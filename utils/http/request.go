package http

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

type PostGetter interface {
	// GetPostRequest
	// サーバのURL,ローカルに作成済みのアップロードファイルとformに入れるパラメータのmapを受け取り
	// それらをPOSTで送信するリクエストを返す。またサーバに送信はしていない。
	GetPostRequest(url, uploadFile string, fields map[string]string) (r *http.Request, err error)
}

func NewPostGetter() PostGetter {
	return &mainPoster{}
}

type mainPoster struct{}

// GetPostRequest
// サーバのURL,ローカルに作成済みのアップロードファイルとformに入れるパラメータのmapを受け取り
// それらをPOSTで送信するリクエストを返す。またサーバに送信はしていない。
func (m *mainPoster) GetPostRequest(url, uploadFile string, fields map[string]string) (r *http.Request, err error) {
	file, err := os.Open(uploadFile)
	if err != nil {
		return nil, fmt.Errorf("MainPost: %v", err)
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			err = fmt.Errorf("MainPost: %v", err)
		}
	}(file)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// mapの中身をformに入れる
	for key, value := range fields {
		err := writer.WriteField(key, value)
		if err != nil {
			return nil, fmt.Errorf("MainPost: %v", err)
		}
	}

	// ファイルをformに入れる
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("MainPost: %v", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("MainPost: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("MainPost: %v", err)
	}

	request, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, fmt.Errorf("MainPost: %v", err)
	}

	request.Header.Set("Content-Type", writer.FormDataContentType())

	return request, nil
}
