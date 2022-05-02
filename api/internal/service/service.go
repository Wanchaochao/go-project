package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
	"project/model"
	"project/pkg/cache"
	"project/pkg/db"
)

type Service struct {
	mysql *gorm.DB
	redis *redis.Client
	//nsq    *nsq.Producer
	single *singleflight.Group
}

type Config struct {
	Mysql db.Mysql
	Redis cache.Redis
	Nsq   struct {
		Producer string
	}
}

func New(cfg *Config) *Service {
	s := &Service{
		mysql: db.NewMysqlDB(&cfg.Mysql),
		redis: cache.NewRedisClient(&cfg.Redis),
		//nsq:    mq.NewNsqProducer(cfg.Nsq.Producer),
		single: &singleflight.Group{},
	}
	return s
}

func (s *Service) WechatToken(ctx context.Context) (string, error) {
	val, err, _ := s.single.Do("WechatToken", func() (any, error) {
		return s.redis.Get(ctx, model.KeyWechatToken).Result()
	})
	return val.(string), err
}
