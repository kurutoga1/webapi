/*
url,ファイルパスを受け取り、ファイルをサーバにアップロードする機能を提供するパッケージ
*/

package upload

import (
	"fmt"
	"strings"
	utils2 "webapi/utils"
)

type Uploader interface {
	// Upload アップロードするファイルをURLを受け取り、アップロードする。
	Upload(url string, uploadFilePath string) error
}

func NewUploader() Uploader {
	return &uploader{}
}

type uploader struct{}

func (u *uploader) Upload(url string, uploadFilePath string) error {
	command := fmt.Sprintf("curl -X POST -F file=@%v %v", uploadFilePath, url)
	stdout, stderr, err := utils2.Exec(command)
	if strings.Contains(stdout, "request body too large") || err != nil {
		return fmt.Errorf("Upload: stdout: %v \n stderr: %v ", stdout, stderr)
	}
	return nil
}
