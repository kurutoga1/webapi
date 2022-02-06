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
	"time"
	"webapi/gw/config"
	gh "webapi/gw/handlers"
	proRouter "webapi/pro/router"
	"webapi/utils/file"
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
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	//gh.Logger.SetOutput(os.Stdout)
	serverSet()
}

func serverSet() {
	// servers.jsonのExpectedAliveServersを見ながらサーバを立てる。
	ports := []string{"8081", "8082", "8083"}
	for _, p := range ports {
		p := p
		go func() {
			if err := http.ListenAndServe(":"+p, proRouter.New().New("fileserver"+p)); err != nil {
				panic(err.Error())
			}
			time.Sleep(1 * time.Second)
		}()
	}
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func TestGetMinimumMemoryServerHandler(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "/MinimumMemoryServer", nil)
	w := httptest.NewRecorder()

	handler := gh.GetMinimumMemoryServerHandler(config.NewServerConfig())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v, body: %v \n", w.Code, http.StatusOK, w.Body.String())
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(w.Body.Bytes(), &d)
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
	r, _ := http.NewRequest(http.MethodGet, "/SuitableServer/convertToJson", nil)
	w := httptest.NewRecorder()

	handler := gh.GetSuitableServerHandler(log.New(os.Stdout, "", log.LstdFlags), config.NewServerConfig())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	type j struct {
		Url string `json:"url"`
	}

	var d j
	err = json.Unmarshal(w.Body.Bytes(), &d)
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
	r, _ := http.NewRequest(http.MethodGet, "/userTop", nil)
	w := httptest.NewRecorder()

	handler := gh.UserTopHandler(log.New(os.Stdout, "", log.LstdFlags), config.NewServerConfig())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func TestGetAliveServersHandler(t *testing.T) {

	r, _ := http.NewRequest(http.MethodGet, "/AliveServers", nil)
	w := httptest.NewRecorder()

	handler := gh.GetAliveServersHandler(log.New(os.Stdout, "", log.LstdFlags), config.NewServerConfig())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	type data struct {
		AliveServers []string `json:"AliveServers"`
	}
	var d data
	err = json.Unmarshal(w.Body.Bytes(), &d)
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
	r, _ := http.NewRequest(http.MethodGet, "/AllServerPrograms", nil)
	w := httptest.NewRecorder()

	handler := gh.GetAllProgramsHandler(log.New(os.Stdout, "", log.LstdFlags), config.NewServerConfig())
	handler.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("got %v, want %v", w.Code, http.StatusOK)
	}

	b := w.Body.String()
	if !strings.Contains(b, "convertToJson") {
		t.Errorf("%v doesn't contain convertToJson", b)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
