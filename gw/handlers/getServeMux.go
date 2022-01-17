package handlers

import (
	"net/http"
	"webapi/utils"
)

func GetServeMux() *http.ServeMux {
	router := http.NewServeMux()

	// ユーザがwebにアクセスした場合はメモリ使用量が一番低いサーバへリダイレクトする。
	router.HandleFunc("/userTop", UserTopHandler)

	router.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// コマンドラインからはここにアクセスし、メモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/MinimumMemoryServer", GetMinimumMemoryServerHandler)

	// コマンドラインからここにアクセスし、プログラムがあるかつメモリ使用量が一番低いサーバのURLを返す。
	router.HandleFunc("/SuitableServer/", GetSuitableServerHandler)

	// 現在稼働しているサーバを返すAPI
	router.HandleFunc("/AliveServers", GetAliveServersHandler)

	// 生きている全てのサーバのプログラムを取得してJSONで表示するAPI
	router.HandleFunc("/AllServerPrograms", GetAllProgramsHandler)

	// このサーバが生きているかを判断するのに使用するハンドラ
	router.HandleFunc("/health", utils.HealthHandler)

	// このサーバプログラムのメモリ状態をJSONで表示するAPI
	router.HandleFunc("/json/health/memory", utils.GetRuntimeHandler)

	return router
}
