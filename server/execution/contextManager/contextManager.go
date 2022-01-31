/*
contextManager.go
入力ファイルや出力ディレクトリ、パラメータなど登録プログラムの実行に必要な
情報等を保持、用意し、executer.Execute()に渡し、実行してもらう。
本来ならばたくさんのパラメータをExecute()に渡さなければいけないがその
たくさんのパラメータを全てこのパッケージのContextManagerインタフェースが保持する。
*/

package contextManager

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"webapi/server/config"
	"webapi/server/handlers/upload"
	"webapi/utils/file"
	log2 "webapi/utils/log"
	utils2 "webapi/utils/string"
)

var Logger *log.Logger

func init() {
	Logger = log.New(new(log2.NullWriter), "", log.Ldate|log.Ltime)
}

// ContextManager はコマンド実行に必要な要素を持つ構造体のインタフェース
// コマンド実行に必要なパラメータ等を一括で管理する構造体のインタフェース
// コマンド実行に必要なものは全てこの中に入れる。
type ContextManager interface {
	// InputFilePath 登録プログラムに処理させる入力ファイルを返す。
	InputFilePath() string
	// SetInputFilePath move fileserver/upload/a.txt fileserver/programOut/(program name)/(random str)/a.txt
	SetInputFilePath() error

	// Command 登録してあるコマンドを返す。
	Command() string
	// SetCommand programConfigHolderを受け取り、登録してあるコマンドを設定する。
	SetCommand() error

	// ProgramName 実行するプログラム名を返す
	ProgramName() string
	// SetProgramName 実行するプログラム名を設定する。
	SetProgramName(string)

	// ProgramTempDir プログラムが作業するテンポラリディレクトリ
	ProgramTempDir() string
	// SetProgramTempDir プログラムが作業するテンポラリディレクトリを取得する
	SetProgramTempDir(proConf config.ProgramConfigHolder, cfg *config.Config) error

	// OutputDir プログラムがファイルを出力するディレクトリを返す。
	OutputDir() string
	// SetOutputDir プログラムがファイルを出力するディレクトリを設定する。
	SetOutputDir(string)

	// UploadedFilePath アップロードされたファイルパスを返す。
	UploadedFilePath() string
	// SetUploadedFilePath アップロードされたファイルパスを設定する
	SetUploadedFilePath(string) // webからの場合に使用する

	// Parameta コマンド実行する際に使用するパラメータを返す。
	Parameta() string
	// SetParameta コマンド実行する際に使用するパラメータを設定する。
	SetParameta(string)

	// ProgramConfig 登録プログラムの情報インターフェースを返す。
	ProgramConfig() config.ProgramConfigHolder
	// SetProgramConfig 登録プログラムの情報インターフェースを設定する。
	SetProgramConfig(holder config.ProgramConfigHolder)

	// Config サーバの設定値を保持する
	Config() *config.Config
	SetConfig(cfg *config.Config)
}

// NewContextManager はcontextManagerを返す。
// プログラム名とプログラム出力ディレクトリはセットする。
// それ以外に必要な要素は定義した後で設定し、executerに渡す感じ
func NewContextManager(w http.ResponseWriter, r *http.Request, cfg *config.Config) (ContextManager, error) {
	fName := "NewContextManager"
	ctx := &contextManager{}

	// file(multi-data)をこのサーバのfileserver/uploadにアップロードする。
	uploadFilePath, err := upload.Upload(w, r, cfg)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	if !file.FileExists(uploadFilePath) {
		return nil, fmt.Errorf("%v: %v", fName, errors.New(uploadFilePath+" is not found."))
	}
	ctx.SetUploadedFilePath(uploadFilePath)

	ctx.SetParameta(r.FormValue("parameta"))
	programName := r.FormValue("proName")

	// プログラム名を取得する処理
	// webからの場合はFormValueにプログラム名が乗ってくるが
	// cliからの場合はURLのエンドポイントに乗ってくるため処理を分ける必要がある。
	if programName == "" {
		programName = r.URL.Path[len("/pro/"):]
	}
	ctx.SetProgramName(programName)

	proConf, err := config.GetProConfByName(programName)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}
	ctx.SetProgramConfig(proConf)

	err = ctx.SetProgramTempDir(proConf, cfg)
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	err = ctx.SetProgramOutDir()
	if err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	if err = ctx.SetInputFilePath(); err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}
	if err = ctx.SetCommand(); err != nil {
		return nil, fmt.Errorf("%v: %v", fName, err)
	}

	ctx.SetConfig(cfg)

	return ctx, nil
}

