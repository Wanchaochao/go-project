package coss

import (
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"log"
	"strings"
)

/*阿里云对象存储OSS（Object Storage Service）*/

type AliOSS interface {
	PutObject(path string, reader io.Reader) error
	GetSignURL(path string, expireSeconds int64) (string, error)
}

type alioss struct {
	bucket *oss.Bucket
}

func NewAliOSS(endpoint, keyID, keySecret, bucketName string) AliOSS {
	client, err := oss.New(endpoint, keyID, keySecret)
	if err != nil {
		log.Fatal(err)
	}
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		log.Fatal(err)
	}
	return &alioss{bucket: bucket}
}

func (s *alioss) PutObject(path string, reader io.Reader) error {
	return s.bucket.PutObject(strings.TrimLeft(path, "/"), reader)
}

func (s *alioss) GetSignURL(path string, expireSeconds int64) (string, error) {
	return s.bucket.SignURL(strings.TrimLeft(path, "/"), oss.HTTPGet, expireSeconds)
}
