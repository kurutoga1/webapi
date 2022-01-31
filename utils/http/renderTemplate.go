package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"text/template"
)

// PrintAsJSON webページにJSONを表示する。
// jsonの元になる構造体を渡す。
func PrintAsJSON(w http.ResponseWriter, jsonStruct interface{}) {
	// jsonに変換
	b, err := json.MarshalIndent(jsonStruct, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// JSONを表示
	_, err = fmt.Fprintf(w, string(b))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// RenderTemplate webページのhtmlにデータを渡し、表示させる
// serverHtmlは絶対パスで記述する。
func RenderTemplate(w http.ResponseWriter, serveHtml string, data interface{}) {
	if !filepath.IsAbs(serveHtml) {
		http.Error(w, "serve html("+serveHtml+") is not found.", http.StatusInternalServerError)
		return
	}
	t, err := template.ParseFiles(serveHtml)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := t.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
