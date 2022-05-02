package service

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	uuid "github.com/satori/go.uuid"
	"project/api/internal/proto"
	"project/model"
	"project/pkg/logger"
	"time"
)

func (s *Service) SaveUser(ctx context.Context, data *model.User) (int, error) {
	err := s.mysql.WithContext(ctx).FirstOrCreate(data, "openid = ?", data.Openid).Error
	if err != nil || data.ID == 0 {
		return 0, err
	}
	b, _ := json.Marshal(data)
	if err := s.redis.Set(ctx, model.UserInfoKey(data.ID), b, time.Hour); err != nil {
		logger.FromContext(ctx).Error("redis.Set error", nil, err)
	}
	return data.ID, nil
}

func (s *Service) SetUserToken(ctx context.Context, data *proto.UserToken) (string, error) {
	h := sha1.New()
	h.Write([]byte(data.Openid))
	h.Write(uuid.NewV4().Bytes())
	h.Write([]byte(data.SessionKey))
	token := base32.StdEncoding.EncodeToString(h.Sum(nil))
	b, _ := json.Marshal(data)
	err := s.redis.Set(ctx, model.UserTokenKey(token), b, time.Hour).Err()
	return token, err
}

func (s *Service) GetUserToken(ctx context.Context, token string) (*proto.UserToken, error) {
	key := model.UserTokenKey(token)
	pipe := s.redis.Pipeline()
	cmd1 := pipe.Expire(ctx, key, time.Hour)
	cmd2 := pipe.Get(ctx, key)
	if _, err := pipe.Exec(ctx); err != nil && err != redis.Nil {
		return nil, err
	}
	var account proto.UserToken
	if ok, err := cmd1.Result(); !ok {
		return &account, err
	}
	b, err := cmd2.Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}
	//b, err := s.redis.GetEx(ctx, key, time.Hour).Bytes() // redis >= 6.2.0 只需一个原子性命令 GETEX
	if len(b) > 0 {
		err = json.Unmarshal(b, &account)
	}
	return &account, err
}

func (s *Service) FindUserByID(ctx context.Context, id int) (*model.User, error) {
	key := model.UserInfoKey(id)
	b, err := s.redis.Get(ctx, key).Bytes()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	var res model.User
	if len(b) > 0 {
		err = json.Unmarshal(b, &res)
		return &res, err
	}

	err = s.mysql.WithContext(ctx).Where("id = ?", id).Take(&res).Error
	if err != nil {
		return nil, err
	}
	b, _ = json.Marshal(res)
	if err := s.redis.Set(ctx, key, b, time.Hour); err != nil {
		logger.FromContext(ctx).Error("redis.Set error", key, err)
	}
	return &res, nil
}

func (s *Service) UpdateUser(ctx context.Context, data *model.User) error {
	opt := s.mysql.WithContext(ctx).Updates(data) // gorm根据ID更新指定非零值字段
	if opt.Error != nil {
		return opt.Error
	}
	if opt.RowsAffected > 0 {
		return s.redis.Del(ctx, model.UserInfoKey(data.ID)).Err()
	}
	return nil
}
