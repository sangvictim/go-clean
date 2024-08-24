package storage

import (
	"fmt"
	apiResponse "go-clean/utils/response"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type UploadController struct {
	Log   *logrus.Logger
	Viper *viper.Viper
}

func NewStorageController(log *logrus.Logger, viper *viper.Viper) *UploadController {
	return &UploadController{
		Log:   log,
		Viper: viper,
	}
}

func (c *UploadController) storageInit() *session.Session {
	sess, _ := session.NewSession(&aws.Config{
		DisableSSL:       aws.Bool(c.Viper.GetBool("s3.disableSSL")),
		Endpoint:         aws.String(c.Viper.GetString("s3.endpoint")),
		S3ForcePathStyle: aws.Bool(c.Viper.GetBool("s3.forcePathStyle")),
		Region:           aws.String(c.Viper.GetString("s3.region")),
		Credentials: credentials.NewStaticCredentials(
			c.Viper.GetString("s3.accessKey"),
			c.Viper.GetString("s3.secretKey"),
			"",
		),
	})

	return sess
}

func (c *UploadController) UploadFile(ctx echo.Context) error {
	// get file from request
	file, err := ctx.FormFile("file")
	filePath := ctx.FormValue("filePath")
	if err != nil {
		return echo.ErrBadRequest
	}

	// check file mime type
	format := file.Header.Get("Content-Type")
	if err := c.checkMimeType(format); err != nil {
		return err
	}

	//open the uploaded file
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// converto to webp
	// fileContent, _ := io.ReadAll(src)
	// imgConvert, _ := bimg.NewImage(fileContent).Convert(bimg.WEBP)

	// Define the S3 key (file name in the bucket)
	key := filePath + "/" + strings.ReplaceAll(file.Filename, " ", "_")

	// initiate session
	sess := c.storageInit()
	client := s3.New(sess)

	_, err = client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(c.Viper.GetString("s3.bucket")),
		Key:    aws.String(key),
		Body:   src, // bytes.NewReader(imgConvert)
		ACL:    aws.String("public-read"),
	})

	if err != nil {
		fmt.Println("Failed Upload: ", err)
	}

	return apiResponse.ResponseJson(ctx, http.StatusCreated, apiResponse.Response{
		Message: "Success Uploaded",
		Data: ResponseBody{
			Url:         c.Viper.GetString("s3.cdn") + key,
			FilePath:    key,
			FileName:    file.Filename,
			ContentType: format,
		},
	})
}

func (c *UploadController) GetFile(ctx echo.Context) error {
	key := ctx.Param("key")
	sess := c.storageInit()
	client := s3.New(sess)

	// Get the object
	object, err := client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(c.Viper.GetString("s3.bucket")),
		Key:    aws.String(key),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchKey:
				return echo.NewHTTPError(http.StatusNotFound, "File not found")
			default:
				return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
			}
		}
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}
	defer object.Body.Close()

	// get file extension & MIME type
	ext := filepath.Ext(key)
	contentType := mime.TypeByExtension(ext)

	// set header
	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=\"%s\"", filepath.Base(key)))
	ctx.Response().Header().Set("Content-Type", contentType)
	ctx.Response().Header().Set("Cache-control", "public, max-age=3600")

	// stream this file, not download
	_, err = io.Copy(ctx.Response().Writer, object.Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Internal server error")
	}

	return nil
}

func (c *UploadController) checkMimeType(format string) error {
	allowedTypes := []string{
		"image/png",
		"image/jpg",
		"image/jpeg",
		"image/pdf",
	}

	if !contains(allowedTypes, format) {
		return echo.NewHTTPError(http.StatusBadRequest, "File type not allowed")
	}

	return nil
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