type contextManager struct {
	programName      string
	programTempDir   string
	outputDir        string
	parameta         string
	uploadedFilePath string
	inputFilePath    string
	command          string
	programConfig    config.ProgramConfigHolder
	stdOutBufferSize int
	stdErrBufferSize int
	cfg              *config.Config
}

// SetProgramOutDir はプログラムが出力するディレクトリを作成、セットする。
// SetProgramTempDirの後に実行されなければいけない
func (c *contextManager) SetProgramOutDir() error {
	if c.programTempDir == "" {
		return fmt.Errorf("c.programTempDir is not set. you shoud SetProgrtamTempDir before me.")
	}
	randomName := utils2.GetRandomString(20) + "_out"
	programOutDir := filepath.Join(c.programTempDir, randomName)
	c.outputDir = programOutDir
	err := os.MkdirAll(programOutDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("SetProgramOutDir: %v", err)
	}
	return nil
}

// SetProgramTempDir 登録プログラムが使用する作業ディレクトリとして使用する日付とランダム文字列から作成されたテンポラリディレクトリを作成、セットする
func (c *contextManager) SetProgramTempDir(proConf config.ProgramConfigHolder, cfg *config.Config) error {
	outDirName := utils2.GetNowTimeStringWithHyphen() + "-" + utils2.GetRandomString(20)
	programDir := filepath.Join(cfg.FileServer.Dir, cfg.FileServer.ProgramOutDir, proConf.Name(), outDirName)
	c.programTempDir = programDir
	err := os.MkdirAll(programDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("SetProgramTempDir: %v", err)
	}
	return nil
}

func (c *contextManager) InputFilePath() string { return c.inputFilePath }

func (c *contextManager) Command() string { return c.command }

func (c *contextManager) ProgramName() string     { return c.programName }
func (c *contextManager) SetProgramName(s string) { c.programName = s }

func (c *contextManager) ProgramTempDir() string { return c.programTempDir }

func (c *contextManager) OutputDir() string     { return c.outputDir }
func (c *contextManager) SetOutputDir(s string) { c.outputDir = s }

func (c *contextManager) Parameta() string     { return c.parameta }
func (c *contextManager) SetParameta(s string) { c.parameta = s }

func (c *contextManager) UploadedFilePath() string     { return c.uploadedFilePath }
func (c *contextManager) SetUploadedFilePath(s string) { c.uploadedFilePath = s }

func (c *contextManager) ProgramConfig() config.ProgramConfigHolder      { return c.programConfig }
func (c *contextManager) SetProgramConfig(pc config.ProgramConfigHolder) { c.programConfig = pc }

func (c *contextManager) Config() *config.Config       { return c.cfg }
func (c *contextManager) SetConfig(cfg *config.Config) { c.cfg = cfg }

// SetInputFilePath uploadの中のファイルをfileserver/programOut/convertToJson/xxxxxx/の中に移動させる。
// move fileserver/upload/a.txt fileserver/programOut/(program name)/(random str)/a.txt
func (c *contextManager) SetInputFilePath() error {
	inputFilePath := filepath.Join(filepath.Dir(c.OutputDir()), filepath.Base(c.UploadedFilePath()))
	if err := os.Rename(c.UploadedFilePath(), inputFilePath); err != nil {
		return fmt.Errorf("SetInputFilePath: %v", err)
	}
	c.inputFilePath = inputFilePath
	if !file.FileExists(c.inputFilePath) {
		return fmt.Errorf("SetInputFilePath: %v", errors.New(c.inputFilePath+"is not found."))
	}
	return nil
}

// SetCommand templateコマンドからINPUTFILE,OUTPUTDIR, PARAMETAなどをreplaceして正規コマンドを作成する。
func (c *contextManager) SetCommand() error {
	if c.inputFilePath == "" {
		return fmt.Errorf("SetCommand: %v", errors.New("c.inputFilePath is empty. should SetInputFilePath() before SetCommand()."))
	}
	c.command = c.programConfig.ReplacedCmd(c.inputFilePath, c.OutputDir(), c.Parameta())
	return nil
}
