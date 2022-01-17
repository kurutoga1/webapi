package config_test

import (
	"reflect"
	"testing"
	"webapi/cli/config"
)

func init() {
	// テスト用のservers.jsonを用意する
	config.SetConfPath("config.json")
}

func TestNewServerConfig(t *testing.T) {
	cfg := config.NewServerConfig()

	servers := cfg.APIGateWayServers
	wantServers := []string{
		"http://127.0.0.1:8001",
		"http://127.0.0.1:8002",
		"http://127.0.0.1:8003",
	}
	if !reflect.DeepEqual(servers, wantServers) {
		t.Errorf("cfg.Servers: %v, want: %v \n", servers, wantServers)
	}
}
