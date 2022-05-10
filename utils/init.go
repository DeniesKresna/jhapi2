package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DeniesKresna/jhapi2/config"
	"github.com/rs/zerolog/log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type PayloadRequest struct {
	Headers map[string]string
	URL     string
	Method  string
	Body    interface{}
}

type Utility struct {
}

func Provide() IUtils {
	return &Utility{}
}

/*
	AddFileToS3 is
*/
func (u *Utility) AddFileToS3(fileDir string) error {
	region := *config.Get().AWS.Region
	secret := *config.Get().AWS.PrivateKey
	public := *config.Get().AWS.PublicKey
	awsBucket := *config.Get().AWS.Bucket
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(public, secret, "")},
	)

	if err != nil {
		log.Error().Err(err)
		return err
	}
	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(fileDir),
		ACL:    aws.String("private"),
		Body:   bytes.NewReader(buffer),
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}

/*
	GetSignedUrl is
*/
func (u *Utility) GetSignedUrl(key string) (string, error) {
	region := *config.Get().AWS.Region
	secret := *config.Get().AWS.PrivateKey
	public := *config.Get().AWS.PublicKey
	bucket := *config.Get().AWS.Bucket
	if bucket == "" || key == "" {
		err := fmt.Errorf("Key / Bucket cannot be empty")
		log.Error().Err(err)
		return "", err
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(public, secret, "")},
	)

	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	urlStr, err := req.Presign(15 * time.Minute)

	if err != nil {
		log.Error().Err(err)
		return "", err
	}

	return urlStr, nil
}

// AddFileToS3 will upload a single file to S3, it will require a pre-built aws session
// and will set file info like content type and encryption on the uploaded file.
func (u *Utility) AddFileToPublicS3(fileDir string) error {
	region := *config.Get().AWS.Region
	secret := *config.Get().AWS.PrivateKey
	public := *config.Get().AWS.PublicKey
	awsBucket := *config.Get().AWS.Bucket
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(public, secret, "")},
	)

	if err != nil {
		log.Error().Err(err)
		return err
	}
	// Open the file for use
	file, err := os.Open(fileDir)
	if err != nil {
		log.Error().Err(err)
		return err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error().Err(err)
		return err
	}

	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
	file.Read(buffer)

	// Config settings: this is where you choose the bucket, filename, content-type etc.
	// of the file you're uploading.
	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(awsBucket),
		Key:    aws.String(fileDir),
		ACL:    aws.String("public-read"),
		Body:   bytes.NewReader(buffer),
	})

	if err != nil {
		log.Error().Err(err)
		return err
	}

	return nil
}

func (u Utility) getObjectNameFromURL(url string, bucketName string) string {
	delimiter := "/"

	objectName := strings.Join(strings.Split(url, bucketName+delimiter)[1:], bucketName)
	return objectName
}

func (u *Utility) ConvertDateStringToIndonesia(monthString string) (namaBulan string) {
	lib := make(map[string]string)
	lib["January"] = "Januari"
	lib["February"] = "Februari"
	lib["March"] = "Maret"
	lib["April"] = "April"
	lib["May"] = "Mei"
	lib["June"] = "Juni"
	lib["July"] = "Juli"
	lib["August"] = "Agustus"
	lib["September"] = "September"
	lib["October"] = "Oktober"
	lib["November"] = "November"
	lib["December"] = "Desember"

	if lib[monthString] == "" {
		namaBulan = monthString
		return
	}
	namaBulan = lib[monthString]
	return
}

//CreateFolder create folder if folder not exists
func (u Utility) CreateFolder(path string) error {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
	}
	return err
}

func (u Utility) CopyFile(src, dst string) error {
	input, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, input, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (u Utility) CopyFileFromMultipart(filename string, originalfile multipart.File) error {
	destFile, err := os.Create(filename)
	if err != nil {
		log.Error().Err(err).Msg("Fail to create temporary file.")
		return err
	}

	if _, err = io.Copy(destFile, originalfile); err != nil {
		log.Error().Err(err).Msg("fail to copy file")
	}
	return err
}

func (u Utility) ConvertDatetimeReadable(dt string) (string, error) {
	var date, err = time.Parse("2006-01-02 15:04:05", dt)
	if err != nil {
		log.Error().Err(err).Interface("data", dt).Msg("fail to convert this date")
		return dt, err
	}

	month := u.ConvertDateStringToIndonesia(date.Month().String())
	return fmt.Sprintf("%v %v %v, pukul %02d:%02d", date.Day(), month,
		date.Year(), date.Hour(), date.Minute()), err
}

func (u Utility) CreateNewHTTPRequest(data PayloadRequest) (*http.Response, error) {
	jsonBody, err := json.Marshal(data.Body)
	if err != nil {
		log.Err(err).Interface("data", data.Body)
		return nil, err
	}

	req, err := http.NewRequest(data.Method, data.URL, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	//Add Additional Headers
	for key, val := range data.Headers {
		req.Header.Set(key, val)
	}

	log.Debug().Interface("HTTP REQUEST", req)
	timeOut := config.Get().Constant.RequestTimeout
	c := &http.Client{Timeout: timeOut * time.Millisecond}
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
