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
	"webapi/utils"

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

	router := handlers.GetServeMux(cfg.FileServer.Dir)

	port := ":" + cfg.ServerPort
	fmt.Printf("web server on %v%v\n", cfg.ServerIP, port)

	rotater := ul.NewLogRotater(utils.KBToByte(cfg.Log.RotateShavingKB), utils.KBToByte(cfg.Log.RotateMaxKB), &logMu, logger, logFile)

	if err := http.ListenAndServe(port, ul.RotateMiddleware(ul.HttpTraceMiddleware(router, logger), rotater)); err != nil {
		panic(fmt.Errorf("[FAILED] start sever. err: %v", err))
	}
}
