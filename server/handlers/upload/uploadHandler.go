/*
アップロードされるハンドラーを定義したファイル。
*/

package upload

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"webapi/server/config"
	msg "webapi/server/messages"
	int2 "webapi/utils/int"
)

var (
	cfg config.Config = *config.Load()
	// maxUploadSize はアップロードするファイルの上限の大きさ
	//maxUploadSize int64 = cfg.MaxUploadSizeMB << 20 // 20がメガ表記になる。
	maxUploadSize int64 = int64(int2.MBToByte(int(cfg.MaxUploadSizeMB)))
)

// UploadHandler はファイルをアップロードするためのハンドラー。
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	_, err := Upload(w, r)
	if err != nil {
		msg := fmt.Sprintf("UploadHandler: %v, err msg: %v", msg.UploadFileSizeExceedError(cfg.MaxUploadSizeMB), err.Error())
		logf(msg)
		http.Error(w, msg, 500)
		return
	}

	_, err = fmt.Fprintf(w, msg.UploadSuccess)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	return
}

// Upload はファイルをアップロードするためのハンドラー。
func Upload(w http.ResponseWriter, r *http.Request) (string, error) {
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Upload: %v ", errors.New(r.Method+" is not allowed."))
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return "", fmt.Errorf("%v, Upload: %v", err, msg.UploadFileSizeExceedError(cfg.MaxUploadSizeMB))
	}

	//FormFileの引数はHTML内のform要素のnameと一致している必要があります
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return "", fmt.Errorf("Upload: %v", err)
	}

	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			return
		}
		return
	}(file)

	// 存在していなければ、保存用のディレクトリを作成します。
	uploadDir := filepath.Join(cfg.FileServer.Dir, cfg.FileServer.UploadDir)
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("Upload: %v", err)
	}

	// 保存用ディレクトリ内に新しいファイルを作成します。
	uploadFilePath := filepath.Join(uploadDir, fileHeader.Filename)
	logf("uploadFilePath: %v", uploadFilePath)
	dst, err := os.Create(uploadFilePath)
	if err != nil {
		return "", fmt.Errorf("Upload: %v", err)
	}

	defer func(dst *os.File) {
		err := dst.Close()
		if err != nil {
			return
		}
		return
	}(dst)

	// アップロードされたファイルを先程作ったファイルにコピーします。
	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("Upload: %v", err)
	}

	logf(msg.UploadSuccess)
	return uploadFilePath, nil

}
