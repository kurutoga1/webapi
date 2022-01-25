package http

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// GetPostRequestWithFileAndFields
// プログラムサーバで実行するために必要なのはfile(multi-part)といくつかのパラメータ
// それらを一気にPOSTで送るリクエストを返す。
func GetPostRequestWithFileAndFields(uploadFile, url string, fields map[string]string) (*http.Request, error) {

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		for field, value := range fields {
			err := form.WriteField(field, value)
			if err != nil {
				panic(err.Error())
			}
		}

		file, err := os.Open(uploadFile)
		if err != nil {
			panic(err.Error())
		}
		w, err := form.CreateFormFile("file", filepath.Base(uploadFile))
		if err != nil {
			panic(err.Error())
		}
		_, err = io.Copy(w, file)
		if err != nil {
			panic(err.Error())
		}
		err = form.Close()
		if err != nil {
			panic(err.Error())
		}
	}()

	r, err := http.NewRequest(http.MethodPost, url, pr)

	if err != nil {
		return nil, fmt.Errorf("GetPostRequestWithFileAndFields: %v", err)
	}

	r.Header.Set("Content-Type", form.FormDataContentType())

	return r, nil
}
