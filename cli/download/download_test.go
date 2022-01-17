package download_test

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"webapi/cli/download"
	u "webapi/cli/upload"
	"webapi/tests"
	"webapi/utils"
)

var (
	currentDir string
	uploadFile string
	addrs      []string
)

func init() {
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	addrs, _, err = tests.GetSomeStartedProgramServer(1)
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestDownload(t *testing.T) {
	// create upload file
	err := utils.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	// upload
	uploader := u.NewUploader()
	err = uploader.Upload(addrs[0]+"/upload", uploadFile)
	if err != nil {
		t.Errorf("err from Upload: %v \n", err.Error())
	}

	url := addrs[0] + "/fileserver/upload/uploadFile"

	downloader := download.NewDownloader()

	// mkdir
	err = os.Mkdir("tmp", os.ModePerm)
	if err != nil {
		t.Errorf("err from Mkdir(): %v \n", err.Error())
	}

	done := make(chan error, 10)
	var wg sync.WaitGroup

	wg.Add(1)
	downloader.Download(url, "tmp", done, &wg)

	t.Cleanup(func() {
		tearDown()
		os.RemoveAll("tmp")
	})
}
