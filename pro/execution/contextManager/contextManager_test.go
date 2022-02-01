package contextManager_test

import (
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"webapi/pro/config"
	ec "webapi/pro/execution/contextManager"
	"webapi/pro/router"
	"webapi/test"
	"webapi/utils/file"
	"webapi/utils/http"
	u "webapi/utils/upload"
)

var (
	currentDir    string
	uploadFile    string
	dummyParameta string
	programName   string
	ctx           ec.ContextManager
	addrs         []string
	ports         []string
)

var cfg *config.Config

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"
	dummyParameta = "dummyParameta"
	programName = "convertToJson"

	addrs, ports, err = test.GetStartedServers(router.New(), 1)
	fmt.Printf("addrs: %v, ports: %v \n", addrs, ports)
	if err != nil {
		panic(err)
	}

	err = file.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		panic(err.Error())
	}

	uploader := u.NewUploader()
	err = uploader.Upload(addrs[0]+"/upload", uploadFile)
	if err != nil {
		panic(err.Error())
	}

	cfg = config.Load()
	ctx, err = GetDummyContextManager(cfg)
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(uploadFile)
}

func GetDummyContextManager(cfg *config.Config) (ec.ContextManager, error) {
	uploadFile = "uploadfile"
	err := file.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]string{
		"proName":  "convertToJson",
		"parameta": "dummyParameta",
	}

	poster := http.NewPostGetter()
	r, err := poster.GetPostRequest("/pro/convertToJson", uploadFile, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	var ctx ec.ContextManager
	ctx, err = ec.NewContextManager(w, r, cfg)
	if err != nil {
		return nil, fmt.Errorf("GetDummyContextManager: %v", err)
	}

	return ctx, nil
}

func TestNewContextManager(t *testing.T) {
	// ctx.SetProgramOutDir, SetUploadedFilePathAndParametaはここで同時に試験できる。
	if ctx.Parameta() != dummyParameta {
		t.Errorf("ctx.Parameta(): %v, want: %v \n", ctx.Parameta(), dummyParameta)
	}
	if filepath.Base(ctx.UploadedFilePath()) != uploadFile {
		t.Errorf("ctx.UploadFilePath(): %v, want: %v \n", filepath.Base(ctx.UploadedFilePath()), uploadFile)
	}
	if ctx.ProgramName() != programName {
		t.Errorf("ctx.ProgramName(): %v, want: %v \n", ctx.ProgramName(), programName)
	}
	if filepath.Base(ctx.InputFilePath()) != uploadFile {
		t.Errorf("ctx.InputFilePath(): %v , want: %v \n", filepath.Base(ctx.InputFilePath()), uploadFile)
	}
	if !file.FileExists(ctx.InputFilePath()) {
		t.Errorf("ctx.InputFilePath(%v) is not found \n", ctx.InputFilePath())
	}

	if !file.FileExists(ctx.OutputDir()) {
		t.Errorf("ctx.OutputDir() is not found.")
	}

	if !reflect.DeepEqual(ctx.Config(), cfg) {
		t.Errorf("ctx.Config(%v) is not equal cfg(%v) \n", ctx.Config(), cfg)
	}

	if !file.FileExists(ctx.ProgramTempDir()) {
		t.Errorf("ctx.ProgramTempDir is not found")
	}

	t.Cleanup(func() {
		tearDown()
	})
}
