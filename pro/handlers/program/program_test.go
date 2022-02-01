package program_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"webapi/pro/config"
	"webapi/pro/execution/outputManager"
	"webapi/pro/handlers/program"
	"webapi/pro/msgs"
	proRouter "webapi/pro/router"
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
		if err := http.ListenAndServe(":"+port, proRouter.New().New("fileserver")); err != nil {
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
	poster := http2.NewPostGetter()
	r, err := poster.GetPostRequest("/pro/"+proName, uploadFile, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	handler := program.ProgramHandler(log.New(os.Stdout, "", log.LstdFlags), config.Load())

	handler.ServeHTTP(w, r)
	var out *outputManager.OutputInfo

	if err := json.Unmarshal(w.Body.Bytes(), &out); err != nil {
		t.Errorf("jsonUnmarshal fail(msg: %v). body is %v", err.Error(), w.Body.String())
	}

	return w, out
}

func TestProgramAllHandler(t *testing.T) {
	set()
	r, err := http.NewRequest("GET", "/json/program/all", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(program.ProgramAllHandler)

	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got: %v want: %v",
			status, http.StatusOK)
	}

	var m map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &m)
	if err != nil {
		t.Errorf("%v is not json format. err msg: %v", w.Body.String(), err.Error())
	}

	expectedProgramNames := []string{"convertToJson", "err", "sleep"}
	for _, name := range expectedProgramNames {
		if !strings.Contains(w.Body.String(), name) {
			t.Errorf("%v is not contaned of %v", name, w.Body.String())
		}
	}

	t.Cleanup(func() {
		tearDown()
	})
}
