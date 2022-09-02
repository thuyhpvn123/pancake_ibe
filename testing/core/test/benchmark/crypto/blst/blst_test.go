package blst_benchmark

import (
	"crypto/rand"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	blst "gitlab.com/meta-node/core/crypto/blst/bindings/go"
	"testing"
)

type PublicKey = blst.P1Affine
type Signature = blst.P2Affine
type AggregatePublicKey = blst.P1Aggregate
type AggregateSignature = blst.P2Aggregate

func BenchmarkBlstGenerateKeyPari(b *testing.B) {
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var ikm [32]byte
		_, _ = rand.Read(ikm[:])
		sk := blst.KeyGen(ikm[:])
		new(PublicKey).From(sk)
	}
}

func BenchmarkBlstSign(b *testing.B) {
	var ikm [32]byte
	_, _ = rand.Read(ikm[:])
	sk := blst.KeyGen(ikm[:])

	var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
	msg := []byte("hello foo")

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		new(Signature).Sign(sk, msg, dst)
	}
}

func BenchmarkBlstVerifySingleSign(b *testing.B) {
	var ikm [32]byte
	_, _ = rand.Read(ikm[:])
	sk := blst.KeyGen(ikm[:])
	pk := new(PublicKey).From(sk)

	var dst = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")
	msg := []byte("hello foo")
	sig := new(Signature).Sign(sk, msg, dst)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sig.Verify(true, pk, true, msg, dst)
	}
}
