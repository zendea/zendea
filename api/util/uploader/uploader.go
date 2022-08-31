package uploader

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/go-resty/resty/v2"
	"github.com/spf13/viper"

	"zendea/util"
	"zendea/util/log"
	"zendea/util/urls"
)

var (
	aliyun = NewAliyun()
	local  = NewLocal()
)

func PutImage(data []byte) (string, error) {
	if viper.GetString("uploader.enable") == "aliyun" {
		return aliyun.PutImage(data)
	} else {
		return local.PutImage(data)
	}
}

func PutObject(key string, data []byte) (string, error) {
	if viper.GetString("uploader.enable") == "aliyun" {
		return aliyun.PutObject(key, data)
	} else {
		return local.PutObject(key, data)
	}
}

func CopyImage(originUrl string) (string, error) {
	if viper.GetString("uploader.enable") == "aliyun" {
		return aliyun.CopyImage(originUrl)
	} else {
		return local.CopyImage(originUrl)
	}
}

func NewAliyun() *aliyunOssUploader {
	return &aliyunOssUploader{
		once:   sync.Once{},
		bucket: nil,
	}
}

func NewLocal() *localUploader {
	return &localUploader{}
}

type uploader interface {
	PutImage(data []byte) (string, error)
	PutObject(key string, data []byte) (string, error)
	CopyImage(originUrl string) (string, error)
}

// 阿里云oss
type aliyunOssUploader struct {
	once   sync.Once
	bucket *oss.Bucket
}

func (aliyun *aliyunOssUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return aliyun.PutObject(key, data)
}

func (aliyun *aliyunOssUploader) PutObject(key string, data []byte) (string, error) {
	bucket := aliyun.getBucket()
	if err := bucket.PutObject(key, bytes.NewReader(data)); err != nil {
		return "", err
	}

	return urls.UrlJoin(viper.GetString("uploader.oss.host"), key), nil
}

func (aliyun *aliyunOssUploader) CopyImage(originUrl string) (string, error) {
	data, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return aliyun.PutImage(data)
}

func (aliyun *aliyunOssUploader) getBucket() *oss.Bucket {
	aliyun.once.Do(func() {
		if client, err := oss.New(viper.GetString("uploader.oss.endpoint"), viper.GetString("uploader.oss.access_id"), viper.GetString("uploader.oss.access_secret")); err != nil {
			log.Error(err.Error())
		} else if aliyun.bucket, err = client.Bucket(viper.GetString("uploader.oss.bucket")); err != nil {
			log.Error(err.Error())
		}
	})
	return aliyun.bucket
}

// 本地文件系统
type localUploader struct{}

func (local *localUploader) PutImage(data []byte) (string, error) {
	key := generateImageKey(data)
	return local.PutObject(key, data)
}

func (local *localUploader) PutObject(key string, data []byte) (string, error) {
	if err := os.MkdirAll("/", os.ModeDir); err != nil {
		return "", err
	}

	filename := filepath.Join(viper.GetString("uploader.local.path"), key)
	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		return "", err
	}
	if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
		return "", err
	}
	return urls.UrlJoin(viper.GetString("uploader.local.host"), key), nil
}

func (local *localUploader) CopyImage(originUrl string) (string, error) {
	data, err := download(originUrl)
	if err != nil {
		return "", err
	}
	return local.PutImage(data)
}

// generateKey 生成图片Key
func generateImageKey(data []byte) string {
	md5 := util.MD5Bytes(data)
	return filepath.Join("images", util.TimeFormat(time.Now(), "2006/01/02/"), md5+".jpg")
}

func download(url string) ([]byte, error) {
	rsp, err := resty.New().R().Get(url)
	if err != nil {
		return nil, err
	}
	return rsp.Body(), nil
}
