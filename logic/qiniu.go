package logic

import (
	"context"
	"manage/setting"
	"mime/multipart"

	"github.com/qiniu/go-sdk/v7/auth/qbox"

	"github.com/qiniu/go-sdk/v7/storage"
)

func UpLoad(file multipart.File, fileSize int64) (url string, err error) {
	putPolicy := storage.PutPolicy{
		Scope: setting.Conf.Bucket,
	}
	mac := qbox.NewMac(setting.Conf.AccessKey, setting.Conf.SecretKey)
	upToken := putPolicy.UploadToken(mac)

	cfg := &storage.Config{
		Zone:          &storage.ZoneHuanan,
		Region:        nil,
		UseHTTPS:      setting.Conf.UseHttps,
		UseCdnDomains: setting.Conf.UseCdnDomains,
	}

	putExtra := storage.PutExtra{}

	formUpLoader := storage.NewFormUploader(cfg)

	ret := storage.PutRet{}
	//自定义key
	//key := "image/" + setting.Conf.Prefix + file.FileName

	err = formUpLoader.PutWithoutKey(context.Background(), &ret, upToken, file, fileSize, &putExtra)
	url = setting.Conf.Domain + ret.Key
	return
}
