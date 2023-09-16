package test

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"gogormlearn/inits"
)

func TestMinio() {
	inits.InitMinio()
	info, err := inits.MinioClient.FPutObject(context.TODO(), "test", "tt.php", "/Users/weifeng/tt.php", minio.PutObjectOptions{ContentType: "application:text"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("上传结果---> %+v", info)

	err = inits.MinioClient.FGetObject(context.TODO(), "test", "tt.php", "/Users/weifeng/tt.phpbak", minio.GetObjectOptions{})
	if err != nil {
		panic(err)
	}
}
