package config_test

import (
	"reflect"
	"testing"
	"webapi/gw/config"
)

func init() {
	// テスト用のservers.jsonを用意する
	config.SetConfPath("config.json")
}

func TestNewServerConfig(t *testing.T) {
	cfg := config.NewServerConfig()

	servers := cfg.Servers
	wantServers := []string{
		"http://127.0.0.1:8081",
		"http://127.0.0.1:8082",
		"http://127.0.0.1:8083",
	}
	if !reflect.DeepEqual(servers, wantServers) {
		t.Errorf("cfg.Servers: %v, want: %v \n", servers, wantServers)
	}

	endPoint := cfg.GetMemoryEndPoint
	wantEndPoint := "/json/health/memory"
	if endPoint != wantEndPoint {
		t.Errorf("cfg.GetMemoryEndPoint: %v, want: %v \n", endPoint, wantEndPoint)
	}

	ip := cfg.LoadBalancerServerIP
	wantIP := "127.0.0.1"
	if ip != wantIP {
		t.Errorf("cfg.LoadBalancerServerIP: %v, want: %v \n", ip, wantIP)
	}

	port := cfg.LoadBalancerServerPort
	wantPort := "8001"
	if port != wantPort {
		t.Errorf("cfg.LoadBalancerServerPort: %v, want: %v \n", port, wantPort)
	}

	logPath := cfg.LogPath
	wantLogPath := "./log.txt"
	if logPath != wantLogPath {
		t.Errorf("cfg.LogPath: %v, want: %v \n", logPath, wantLogPath)
	}

	rotateMaxKB := cfg.RotateMaxKB
	wantRotateMaxKB := 10
	if rotateMaxKB != wantRotateMaxKB {
		t.Errorf("got: %v, want: %v \n", rotateMaxKB, wantRotateMaxKB)
	}

	rotateShavingKB := cfg.RotateShavingKB
	wantRotateShavingKB := 5
	if rotateShavingKB != wantRotateShavingKB {
		t.Errorf("got: %v, want: %v \n", rotateShavingKB, wantRotateShavingKB)
	}

}
