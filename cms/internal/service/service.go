package service

import (
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"project/pkg/cache"
	"project/pkg/db"
)

type Service struct {
	mysql *gorm.DB
	redis *redis.Client
	//nsq   *nsq.Producer
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
		//nsq:   mq.NewNsqProducer(cfg.Nsq.Producer),
	}
	return s
}
