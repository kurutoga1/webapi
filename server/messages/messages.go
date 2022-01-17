/*
ユーザにレスポンスするメッセージを定義したパッケージ
ここにわかりやすく書いておいて、他で使用する。
*/

package messages

import (
	"fmt"
)

var (
	NotAllowMethodError = "許可されていないメソッドです."
	UploadSuccess       = "アップロードが成功しました。"
)

func UploadFileSizeExceedError(maxUploadFileSize int64) string {
	return fmt.Sprintf("アップロードされたファイルが大きすぎます。%vMB以下のファイルを指定してください", maxUploadFileSize)
}

func PassDataToHtml(d interface{}) string {
	return fmt.Sprintf("Pass data to html: %v\n", d)
}

func ServeHtml(htmlFile string) string {
	return fmt.Sprintf("Serve Html: %v\n", htmlFile)
}
