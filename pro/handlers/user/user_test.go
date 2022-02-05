package user_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"webapi/pro/config"
	"webapi/pro/handlers/user"
	"webapi/utils/file"
	http2 "webapi/utils/http"
)

var (
	uploadFile string
)

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestUserTopHandler(t *testing.T) {
	r, _ := http.NewRequest(http.MethodGet, "/user/top", nil)
	w := httptest.NewRecorder()

	handler := user.TopHandler(log.New(os.Stdout, "", log.LstdFlags), config.Load())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	programName := "convertToJson"
	if !strings.Contains(w.Body.String(), programName) {
		t.Errorf("html doesn't have %v", programName)
	}
}

func TestPrepareExecHandler(t *testing.T) {
	form := url.Values{
		"proName": []string{"convertToJson"},
	}
	r, err := http.NewRequest("POST", "/user/perpareExec", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err.Error())
	}

	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()

	handler := user.PrepareExecHandler(config.Load())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}
}

func TestExecHandler(t *testing.T) {
	// TODO: most important test
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
	r, err := poster.GetPostRequest("/pro/convertToJson", uploadFile, fields)
	if err != nil {
		panic(err.Error())
	}
	w := httptest.NewRecorder()

	handler := user.ExecHandler(log.New(os.Stdout, "", log.LstdFlags), config.Load())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	expected := "<p>結果: ok</p>"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("response.Body doesn't contain %v", expected)
	}

	tearDown()
}
