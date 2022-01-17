package processFileOnServer_test

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	p "webapi/cli/processFileOnServer"
	u "webapi/cli/upload"
	sh "webapi/server/handlers"
	"webapi/utils"
)

var (
	currentDir string
	uploadFile string
)

func init() {
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	go func() {
		if err := http.ListenAndServe(":8881", sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func TestProcessFileOnServer(t *testing.T) {
	// create upload file
	err := utils.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		t.Errorf("err from CreateSpecifiedFile: %v \n", err.Error())
	}

	// upload
	uploader := u.NewUploader()
	err = uploader.Upload("http://127.0.0.1:8881/upload", uploadFile)
	if err != nil {
		t.Errorf("err from Upload: %v \n", err.Error())
	}

	// ファイル上で処理させる
	basename := filepath.Base(uploadFile)
	programName := "convertToJson"
	serverURL := "http://127.0.0.1:8881"
	parameta := ""

	// 処理させるためのurl
	url := fmt.Sprintf("%v/pro/%v", serverURL, programName)

	processor := p.NewFileProcessor()

	res, err := processor.Process(url, basename, parameta)
	if err != nil {
		t.Errorf("err from Process: %v \n", err.Error())
	}

	outBase := filepath.Base(res.OutURLs()[0])
	if outBase != "uploadfile.json" {
		t.Errorf("output file is not %v \n", "uploadfile.json")
	}

	t.Cleanup(func() {
		tearDown()
	})
}
