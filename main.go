package main

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <file_path>")
		os.Exit(1)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		os.Exit(1)
	}

	filePath := os.Args[1]
	bucketName := viper.GetString("bucket_name")

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	fileExtension := filepath.Ext(filePath)
	fileName := generateRandomString(10)
	newFileName := fileName + fileExtension

	sess, err := session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials(viper.GetString("minio_access_key"), viper.GetString("minio_secret_key"), ""),
		Endpoint:         aws.String("http://s3.myinfra.lol"),
		Region:           aws.String("us-kanc"),
		S3ForcePathStyle: aws.Bool(true),
	})
	if err != nil {
		fmt.Printf("Failed to initialize AWS session: %v\n", err)
		os.Exit(1)
	}

	svc := s3.New(sess)

	_, err = svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(newFileName),
		Body:   file,
	})
	if err != nil {
		fmt.Printf("Failed to upload file: %v\n", err)
		os.Exit(1)
	}

	url := fmt.Sprintf("https://s3.tritan.gg/%s/%s", bucketName, newFileName)
	fmt.Printf("File uploaded successfully. URL: %s\n", url)
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var builder strings.Builder
	for i := 0; i < length; i++ {
		builder.WriteByte(charset[rand.Intn(len(charset))])
	}
	return builder.String()
}
