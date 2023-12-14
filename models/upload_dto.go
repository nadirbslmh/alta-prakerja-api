package models

import "mime/multipart"

type UploadDTO struct {
	File           *multipart.FileHeader
	UploadFormData UploadForm
}
