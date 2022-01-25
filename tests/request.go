package tests

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

// MainRequest web,cliのどちらともからこのリクエストをプログラムサーバへ送信する。
// toパラメータは"web"か"cli"のどちらかを指定する。
// 主にtestで使用する。
func MainRequest(uploadFile, proName, parameta, from string) (*httptest.ResponseRecorder, *http.Request, error) {
	if from != "web" && from != "cli" {
		return nil, nil, fmt.Errorf("MainRequest: from(%v) is not valid. only web or cli.", from)
	}

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		if from == "web" {
			err := form.WriteField("proName", proName)
			if err != nil {
				panic(err.Error())
			}
		}
		err := form.WriteField("parameta", parameta)
		if err != nil {
			panic(err.Error())
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

	var r *http.Request
	var err error

	if from == "cli" {
		r, err = http.NewRequest(http.MethodPost, "/pro/"+proName, pr)
	} else if from == "web" {
		r, err = http.NewRequest(http.MethodPost, "/user/exec", pr)
	}

	if err != nil {
		panic(err.Error())
	}

	r.Header.Set("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()

	return w, r, nil
}
