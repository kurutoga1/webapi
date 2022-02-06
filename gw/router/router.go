package router

import (
	"log"
	"net/http"
	"path/filepath"
	"webapi/gw/config"
	"webapi/gw/handlers"
	http2 "webapi/utils/http"
	ul "webapi/utils/log"
)

var (
	cfg     *config.Config = config.NewServerConfig()
	logFile                = filepath.Join(cfg.LogPath)
	logger  *log.Logger    = ul.GetLogger(logFile)
)

func New() *apiGwServerMux {
	return &apiGwServerMux{}
}

type apiGwServerMux struct{}

func (a *apiGwServerMux) New(fileServerDir string) *http.ServeMux {
	router := http.NewServeMux()

	// ファイルサーバーの機能のハンドラ
	fileServer := "/" + fileServerDir + "/"
	router.Handle(fileServer, http.StripPrefix(fileServer, http.FileServer(http.Dir(fileServerDir))))

	// ユーザがこのハンドラにアクセスした場合は全てのサーバにアクセスし、全てのプログラムを表示する。
	router.HandleFunc("/userTop", handlers.UserTopHandler(logger, cfg))

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// コマンドラインからはここにアクセスし、メモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/MinimumMemoryServer", handlers.GetMinimumMemoryServerHandler(cfg))

	// コマンドラインからここにアクセスし、プログラムがあるかつメモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/SuitableServer/", handlers.GetSuitableServerHandler(logger, cfg))

	// 現在稼働しているサーバを返すAPI
	router.HandleFunc("/AliveServers", handlers.GetAliveServersHandler(logger, cfg))

	// 生きている全てのサーバのプログラムを取得してJSONで表示するAPI
	router.HandleFunc("/AllServerPrograms", handlers.GetAllProgramsHandler(logger, cfg))

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", http2.HealthHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/json/health/memory", http2.GetRuntimeHandler)

	return router
}
