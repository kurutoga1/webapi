package download_test

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"webapi/cli/download"
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
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestDownload(t *testing.T) {
	// create upload file
	err := file.CreateSpecifiedFile(uploadFile, 200)
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
	downloader.Download(url, "tmp", done, &wg, file.NewMover())

	t.Cleanup(func() {
		tearDown()
		os.RemoveAll("tmp")
	})
}
