package dispatcher

import (
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/buf"
	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/features/stats"
)

type SizeStatWriter struct {
	Counter stats.Counter
	Writer  buf.Writer
}

func (w *SizeStatWriter) WriteMultiBuffer(mb buf.MultiBuffer) error {
	w.Counter.Add(int64(mb.Len()))
	return w.Writer.WriteMultiBuffer(mb)
}

func (w *SizeStatWriter) Close() error {
	return common.Close(w.Writer)
}

func (w *SizeStatWriter) Interrupt() {
	common.Interrupt(w.Writer)
}
