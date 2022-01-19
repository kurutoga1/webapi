/*
servers.jsonを読み込み、構造体に保持する。
*/

package config

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"webapi/utils/file"
)

var serversConfigPath string

func SetConfPath(confName string) {
	currentDir, err := file.GetCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	serversConfigPath = filepath.Join(currentDir, confName)
}

type serversConfig struct {
	Servers                []string `json:"ExpectedAliveServers"`
	GetMemoryEndPoint      string   `json:"GetMemoryEndPoint"`
	LoadBalancerServerIP   string   `json:"LoadBalancerServerIP"`
	LoadBalancerServerPort string   `json:"LoadBalancerServerPort"`
	LogPath                string   `json:"LogPath"`
	RotateMaxKB            int      `json:"RotateMaxKB"`
	RotateShavingKB        int      `json:"RotateShavingKB"`
}

// NewServerConfig はservers.jsonの中身をserversConfig構造体にセットし、返す
func NewServerConfig() *serversConfig {
	if serversConfigPath == "" {
		SetConfPath("config.json")
	}
	// 構造体を初期化
	conf := &serversConfig{}

	// 設定ファイルを読み込む
	cValue, err := ioutil.ReadFile(serversConfigPath)
	if err != nil {
		panic(err.Error())
	}

	// 読み込んだjson文字列をデコードし構造体にマッピング
	err = json.Unmarshal([]byte(cValue), conf)
	if err != nil {
		panic(err.Error())
	}

	return conf
}
