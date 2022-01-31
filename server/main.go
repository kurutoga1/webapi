/*
プログラムサーバの開始
*/

package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	int2 "webapi/utils/int"

	"webapi/server/config"
	"webapi/server/handlers"
	ul "webapi/utils/log"
)

var (
	cfg     *config.Config = config.Load()
	logFile                = filepath.Join(cfg.Log.Dir, cfg.Log.GoLog)

	logger *log.Logger = ul.GetLogger(logFile)
	logMu  sync.Mutex
)

func init() {}

func main() {

	router := handlers.NewRouter(cfg.FileServer.Dir)

	port := ":" + cfg.ServerPort
	fmt.Printf("web server on %v%v\n", cfg.ServerIP, port)

	rotater := ul.NewLogRotater(int2.KBToByte(cfg.Log.RotateShavingKB), int2.KBToByte(cfg.Log.RotateMaxKB), &logMu, logger, logFile)

	if err := http.ListenAndServe(port, ul.RotateMiddleware(ul.HttpTraceMiddleware(router, logger), rotater)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
