/*
ユーザにレスポンスするメッセージを定義したパッケージ
ここにわかりやすく書いておいて、他で使用する。
*/

package msgs

import (
	"fmt"
)

var (
	NotAllowMethodError = "許可されていないメソッドです."
	UploadSuccess       = "アップロードが成功しました。"
	UploadSizeTooBig    = "アップロードされたファイルが大きすぎます"
)

func UploadFileSizeExceedError(maxUploadFileSize int64) string {
	return fmt.Sprintf("%v %vMB以下のファイルを指定してください", UploadSizeTooBig, maxUploadFileSize)
}

func PassDataToHtml(d interface{}) string {
	return fmt.Sprintf("Pass data to html: %v\n", d)
}

func ServeHtml(htmlFile string) string {
	return fmt.Sprintf("Serve Html: %v\n", htmlFile)
}
