package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"webapi/gw/config"
	"webapi/gw/router"
	int2 "webapi/utils/int"
	log2 "webapi/utils/log"
)

func init() {
	// ライブラリでLoggerを使用する場合、ここでライブラリのLoggerにロガーをセットする。
	// なるだけライブラリではLoggerを使用しない設計にする。
	l := log2.GetLogger(cfg.LogPath)
	l.SetFlags(log.LstdFlags)
}

var (
	logger = log2.GetLogger(cfg.LogPath)
	cfg    = config.NewServerConfig()
	logMu  sync.Mutex
)

func main() {
	addr := cfg.LoadBalancerServerIP + ":" + cfg.LoadBalancerServerPort
	fmt.Printf("web pro on: %v \n", addr)

	r := router.New().New("gwFileServer")

	rotater := log2.NewLogRotater(int2.KBToByte(cfg.RotateShavingKB), int2.KBToByte(cfg.RotateMaxKB), &logMu, logger, cfg.LogPath)

	if err := http.ListenAndServe(addr, log2.RotateMiddleware(log2.HttpTraceMiddleware(r, logger), rotater)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
