package service

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"project/cms/internal/acl"
	"project/cms/internal/proto"
	"project/model"
	"project/pkg/logger"
	"time"
)

const (
	tokenTTL = 6 * time.Hour
	minTTL   = 5 * time.Hour
)

func (s *Service) SetAdminToken(ctx context.Context, data *acl.AdminToken) (string, error) {
	ssoKey := model.AdminSSOKey(data.ID)
	oldToken, err := s.redis.Get(ctx, ssoKey).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(data.Username))
	h.Write(uuid.NewV4().Bytes())
	newToken := base32.StdEncoding.EncodeToString(h.Sum(nil))
	b, _ := json.Marshal(data)
	tx := s.redis.TxPipeline()
	if oldToken != "" {
		tx.Del(ctx, model.AdminTokenKey(oldToken))
	}
	tx.Set(ctx, model.AdminTokenKey(newToken), b, tokenTTL)
	tx.Set(ctx, ssoKey, newToken, tokenTTL)
	_, err = tx.Exec(ctx)
	return newToken, err
}

func (s *Service) GetAdminToken(ctx context.Context, token string) (*acl.AdminToken, error) {
	key := model.AdminTokenKey(token)
	pipe := s.redis.Pipeline()
	cmd1 := pipe.TTL(ctx, key)
	cmd2 := pipe.Get(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}
	ttl, err1 := cmd1.Result()
	if err1 != nil {
		return nil, err1
	}
	b, err2 := cmd2.Bytes()
	if err2 != nil && err2 != redis.Nil {
		return nil, err2
	}
	var result acl.AdminToken
	if (ttl > time.Second || ttl == -1) && len(b) > 0 {
		_ = json.Unmarshal(b, &result)
		if ttl < minTTL && result.ID > 0 {
			tx := s.redis.TxPipeline()
			tx.Expire(ctx, model.AdminSSOKey(result.ID), tokenTTL)
			tx.Expire(ctx, key, tokenTTL)
			if _, err := tx.Exec(ctx); err != nil {
				logger.FromContext(ctx).Error("redis.TxPipeline.Exec error", nil, err)
			}
		}
	}
	return &result, nil
}

func (s *Service) DelAdminToken(ctx context.Context, token string) error {
	return s.redis.Del(ctx, model.AdminTokenKey(token)).Err()
}

func (s *Service) FindAdminUserByName(ctx context.Context, name string) (*acl.AdminUser, error) {
	var user acl.AdminUser
	err := s.mysql.WithContext(ctx).Where("username = ?", name).
		Preload("AdminRole").Take(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func (s *Service) FindAdminUserByID(ctx context.Context, id int) (*acl.AdminUser, error) {
	var user acl.AdminUser
	err := s.mysql.WithContext(ctx).Where("id = ? AND username <> ?", id, acl.Super).Take(&user).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &user, nil
}

func (s *Service) UpdateAdminUser(ctx context.Context, data *acl.AdminUser) error {
	opt := s.mysql.WithContext(ctx).Updates(data) // gorm自动根据ID更新非零字段
	if opt.Error != nil {
		return opt.Error
	}
	if opt.RowsAffected > 0 && (data.Status == model.StatusOff || data.RoleID != 0) {
		return s.LogoutAdminUser(ctx, data.ID) //禁用或变更角色强制退出登录
	}
	return nil
}

func (s *Service) LogoutAdminUser(ctx context.Context, id int) error {
	ssoKey := model.AdminSSOKey(id)
	token, err := s.redis.Get(ctx, ssoKey).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if token != "" {
		return s.redis.Del(ctx, ssoKey, model.AdminTokenKey(token)).Err()
	}
	return nil
}

func (s *Service) PaginateAdminRole(ctx context.Context,
	p *proto.ListArgs) (total int64, list []*acl.AdminRole, err error) {
	query := s.mysql.WithContext(ctx).Model(&acl.AdminRole{})
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) AllAdminRole(ctx context.Context) ([]*acl.AdminRole, error) {
	var data []*acl.AdminRole
	err := s.mysql.WithContext(ctx).Select("id", "name").
		Order("id").Find(&data).Error
	return data, err
}

func (s *Service) FindAdminRoleByID(ctx context.Context, id int) (*acl.AdminRole, error) {
	var role acl.AdminRole
	err := s.mysql.WithContext(ctx).Where("id = ?", id).Take(&role).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &role, nil
}

func (s *Service) CreateAdminRole(ctx context.Context, data *acl.AdminRole) error {
	return s.mysql.WithContext(ctx).Create(data).Error
}

func (s *Service) UpdateAdminRole(ctx context.Context, data *acl.AdminRole) error {
	return s.mysql.WithContext(ctx).Updates(data).Error // gorm默认根据ID更新非零值字段
}

func (s *Service) PaginateAdminUser(ctx context.Context,
	p *proto.AdminUserListArgs) (total int64, list []*acl.AdminUser, err error) {
	query := s.mysql.WithContext(ctx).
		Model(&acl.AdminUser{}).Where("username <> ?", acl.Super)
	if p.RoleID > 0 {
		query = query.Where("role_id = ?", p.RoleID)
	}
	if p.Username != "" {
		query = query.Where("username LIKE ?", p.Username+"%")
	}
	if p.Status != 0 {
		query = query.Where("status = ?", p.Status)
	}
	err = query.Count(&total).Error
	offset := p.Size * (p.Page - 1)
	if err != nil || total == 0 || offset >= int(total) {
		return
	}
	err = query.Order("id DESC").Limit(p.Size).Offset(offset).Find(&list).Error
	return
}

func (s *Service) CreateAdminUser(ctx context.Context, data *acl.AdminUser) (bool, error) {
	opt := s.mysql.WithContext(ctx).FirstOrCreate(data, "username = ?", data.Username)
	return opt.RowsAffected > 0, opt.Error
}
