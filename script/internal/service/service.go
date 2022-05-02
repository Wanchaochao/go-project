package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/nsqio/go-nsq"
	"gorm.io/gorm"
	"project/pkg/cache"
	"project/pkg/db"
	"project/pkg/mq"
)

type Service struct {
	mysql    *gorm.DB
	redis    *redis.Client
	producer *nsq.Producer
}

type Option func(*Service)

func NewMysql(cfg *db.Mysql) Option {
	return func(s *Service) {
		if s.mysql == nil {
			s.mysql = db.NewMysqlDB(cfg)
		}
	}
}

func NewRedis(cfg *cache.Redis) Option {
	return func(s *Service) {
		if s.redis == nil {
			s.redis = cache.NewRedisClient(cfg)
		}
	}
}

func NewProducer(addr string) Option {
	return func(s *Service) {
		if s.producer == nil {
			s.producer = mq.NewNsqProducer(addr)
		}
	}
}

func NewService(options ...Option) *Service {
	s := &Service{}
	for _, opt := range options {
		opt(s)
	}
	return s
}
