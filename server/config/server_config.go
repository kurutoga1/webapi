package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"webapi/utils"
)

var (
	configFilePath string

	// Load()が何回もされるのでattentionMsgFlagで一回にする。
	attentionMsgFlag bool
)

func SetConfPath(confName string) {
	currentDir, err := getCurrentDir()
	if err != nil {
		panic("err msg: " + err.Error())
	}
	configFilePath = filepath.Join(currentDir, "config", confName)
}

// LogConf is log information struct.
type LogConf struct {
	Dir             string `json:"dir"`
	GoLog           string `json:"goLog"`
	RotateMaxKB     int    `json:"rotateMaxKB"`     // ログファイルがこのサイズになったら上の行から削る
	RotateShavingKB int    `json:"rotateShavingKB"` // ログファイルをこのサイズを下回るまで上の行から削る
}

// FileServerConf is relevant fileserver information struct.
type FileServerConf struct {
	Dir           string `json:"dir"`
	UploadDir     string `json:"uploadDir"`
	ProgramOutDir string `json:"programOutDir"`
}

// Config has all config that are LogConf, FileServerConf,,etc.
type Config struct {
	TemplatesDir                       string         `json:"templatesDir"`
	ProgramsDir                        string         `json:"programsDir"`
	ProgramsJSON                       string         `json:"programsJson"`
	Log                                LogConf        `json:"log"`
	FileServer                         FileServerConf `json:"fileServer"`
	DeleteProcessedFileLimitSecondTime int            `json:"DeleteProcessedFileLimitSecondTime"`
	ServerIP                           string         `json:"ServerIP"`
	ServerPort                         string         `json:"ServerPort"`
	ProgramTimeOut                     int            `json:"programTimeOutSecond"`
	StdoutBufferSize                   int            `json:"StdoutBufferSize"`
	StderrBufferSize                   int            `json:"StderrBufferSize"`
	MaxUploadSizeMB                    int64          `json:"MaxUploadSizeMB"`
}

// Load はconfig.jsonの中身をstructに入れたものを返す
func Load() *Config {
	// configFilePath がセットされていない場合は値をセットする。
	if configFilePath == "" {
		SetConfPath("config.json")
	}

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	// 構造体を初期化
	conf := new(Config)

	if !utils.FileExists(configFilePath) {
		// configFilePathが存在しない場合は２通りある。
		// 1つはシンプルにプログラムサーバにconfig.jsonを配置していない場合
		// もう一つはgwをビルドしてこのLoad()関数が呼ばれたけどgwのため、同じディレクトリにconfig.jsonがない場合
		// gwをビルドするときにこのLoad()関数が呼ばれる。理由はGetSomeStartedServerをgwがテストで頻繁に使用しているため
		// GetSomeStartedServerはhandlersパッケージを呼び出し、handlersパッケージは至る箇所でLoad()を実行している
		// gwのテストではプログラムサーバを立てることが必須である。
		// gwはテストはうまくいくがビルドした後、実行するときにconfig.jsonがないので失敗する。
		// なのでgwの実行する際にエラーが出る問題を解決する必要がある。
		// なので警告文を発すると同時に、config.Configに偽の値をセットすることでgwを実行する際にエラーをすり抜ける。
		// Load()が何回もされるのでattentionMsgFlagで一回にする。
		if !attentionMsgFlag {
			fmt.Println(`注意
あなたがプログラムサーバを実行した場合、このエラーが出た場合はconfigディレクトリの中にconfig.jsonがあるか確認してください。
ない場合は再度実行し、このエラーが出なかったら起動成功です。
あなたがゲートウェイサーバ、またはコマンドラインツールを実行したのなら、このエラーは無視してください。`)
		}
		attentionMsgFlag = true
		dummyCfg := &Config{
			TemplatesDir: "aaaa",
		}
		return dummyCfg
	}

	// 設定ファイルを読み込む
	cValue, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// 読み込んだjson文字列をデコードし構造体にマッピング
	err = json.Unmarshal([]byte(cValue), conf)
	if err != nil {
		log.Fatal(err)
	}

	// templatesDir(html)が入っているディレクトリにフルパスを記述する。
	currentDir, err := getCurrentDir()
	if err != nil {
		panic("panic: " + err.Error())
	}
	conf.TemplatesDir = filepath.Join(currentDir, conf.TemplatesDir)

	return conf
}
