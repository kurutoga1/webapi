package http

import (
	"errors"
	"fmt"
	"net"
	"path/filepath"
	"strconv"
	"strings"
)

// GetURLFromFilePath はファイルサーバとして公開しているディレクトリの下にあるファイルのパスを受け取り
// 外部からそのファイルパスにアクセスするためのURLを返す。
// input: fileserver/something/a.txt -> output: http://localhost:8082/fileserver/something/a.txt
func GetURLFromFilePath(filePath string, ip, port, fileServerDir string) (fileURLPath string, err error) {

	// ファイルサーバディレクトリはファイルサーバとして公開しているため
	// filePathにファイルサーバディレクトリの文字列が入っていないとエラーを出す。
	fileServerStr := filepath.Join(fileServerDir, "")
	if ok := strings.Contains(filePath, fileServerStr); !ok {
		err := errors.New(filePath + "doesn't contain " + fileServerStr)
		return "", fmt.Errorf("GetURLFromFilePath: %v", err)
	}

	basename := filepath.Base(filePath)
	port = ":" + port + "/"

	fileURLPath = "http://" + ip + port + filepath.Join(filepath.Dir(filePath), basename)

	return fileURLPath, nil
}

// GetKeysFromMap マップからキーを取り出し、リストで返す
func GetKeysFromMap(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetUnUsedPort は使用していないポートを取得し、ストリングで返す。
func GetUnUsedPort() (string, error) {
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return "", err
	}
	portInt := listener.Addr().(*net.TCPAddr).Port
	port := strconv.Itoa(portInt)
	err = listener.Close()
	if err != nil {
		return "", err
	}
	return port, nil
}

// GetLoopBackAddrWithUnUsedPort 127.0.0.1:(使用していないport)を返す。
func GetLoopBackAddrWithUnUsedPort() (string, error) {
	ip := "127.0.0.1"
	port, err := GetUnUsedPort()
	if err != nil {
		return "", fmt.Errorf("GetLoopBackAddrWithUnUsedPort: %v", err)
	}
	port = ":" + port
	addr := ip + port
	return addr, nil
}

// GetLoopBackURL http://127.0.0.1:(使用していないport)を返す。
func GetLoopBackURL() (string, error) {
	scheme := "http://"
	addr, err := GetLoopBackAddrWithUnUsedPort()
	if err != nil {
		return "", fmt.Errorf("GetLoopBackURL: %v", err)
	}
	return scheme + addr, nil
}

func GetPortFromURL(url string) string {
	return url[strings.LastIndex(url, ":")+1:]
}
