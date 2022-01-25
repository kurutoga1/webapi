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
// いくつかのパラメータはfieldsパラメータにmapで渡す。
// それらを一気にPOSTで送るリクエストを返す。
func GetPostRequestWithFileAndFields(uploadFile, url string, fields map[string]string) (r *http.Request, err error) {

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer func(pw *io.PipeWriter) {
			err = pw.Close()
		}(pw)

		// フォームにフィールドを追加
		for field, value := range fields {
			err = form.WriteField(field, value)
		}

		var file *os.File
		file, err = os.Open(uploadFile)

		var w io.Writer
		w, err = form.CreateFormFile("file", filepath.Base(uploadFile))
		_, err = io.Copy(w, file)
		err = form.Close()

		if err != nil {
			err = fmt.Errorf("GetPostRequestWithFileAndFields: %v", err)
		}
	}()

	r, err = http.NewRequest(http.MethodPost, url, pr)

	if err != nil {
		return nil, fmt.Errorf("GetPostRequestWithFileAndFields: %v", err)
	}

	r.Header.Set("Content-Type", form.FormDataContentType())

	return r, nil
}
