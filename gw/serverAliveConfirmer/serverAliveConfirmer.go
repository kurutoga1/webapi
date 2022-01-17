/*
サーバが疎通確認ができるか確認する機能を提供するパッケージ
*/

package serverAliveConfirmer

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

// ServerAliveConfirmer サーバのリストを受け取り、実際に疎通確認をし、サーバが起動していればtrue, していなければfalseを返す
// IsAliveを実装したインターフェース。
type ServerAliveConfirmer interface {
	// IsAlive はサーバのURLのを受け取り、
	// 実際にhttp getで疎通確認をし、サーバが起動していればtrue, していなければfalseを返す。
	// arg eg. serverAddr -> http://127.0.0.1:8081
	// arg eg. endPoint -> /user/top
	IsAlive(addr string, endPoint string) (bool, error)
}

func NewServerAliveConfirmer() ServerAliveConfirmer {
	return &confirmer{}
}

type confirmer struct{}

// IsAlive はサーバのURLのを受け取り、
// 実際にhttp getで疎通確認をし、サーバが起動していればtrue, していなければfalseを返す。
// arg eg. serverAddr -> http://127.0.0.1:8081
// arg eg. endPoint -> /user/top
func (c *confirmer) IsAlive(serverAddr string, endPoint string) (bool, error) {
	url := serverAddr + endPoint // eg. http:127.0.0.1:8081/user/top

	resp, err := http.Get(url)
	if err != nil {
		// Getでアクセスしてみて、サーバが立ち上がっていない場合はpanicを起こし、
		// プログラムがエラーになるため、recoverし処理を継続する
		defer func() {
			err := recover()
			if err != nil {
				panic("Panic!!. from IsAlive. recover failed of http get. program exit.")
			}
		}()
		return false, nil
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			err = errors.New("response body failed to close. err msg: " + err.Error())
		}
	}(resp.Body)
	if err != nil {
		return false, fmt.Errorf("IsAlive: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		// server is alive.
		return true, nil
	} else {
		// server is not alive.
		return false, nil
	}

}

// GetAliveServers は[http://127.0.0.1:8081,....]みたいなURLのリストを受け取る
// サーバのリストを受け取り、実際に疎通確認をし、実際に生きているサーバをリストで返す
func GetAliveServers(servers []string, endPoint string, confirmer ServerAliveConfirmer) (aliveServers []string, err error) {
	aliveServers = make([]string, 0, 20)
	for _, addr := range servers {
		alive, err := confirmer.IsAlive(addr, endPoint)
		if err != nil {
			return nil, err
		}

		if alive {
			aliveServers = append(aliveServers, addr)
		}
	}

	if len(aliveServers) == 0 {
		err := errors.New("alive servers are nothing!!.")
		return nil, fmt.Errorf("GetAliveServers: %v", err)
	}
	return
}
