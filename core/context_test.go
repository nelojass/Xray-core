package core_test

import (
	"context"
	"testing"
	_ "unsafe"

	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/core"
)

func TestFromContextPanic(t *testing.T) {
	defer func() {
		r := recover()
		if r == nil {
			t.Error("expect panic, but nil")
		}
	}()

	MustFromContext(context.Background())
}
