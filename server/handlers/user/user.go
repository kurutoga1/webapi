/*
ユーザーがWebページを使用し、登録プログラムの実行を行う。
*/

package user

import (
	"net/http"
	"path/filepath"
	"strings"
	"text/template"
	"webapi/server/config"
	"webapi/server/execution/contextManager"
	"webapi/server/execution/executer"
	"webapi/server/outputManager"
)

var cfg = config.Load()

// UserFileUploadHandler はtemplates/index.htmlを返す。
func UserFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	serveHtml := filepath.Join(cfg.TemplatesDir, "upload.html")

	w.Header().Add("Content-Type", "text/html")
	http.ServeFile(w, r, serveHtml)
}

// convertNewLine は改行コードを任意の文字列に置換する。
func convertNewLine(str, newStr string) string {
	return strings.NewReplacer(
		"\r\n", newStr,
		"\r", newStr,
		"\n", newStr,
	).Replace(str)
}

// UserTopHandler はtemplates/index.htmlを返す。
func UserTopHandler(w http.ResponseWriter, r *http.Request) {
	serveHtml := filepath.Join(cfg.TemplatesDir, "userTop.html")

	_ = r.Body
	cfg := *config.Load()
	w.Header().Add("Content-Type", "text/html")

	proConfList, err := config.GetPrograms()
	if err != nil {
		logf("get programs err: %v", err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
	type ProInfo struct {
		Name string
		Help string
	}
	pList := make([]ProInfo, 0, 20)
	for _, p := range proConfList {
		help, err := p.Help()
		if err != nil {
			logf("get help err: %v", err)
			http.Error(w, err.Error(), 500)
			return
		}
		// テキストの改行をhtmlのbrタグに変換する。
		helpWithBRtag := convertNewLine(help, "<br />")
		proinfo := ProInfo{Name: p.Name(), Help: helpWithBRtag}
		pList = append(pList, proinfo)
	}

	type data struct {
		ProInfos   []ProInfo
		ServerIP   string
		ServerPort string
	}
	d := data{ProInfos: pList, ServerIP: cfg.ServerIP, ServerPort: cfg.ServerPort}

	t, err := template.ParseFiles(serveHtml)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.Execute(w, d); err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

// PrepareExecHandler はプログラムを実行するためのwebページ
// を表示するハンドラ。
func PrepareExecHandler(w http.ResponseWriter, r *http.Request) {
	serveHtml := filepath.Join(cfg.TemplatesDir, "prepareExec.html")

	proName := r.FormValue("proName")
	p, err := config.GetProConfByName(proName)
	type data struct {
		Name    string
		Command string
	}
	d := data{Name: p.Name(), Command: p.Command()}

	t, err := template.ParseFiles(serveHtml)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.Execute(w, d); err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}

// ExecHandler はプログラムを実行するためのwebページを表示するハンドラ。
// webページからプログラム名、ファイル(multi-part)、パラメタが送られてくる。
// 上記の情報を使用し、コマンドを実行。
// 実行結果をwebページに挿入し、webページを返す。
func ExecHandler(w http.ResponseWriter, r *http.Request) {
	fName := "ExecHandler"
	serveHtml := filepath.Join(cfg.TemplatesDir, "execResult.html")

	ctx, err := contextManager.NewContextManager(w, r, cfg)
	if err != nil {
		logf("%v: %v", fName, err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	var outputInfo outputManager.OutputManager
	var executer executer.Executer = executer.NewExecuter()

	logf("programName: %v", ctx.ProgramName())
	logf("uploadFilePath: %v", ctx.UploadedFilePath())
	logf("parameta: %v", ctx.Parameta())
	logf("command: %v", ctx.Command())
	outputInfo = executer.Execute(ctx)
	logf("Status: %v", outputInfo.Status())
	logf("ErrMsg: %v", outputInfo.ErrorMsg())

	type data struct {
		Name                    string
		OutURLs                 []string
		DownloadLimitSecondTime int
		Result                  string
		Errmsg                  string
		Stdout                  string
		Stderr                  string
	}
	d := data{
		Name:                    ctx.ProgramConfig().Name(),
		OutURLs:                 outputInfo.OutURLs(),
		DownloadLimitSecondTime: cfg.DeleteProcessedFileLimitSecondTime,
		Result:                  outputInfo.Status(),
		Errmsg:                  outputInfo.ErrorMsg(),
		Stdout:                  outputInfo.StdOut(),
		Stderr:                  outputInfo.StdErr(),
	}

	t, err := template.ParseFiles(serveHtml)
	if err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}

	if err := t.Execute(w, d); err != nil {
		logf(err.Error())
		http.Error(w, err.Error(), 500)
		return
	}
}
