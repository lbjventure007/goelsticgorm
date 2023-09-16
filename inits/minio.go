package inits

import (
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"sync"
)

var MinioClient *minio.Client
var once1 sync.Once

func InitMinio() {
	once1.Do(func() {
		client, err := minio.New("127.0.0.1:9000", &minio.Options{

			Creds: credentials.NewStaticV4("7mmIi11ZiwkUefMT6yDv",
				"Wnt33apGw2Igw5xYOQmlUydECWJwMpkcxUq5WR9e", ""),
			Secure: false,
		})
		if err != nil {
			panic(err)
		}
		MinioClient = client
	})

}
