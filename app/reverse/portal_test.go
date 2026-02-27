package reverse_test

import (
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/app/reverse"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
)

func TestStaticPickerEmpty(t *testing.T) {
	picker, err := reverse.NewStaticMuxPicker()
	common.Must(err)
	worker, err := picker.PickAvailable()
	if err == nil {
		t.Error("expected error, but nil")
	}
	if worker != nil {
		t.Error("expected nil worker, but not nil")
	}
}
