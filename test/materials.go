package test

import (
	"fmt"
	"net/http/httptest"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	"webapi/pro/handlers/upload"
	"webapi/pro/msgs"
	"webapi/utils/file"
	http2 "webapi/utils/http"
)

/*
プログラムサーバのテストを行うパッケージ
タイムアウト、エラー、OK,アップロードの失敗等
アップロードファイル、プログラムネームにスペースがある場合
stdoutBufferSiZe等もテストしなければいけない
入力ファイルにスペースがある場合
*/

type Struct struct {
	TestName              string
	IsSkip                bool
	Config                *config.Config
	ProgramName           string
	UploadFilePath        string
	UploadFileSize        int64
	Parameta              string
	ExpectedOutFileNames  []string
	ExpectedStdOutIsEmpty bool
	ExpectedStdErrIsEmpty bool
	ExpectedStatus        string
	ExpectedErrMsgIsEmpty bool
	ExpectedError         error
}

func (s *Struct) Setup() (contextManager.ContextManager, error) {
	if s.Config == nil {
		s.Config = config.Load()
	}
	err := file.CreateSpecifiedFile(s.UploadFilePath, s.UploadFileSize)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]string{
		"proName":  s.ProgramName,
		"parameta": s.Parameta,
	}

	poster := http2.NewPostGetter()
	r, err := poster.GetPostRequest("/pro/"+s.ProgramName, s.UploadFilePath, fields)
	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}
	w := httptest.NewRecorder()

	var ctx contextManager.ContextManager
	ctx, err = contextManager.NewContextManager(w, r, s.Config)
	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}

	return ctx, nil
}

func GetMaterials() []Struct {
	tests := []Struct{
		// 普通にプログラムを実行し、出力させる一番典型的なパターン
		{
			TestName:              "usually convertToJson",
			IsSkip:                false,
			ProgramName:           "convertToJson",
			UploadFilePath:        "uploadfile1",
			UploadFileSize:        200,
			Parameta:              "dummyParameta",
			ExpectedOutFileNames:  []string{"uploadfile1.json"},
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			ExpectedError:         nil,
		},
		{ // アップロードファイルにスペースがある場合
			TestName:              "upload file with space. convertToJson",
			IsSkip:                false,
			ProgramName:           "convertToJson",
			UploadFilePath:        "upload file1",
			UploadFileSize:        200,
			Parameta:              "dummyParameta",
			ExpectedOutFileNames:  []string{"uploadfile1.json"},
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			ExpectedError:         nil,
		},
		// アップロードファイルのサイズがconfigのMaxUploadSizeMBを超えている場合
		// アップロードエラーが発生する。
		{
			TestName:              "upload file too large. convertToJson",
			IsSkip:                false,
			ProgramName:           "convertToJson",
			Parameta:              "dummyParameta",
			UploadFilePath:        "uploadfile2",
			UploadFileSize:        400,
			ExpectedOutFileNames:  []string{},
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			ExpectedError:         upload.FileSizeTooBigError,
		},
		// プログラムがエラーを起こす場合
		{
			TestName:              "err raise.",
			IsSkip:                false,
			ProgramName:           "err",
			UploadFilePath:        "uploadfile3",
			UploadFileSize:        200,
			Parameta:              "dummyParameta",
			ExpectedOutFileNames:  []string{},
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: false,
			ExpectedStatus:        msgs.PROGRAMERROR,
			ExpectedErrMsgIsEmpty: false,
			ExpectedError:         nil,
		},
		// プログラムがタイムアウトを起こす場合
		{
			TestName:              "sleep. time out",
			IsSkip:                false,
			ProgramName:           "sleep",
			UploadFilePath:        "uploadfile4",
			UploadFileSize:        200,
			Parameta:              "10",
			ExpectedOutFileNames:  []string{},
			ExpectedStdOutIsEmpty: true,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.PROGRAMTIMEOUT,
			ExpectedErrMsgIsEmpty: false,
			ExpectedError:         nil,
		},
		// 実行するプログラムをlinuxのmvコマンドでやってみる。
		{
			TestName:              "move success",
			IsSkip:                false,
			ProgramName:           "move",
			UploadFilePath:        "uploadfile",
			UploadFileSize:        200,
			Parameta:              "10",
			ExpectedOutFileNames:  []string{"moved.txt"},
			ExpectedStdOutIsEmpty: true,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			ExpectedError:         nil,
		},
		// 存在しないプログラム名の場合
		{
			TestName:              "no program name",
			IsSkip:                false,
			ProgramName:           "nothingProgram",
			UploadFilePath:        "uploadfile",
			UploadFileSize:        200,
			Parameta:              "10",
			ExpectedOutFileNames:  []string{},
			ExpectedStdOutIsEmpty: true,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.PROGRAMERROR,
			ExpectedErrMsgIsEmpty: true,
			ExpectedError:         config.ProgramNotFoundError,
		},
	}
	return tests
}
