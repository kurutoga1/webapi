/*
ゲートウェイにアクセスし、プログラムがあるサーバの中で一番メモリ消費が少ないサーバを取得する機能を提供するパッケージ
*/

package selectServer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type ServerSelector interface {
	// Select はロードバランサにプログラム名をアクセスし、プログラムがあるサーバたちの中で
	// 一番メモリ消費が少ないサーバを選択し、返す。
	// eg url: http://127.0.0.1:8001/SuitableServer/convertToJson
	Select(url string) (addr string, err error)
}

func NewServerSelector() ServerSelector {
	return &selector{}
}

type selector struct{}

// Select はロードバランサにプログラム名をアクセスし、プログラムがあるサーバたちの中で
// 一番メモリ消費が少ないサーバを選択し、返す。
// eg url: http://127.0.0.1:8001/SuitableServer/convertToJson
func (s *selector) Select(url string) (addr string, err error) {
	resp, _ := http.Get(url)

	defer func(Body io.ReadCloser) {
		if err == nil {
			err = Body.Close()
		}
	}(resp.Body)

	byteArray, _ := ioutil.ReadAll(resp.Body)

	type j struct {
		Url string `json:"url"`
	}

	var decoded j
	err = json.Unmarshal(byteArray, &decoded)
	if err != nil {
		return "", fmt.Errorf("Select: %v, response from APIGW server: %v", err, string(byteArray))
	}

	return decoded.Url, nil
}
