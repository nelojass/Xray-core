package protocol_test

import (
	"testing"

	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/protocol"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/uuid"
)

func TestIdEquals(t *testing.T) {
	id1 := NewID(uuid.New())
	id2 := NewID(id1.UUID())

	if !id1.Equals(id2) {
		t.Error("expected id1 to equal id2, but actually not")
	}

	if id1.String() != id2.String() {
		t.Error(id1.String(), " != ", id2.String())
	}
}
