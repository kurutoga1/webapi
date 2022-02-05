/*
全てのハンドラをセットしたrouterを返すGetServeMuxを定義している。
*/

package router

import (
	"log"
	"net/http"
	"path/filepath"
	"webapi/pro/config"
	"webapi/pro/handlers/program"
	"webapi/pro/handlers/upload"
	"webapi/pro/handlers/user"
	http2 "webapi/utils/http"
	ul "webapi/utils/log"
)

var (
	cfg     *config.Config = config.Load()
	logFile                = filepath.Join(cfg.Log.Dir, cfg.Log.GoLog)

	logger *log.Logger = ul.GetLogger(logFile)
)

func New() *programServerMux {
	return &programServerMux{}
}

type programServerMux struct{}

// New ハンドラをセットしたrouterを返す。
func (p *programServerMux) New(fileServerDir string) *http.ServeMux {
	router := http.NewServeMux()

	// ファイルサーバーの機能のハンドラ
	// cfg.FileServer.Dir以下のファイルをwebから見ることができる。
	fileServer := "/" + fileServerDir + "/"
	router.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(fileServerDir))))

	// staticをcss,js等を格納するディレクトリとする。,favicon.icoも格納する。
	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// 登録プログラムを実行させるAPI
	router.HandleFunc("/pro/", program.Handler(logger, cfg))

	// ファイルをアップロードするAPI
	router.HandleFunc("/upload", upload.UploadHandler(logger, cfg))

	// /user....の場合は全てAPIではなく、ユーザーが実際にwebにアクセスし、
	// webページのように使用する。
	router.HandleFunc("/user/top", user.UserTopHandler)
	router.HandleFunc("/user/fileUpload", user.UserFileUploadHandler)
	router.HandleFunc("/user/prepareExec", user.PrepareExecHandler)
	router.HandleFunc("/user/exec", user.ExecHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/json/health/memory", http2.GetRuntimeHandler)

	// プログラムサーバに登録してあるプログラム一覧をJSONで表示するAPI
	router.HandleFunc("/json/program/all", program.AllHandler(logger))

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", http2.HealthHandler)

	return router
}
