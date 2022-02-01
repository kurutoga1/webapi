package executer_test

import (
	"fmt"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	executer2 "webapi/pro/execution/executer"
	"webapi/pro/execution/outputManager"
	"webapi/pro/msgs"
	"webapi/test"
	"webapi/utils/file"
	http2 "webapi/utils/http"
)

var (
	executer executer2.Executer
	cfg      *config.Config
)

func init() {
	cfg = config.Load()
	executer = executer2.NewExecuter()
}

func setup(cfg *config.Config, uploadFileName, proName, parameta string, uploadFileKBSize int64) (contextManager.ContextManager, error) {
	err := file.CreateSpecifiedFile(uploadFileName, uploadFileKBSize)
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
	r, err := poster.GetPostRequest("/pro/"+proName, uploadFileName, fields)
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
}

func TestFileExecuter_Execute(t *testing.T) {
	// コンテキストを用意する→実行→出力を判定、テストする

	tests := test.GetMaterials()

	for _, tt := range tests {
		if tt.IsSkip {
			continue
		}
		ctx, err := tt.Setup()
		if err != nil {
			if strings.Contains(err.Error(), msgs.UploadSizeIsTooBig) != tt.ExpectedUploadIsError {
				t.Errorf("got: %v, want: upload err.", err.Error())
				continue
			} else {
				// ファイルのアップロードテストが予想通りにいった場合。しかし他にエラーの可能性もある。
				// とりあえずアップロードテストは成功。他のエラーハンドリングは
				t.Run(tt.TestName, func(j *testing.T) {
				})
				fmt.Printf("err: %v \n", err.Error())
				os.Remove(tt.UploadFilePath)
				continue
			}
		}

		out := executer.Execute(ctx)
		t.Run(tt.TestName, func(j *testing.T) {
			testExecute(t, out, tt, cfg)
		})
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func testExecute(t *testing.T, out outputManager.OutputManager, tt test.Struct, cfg *config.Config) {
	t.Helper()
	if len(out.OutURLs()) != tt.ExpectedLenOfOutURLs {
		t.Errorf("name: %v, len of out.OutURLs() is more than 0. got: %v", tt.TestName, len(out.OutURLs()))
	}
	if (out.StdOut() == "") != tt.ExpectedStdOutIsEmpty {
		t.Errorf("name: %v, out.StdOut() is empty. stdout: %v", tt.TestName, out.StdOut())
	}
	if (out.StdErr() == "") != tt.ExpectedStdErrIsEmpty {
		t.Errorf("name: %v, out.StdErr() is not empty. stdout: %v", tt.TestName, out.StdErr())
	}
	if out.Status() != tt.ExpectedStatus {
		t.Errorf("name: %v, out.ExpectedStatus() is not %v. got: %v", tt.TestName, msgs.OK, out.Status())
	}
	if (out.ErrorMsg() == "") != tt.ExpectedErrMsgIsEmpty {
		t.Errorf("name: %v, out.ErrorMsg is not empty. got: %v", tt.TestName, out.ErrorMsg())
	}

	// out.Stdout,errはcfg（設定ファイル）の値より小さくなくてはならない。設定値がマックスなので。
	if len(out.StdOut()) > cfg.StdoutBufferSize {
		t.Errorf("name: %v, len(out.StdOut()):%v is not more less cfg.StdoutBufferSize: %v \n", tt.TestName, len(out.StdOut()), cfg.StdoutBufferSize)
	}
	if len(out.StdErr()) > cfg.StderrBufferSize {
		t.Errorf("name: %v, len(out.StdErr()):%v is not more less cfg.StderrBufferSize: %v. ", tt.TestName, len(out.StdErr()), cfg.StderrBufferSize)
	}

	t.Cleanup(func() {
		tearDown()
		os.Remove(tt.UploadFilePath)
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
