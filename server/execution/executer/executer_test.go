package executer_test

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"webapi/server/config"
	"webapi/server/execution/contextManager"
	executer2 "webapi/server/execution/executer"
	"webapi/server/execution/msgs"
	sh "webapi/server/handlers"
	"webapi/utils/file"
	u "webapi/utils/upload"
)

var (
	executer      executer2.Executer
	currentDir    string
	dummyParameta string
	programName   string
	uploadFile    string
	ctx           contextManager.ContextManager
	cfg           *config.Config
)

func init() {
	cfg = config.Load()
	executer = executer2.NewExecuter()
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"
	dummyParameta = "-a mike"
	programName = "convertToJson"

	go func() {
		if err := http.ListenAndServe(":8882", sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()

	err = file.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		panic(err.Error())
	}

	uploader := u.NewUploader()
	err = uploader.Upload("http://127.0.0.1:8882/upload", uploadFile)
	if err != nil {
		panic(err.Error())
	}

	cfg.StdoutBufferSize = 3000
	cfg.StderrBufferSize = 2000
	ctx, err = GetDummyContextManager(cfg)
	if err != nil {
		panic(err)
	}
}

func GetDummyContextManager(cfg *config.Config) (contextManager.ContextManager, error) {
	uploadFile = "uploadfile"
	err := file.CreateSpecifiedFile(uploadFile, 2)
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

	r, err := http.NewRequest(http.MethodPost, "/pro/"+programName, pr)
	if err != nil {
		panic(err.Error())
	}
	r.Header.Set("Content-Type", form.FormDataContentType())

	w := httptest.NewRecorder()

	var ctx contextManager.ContextManager
	ctx, err = contextManager.NewContextManager(w, r, cfg)
	if err != nil {
		return nil, fmt.Errorf("GetDummyContextManager: %v", err)
	}

	return ctx, nil
}

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestFileExecuter_Execute(t *testing.T) {

	_, _ = http.NewRequest(http.MethodPost, "/pro/convertToJson", nil)
	out := executer.Execute(ctx)
	fmt.Println(out)

	wantOut := uploadFile + ".json"
	if filepath.Base(out.OutURLs()[0]) != wantOut {
		t.Errorf("out.OutURLs()[0] is not %v.", wantOut)
	}
	if out.StdOut() == "" {
		t.Errorf("out.StdOut() is empty.")
	}
	if out.StdErr() != "" {
		t.Errorf("out.StdErr() is not empty.")
	}
	if out.Status() != msgs.OK {
		t.Errorf("out.Status() is not %v", msgs.OK)
	}
	if out.ErrorMsg() != "" {
		t.Errorf("out.ErrorMsg is not empty.")
	}

	// out.Stdout,errはcfg（設定ファイル）の値より小さくなくてはならない。設定値がマックスなので。
	if len(out.StdOut()) > cfg.StdoutBufferSize {
		t.Errorf("len(out.StdOut()):%v is not more less cfg.StdoutBufferSize: %v \n", len(out.StdOut()), cfg.StdoutBufferSize)
	}
	if len(out.StdErr()) > cfg.StderrBufferSize {
		t.Errorf("len(out.StdErr()):%v is not more less cfg.StderrBufferSize: %v", len(out.StdErr()), cfg.StderrBufferSize)
	}

	t.Cleanup(func() {
		tearDown()
	})

}

func TestFileExecuter_DeleteOutputDir(t *testing.T) {
	contextManager.Logger.SetOutput(os.Stdout)

	err := os.MkdirAll("tmpDir", os.ModePerm)
	if err != nil {
		t.Errorf("err from os.MkdirAll(): %v \n", err.Error())
	}

	err = executer2.DeleteDirSomeTimeLater("tmpDir", 1)
	if err != nil {
		t.Errorf("err from DeleteDirSomeTimeLater() : %v \n", err.Error())
	}

	t.Cleanup(func() {
		err := os.RemoveAll("tmpDir")
		if err != nil {
			t.Errorf("err from RemoveAll(): %v \n", err.Error())
		}
		tearDown()
	})
}
