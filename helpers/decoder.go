package helpers

import (
	"mime/multipart"
)

// DecodeImage trsanslate image from multipart.FileHeader to []byte
func DecodeImage(image *multipart.FileHeader) ([]byte, error) {
	src, err := image.Open()
	if err != nil {
		return nil, err
	}

	defer src.Close()

	data := make([]byte, image.Size)
	_, err = src.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
