package user_test

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"webapi/server/handlers/user"
	"webapi/utils"
)

var (
	uploadFile string
)

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestUserTopHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/user/top", nil)
	response := httptest.NewRecorder()

	user.UserTopHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	programName := "convertToJson"
	if !strings.Contains(response.Body.String(), programName) {
		t.Errorf("html doesn't have %v", programName)
	}
}

func TestPrepareExecHandler(t *testing.T) {
	form := url.Values{
		"proName": []string{"convertToJson"},
	}
	req, err := http.NewRequest("POST", "/user/perpareExec", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	user.PrepareExecHandler(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}
}

func TestExecHandler(t *testing.T) {
	uploadFile = "uploadfile"
	err := utils.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		err = form.WriteField("proName", "convertToJson")
		if err != nil {
			panic(err.Error())
		}

		err = form.WriteField("parameta", "dummyParameta")
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

	req, err := http.NewRequest(http.MethodPost, "/user/exec", pr)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

	response := httptest.NewRecorder()

	user.ExecHandler(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	expected := "<p>Result: ok</p>"
	if !strings.Contains(response.Body.String(), expected) {
		t.Errorf("response.Body doesn't contain %v", expected)
	}

	tearDown()
}
