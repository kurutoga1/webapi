package executer_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	executer2 "webapi/pro/execution/executer"
	"webapi/pro/execution/outputManager"
	"webapi/pro/msgs"
	"webapi/test"
)

var (
	executer executer2.Executer
	cfg      *config.Config
)

func init() {
	cfg = config.Load()
	executer = executer2.NewExecuter()
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
}

func testExecute(t *testing.T, out outputManager.OutputManager, tt test.Struct, cfg *config.Config) {
	t.Helper()

	if len(tt.ExpectedOutFileNames) > 0 {
		for _, ExpectedOutFileName := range tt.ExpectedOutFileNames {
			isExists := false
			for _, outPath := range out.OutURLs() {
				if filepath.Base(outPath) == ExpectedOutFileName {
					isExists = true
				}
			}
			if !isExists {
				t.Errorf("test name: %v,\"%v\" is not exists of %v", tt.TestName, ExpectedOutFileName, out.OutURLs())
			}
		}
	}

	if (out.StdOut() == "") != tt.ExpectedStdOutIsEmpty {
		t.Errorf("test name: %v, out.StdOut() is empty. stdout: %v", tt.TestName, out.StdOut())
	}
	if (out.StdErr() == "") != tt.ExpectedStdErrIsEmpty {
		t.Errorf("test name: %v, out.StdErr() is not empty. stdout: %v", tt.TestName, out.StdErr())
	}
	if out.Status() != tt.ExpectedStatus {
		t.Errorf("test name: %v, out.ExpectedStatus() is not %v. got: %v", tt.TestName, msgs.OK, out.Status())
	}
	if (out.ErrorMsg() == "") != tt.ExpectedErrMsgIsEmpty {
		t.Errorf("test name: %v, out.ErrorMsg is not empty. got: %v", tt.TestName, out.ErrorMsg())
	}

	// out.Stdout,errはcfg（設定ファイル）の値より小さくなくてはならない。設定値がマックスなので。
	if len(out.StdOut()) > cfg.StdoutBufferSize {
		t.Errorf("test name: %v, len(out.StdOut()):%v is not more less cfg.StdoutBufferSize: %v \n", tt.TestName, len(out.StdOut()), cfg.StdoutBufferSize)
	}
	if len(out.StdErr()) > cfg.StderrBufferSize {
		t.Errorf("test name: %v, len(out.StdErr()):%v is not more less cfg.StderrBufferSize: %v. ", tt.TestName, len(out.StdErr()), cfg.StderrBufferSize)
	}

	t.Cleanup(func() {
		os.RemoveAll("fileserver")
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
	})
}
