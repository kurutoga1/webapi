/*
サーバにアクセスし、プログラムが保持されているか確認する機能を提供するパッケージ
*/

package programHasServers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ProgramHasServersGetter interface {
	// Get 生きているサーバリストとプログラム名を受け取り、生きているサーバの中から
	// そのプログラムを持っているサーバを見つけて、リストで返す。
	Get(aliveServers []string, endPoint string, programName string) ([]string, error)
}

func GetProgramHasServersGetter() ProgramHasServersGetter {
	return &programHasServersGetter{}
}

type programHasServersGetter struct{}

// Get 生きているサーバリストとプログラム名を受け取り、生きているサーバの中から
// そのプログラムを持っているサーバを見つけて、リストで返す。
func (p *programHasServersGetter) Get(aliveServers []string, endPoint, programName string) ([]string, error) {
	programHasServers := make([]string, 0, 20)
	// 生きているサーバたちにプログラムを持っているか聞いていく
	for _, url := range aliveServers {
		URLToGetProgramInfo := url + endPoint
		ok, err := IsProgramHasServer(URLToGetProgramInfo, programName)
		if err != nil {
			return nil, fmt.Errorf("Get: %v \n", err)
		}
		// サーバが該当プログラムを保持している場合
		if ok {
			programHasServers = append(programHasServers, url)
		}
	}
	return programHasServers, nil
}

// IsProgramHasServer はプログラム名を受け取り、そのサーバにプログラム名が登録されているか確認する。
// eg. url: http://127.0.0.1:8000/json/program/all
func IsProgramHasServer(url, programName string) (bool, error) {

	// サーバにアクセス
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("IsProgramhasServer: %v", err)
	}

	byteArray, _ := ioutil.ReadAll(resp.Body)

	// レスポンスをパース
	decoded := map[string]interface{}{}
	err = json.Unmarshal(byteArray, &decoded)
	if err != nil {
		return false, fmt.Errorf("IsProgramHasServer: %v", err)
	}

	// プログラム名はあるのか
	for k, _ := range decoded {
		if k == programName {
			return true, nil
		}
	}

	return false, nil
}
