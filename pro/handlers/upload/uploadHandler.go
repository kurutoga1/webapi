/*
アップロードされるハンドラーを定義したファイル。
*/

package upload

import (
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"webapi/pro/config"
	msg "webapi/pro/msgs"
	int2 "webapi/utils/int"
	utilString "webapi/utils/string"
)

// Handler はファイルをアップロードするためのハンドラー。
func Handler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := Upload(w, r, cfg)
		if err != nil {
			msg := fmt.Sprintf("Handler: %v, err msg: %v", msg.UploadFileSizeExceedError(cfg.MaxUploadSizeMB), err.Error())
			l.Printf(msg)
			http.Error(w, msg, 500)
			return
		}

		_, err = fmt.Fprintf(w, msg.UploadSuccess)
		if err != nil {
			l.Printf(err.Error())
			http.Error(w, err.Error(), 500)
			return
		}
		return
	}
}

var FileSizeTooBigError = errors.New("upload file size is too big.")

// Upload はファイルをアップロードするためのハンドラー。
func Upload(w http.ResponseWriter, r *http.Request, cfg *config.Config) (string, error) {
	maxUploadSize := int64(int2.MBToByte(int(cfg.MaxUploadSizeMB)))
	if r.Method != http.MethodPost {
		return "", fmt.Errorf("Upload: %v ", errors.New(r.Method+" is not allowed."))
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return "", fmt.Errorf("%w, file size: %v", FileSizeTooBigError, msg.UploadFileSizeExceedError(cfg.MaxUploadSizeMB))
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
	// アップロードファイルに半角や全角のスペースがある場合は削除する。
	spaceRemovedUploadFileName := utilString.RemoveSpace(fileHeader.Filename)
	uploadFilePath := filepath.Join(uploadDir, spaceRemovedUploadFileName)
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
		return "", fmt.Errorf("Upload: %w", err)
	}

	return uploadFilePath, nil

}
