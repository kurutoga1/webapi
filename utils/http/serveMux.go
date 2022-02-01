package http

import "net/http"

// MyServeMux fileServerDirをファイルサーバとしたmuxを返すインタフェース。
// プログラムサーバとAPIGWサーバはこのインターフェースを実装している。
type MyServeMux interface {
	New(fileServerDir string) *http.ServeMux
}
