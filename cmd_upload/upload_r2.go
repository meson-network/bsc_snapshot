package cmd_upload

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/meson-network/bsc-data-file-utils/src/split_file"
	"github.com/urfave/cli/v2"
)

func Upload_r2(clictx *cli.Context) error {

	//read params
	originDir := clictx.String("dir")
	// thread := clictx.Int("thread")
	bucketName := clictx.String("bucket_name")
	additional_path := clictx.String("additional_path")
	accountId := clictx.String("account_id")
	accessKeyId := clictx.String("access_key_id")
	accessKeySecret := clictx.String("access_key_secret")

	// read json from originDir
	configFilePath := filepath.Join(originDir, split_file.FILES_CONFIG_JSON_NAME)
	jsonContent, err := os.ReadFile(configFilePath)
	if err != nil {
		return err
	}
	fileConfig := split_file.FileSplitConfig{}
	err = json.Unmarshal(jsonContent, &fileConfig)
	if err != nil {
		return err
	}

	client, err := genR2Client(accountId, accessKeyId, accessKeySecret)
	if err != nil {
		return err
	}

	fileList := fileConfig.ChunkedFileList
	errorFiles := []*split_file.ChunkedFileInfo{}
	for index, v := range fileList {
		localFilePath := filepath.Join(originDir, v.FileName)
		err := uploadFile(client, bucketName, additional_path, v.FileName, localFilePath, v.Md5)
		if err != nil {
			log.Println(index+1, "/", len(fileList), v.FileName, "upload err:", err)
			errorFiles = append(errorFiles, &v)
		} else {
			log.Println(index+1, "/", len(fileList), v.FileName, "uploaded")
		}
	}

	if len(errorFiles) > 0 {
		log.Println("the following files upload failed, please try again:")
		for _, v := range errorFiles {
			log.Println(v.FileName)
		}
		return errors.New("upload error")
	}

	//upload config
	localFilePath := filepath.Join(originDir, split_file.FILES_CONFIG_JSON_NAME)
	err = uploadFile(client, bucketName, additional_path, split_file.FILES_CONFIG_JSON_NAME, localFilePath, "")
	if err != nil {
		log.Println(split_file.FILES_CONFIG_JSON_NAME, "upload err:", err)
	} else {
		log.Println(split_file.FILES_CONFIG_JSON_NAME, "uploaded")
	}

	return nil

}

func genR2Client(accountId string, accessKeyId string, accessKeySecret string) (*s3.Client, error) {
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

func uploadFile(client *s3.Client, bucketName string, additional_path string, fileName string, localFilePath string, localFileMd5 string) error {
	// get fileInfo from bucket
	// ignore err, if info==nil, just reupload file
	keyInRemote := fileName
	if additional_path != "" {
		keyInRemote = additional_path + "/" + fileName
	}

	if localFileMd5 != "" {
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
				log.Println("file already uploaded")
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

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(keyInRemote),
		Body:          uploadFile,
		ContentLength: fileInfo.Size(),
	})
	if err != nil {
		return err
	}

	return nil
}
