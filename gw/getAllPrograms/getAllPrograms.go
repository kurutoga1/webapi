/*
生きているサーバの全てのプログラム一覧を取得する。
{
	"programName" : {
		"command": "xxxxx",
		"help": "xxxxxxx"
	},
	......
}
*/

package getAllPrograms

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type AllProgramGetter interface {
	// Get は生きているサーバたちのリストを入れて、全てのサーバにアクセスし、
	// 全てのプログラム情報を取得しmapで返す。
	Get(aliveServers []string, endPoint string) (map[string]interface{}, error)
}

func NewAllProgramGetter() AllProgramGetter {
	return &allProgramGetter{}
}

type allProgramGetter struct{}

func (a *allProgramGetter) Get(aliveServers []string, endPoint string) (map[string]interface{}, error) {
	allServerMaps := map[string]interface{}{}

	// aliveServersにアクセスしていき、プログラム情報を取得、
	// 全てのプログラム情報を取得し、allServerMapsに格納する。
	for _, server := range aliveServers {
		url := server + endPoint

		// サーバにアクセス
		resp, _ := http.Get(url)
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(resp.Body)

		// レスポンスをパース
		resBody, _ := ioutil.ReadAll(resp.Body)
		mapForResBody := map[string]interface{}{}
		err := json.Unmarshal(resBody, &mapForResBody)
		if err != nil {
			return nil, fmt.Errorf("Get: %v", err)
		}

		// mapForResBodyに一つのサーバから取り出したプログラム一覧(map[string]interface{}{})がある。
		for proName, proInfo := range mapForResBody {
			// allServerMapsにプログラムネームのキーが入っていなかったら追加する
			if _, ok := allServerMaps[proName]; !ok {
				m, _ := proInfo.(map[string]interface{})
				allServerMaps[proName] = m
			}
		}
	}

	return allServerMaps, nil
}
