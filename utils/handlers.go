package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webapi/server/myRuntime"
)

// HealthHandler はサーバが生きているか確認するだけのハンドラ
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	type h struct {
		Health string `json:"health"`
	}
	healthStr := h{Health: "ok"}
	bytes, _ := json.MarshalIndent(&healthStr, "", "    ")
	_, err := fmt.Fprintf(w, string(bytes))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// GetRuntimeHandler はこのサーバプログラムのメモリの状態をJSONで表示する
func GetRuntimeHandler(w http.ResponseWriter, r *http.Request) {
	runtimeGetter := myRuntime.NewRuntimeGetter()
	runtimeJSON, err := runtimeGetter.GetRuntimeAsJSON()

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = fmt.Fprintf(w, runtimeJSON)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	return
}
