/*
ランタイム(メモリ状況)を取得するパッケージ
*/

package runtime

import (
	"encoding/json"
	"runtime"
)

// RuntimeGetter はランタムをゲットするためのインターフェース
type RuntimeGetter interface {
	GetRuntime() runtime.MemStats
	GetRuntimeAsJSON() (string, error)
}

// NewRuntimeGetter はmyRuntimeストラクトを取得する。
// myRuntimeストラクトはRuntimeGetter を実装している。
func NewRuntimeGetter() RuntimeGetter {
	return &myRuntime{}
}

type myRuntime struct {
	memStatus runtime.MemStats
}

//GetRuntime はruntime.MemStatsを返す
func (m myRuntime) GetRuntime() runtime.MemStats {
	runtime.ReadMemStats(&m.memStatus)
	return m.memStatus
}

// GetRuntimeAsJSON はメモリの状態をJSONにして返す。
func (m myRuntime) GetRuntimeAsJSON() (string, error) {
	r := m.GetRuntime()
	bytes, err := json.MarshalIndent(&r, "", "    ")
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
