/*
プログラムサーバにアクセスし、メモリの状態を把握し、一番負荷が少ないサーバを選択する機能を提供するパッケージ
*/

package minimumServerSelector

import (
	"fmt"
	"strings"
	"webapi/gw/memoryGetter"
	"webapi/gw/serverAliveConfirmer"
)

// MinimumMemoryServerSelector はサーバのURLリストを受け取り、メモリなどの状態を判断し、一番負荷が少ないサーバのURLを返す。
// arg -> [http:127.0.0.1:8081,http:127.0.0.1:8082], return -> http://127.0.0.1:8082
type MinimumMemoryServerSelector interface {
	// Select はサーバのリストを受け取り、最適なサーバを返す。
	// 最初にサーバが生きているか確認する。
	Select(expectedAliveServers []string, confirmer serverAliveConfirmer.ServerAliveConfirmer, memoryGetter memoryGetter.Getter, getMemoryEndPoint, healthEndPoint string) (string, error)
}

func NewMinimumMemoryServerSelector() MinimumMemoryServerSelector {
	return &minimumMemoryServerSelector{}
}

type minimumMemoryServerSelector struct{}

// Select はサーバのリストを受け取り、最適なサーバを返す。
// 最初にサーバが生きているか確認する。
func (s *minimumMemoryServerSelector) Select(expectedAliveServers []string, confirmer serverAliveConfirmer.ServerAliveConfirmer, memoryGetter memoryGetter.Getter, getMemoryEndPoint, healthEndPoint string) (string, error) {

	for i, serverAddr := range expectedAliveServers {
		if !strings.Contains(serverAddr, "http://") {
			expectedAliveServers[i] = "http://" + serverAddr
		}
	}

	// 生きているサーバを抜き出す。
	aliveServers, err := serverAliveConfirmer.GetAliveServers(expectedAliveServers, healthEndPoint, confirmer)
	if err != nil {
		return "", fmt.Errorf("Select: %v", err)
	}

	/*
		{
			"127.0.0.1:8081": 31212,
			.....
		}
		みたいなサーバとメモリのマップをもらう
	*/
	serverMemoryMap, err := GetServerMemoryMap(aliveServers, getMemoryEndPoint, memoryGetter)
	if err != nil {
		return "", fmt.Errorf("Select: %v", err)
	}

	minUrl := GetMinimumMemoryServer(serverMemoryMap)

	return minUrl, nil
}

// GetServerMemoryMap のサーバとサーバのメモリ使用量をマップにして返す。
// 引数は
// getMemoryInfoUrlsは[http://127.0.0.1:8082/json/health/memory,...]みたいな
// serversは[127.0.0.1:8082,...]みたいな
// 返り値はurl(key): 使用メモリ量(value)のマップを返す。
// eg. http://127.0.0.1:8082: 5120(メモリ量)
// メモリはヒープ上に割り当てられた累積メモリ量をバリューに入れたものを返す。
func GetServerMemoryMap(servers []string, getMemoryEndPoint string, memoryGetter memoryGetter.Getter) (map[string]uint64, error) {
	m := make(map[string]uint64)
	for _, addr := range servers {
		url := addr + getMemoryEndPoint
		r, err := memoryGetter.Get(url) // eg. url -> http://127.0.0.1:8082/json/health/memory
		if err != nil {
			return m, err
		}

		for _, u := range servers {
			// eg. u -> http://127.0.0.1:8082
			if strings.Contains(url, u) {
				m[u] = r.Mallocs
				break
			}
		}
	}
	return m, nil
}

// GetMinimumMemoryServer はurl(key): memory(value)のようなマップを受け取り
// 使用メモリが最小のサーバを判断し,URLを返す。
func GetMinimumMemoryServer(serverMemoryMap map[string]uint64) string {
	var minUrl string
	var minMemory uint64 = 10000000000000000
	for url, memory := range serverMemoryMap {
		if memory < minMemory {
			minUrl = url
			minMemory = memory
		}
	}
	return minUrl
}
