package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	ctx := context.Background()
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
		log.Fatalln(err)
	}

	log.Printf("client show up!")
	log.Printf("%#v\n", minioClient) // minioClient is now set up

	// Make a new bucket called mybucket
	bucketName := "mybucket"
	location := "us-east-1"
	createBucket(minioClient, ctx, bucketName, location)

	// Upload the zip file
	objectName := "COEXODDS.png"
	filePath := "./COEXODDS.png"
	// contentType := "application/zip" //zip file
	contentType := "" //for image file
	uploadFileToBucket(minioClient, ctx, bucketName, objectName, filePath, contentType)

	// expiryTime := time.Second * 24 * 60 * 60 * 7
	// getURLObject(minioClient, ctx, bucketName, objectName, expiryTime)
	policy, err := minioClient.GetBucketPolicy(ctx, bucketName)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(fmt.Sprintf("policy : %s", policy))

	url := getImageURL(endpoint, bucketName, objectName)
	log.Println("URL: ", url)
}

func createBucket(minioClient *minio.Client, ctx context.Context, bucketName string, location string) {
	err := minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Printf("Fatal!!")
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
}

func getImageURL(endpoint string, bucketName string, objectName string) string {
	return fmt.Sprintf("%s/%s/%s", endpoint, bucketName, objectName)
}

func uploadFileToBucket(minioClient *minio.Client, ctx context.Context, bucketName string, objectName string, filePath string, contentType string) {

	userMetaData := map[string]string{"x-amz-acl": "public-read"} // make it public
	info, err := minioClient.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType, UserMetadata: userMetaData})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Successfully uploaded %s of size %d\n", objectName, info.Size)
	log.Printf("info: %s", info.ETag)
}

func getURLObject(minioClient *minio.Client, ctx context.Context, bucketName string, objectName string, expiryTime time.Duration) {
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf("attachment; filename=\"%s\"", objectName))

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), bucketName, objectName, expiryTime, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)
}
