package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"sort"
	"text/template"
	"webapi/gw/config"
	gg "webapi/gw/getAllPrograms"
	mg "webapi/gw/memoryGetter"
	"webapi/gw/minimumServerSelector"
	gp "webapi/gw/programHasServers"
	sc "webapi/gw/serverAliveConfirmer"
	http2 "webapi/utils/http"
)

var (
	cfg                         = config.NewServerConfig()
	minimumMemoryServerSelector = minimumServerSelector.NewMinimumMemoryServerSelector()
	memoryGetter                = mg.NewMemoryGetter()
	serverAliveConfirmer        = sc.NewServerAliveConfirmer()
	getJsonInfoEndPoint         = "/json/program/all"
)

func PrintJSONToWeb(w http.ResponseWriter, JSONStruct interface{}) {
	bytes, _ := json.MarshalIndent(&JSONStruct, "", "    ")

	_, err := fmt.Fprintf(w, string(bytes))
	if err != nil {
		logf("err from Fprintf(): %v \n", err)
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// GetMinimumMemoryServerHandler は実際に疎通できるサーバの中から使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func GetMinimumMemoryServerHandler(w http.ResponseWriter, r *http.Request) {
	url, err := minimumMemoryServerSelector.Select(cfg.Servers, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	type j struct {
		Url string `json:"url"`
	}

	jsonStr := j{Url: url}

	PrintJSONToWeb(w, jsonStr)
}

// GetSuitableServerHandler は実際にプログラムがあるサーバかつ、使用メモリが最小の
// サーバのURLをJSONで表示するAPI
func GetSuitableServerHandler(w http.ResponseWriter, r *http.Request) {
	programName := r.URL.Path[len("/SuitableServer/"):]
	logf("programName: %v ", programName)

	aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
	if err != nil {
		logf("GetSuitableServerHandler: %v ", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	logf("aliveservers: %v", aliveServers)

	programHasServersGetter := gp.GetProgramHasServersGetter()
	programHasServers, err := programHasServersGetter.Get(aliveServers, getJsonInfoEndPoint, programName)
	if err != nil {
		logf("err from Get(): %v ", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	logf("programHasServers: %v ", programHasServers)

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

	PrintJSONToWeb(w, jsonStr)
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
func UserTopHandler(w http.ResponseWriter, r *http.Request) {
	// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
	allServerMaps := map[string]interface{}{}
	aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
	if err != nil {
		logf("UserTopHandler: %v ", err.Error())
		http.Error(w, err.Error()+"全てのプログラムサーバが生きていない可能性があります。", 500)
		return
	}

	allProgramGetter := gg.NewAllProgramGetter()
	endPoint := "/json/program/all"
	allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

	logf("programs of all pro: %v ", http2.GetKeysFromMap(allServerMaps))

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
			logf("UserTopHandler: %v ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		programHasServersGetter := gp.GetProgramHasServersGetter()
		programHasServers, err := programHasServersGetter.Get(aliveServers, getJsonInfoEndPoint, proName)
		if err != nil {
			logf("UserTopHandler: %v ", err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		// プログラムを保持しているサーバたちの中で一番使用メモリが少ないサーバを選択する。
		url, err := minimumMemoryServerSelector.Select(programHasServers, serverAliveConfirmer, memoryGetter, cfg.GetMemoryEndPoint, "/health")
		if err != nil {
			logf("UserTopHandler: %v ", err.Error())
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

	// 名前順でソートする
	sort.Slice(dataToHtml.ProInfos, func(i, j int) bool { return dataToHtml.ProInfos[i].Name < dataToHtml.ProInfos[j].Name })

	t, err := template.ParseFiles(serveHtml)
	if err != nil {
		logf("UserTopHandler: %v ", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.Execute(w, dataToHtml); err != nil {
		logf("UserTopHandler: %v ", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

func GetAliveServersHandler(w http.ResponseWriter, r *http.Request) {
	aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
	if err != nil {
		logf("err from GetAliveServers(): %v \n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	type data struct {
		AliveServers []string `json:"AliveServers"`
	}

	jsonStr := data{AliveServers: aliveServers}

	PrintJSONToWeb(w, jsonStr)
}

func GetAllProgramsHandler(w http.ResponseWriter, r *http.Request) {
	// allServerMapsはキーにプログラム名が入る。値はプログラム情報のmapが入る。
	allServerMaps := map[string]interface{}{}
	aliveServers, err := sc.GetAliveServers(cfg.Servers, "/health", serverAliveConfirmer)
	if err != nil {
		logf("err from GetAliveServers(): %v \n", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	allProgramGetter := gg.NewAllProgramGetter()
	endPoint := "/json/program/all"
	allServerMaps, err = allProgramGetter.Get(aliveServers, endPoint)

	PrintJSONToWeb(w, allServerMaps)
}
