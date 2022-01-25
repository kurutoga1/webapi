/*
プロセスサーバでプログラムを実行させる機能を提供するパッケージ
*/

package processFileOnServer

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
	Process(url, uploadFilePath, parameta string) (outLib.OutputManager, error)
}

func NewFileProcessor() FileProcessor {
	return &fileProcessor{}
}

type fileProcessor struct{}

// Process はリクエストの中にfile(multi-part)とパラメータを付与し、サーバへ送信する。
// サーバのurl, アップロードしたuploadedFile、サーバ上でコマンドを実行するためのparametaを受け取る
// 返り値はoutLib.OutputManagerインタフェースを返す。
func (f *fileProcessor) Process(url, uploadFilePath, parameta string) (outLib.OutputManager, error) {
	pr, pw := io.Pipe()
	form := multipart.NewWriter(pw)

	go func() {
		defer pw.Close()

		err := form.WriteField("parameta", parameta)
		if err != nil {
			panic(err.Error())
		}

		file, err := os.Open(uploadFilePath)
		if err != nil {
			panic(err.Error())
		}
		w, err := form.CreateFormFile("file", filepath.Base(uploadFilePath))
		if err != nil {
			panic(err.Error())
		}
		_, err = io.Copy(w, file)
		if err != nil {
			panic(err.Error())
		}
		err = form.Close()
		if err != nil {
			panic(err.Error())
		}
	}()

	req, err := http.NewRequest(http.MethodPost, url, pr)
	if err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Process: %v", err)
	}
	req.Header.Set("Content-Type", form.FormDataContentType())

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
