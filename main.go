package main

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func mainConnect() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "minio"
	secretAccessKey := "12345678"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("error!")
		log.Fatalln(err)
	}

	log.Printf("client show up!")
	log.Printf("%#v\n", minioClient) // minioClient is now set up
}
