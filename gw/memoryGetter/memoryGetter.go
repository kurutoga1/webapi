/*
サーバにアクセスし、ランタイム構造体を取得する機能を提供するパッケージ
*/

package memoryGetter

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
)

type Getter interface {
	Get(string) (runtime.MemStats, error)
}

func NewMemoryGetter() Getter {
	return &getter{}
}

type getter struct{}

// Get はメモリ状況のjsonAPIを公開しているサーバのURLを受け取り、
// get(http)し、そのjsonをruntime.MemStatsにデコードし、runtime.MemStatsを返す
// eg url: "http://127.0.0.1:8093/json/health/memory"
func (g *getter) Get(url string) (runtime.MemStats, error) {
	resp, _ := http.Get(url)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalln(err.Error())
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)

	var d runtime.MemStats
	err := json.Unmarshal(body, &d)
	if err != nil {
		return runtime.MemStats{}, fmt.Errorf("Get: %v", err)
	}

	return d, nil
}
