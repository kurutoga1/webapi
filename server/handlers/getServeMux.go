/*
全てのハンドラをセットしたrouterを返すGetServeMuxを定義している。
*/

package handlers

import (
	"net/http"
	"webapi/server/handlers/program"
	"webapi/server/handlers/upload"
	"webapi/server/handlers/user"
	"webapi/utils/hanlders"
)

// GetServeMux ハンドラをセットしたrouterを返す。
func GetServeMux(fileServerDir string) *http.ServeMux {
	router := http.NewServeMux()

	// ファイルサーバーの機能のハンドラ
	// cfg.FileServer.Dir以下のファイルをwebから見ることができる。
	fileServer := "/" + fileServerDir + "/"
	router.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(fileServerDir))))

	// staticをcss,js等を格納するディレクトリとする。,favicon.icoも格納する。
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 登録プログラムを実行させるAPI
	router.HandleFunc("/pro/", program.ProgramHandler)

	// ファイルをアップロードするAPI
	router.HandleFunc("/upload", upload.UploadHandler)

	// /user....の場合は全てAPIではなく、ユーザーが実際にwebにアクセスし、
	// webページのように使用する。
	router.HandleFunc("/user/top", user.UserTopHandler)
	router.HandleFunc("/user/fileUpload", user.UserFileUploadHandler)
	router.HandleFunc("/user/prepareExec", user.PrepareExecHandler)
	router.HandleFunc("/user/exec", user.ExecHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/json/health/memory", hanlders.GetRuntimeHandler)

	// プログラムサーバに登録してあるプログラム一覧をJSONで表示するAPI
	router.HandleFunc("/json/program/all", program.ProgramAllHandler)

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", hanlders.HealthHandler)

	return router
}
