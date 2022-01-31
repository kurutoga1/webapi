package selectServer_test

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"webapi/cli/selectServer"
	gw "webapi/gw/handlers"
	sh "webapi/server/handlers"
	"webapi/utils/file"
)

var (
	currentDir string
)

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c

	// 下の３つのポートはgw/config/config.jsonの設定値と同じにしなければならない。
	go func() {
		if err := http.ListenAndServe(":8081", sh.NewRouter("fileserver")); err != nil {
			panic(err.Error())
		}
	}()
	go func() {
		if err := http.ListenAndServe(":8082", sh.NewRouter("fileserver")); err != nil {
			panic(err.Error())
		}
	}()
	go func() {
		if err := http.ListenAndServe(":8083", sh.NewRouter("fileserver")); err != nil {
			panic(err.Error())
		}
	}()
	go func() {
		if err := http.ListenAndServe(":8005", gw.GetServeMux()); err != nil {
			panic(err.Error())
		}
	}()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.Remove(filepath.Join(currentDir, "log.txt"))
}

func contains(li []string, s string) bool {
	for _, a := range li {
		if strings.Contains(s, a) {
			return true
		}
	}
	return false
}

func TestSelectServer(t *testing.T) {
	selector := selectServer.NewServerSelector()

	loadBalancerURL := "http://127.0.0.1:8005"
	programName := "convertToJson"

	url := loadBalancerURL + "/SuitableServer/" + programName
	serverURL, err := selector.Select(url)
	if err != nil {
		t.Errorf("err from Select(): %v \n", err.Error())
	}

	serverPorts := []string{"8081", "8082", "8083"}
	if !contains(serverPorts, serverURL) {
		t.Errorf("%v doesn't contain of %v \n", serverURL, serverPorts)
	}

	t.Cleanup(func() {
		tearDown()
	})
}
