/*
プログラムサーバ内で登録プログラムを実行するが、その出力を格納するライブラリを記述している。
*/

package outputManager

import (
	"fmt"
	"io"
	"webapi/utils/file"
)

// OutputManager はコマンド実行結果を格納した構造体のインターフェース
// 登録プログラムの実行はエラーも含めてこのインターフェースに入れて管理する。そのまま出力のjsonにする。
type OutputManager interface {
	// OutURLs 登録プログラムにより出力されたファイルのリスト(URL)
	OutURLs() []string
	SetOutURLs([]string)

	// StdOut 登録プログラムの標準出力(string)
	StdOut() string
	SetStdOut(io.Reader, int) error

	// StdErr 登録プログラムの標準出力(string)
	StdErr() string
	SetStdErr(io.Reader, int) error

	// Status プログラムサーバ内の挙動、または登録プログラムの挙動によって、ユーザに返すステータスを定義したもの
	Status() string
	SetStatus(string)

	// ErrorMsg プログラムサーバ内でエラーが起きた場合のメッセージ
	ErrorMsg() string
	SetErrorMsg(string)
}

// NewOutputManager はOutputManager(interface)を返す。
func NewOutputManager() OutputManager {
	return &OutputInfo{}
}

// OutputInfo OutputManager(interface)を実装したもの。
// クライアントプログラムで使用しているため外部に公開している。
type OutputInfo struct {
	OutputURLs []string `json:"outURLs"`
	Stdout     string   `json:"stdout"`
	Stderr     string   `json:"stderr"`
	StaTus     string   `json:"status"`
	Errormsg   string   `json:"errmsg"`
}

func (o *OutputInfo) OutURLs() []string { return o.OutputURLs }
func (o *OutputInfo) StdOut() string    { return o.Stdout }
func (o *OutputInfo) StdErr() string    { return o.Stderr }
func (o *OutputInfo) Status() string    { return o.StaTus }
func (o *OutputInfo) ErrorMsg() string  { return o.Errormsg }

func (o *OutputInfo) SetOutURLs(s []string) { o.OutputURLs = s }

// SetStdOut io.ReaderとbufferSizeをもらい、bufferSizeがマックスとして読み込み、
// セットする。
func (o *OutputInfo) SetStdOut(r io.Reader, bufferSize int) error {
	stdout, err := file.ReadBytesWithSize(r, bufferSize)
	if err != nil {
		return fmt.Errorf("SetStdOut: %v", err)
	}
	o.Stdout = stdout
	return nil
}

// SetStdErr io.ReaderとbufferSizeをもらい、bufferSizeがマックスとして読み込み、
// セットする。
func (o *OutputInfo) SetStdErr(r io.Reader, bufferSize int) error {
	stderr, err := file.ReadBytesWithSize(r, bufferSize)
	if err != nil {
		return fmt.Errorf("SetStdErr: %v", err)
	}
	o.Stderr = stderr
	return nil
}

func (o *OutputInfo) SetStatus(s string)   { o.StaTus = s }
func (o *OutputInfo) SetErrorMsg(s string) { o.Errormsg = s }
