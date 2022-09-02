package crypto_benchmark

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/herumi/bls-eth-go-binary/bls"
	cc "gitlab.com/meta-node/core/controllers"
	ccrypto "gitlab.com/meta-node/core/crypto"
	"google.golang.org/protobuf/proto"
)

func BenchmarkVeryLongDataHash(b *testing.B) {
	rb := make([]byte, 256)
	rand.Read(rb)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crypto.Keccak256(rb)
	}
}

func BenchmarkHerumiGetTransactionHash(b *testing.B) {
	transaction := cc.GetEmptyTransaction()
	transaction.PreviousTransaction = cc.GetEmptyTransaction()
	by, _ := proto.Marshal(transaction)
	fmt.Printf("transaction length %v", len(by))
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		cc.GetTransactionHash(transaction)
	}
}

func BenchmarkHerumiSign(b *testing.B) {
	ccrypto.Init()
	sec, _, _ := ccrypto.GenerateKeyPair()
	b.ResetTimer()
	testMessage := crypto.Keccak256([]byte("test"))
	for n := 0; n < b.N; n++ {
		ccrypto.Sign(sec, testMessage)
	}
}

func BenchmarkHerumiVerifySign(b *testing.B) {
	ccrypto.Init()
	sec, pub, _ := ccrypto.GenerateKeyPair()
	testMessage := crypto.Keccak256([]byte("test"))
	sign := ccrypto.Sign(sec, testMessage)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		ccrypto.VerifySign(pub, sign, testMessage)
	}
}

func BenchmarkHerumiAggregate(b *testing.B) {
	ccrypto.Init()
	total := 10000
	signs := [][]byte{}
	for n := 0; n < total; n++ {
		sec, _, _ := ccrypto.GenerateKeyPair()
		testMessage := crypto.Keccak256([]byte(fmt.Sprintf("%d", n)))
		sign := ccrypto.Sign(sec, testMessage)
		signs = append(signs, sign)
	}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var aggSign bls.Sign
		signArr := make([]bls.Sign, len(signs))
		var sign bls.Sign
		for i, v := range signs {
			if sign.Deserialize(v) == nil {
				signArr[i] = sign
			}
		}
		aggSign.Aggregate(signArr)
	}

}

func BenchmarkHerumiAggregateVerify(b *testing.B) {
	ccrypto.Init()

	signs := [][]byte{}
	pubs := [][]byte{}
	hashes := [][]byte{}
	total := 10000
	for n := 0; n < total; n++ {
		sec, pub, _ := ccrypto.GenerateKeyPair()
		pubs = append(pubs, pub)
		testMessage := crypto.Keccak256([]byte(fmt.Sprintf("%d", n)))
		sign := ccrypto.Sign(sec, testMessage)
		signs = append(signs, sign)
		hashes = append(hashes, testMessage)
	}
	var aggSign bls.Sign
	signArr := make([]bls.Sign, len(signs))
	var sign bls.Sign
	for i, v := range signs {
		if sign.Deserialize(v) == nil {
			signArr[i] = sign
		}
	}
	aggSign.Aggregate(signArr)

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		pubArr := make([]bls.PublicKey, total)
		hashArr := make([][]byte, total)

		for index := 0; index < total; index++ {
			var pub bls.PublicKey
			hashArr[index] = hashes[index]
			if pub.Deserialize(pubs[index]) == nil {
				pubArr[index] = pub
			}
		}

		aggSign.VerifyAggregateHashes(pubArr, hashArr)
	}

}

var (
	testmsg     = hexutil.MustDecode("0xce0677bb30baa8cf067c88db9811f4333d131bf8bcf12fe7065d211dce971009")
	testsig     = hexutil.MustDecode("0x90f27b8b488db00b00606796d2987f6a5f59ae62ea05effe84fef5b8b0e549984a691139ad57a3f0b906637673aa2f63d1f55cb1a69199d4009eea23ceaddc9301")
	testpubkey  = hexutil.MustDecode("0x04e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a0a2b2667f7e725ceea70c673093bf67663e0312623c8e091b13cf2c0f11ef652")
	testpubkeyc = hexutil.MustDecode("0x02e32df42865e97135acfb65f3bae71bdc86f4d49150ad6a440b6f15878109880a")
)

func BenchmarkVerifySignatureEthereum(b *testing.B) {
	sig := testsig[:len(testsig)-1]
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		crypto.VerifySignature(testpubkey, testmsg, sig)
	}

}
