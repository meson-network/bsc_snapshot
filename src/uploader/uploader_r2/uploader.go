package uploader_r2

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gosuri/uiprogress"
	"github.com/meson-network/bsc-data-file-utils/src/common/custom_reader"
)

func GenR2Client(accountId string, accessKeyId string, accessKeySecret string) (*s3.Client, error) {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId),
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(r2Resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)
	return client, nil
}

func UploadFile(client *s3.Client, bucketName string, additional_path string, fileName string, localFilePath string, localFileMd5 string, bar *uiprogress.Bar) error {

	keyInRemote := fileName
	if additional_path != "" {
		keyInRemote = additional_path + "/" + fileName
	}

	if localFileMd5 != "" {
		// if provide md5, check remote file md5 first
		// get fileInfo from bucket

		// ignore err, if info==nil, just reupload file
		info, _ := client.GetObject(context.Background(), &s3.GetObjectInput{
			Bucket: &bucketName,
			Key:    aws.String(keyInRemote),
		})

		// if exist check md5
		if info != nil && info.ETag != nil {
			remoteMd5 := *info.ETag
			localMd5 := "\"" + localFileMd5 + "\""

			if remoteMd5 == localMd5 {
				// if same file, upload success
				// bar.PrintBar(100)
				return nil
			}
		}
	}

	//  upload new one
	fileInfo, err := os.Stat(localFilePath)
	if err != nil {
		return err
	}

	uploadFile, err := os.OpenFile(localFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer uploadFile.Close()

	// use custom reader to show upload progress
	reader := &custom_reader.CustomReader{
		Reader: uploadFile,
		Size: fileInfo.Size(),
		Bar:  bar,
	}

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(keyInRemote),
		Body:          reader,
		ContentLength: fileInfo.Size(),
	})
	if err != nil {
		return err
	}

	return nil
}
