package reverse

import (
	"crypto/rand"
	"io"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/dice"
)

func (c *Control) FillInRandom() {
	randomLength := dice.Roll(64)
	randomLength++
	c.Random = make([]byte, randomLength)
	io.ReadFull(rand.Reader, c.Random)
}
