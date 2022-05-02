package coss

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

/* 腾讯云对象存储（Cloud Object Storage，COS）*/

type TCOS interface {
	PutObject(ctx context.Context, path string, reader io.Reader) error
	GetSignURL(ctx context.Context, path string, expired time.Duration) (string, error)
}

type tcos struct {
	secretID  string
	secretKey string
	client    *cos.Client
}

func NewTCOS(bucketURL, serviceURL, secretID, secretKey string) TCOS {
	bu, _ := url.Parse(bucketURL)
	su, _ := url.Parse(serviceURL)
	baseURL := &cos.BaseURL{
		BucketURL:  bu,
		ServiceURL: su,
	}
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})
	return &tcos{
		secretID:  secretID,
		secretKey: secretKey,
		client:    client,
	}
}

func (s *tcos) PutObject(ctx context.Context, path string, reader io.Reader) error {
	_, err := s.client.Object.Put(ctx, strings.TrimLeft(path, "/"), reader, nil)
	return err
}

func (s *tcos) GetSignURL(ctx context.Context, path string, expired time.Duration) (string, error) {
	u, err := s.client.Object.GetPresignedURL(ctx, http.MethodGet, strings.TrimLeft(path, "/"),
		s.secretID, s.secretKey, expired, nil)
	return u.String(), err
}
