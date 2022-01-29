package executer_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"webapi/server/config"
	"webapi/server/execution/contextManager"
	executer2 "webapi/server/execution/executer"
	"webapi/server/execution/msgs"
	sh "webapi/server/handlers"
	"webapi/server/outputManager"
	"webapi/utils/file"
	http2 "webapi/utils/http"
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

	cfg.StdoutBufferSize = 3000
	cfg.StderrBufferSize = 2000

}

func setup(cfg *config.Config, proName, parameta string, uploadFileKBSize int64) (contextManager.ContextManager, error) {
	uploadFile = "uploadfile"
	err := file.CreateSpecifiedFile(uploadFile, uploadFileKBSize)
	if err != nil {
		panic(err.Error())
	}
	// TODO: urlパラメータに限界文字数はあるのか
	// urlパラメータで情報を送る場合はquery.getみたいなので取得しないといけない。またhtmlをそっちに合わせないといけない

	fields := map[string]string{
		"proName":  proName,
		"parameta": parameta,
	}

	poster := http2.NewPostGetter()
	r, err := poster.GetPostRequest("/pro/"+proName, uploadFile, fields)
	w := httptest.NewRecorder()

	var ctx contextManager.ContextManager
	ctx, err = contextManager.NewContextManager(w, r, cfg)
	if err != nil {
		return nil, fmt.Errorf("setup: %v", err)
	}

	return ctx, nil
}

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestFileExecuter_Execute(t *testing.T) {
	// コンテキストを用意する→実行→出力を判定、テストする
	// TODO: ここはシステムの要なのでしっかりテストする必要がある。
	// TODO: タイムアウト、エラー、OK,アップロードの失敗等
	// TODO: stdoutBufferSiZe等もテストしなければいけない

	tests := []struct {
		programName    string
		parameta       string
		lenOfOutURLs   int
		stdOutIsEmpty  bool
		stdErrIsEmpty  bool
		status         string
		errMsgIsEmpty  bool
		uploadFileSize int64
		uploadIsError  bool
	}{
		{
			programName:    "convertToJson",
			parameta:       "dummyParameta",
			lenOfOutURLs:   1,
			stdOutIsEmpty:  false,
			stdErrIsEmpty:  true,
			status:         msgs.OK,
			errMsgIsEmpty:  true,
			uploadFileSize: 200,
			uploadIsError:  false,
		},
		{
			programName:    "convertToJson",
			parameta:       "dummyParameta",
			lenOfOutURLs:   1,
			stdOutIsEmpty:  false,
			stdErrIsEmpty:  true,
			status:         msgs.OK,
			errMsgIsEmpty:  true,
			uploadFileSize: 400,
			uploadIsError:  true,
		},
		{
			programName:    "err",
			parameta:       "dummyParameta",
			lenOfOutURLs:   0,
			stdOutIsEmpty:  false,
			stdErrIsEmpty:  false,
			status:         msgs.PROGRAMERROR,
			errMsgIsEmpty:  false,
			uploadFileSize: 200,
			uploadIsError:  false,
		},
		{
			programName:    "sleep",
			parameta:       "10",
			lenOfOutURLs:   0,
			stdOutIsEmpty:  true,
			stdErrIsEmpty:  true,
			status:         msgs.PROGRAMTIMEOUT,
			errMsgIsEmpty:  false,
			uploadFileSize: 200,
			uploadIsError:  false,
		},
	}

	for _, tt := range tests {
		ctx, err := setup(cfg, tt.programName, tt.parameta, tt.uploadFileSize)
		if err != nil {
			if strings.Contains(err.Error(), "アップロードされたファイルが大きすぎます。") != tt.uploadIsError {
				t.Errorf("got: %v, want: upload err.", err.Error())
			}
		}

		if ctx != nil {
			out := executer.Execute(ctx)
			t.Run("test "+tt.programName, func(j *testing.T) {
				testExecute(t, out, tt.lenOfOutURLs, tt.stdOutIsEmpty, tt.stdErrIsEmpty, tt.status, tt.errMsgIsEmpty, cfg.StdoutBufferSize, cfg.StderrBufferSize)
			})
		}
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func testExecute(t *testing.T, out outputManager.OutputManager, lenOfOutURLs int, stdOutIsEmpty, stdErrIsEmpty bool, status string, errMsgIsEmpty bool, stdOutBufferSize, stdErrBufferSize int) {
	if len(out.OutURLs()) != lenOfOutURLs {
		t.Errorf("len of out.OutURLs() is more than 0 ")
	}
	if (out.StdOut() == "") != stdOutIsEmpty {
		t.Errorf("out.StdOut() is empty.")
	}
	if (out.StdErr() == "") != stdErrIsEmpty {
		t.Errorf("out.StdErr() is empty.")
	}
	if out.Status() != status {
		t.Errorf("out.Status() is not %v", msgs.OK)
	}
	if (out.ErrorMsg() == "") != errMsgIsEmpty {
		t.Errorf("out.ErrorMsg is empty.")
	}

	// out.Stdout,errはcfg（設定ファイル）の値より小さくなくてはならない。設定値がマックスなので。
	if len(out.StdOut()) > stdOutBufferSize {
		t.Errorf("len(out.StdOut()):%v is not more less cfg.StdoutBufferSize: %v \n", len(out.StdOut()), cfg.StdoutBufferSize)
	}
	if len(out.StdErr()) > stdErrBufferSize {
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
