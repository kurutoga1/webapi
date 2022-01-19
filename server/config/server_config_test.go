package config_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"
	"webapi/server/config"
	"webapi/utils/file"
)

var currentDir string

func init() {
	c, err := file.GetCurrentDir()
	if err != nil {
		log.Fatalln(err.Error())
	}
	currentDir = c
}

func TestSetConfPath(t *testing.T) {
	config.SetConfPath("config.json")
	t.Cleanup(func() {
		tearDown()
	})
}

func TestLoad(t *testing.T) {
	// config.jsonを追加するごとに追加する
	cfg := config.Load()

	wantIP := "127.0.0.1"
	if cfg.ServerIP != wantIP {
		t.Errorf("ServerIP: %v, want: %v \n", cfg.ServerIP, wantIP)
	}

	wantPort := "8081"
	if cfg.ServerPort != wantPort {
		t.Errorf("ServerPort: %v, want: %v \n", cfg.ServerPort, wantPort)
	}

	wantRotateMaxKB := 6
	if cfg.Log.RotateMaxKB != wantRotateMaxKB {
		t.Errorf("got: %v, want: %v \n", cfg.Log.RotateMaxKB, wantRotateMaxKB)
	}

	wantRotateShavingSize := 3
	if cfg.Log.RotateShavingKB != wantRotateShavingSize {
		t.Errorf("got: %v, want: %v \n", cfg.Log.RotateShavingKB, wantRotateShavingSize)
	}

	t.Cleanup(func() {
		tearDown()
	})
}

func tearDown() {
	os.RemoveAll(filepath.Join(currentDir, "fileserver"))
}
