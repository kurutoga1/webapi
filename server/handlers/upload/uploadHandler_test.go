package upload_test

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"webapi/server/handlers/upload"
	uf "webapi/utils/file"
)

var (
	uploadFile string
)

func init() {
	uploadFile = "200MB.txt"
	if !uf.FileExists(uploadFile) {
		err := uf.CreateSpecifiedFile(uploadFile, 200000)
		if err != nil {
			panic(err.Error())
		}
	}
}

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestUploadHandler(t *testing.T) {
	file, err := os.Open(uploadFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		log.Fatal(err)
	}

	io.Copy(part, file)
	writer.Close()
	request, err := http.NewRequest("POST", "/upload", body)
	if err != nil {
		log.Fatal(err)
	}

	request.Header.Add("Content-Type", writer.FormDataContentType())

	response := httptest.NewRecorder()

	upload.UploadHandler(response, request)

	// アップロードされているか
	uploadedPath := filepath.Join("fileserver", "upload", uploadFile)
	if !uf.FileExists(uploadedPath) {
		t.Errorf("uploadedPath(%v) is not exist.", uploadedPath)
	}

	t.Cleanup(func() {
		tearDown()
	})

}
