/*
プロセスサーバでプログラムを実行させる機能を提供するパッケージ
*/

package post

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	outLib "webapi/server/outputManager"
	http2 "webapi/utils/http"
)

type Poster interface {
	Post(proName, url, uploadFilePath, parameta string) (outLib.OutputManager, error)
}

func NewPoster() Poster {
	return &fileProcessor{}
}

type fileProcessor struct{}

// Post はリクエストの中にfile(multi-part)とパラメータを付与し、サーバへ送信する。
// サーバのurl, アップロードしたuploadedFile、サーバ上でコマンドを実行するためのparametaを受け取る
// 返り値はoutLib.OutputManagerインタフェースを返す。
func (f *fileProcessor) Post(proName, url, uploadFilePath, parameta string) (outlib outLib.OutputManager, err error) {

	fields := map[string]string{
		"proName":  proName,
		"parameta": parameta,
	}

	poster := http2.NewPostGetter()
	req, err := poster.GetPostRequest(url, uploadFilePath, fields)
	if err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Post: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Post: %v", err)
	}

	defer func(Body io.ReadCloser) {
		err1 := Body.Close()
		if err == nil {
			err = err1
		}
	}(resp.Body)

	// レスポンスを受け取り、格納する。
	var res *outLib.OutputInfo
	b, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(b, &res); err != nil {
		return &outLib.OutputInfo{}, fmt.Errorf("Post: %v, response body from server: %v", err, string(b))
	}

	if res.OutputURLs == nil {
		res.OutputURLs = []string{}
	}

	return res, err
}
