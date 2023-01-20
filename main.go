package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	AWS_S3_REGION = "us-east-1" //Region
)

var s3Client *s3.Client

func fileContent(w http.ResponseWriter, r *http.Request) {

	// Get the bucket name and file key from the URL
	bucketName := r.URL.Query().Get("bucket")
	fileKey := r.URL.Query().Get("file")

	// Retrieve the file from the bucket
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	}
	result, err := s3Client.GetObject(context.TODO(), input)
	if err != nil {
		fmt.Fprintln(w, "No such file is present....")
		log.Println(err)
		return
	}

	// Write the file contents to the response
	_, err = io.Copy(w, result.Body)
	if err != nil {
		fmt.Fprintln(w, "File Contents not written OOPS!!!")
		log.Println(err)
		return
	}
	defer result.Body.Close()
}

func bucketContent(w http.ResponseWriter, r *http.Request) {

	output, err := s3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(string(r.URL.Path[12:])),
	})
	if err != nil {
		fmt.Fprintln(w, "No Bucket with such name....")
		log.Println(err)
		return
	}
	fmt.Fprintf(w, "Total Files are: %d\n", len(output.Contents))
	for _, object := range output.Contents {
		fmt.Fprintln(w, "File Name:", *object.Key)
		fmt.Fprintln(w, "File size:", object.Size)
		fmt.Fprintln(w, "File last modified:", *object.LastModified)
	}
}

func main() {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("personal-puneeth"),
		config.WithRegion("us-east-1"))
	if err != nil {
		log.Printf("error: %v", err)
		return
	}

	s3Client = s3.NewFromConfig(cfg)

	http.HandleFunc("/averlon/s3", fileContent)
	http.HandleFunc("/averlon/s3/", bucketContent)
	http.ListenAndServe(":8080", nil)
}
