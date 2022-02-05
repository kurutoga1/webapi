/*
ユーザーがWebページを使用し、登録プログラムの実行を行う。
*/

package user

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"webapi/pro/config"
	"webapi/pro/execution/contextManager"
	"webapi/pro/execution/executer"
	"webapi/pro/execution/outputManager"
	http2 "webapi/utils/http"
)

// FileUploadHandler はtemplates/index.htmlを返す。
func FileUploadHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveHtml := filepath.Join(cfg.TemplatesDir, "upload.html")

		w.Header().Add("Content-Type", "text/html")
		http.ServeFile(w, r, serveHtml)
	}
}

// convertNewLine は改行コードを任意の文字列に置換する。
func convertNewLine(str, newStr string) string {
	return strings.NewReplacer(
		"\r\n", newStr,
		"\r", newStr,
		"\n", newStr,
	).Replace(str)
}

// TopHandler はtemplates/index.htmlを返す。
func TopHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveHtml := filepath.Join(cfg.TemplatesDir, "userTop.html")

		_ = r.Body
		w.Header().Add("Content-Type", "text/html")

		proConfList, err := config.GetPrograms()
		if err != nil {
			l.Printf("get programs err: %v", err.Error())
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
				l.Printf("get help err: %v", err)
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

		http2.RenderTemplate(w, serveHtml, d)
		return
	}
}

// PrepareExecHandler はプログラムを実行するためのwebページ
// を表示するハンドラ。
func PrepareExecHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serveHtml := filepath.Join(cfg.TemplatesDir, "prepareExec.html")

		proName := r.FormValue("proName")
		p, err := config.GetProConfByName(proName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		type data struct {
			Name    string
			Command string
		}
		d := data{Name: p.Name(), Command: p.Command()}

		http2.RenderTemplate(w, serveHtml, d)
		return
	}
}

// ExecHandler はプログラムを実行するためのwebページを表示するハンドラ。
// webページからプログラム名、ファイル(multi-part)、パラメタが送られてくる。
// 上記の情報を使用し、コマンドを実行。
// 実行結果をwebページに挿入し、webページを返す。
func ExecHandler(l *log.Logger, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fName := "ExecHandler"

		ctx, err := contextManager.NewContextManager(w, r, cfg)
		if err != nil {
			l.Printf("%v: %v", fName, err.Error())
			http.Error(w, err.Error(), 500)
			return
		}

		var outputInfo outputManager.OutputManager
		var executer executer.Executer = executer.NewExecuter()

		l.Printf("programName: %v", ctx.ProgramName())
		l.Printf("uploadFilePath: %v", ctx.UploadedFilePath())
		l.Printf("parameta: %v", ctx.Parameta())
		l.Printf("command: %v", ctx.Command())
		outputInfo = executer.Execute(ctx)
		l.Printf("ExpectedStatus: %v", outputInfo.Status())
		l.Printf("ErrMsg: %v", outputInfo.ErrorMsg())

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

		serveHtml := filepath.Join(cfg.TemplatesDir, "execResult.html")
		http2.RenderTemplate(w, serveHtml, d)
		return
	}
}
