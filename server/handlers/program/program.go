/*
クライアントプログラムからAPIにアクセスし、プログラムを実行するハンドラを定義している。
*/

package program

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	cp "webapi/cli/processFileOnServer"
	"webapi/server/config"
	"webapi/server/execution/contextManager"
	"webapi/server/execution/executer"
	"webapi/server/execution/msgs"
	"webapi/server/outputManager"
	"webapi/utils"
)

var cfg = config.Load()

// ProgramHandler はプログラムを実行するためのハンドラー。処理結果をJSON文字列で返す。
// cliからアクセスされる。cliの場合はこのハンドラにリクエストがくる前にファイルのアップロードは
// 完了し,アップロードディレクトリに格納されている。bodyにファイルベース名とパラメタを格納し、リクエストとしてこのハンドラ
// に送られる。
// アップロードファイルやパラメータ等を使用し、コマンド実行する。
func ProgramHandler(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		err := errors.New("Body of request is nil. should set body that are Filename,parameta.")
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	programName := r.URL.Path[len("/pro/"):]

	var out outputManager.OutputManager = outputManager.NewOutputManager()
	var newExecuter executer.Executer = executer.NewExecuter()

	_, err := config.GetProConfByName(programName)

	// プログラムがこのサーバになかった場合
	if err != nil {
		msg := fmt.Sprintf("%v is not found.", programName)

		out.SetErrorMsg(msg)
		out.SetStatus(msgs.SERVERERROR)
		logf(msg)

		// jsonに変換
		b, err := json.Marshal(out)
		if err != nil {
			logf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// JSONを表示
		_, err = fmt.Fprintf(w, string(b))
		if err != nil {
			logf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}

	// bodyを構造体に代入
	var u cp.UploadedInfo
	err = json.NewDecoder(r.Body).Decode(&u)
	if err != nil && err != io.EOF {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	logf("programName: %v, parameta: %v", programName, u.Parameta)

	// このリクエストがくる前にuploadディレクトリにアップロードは完了しているのでbasenameを使用しパスをつなげる。
	uploadedDir := filepath.Join(cfg.FileServer.Dir, cfg.FileServer.UploadDir)
	uploadedFilePath := filepath.Join(uploadedDir, u.Filename)
	if !utils.FileExists(uploadedFilePath) {
		err := errors.New("アップロードファイルが見つかりません。path: " + uploadedFilePath)
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	ctx, err := contextManager.NewContextManager(programName, uploadedFilePath, u.Parameta, cfg)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	logf("command: %v", ctx.Command())
	out = newExecuter.Execute(ctx)
	logf("Status: %v", out.Status())
	logf("ErrMsg: %v", out.ErrorMsg())

	// JSONに変換
	b, err := json.Marshal(out)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	// JSONを表示
	_, err = fmt.Fprintf(w, string(b))
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
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

	// JSONに変換
	b, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	// JSONを表示
	_, err = fmt.Fprintf(w, string(b))
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
