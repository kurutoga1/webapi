/*
プロセスサーバでプログラムを実行させる機能を提供するパッケージ
*/

package processFileOnServer

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	outLib "webapi/server/outputManager"
	log2 "webapi/utils/log"
)

var Logger *log.Logger

func init() {
	Logger = log.New(new(log2.NullWriter), "", log.Ldate|log.Ltime)
}

// UploadedInfo はアップロードされた情報をサーバーに送る際の構造体
type UploadedInfo struct {
	Filename string `json:"filename"`
	Parameta string `json:"parameta"`
}

type FileProcessor interface {
	Process(url, uploadedBaseName, parameta string) (outLib.OutputManager, error)
}

func NewFileProcessor() FileProcessor {
	return &fileProcessor{}
}

type fileProcessor struct{}

// Process はサーバにアップロードしたファイルを処理させる。
// サーバのurl, アップロードしたuploadedFile、サーバ上でコマンドを実行するためのparametaを受け取る
// 返り値はサーバー内で出力したファイルを取得するためのURLパスのリストを返す。
func (f *fileProcessor) Process(url string, uploadedBasename string, parameta string) (outLib.OutputManager, error) {

	Logger.Printf("url: %v\n", url)
	Logger.Printf("uploadFile: %v\n", uploadedBasename)

	// 値をリクエストボディにセットする
	reqBody := UploadedInfo{Filename: uploadedBasename, Parameta: parameta}

	// リクエストボディをjsonに変換
	requestBody, err := json.Marshal(reqBody)
	Logger.Printf("requestBody: %v\n", string(requestBody))
	if err != nil {
		return &outLib.OutputInfo{}, err
	}
	body := bytes.NewReader(requestBody)

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return &outLib.OutputInfo{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &outLib.OutputInfo{}, err
	}

	defer resp.Body.Close()

	// レスポンスを受け取り、格納する。
	var res *outLib.OutputInfo
	b, err := ioutil.ReadAll(resp.Body)
	Logger.Printf("Response body: %v\r", string(b))
	if err := json.Unmarshal(b, &res); err != nil {
		log.Fatal(err)
	}
	Logger.Printf("res: %v\n", res)

	if res.OutputURLs == nil {
		res.OutputURLs = []string{}
	}

	return res, err
}
