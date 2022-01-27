/*
プロセスサーバでプログラムを実行させる機能を提供するパッケージ
*/

package processFileOnServer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	outLib "webapi/server/outputManager"
	http2 "webapi/utils/http"
)

type FileProcessor interface {
	Process(proName, url, uploadFilePath, parameta string) (outLib.OutputManager, error)
}

func NewFileProcessor() FileProcessor {
	return &fileProcessor{}
}

type fileProcessor struct{}

// Process はリクエストの中にfile(multi-part)とパラメータを付与し、サーバへ送信する。
// サーバのurl, アップロードしたuploadedFile、サーバ上でコマンドを実行するためのparametaを受け取る
// 返り値はoutLib.OutputManagerインタフェースを返す。
func (f *fileProcessor) Process(proName, url, uploadFilePath, parameta string) (outlib outLib.OutputManager, err error) {

	fields := map[string]string{
		"proName":  proName,
		"parameta": parameta,
	}
	req, err := http2.GetPostRequestWithFileAndFields(uploadFilePath, url, fields)
	if err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Process: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Process: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	// レスポンスを受け取り、格納する。
	var res *outLib.OutputInfo
	b, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Process: %v", err)
	}

	if res.OutputURLs == nil {
		res.OutputURLs = []string{}
	}

	return res, err
}
