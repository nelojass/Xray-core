package wechat_test

import (
	"context"
	"testing"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/buf"
	. "v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/transport/internet/headers/wechat"
)

func TestUTPWrite(t *testing.T) {
	videoRaw, err := NewVideoChat(context.Background(), &VideoConfig{})
	common.Must(err)

	video := videoRaw.(*VideoChat)

	payload := buf.New()
	video.Serialize(payload.Extend(video.Size()))

	if payload.Len() != video.Size() {
		t.Error("expected payload size ", video.Size(), " but got ", payload.Len())
	}
}
