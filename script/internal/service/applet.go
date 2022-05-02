package service

import (
	"context"
	"project/model"
	"time"
)

func (s *Service) GetWechatToken(ctx context.Context) (string, error) {
	return s.redis.Get(ctx, model.KeyWechatToken).Result()
}

func (s *Service) TtlWechatToken(ctx context.Context) (time.Duration, error) {
	return s.redis.TTL(ctx, model.KeyWechatToken).Result()
}

func (s *Service) SetWechatToken(ctx context.Context, tk string, ttl time.Duration) error {
	return s.redis.Set(ctx, model.KeyWechatToken, tk, ttl).Err()
}

func (s *Service) SaveWechatAnalysis(ctx context.Context, data *model.WechatAnalysis) error {
	return s.mysql.WithContext(ctx).FirstOrCreate(data, "ref_date = ?", data.RefDate).Error
}
