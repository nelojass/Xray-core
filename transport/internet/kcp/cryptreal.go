package kcp

import (
	"crypto/cipher"
	"crypto/sha256"

	"v12w.x34y.com/flyfishLib/forkHub/xtls/xray-core/common/crypto"
)

func NewAEADAESGCMBasedOnSeed(seed string) cipher.AEAD {
	hashedSeed := sha256.Sum256([]byte(seed))
	return crypto.NewAesGcm(hashedSeed[:])
}
