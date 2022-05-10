package utils

import (
	"mime/multipart"
)

// IUtils is utils interface
type IUtils interface {
	AddFileToS3(fileDir string) error
	GetSignedUrl(key string) (string, error)
	AddFileToPublicS3(fileDir string) error
	ConvertDateStringToIndonesia(monthString string) (namaBulan string)
	CreateFolder(path string) error
	CopyFile(src, dst string) error
	CopyFileFromMultipart(filename string, originalfile multipart.File) error
	ConvertDatetimeReadable(dt string) (string, error)
}
