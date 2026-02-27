package signal_test

import (
	"testing"

	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/signal"
)

func TestNotifierSignal(t *testing.T) {
	n := NewNotifier()

	w := n.Wait()
	n.Signal()

	select {
	case <-w:
	default:
		t.Fail()
	}
}
