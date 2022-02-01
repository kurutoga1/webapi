package test

import (
	"fmt"
	"net/http/httptest"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	"webapi/pro/msgs"
	"webapi/utils/file"
	http2 "webapi/utils/http"
)

/*
プログラムサーバのテストを行うパッケージ
TODO: タイムアウト、エラー、OK,アップロードの失敗等
TODO: アップロードファイル、プログラムネームにスペースがある場合
TODO: stdoutBufferSiZe等もテストしなければいけない
*/

type Struct struct {
	TestName              string
	IsSkip                bool
	Config                *config.Config
	ProgramName           string
	UploadFilePath        string
	UploadFileSize        int64
	Parameta              string
	ExpectedLenOfOutURLs  int
	ExpectedStdOutIsEmpty bool
	ExpectedStdErrIsEmpty bool
	ExpectedStatus        string
	ExpectedErrMsgIsEmpty bool
	ExpectedUploadIsError bool
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
	w := httptest.NewRecorder()

	var ctx contextManager.ContextManager
	ctx, err = contextManager.NewContextManager(w, r, s.Config)
	if err != nil {
		return nil, fmt.Errorf("setup: %v", err)
	}

	return ctx, nil
}

func GetMaterials() []Struct {
	tests := []Struct{
		{
			TestName:              "usually convertToJson",
			IsSkip:                false,
			ProgramName:           "convertToJson",
			Parameta:              "dummyParameta",
			ExpectedLenOfOutURLs:  1,
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			UploadFilePath:        "uploadfile1",
			UploadFileSize:        200,
			ExpectedUploadIsError: false,
		},
		{ // アップロードファイルネームにスペースがある場合
			TestName:              "upload file with space. convertToJson",
			IsSkip:                true,
			ProgramName:           "convertToJson",
			Parameta:              "dummyParameta",
			ExpectedLenOfOutURLs:  1,
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			UploadFilePath:        "upload file1",
			UploadFileSize:        200,
			ExpectedUploadIsError: false,
		},
		{
			TestName:              "upload file too large. convertToJson",
			IsSkip:                false,
			ProgramName:           "convertToJson",
			Parameta:              "dummyParameta",
			ExpectedLenOfOutURLs:  1,
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.OK,
			ExpectedErrMsgIsEmpty: true,
			UploadFilePath:        "uploadfile2",
			UploadFileSize:        400,
			ExpectedUploadIsError: true,
		},
		{
			TestName:              "err raise.",
			IsSkip:                false,
			ProgramName:           "err",
			Parameta:              "dummyParameta",
			ExpectedLenOfOutURLs:  0,
			ExpectedStdOutIsEmpty: false,
			ExpectedStdErrIsEmpty: false,
			ExpectedStatus:        msgs.PROGRAMERROR,
			ExpectedErrMsgIsEmpty: false,
			UploadFilePath:        "uploadfile3",
			UploadFileSize:        200,
			ExpectedUploadIsError: false,
		},
		{
			TestName:              "sleep. time out",
			IsSkip:                false,
			ProgramName:           "sleep",
			Parameta:              "10",
			ExpectedLenOfOutURLs:  0,
			ExpectedStdOutIsEmpty: true,
			ExpectedStdErrIsEmpty: true,
			ExpectedStatus:        msgs.PROGRAMTIMEOUT,
			ExpectedErrMsgIsEmpty: false,
			UploadFilePath:        "uploadfile4",
			UploadFileSize:        200,
			ExpectedUploadIsError: false,
		},
	}
	return tests
}
