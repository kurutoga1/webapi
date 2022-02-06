package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"sort"
	"webapi/gw/config"
	gg "webapi/gw/getAllPrograms"
	mg "webapi/gw/memoryGetter"
	"webapi/gw/minimumServerSelector"
	gp "webapi/gw/programHasServers"
	sc "webapi/gw/serverAliveConfirmer"
	http2 "webapi/utils/http"
)

// GetMinimumMemoryServerHandler は実際に疎通できるサーバの中から使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func GetMinimumMemoryServerHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		minimumMemoryServerSelector := minimumServerSelector.NewMinimumMemoryServerSelector()
		memoryGetter := mg.NewMemoryGetter()
		serverAliveConfirmer := sc.NewServerAliveConfirmer()
		url, err := minimumMemoryServerSelector.Select(cfg.Servers, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		type j struct {
			Url string `json:"url"`
		}

		jsonStr := j{Url: url}

		http2.PrintAsJSON(w, jsonStr)
	}
}

// GetSuitableServerHandler は実際にプログラムがあるサーバかつ、使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func GetSuitableServerHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		programName := r.URL.Path[len("/SuitableServer/"):]
		l.Printf("programName: %v ", programName)
		serverAliveConfirmer := sc.NewServerAliveConfirmer()
		aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
		if err != nil {
			l.Printf("GetSuitableServerHandler: %v ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		l.Printf("aliveservers: %v", aliveServers)

		programHasServersGetter := gp.GetProgramHasServersGetter()
		programHasServers, err := programHasServersGetter.Get(aliveServers, "/json/program/all", programName)
		if err != nil {
			l.Printf("err from Get(): %v ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		l.Printf("programHasServers: %v ", programHasServers)

		minimumMemoryServerSelector := minimumServerSelector.NewMinimumMemoryServerSelector()
		memoryGetter := mg.NewMemoryGetter()
		url, err := minimumMemoryServerSelector.Select(programHasServers, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		type j struct {
			Url string `json:"url"`
		}

		var jsonStr j
		if len(programHasServers) == 0 {
			jsonStr.Url = fmt.Sprintf("%v is not found in all pro.", programName)
		} else {
			jsonStr.Url = url
		}

		http2.PrintAsJSON(w, jsonStr)
	}
}

// mapToStruct はmapからstructに変換する。
func mapToStruct(m interface{}, val interface{}) error {
	tmp, err := json.Marshal(m)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	err = json.Unmarshal(tmp, val)
	if err != nil {
		return fmt.Errorf("mapToStruct: %v", err)
	}
	return nil
}

// UserTopHandler はユーザがwebにアクセスした場合はに
// 全サーバにアクセスし、全てのプログラムリストを取得,
// プログラムがあるサーバなおかつメモリ使用量が最小のサーバのURLをセットし、
// webpageに表示する.ボタンが押されたらそのサーバのプログラム実行準備画面に飛ぶ
func UserTopHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
		allServerMaps := map[string]interface{}{}
		serverAliveConfirmer := sc.NewServerAliveConfirmer()
		aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
		if err != nil {
			l.Printf("UserTopHandler: %v ", err.Error())
			http.Error(w, err.Error()+"全てのプログラムサーバが生きていない可能性があります。", 500)
			return
		}

		allProgramGetter := gg.NewAllProgramGetter()
		endPoint := "/json/program/all"
		allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

		l.Printf("programs of all pro: %v ", http2.GetKeysFromMap(allServerMaps))

		type tmpProInfo struct {
			Help    string `json:"help"`
			Command string `json:"command"`
		}
		type proInfo struct {
			Name    string
			Help    string
			Command string
			URL     string
		}
		// proInfoのリストをhtmlに与えるだけで良いが今後さらに項目を渡す場合に備えて構造体を定義しておく
		type htmlData struct {
			ProInfos []proInfo
		}

		var dataToHtml htmlData

		// 生きているサーバたちにプログラムを持っているか聞いていき、持っていた何台かがあった場合
		// 一番消費メモリの少ないサーバを選択する。
		// htmlに渡すhtmlData(struct)をこのループで完成させる。
		for proName, programInfoMap := range allServerMaps {
			var tmpInfo tmpProInfo
			err := mapToStruct(programInfoMap, &tmpInfo)
			if err != nil {
				l.Printf("UserTopHandler: %v ", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}

			programHasServersGetter := gp.GetProgramHasServersGetter()
			programHasServers, err := programHasServersGetter.Get(aliveServers, "/json/program/all", proName)
			if err != nil {
				l.Printf("UserTopHandler: %v ", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}

			// プログラムを保持しているサーバたちの中で一番使用メモリが少ないサーバを選択する。
			minimumMemoryServerSelector := minimumServerSelector.NewMinimumMemoryServerSelector()
			memoryGetter := mg.NewMemoryGetter()
			serverAliveConfirmer := sc.NewServerAliveConfirmer()
			url, err := minimumMemoryServerSelector.Select(programHasServers, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
			if err != nil {
				l.Printf("UserTopHandler: %v ", err.Error())
				http.Error(w, err.Error(), 500)
				return
			}

			p := proInfo{
				Name:    proName,
				Help:    tmpInfo.Help,
				Command: tmpInfo.Command,
				URL:     url,
			}
			dataToHtml.ProInfos = append(dataToHtml.ProInfos, p)
		}

		w.Header().Add("Content-Type", "text/html")
		serveHtml := filepath.Join("./templates", "top.html")
		absHtml, err := filepath.Abs(serveHtml)
		if err != nil {
			l.Printf("UserTopHandler: %v ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// 名前順でソートする
		sort.Slice(dataToHtml.ProInfos, func(i, j int) bool { return dataToHtml.ProInfos[i].Name < dataToHtml.ProInfos[j].Name })

		http2.RenderTemplate(w, absHtml, dataToHtml)
	}
}

func GetAliveServersHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serverAliveConfirmer := sc.NewServerAliveConfirmer()
		aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
		if err != nil {
			l.Printf("err from GetAliveServers(): %v \n", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		type data struct {
			AliveServers []string `json:"AliveServers"`
		}

		jsonStr := data{AliveServers: aliveServers}

		http2.PrintAsJSON(w, jsonStr)
	}
}

func GetAllProgramsHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
		allServerMaps := map[string]interface{}{}
		serverAliveConfirmer := sc.NewServerAliveConfirmer()
		aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
		if err != nil {
			l.Printf("err from GetAliveServers(): %v \n", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		allProgramGetter := gg.NewAllProgramGetter()
		endPoint := "/json/program/all"
		allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

		http2.PrintAsJSON(w, allServerMaps)
	}
}
