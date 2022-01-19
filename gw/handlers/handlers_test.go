package handlers_test

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	gh "webapi/gw/handlers"
	sh "webapi/server/handlers"
	"webapi/utils"
)

/*
ゲートウェイのハンドラーのテストはまずサーバを立てて、
そこから初めてハンドラーのテストをしなければならない。
またテストサーバのIPとconfのservers.jsonのサーバIPを同じにしなければならない。
*/

var (
	currentDir    string
	deletes       []string
	ts1, ts2, ts3 *httptest.Server
	err           error
)

func init() {
	// サーバを立てるとカレントディレクトリにfileserverディレクトリとlog.txt
	// ができるのでそれを削除する。
	c, err := utils.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	//gh.Logger.SetOutput(os.Stdout)
	serverSet()
}

func serverSet() {
	// servers.jsonのExpectedAliveServersを見ながらサーバを立てる。
	go func() {
		if err := http.ListenAndServe(":8081", sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()

	go func() {
		if err := http.ListenAndServe(":8082", sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()

	go func() {
		if err := http.ListenAndServe(":8083", sh.GetServeMux("fileserver")); err != nil {
			panic(err.Error())
		}
	}()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGetMinimumMemoryServerHandler(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/MinimumMemoryServer", nil)
	response := httptest.NewRecorder()

	gh.GetMinimumMemoryServerHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v, body: %v \n", response.Code, http.StatusOK, response.Body.String())
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(response.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error())
	}

	if d.Url == "" {
		t.Errorf("d.Url is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetSuitableServerHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/SuitableServer/convertToJson", nil)
	response := httptest.NewRecorder()

	gh.GetSuitableServerHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(response.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error())
	}

	if d.Url == "" {
		t.Errorf("d.Url is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestUserTopHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/userTop2", nil)
	response := httptest.NewRecorder()

	gh.UserTopHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetAliveServersHandler(t *testing.T) {

	request, _ := http.NewRequest(http.MethodGet, "/AliveServers", nil)
	response := httptest.NewRecorder()

	gh.GetAliveServersHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	type data struct {
		AliveServers []string `json:"AliveServers"`
	}
	var d data
	err = json.Unmarshal(response.Body.Bytes(), &d)
	if err != nil {
		t.Errorf(err.Error())
	}

	if len(d.AliveServers) == 1 {
		t.Errorf("AliveServers is empty.")
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetAllProgramsHandler(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/AllServerPrograms", nil)
	response := httptest.NewRecorder()

	gh.GetAllProgramsHandler(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("got %v, want %v", response.Code, http.StatusOK)
	}

	b := response.Body.String()
	if !strings.Contains(b, "convertToJson") {
		t.Errorf("%v doesn't contain convertToJson", b)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
