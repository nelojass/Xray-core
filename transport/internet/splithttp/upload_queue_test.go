package splithttp_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/splithttp"
)

func Test_regression_readzero(t *testing.T) {
	q := NewUploadQueue(10)
	q.Push(Packet{
		Payload: []byte("x"),
		Seq:     0,
	})
	buf := make([]byte, 20)
	n, err := q.Read(buf)
	common.Must(err)
	if n != 1 {
		t.Error("n=", n)
	}
}
