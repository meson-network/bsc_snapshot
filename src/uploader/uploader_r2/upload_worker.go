package uploader_r2

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vbauerster/mpb/v8"

	"github.com/meson-network/bsc-snapshot/src/utils/custom_reader"
)

type UploaderWorker struct {
	client          *s3.Client
	bucketName      string
	additional_path string
	fileName        string
	localFileMd5    string
	bar             *mpb.Bar
}

func NewUploadWorker(client *s3.Client, bucketName string, additional_path string,
	fileName string, localFileMd5 string, bar *mpb.Bar) *UploaderWorker {

	return &UploaderWorker{
		client: client, bucketName: bucketName, additional_path: additional_path,
		fileName: fileName, localFileMd5: localFileMd5, bar: bar,
	}
}

func (u *UploaderWorker) UploadFile(localFilePath string) error {

	keyInRemote := u.fileName
	if u.additional_path != "" {
		keyInRemote = u.additional_path + "/" + u.fileName
	}

	if u.localFileMd5 != "" {
		if exists := validateRemoteChunk(u.client, u.bucketName, keyInRemote, u.localFileMd5); exists {
			return nil
		}
	}

	// upload new one
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
		Size:   fileInfo.Size(),
		Bar:    u.bar,
	}

	_, err = u.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:        aws.String(u.bucketName),
		Key:           aws.String(keyInRemote),
		Body:          reader,
		ContentLength: fileInfo.Size(),
	})
	if err != nil {
		return err
	}

	return nil
}

func validateRemoteChunk(client *s3.Client, bucketName string, keyInRemote string,
	localFileMd5 string) bool {
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

		if strings.EqualFold(remoteMd5, localMd5) {
			// if same file, upload success
			return true
		}
	}

	return false
}
