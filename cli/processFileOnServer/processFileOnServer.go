/*
プロセスサーバでプログラムを実行させる機能を提供するパッケージ
*/

package processFileOnServer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	outLib "webapi/server/outputManager"
	http2 "webapi/utils/http"
	log2 "webapi/utils/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New(new(log2.NullWriter), "", log.Ldate|log.Ltime)
}

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
func (f *fileProcessor) Process(proName, url, uploadFilePath, parameta string) (outLib.OutputManager, error) {

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

	defer resp.Body.Close()

	// レスポンスを受け取り、格納する。
	var res *outLib.OutputInfo
	b, err := ioutil.ReadAll(resp.Body)
	Logger.Printf("Response body: %v\r", string(b))
	if err := json.Unmarshal(b, &res); err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Process: %v", err)
	}
	Logger.Printf("res: %v\n", res)

	if res.OutputURLs == nil {
		res.OutputURLs = []string{}
	}

	return res, err
}
