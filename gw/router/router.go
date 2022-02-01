package router

import (
	"net/http"
	"webapi/gw/handlers"
	http2 "webapi/utils/http"
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
	router.HandleFunc("/userTop", handlers.UserTopHandler)

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// コマンドラインからはここにアクセスし、メモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/MinimumMemoryServer", handlers.GetMinimumMemoryServerHandler)

	// コマンドラインからここにアクセスし、プログラムがあるかつメモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/SuitableServer/", handlers.GetSuitableServerHandler)

	// 現在稼働しているサーバを返すAPI
	router.HandleFunc("/AliveServers", handlers.GetAliveServersHandler)

	// 生きている全てのサーバのプログラムを取得してJSONで表示するAPI
	router.HandleFunc("/AllServerPrograms", handlers.GetAllProgramsHandler)

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", http2.HealthHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/json/health/memory", http2.GetRuntimeHandler)

	return router
}
