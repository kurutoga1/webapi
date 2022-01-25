package program_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"webapi/server/execution/msgs"
	sh "webapi/server/handlers"
	"webapi/server/handlers/program"
	"webapi/server/outputManager"
	"webapi/utils/file"
	http2 "webapi/utils/http"
	u "webapi/utils/upload"
)

var (
	currentDir string
	uploadFile string
)

func set() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	uploadFile = "uploadfile"

	addr, err := http2.GetLoopBackURL()
	if err != nil {
		log.Fatalln(err.Error())
	}
	port := http2.GetPortFromURL(addr)

	go func() {
		if err := http.ListenAndServe(":"+port, sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()

	err = file.CreateSpecifiedFile(uploadFile, 200)
	if err != nil {
		panic(err.Error())
	}

	uploader := u.NewUploader()
	err = uploader.Upload(addr+"/upload", uploadFile)
	if err != nil {
		panic(err.Error())
	}
}

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestProgramHandler(t *testing.T) {
	set()
	// 保持しているプログラムの場合
	programName := "convertToJson"
	t.Run("Success test", func(t *testing.T) {
		rr, out := testProgramHandler(t, programName)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if len(out.OutputURLs) != 1 {
			t.Errorf("len(out.OutputURLs): %v, want: 1", len(out.OutputURLs))
		}

		if out.Stdout == "" {
			t.Errorf("out.Stdout is empty.")
		}

		if out.Stderr != "" {
			t.Errorf("out.Stderr is not empty.")
		}

		if out.StaTus != msgs.OK {
			t.Errorf("out.StaTus: %v, want : %v", out.StaTus, msgs.SERVERERROR)
		}

		if out.Errormsg != "" {
			t.Errorf("out.Errormsg: %v, want: %v", out.Errormsg, "")
		}
	})

	// 保持していないプログラム名の場合
	programName = "nothingProgram"
	t.Run("fail test", func(t *testing.T) {
		rr, out := testProgramHandler(t, programName)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, http.StatusOK)
		}

		if len(out.OutputURLs) != 0 {
			t.Errorf("len(out.OutputURLs): %v, want: 0", len(out.OutputURLs))
		}

		if out.Stdout != "" {
			t.Errorf("out.Stdout is not empty.")
		}

		if out.Stderr != "" {
			t.Errorf("out.Stderr is not empty.")
		}

		if out.StaTus != msgs.SERVERERROR {
			t.Errorf("out.StaTus: %v, want : %v", out.StaTus, msgs.SERVERERROR)
		}

		expected := programName + " is not found."
		if out.Errormsg != expected {
			t.Errorf("out.Errormsg: %v, want: %v", out.Errormsg, expected)
		}
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testProgramHandler(t *testing.T, proName string) (*httptest.ResponseRecorder, *outputManager.OutputInfo) {
	uploadFile = "uploadfile"
	err := file.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	fields := map[string]string{
		"proName":  "convertToJson",
		"parameta": "dummyParameta",
	}
	r, err := http2.GetPostRequestWithFileAndFields(uploadFile, "/pro/"+proName, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	handler := http.HandlerFunc(program.ProgramHandler)

	handler.ServeHTTP(w, r)
	var out *outputManager.OutputInfo

	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Errorf("jsonUnmarshal fail(msg: %v). body is %v", err.Error(), w.Body.String())
	}

	return w, out
}

func TestProgramAllHandler(t *testing.T) {
	set()
	req, err := http.NewRequest("GET", "/json/program/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(program.ProgramAllHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	expected := `{
    "convertToJson": {
        "command": "python3 /Users/hibiki/go/src/webapi/server/config/programs/convertToJson/convert_json.py INPUTFILE OUTPUTDIR /Users/hibiki/go/src/webapi/server/config/programs/convertToJson/config.json PARAMETA",
        "help": "入力ファイルに拡張子をつけて出力します\nあまりサイズが大きすぎるファイルを処理できません。"
    }
}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
