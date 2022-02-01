/*
クライアントプログラムからAPIにアクセスし、プログラムを実行するハンドラを定義している。
*/

package program

import (
	"fmt"
	"log"
	"net/http"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	"webapi/pro/execution/executer"
	"webapi/pro/execution/outputManager"
	"webapi/pro/msgs"
	http2 "webapi/utils/http"
)

// ProgramHandler はプログラムを実行するためのハンドラー。処理結果をJSON文字列で返す。
// cliからアクセスされる。cliの場合はこのハンドラにリクエストがくる前にファイルのアップロードは
// 完了し,アップロードディレクトリに格納されている。bodyにファイルベース名とパラメタを格納し、リクエストとしてこのハンドラ
// に送られる。
// アップロードファイルやパラメータ等を使用し、コマンド実行する。
func ProgramHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		programName := r.URL.Path[len("/pro/"):]

		var out outputManager.OutputManager = outputManager.NewOutputManager()
		var newExecuter executer.Executer = executer.NewExecuter()

		// TODO:
		// ctx, err := contextManager.NewContextManager(w, r, cfg)
		// errors.Is(err, config.ProgramNotFoundError)みたいなので判定した方がいいかも
		_, err := config.GetProConfByName(programName)

		// プログラムがこのサーバになかった場合
		if err != nil {
			msg := fmt.Sprintf("%v is not found.", programName)

			out.SetErrorMsg(msg)
			out.SetStatus(msgs.SERVERERROR)
			logf(msg)

			http2.PrintAsJSON(w, out)
			return
		}

		ctx, err := contextManager.NewContextManager(w, r, cfg)
		if err != nil {
			logf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		logf("programName: %v", ctx.ProgramName())
		logf("uploadFilePath: %v", ctx.UploadedFilePath())
		logf("parameta: %v", ctx.Parameta())
		logf("command: %v", ctx.Command())
		out = newExecuter.Execute(ctx)
		logf("ExpectedStatus: %v", out.Status())
		logf("ErrMsg: %v", out.ErrorMsg())

		http2.PrintAsJSON(w, out)
		return
	}
}

// ProgramAllHandler は登録されているプログラムの全てをJSONで返す。
func ProgramAllHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.Body
	w.Header().Add("Content-Type", "application/json")

	proConfList, err := config.GetPrograms()
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	m := map[string]interface{}{}
	// proConfListから代入していく
	for _, ele := range proConfList {
		m1 := map[string]string{}
		m[ele.Name()] = m1
		m1["command"] = ele.Command()
		help, err := ele.Help()
		if err != nil {
			logf("ProgramAllHandler: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		m1["help"] = help
	}

	http2.PrintAsJSON(w, m)
	return
}
