package arrangeAllProgramJSON_test

import (
	"testing"
	cp "webapi/cli/arrangeAllProgramJSON"
)

func TestAllserverProgramsParser_Parse(t *testing.T) {

	parser := cp.NewAllServerProgramsParser()
	jsonStr := `{
    "toJson": {
        "command": "python3 /Users/hibiki/go/src/webapi/build/server1/programs/convertToJson/convert_json.py INPUTFILE OUTPUTDIR /Users/hibiki/go/src/webapi/build/server1/programs/convertToJson/config.json PARAMETA",
        "help": "拡張子に.jsonをつけて出力します。\n"
    },
    "toZip": {
        "command": "python3 /Users/hibiki/go/src/webapi/build/server2/programs/convertToJson/convert_json.py INPUTFILE OUTPUTDIR /Users/hibiki/go/src/webapi/build/server2/programs/convertToJson/config.json PARAMETA",
        "help": "拡張しに.zipをつけて出力します。\n"
    }
}`

	_, err := parser.Parse(jsonStr)
	if err != nil {
		t.Errorf("err from Parse(): %v \n", err.Error())
	}
}
