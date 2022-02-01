package serverAliveConfirmer_test

import (
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"webapi/gw/serverAliveConfirmer"
	"webapi/pro/router"
	"webapi/test"
	"webapi/utils/file"
)

var currentDir string
var confirmer serverAliveConfirmer.ServerAliveConfirmer

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
	confirmer = serverAliveConfirmer.NewServerAliveConfirmer()
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
	os.RemoveAll(filepath.Join(currentDir, "log.txt"))
}

func TestIsAlive(t *testing.T) {
	addrs, _, err := test.GetStartedServers(router.New(), 1)
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("pro is alive.", func(t *testing.T) {
		testIsAlive(t, addrs[0], "/user/top", true)
	})

	t.Run("pro is not alive.", func(t *testing.T) {
		testIsAlive(t, "http://127.0.0.1:8052", "/user/top", false)
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testIsAlive(t *testing.T, addr, endPoint string, expect bool) {
	t.Helper()
	alive, err := confirmer.IsAlive(addr, endPoint)
	if err != nil {
		t.Errorf("err occur: %v \n", err.Error())
	}
	if alive != expect {
		t.Errorf("got: %v, want false.", expect)
	}
}

func TestGetAliveServers(t *testing.T) {
	addrs, _, err := test.GetStartedServers(router.New(), 3)
	if err != nil {
		log.Fatalln(err)
	}

	t.Run("get alive servers 1", func(t *testing.T) {
		testGetAliveServers(t, addrs, "/health", addrs)
	})

	t.Cleanup(func() {
		tearDown()
	})
}

func testGetAliveServers(t *testing.T, servers []string, endPoint string, expectServers []string) {
	t.Helper()
	aliveServers, err := serverAliveConfirmer.GetAliveServers(servers, endPoint, confirmer)
	if err != nil {
		t.Errorf(err.Error())
	}

	if !reflect.DeepEqual(aliveServers, expectServers) {
		t.Errorf("got: %v, want: %v \n", aliveServers, expectServers)
	}
}
