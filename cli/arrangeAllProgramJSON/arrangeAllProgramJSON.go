package arrangeAllProgramJSON

type AllServerProgramsParser interface {
	// Parse ロードバランサから全てのサーバのプログラム情報JSONを取得してそれをうまく表みたいにパースし
	// 文字列を返す。
	Parse(string) (string, error)
}

/*
{
    "toJson": {
        "command": "python3 /Users/hibiki/go/src/webapi/build/server1/programs/convertToJson/convert_json.py INPUTFILE OUTPUTDIR /Users/hibiki/go/src/webapi/build/server1/programs/convertToJson/config.json PARAMETA",
        "help": "拡張子に.jsonをつけて出力します。\n"
    },
    "toZip": {
        "command": "python3 /Users/hibiki/go/src/webapi/build/server2/programs/convertToJson/convert_json.py INPUTFILE OUTPUTDIR /Users/hibiki/go/src/webapi/build/server2/programs/convertToJson/config.json PARAMETA",
        "help": "拡張しに.zipをつけて出力します。\n"
    }
}
*/

func NewAllServerProgramsParser() AllServerProgramsParser {
	return &allserverProgramsParser{}
}

type allserverProgramsParser struct{}

func (a *allserverProgramsParser) Parse(programInfoJSON string) (string, error) {
	/*
		TODO: 時間があるときにこんな感じを作りたい
		+------+-----------------------+--------+
		| NAME |         SIGN          | RATING |
		+------+-----------------------+--------+
		|  A   |       The Good        |    500 |
		|  B   | The Very very Bad Man |    288 |
		|  C   |       The Ugly        |    120 |
		|  D   |      The Gopher       |    800 |
		+------+-----------------------+--------+
	*/

	return programInfoJSON, nil
}
