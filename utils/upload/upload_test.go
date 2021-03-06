package upload_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"webapi/pro/router"
	"webapi/test"
	"webapi/utils/file"
	u "webapi/utils/upload"
)

var (
	currentDir string
	uploadFile string
	addrs      []string
)

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	addrs, _, err = test.GetStartedServers(router.New(), 1)
	if err != nil {
		log.Fatalln(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestUpload(t *testing.T) {
	err := file.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	uploader := u.NewUploader()
	err = uploader.Upload(addrs[0]+"/upload", uploadFile)
	if err != nil {
		t.Errorf("err from Upload: %v \n", err.Error())
	}

	uploadedFilePath := filepath.Join(currentDir, "fileserver", "upload", uploadFile)
	if !file.FileExists(uploadedFilePath) {
		t.Errorf("uploadedPath(%v) is not exists. \n", uploadedFilePath)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
