package storage

import "mime/multipart"

type RequestBody struct {
	File     multipart.File `json:"file" validate:"required"`
	FilePath string         `json:"filePath" validate:"required"`
}

type ResponseBody struct {
	Url         string `json:"url"`
	FilePath    string `json:"filePath"`
	FileName    string `json:"fileName"`
	ContentType string `json:"contentType"`
}
