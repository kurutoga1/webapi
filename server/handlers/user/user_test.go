package user_test

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"webapi/server/handlers/user"
	http2 "webapi/tests"
	"webapi/utils/file"
)

var (
	uploadFile string
)

func tearDown() {
	os.RemoveAll("fileserver")
	os.Remove(uploadFile)
}

func TestUserTopHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/user/top", nil)
	response := httptest.NewRecorder()

	user.UserTopHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	programName := "convertToJson"
	if !strings.Contains(response.Body.String(), programName) {
		t.Errorf("html doesn't have %v", programName)
	}
}

func TestPrepareExecHandler(t *testing.T) {
	form := url.Values{
		"proName": []string{"convertToJson"},
	}
	req, err := http.NewRequest("POST", "/user/perpareExec", strings.NewReader(form.Encode()))
	if err != nil {
		panic(err.Error())
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	user.PrepareExecHandler(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}
}

func TestExecHandler(t *testing.T) {
	uploadFile = "uploadfile"
	err := file.CreateSpecifiedFile(uploadFile, 2)
	if err != nil {
		panic(err.Error())
	}

	w, r, err := http2.MainRequest(uploadFile, "convertToJson", "dummyParameta", "web")
	if err != nil {
		panic(err.Error())
	}

	user.ExecHandler(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	expected := "<p>結果: ok</p>"
	if !strings.Contains(w.Body.String(), expected) {
		t.Errorf("response.Body doesn't contain %v", expected)
	}

	tearDown()
}
