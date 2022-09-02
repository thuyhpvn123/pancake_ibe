package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"runtime"

	"github.com/ethereum/go-ethereum/crypto"
	log "github.com/sirupsen/logrus"
	blst "gitlab.com/meta-node/core/crypto/blst/bindings/go"
	pb "gitlab.com/meta-node/core/proto"
)

type blstPublicKey = blst.P1Affine
type blstSignature = blst.P2Affine
type blstAggregateSignature = blst.P2Aggregate
type blstAggregatePublicKey = blst.P1Aggregate
type blstSecretKey = blst.SecretKey

var dstMinPk = []byte("BLS_SIG_BLS12381G2_XMD:SHA-256_SSWU_RO_POP_")

func Init() {
	blst.SetMaxProcs(runtime.GOMAXPROCS(0))
}

func Sign(bPri []byte, bMessage []byte) []byte {
	sk := new(blstSecretKey).Deserialize(bPri)
	sign := new(blstSignature).Sign(sk, bMessage, dstMinPk)
	return sign.Compress()
}

func GetByteAddress(pubkey []byte) []byte {
	hash := crypto.Keccak256(pubkey)
	address := hash[12:]
	return address
}

func VerifySign(bPub []byte, bSig []byte, bMsg []byte) bool {
	return new(blstSignature).VerifyCompressed(bSig, true, bPub, false, bMsg, dstMinPk)
}

func VerifyAggregateSign(bPubs [][]byte, bSig []byte, bMsgs [][]byte) bool {
	return new(blstSignature).AggregateVerifyCompressed(bSig, true, bPubs, false, bMsgs, dstMinPk)
}

func GenerateKeyPairFromSecretKey(hexSecretKey string) ([]byte, []byte, []byte) {
	secByte, _ := hex.DecodeString(hexSecretKey)
	sec := new(blstSecretKey).Deserialize(secByte)
	pub := new(blstPublicKey).From(sec).Compress()
	hash := crypto.Keccak256([]byte(pub))
	return sec.Serialize(), pub, hash[12:]
}

func randBLSTSecretKey() *blstSecretKey {
	var t [32]byte
	_, _ = rand.Read(t[:])
	secretKey := blst.KeyGen(t[:])
	return secretKey
}

func GenerateKeyPair() ([]byte, []byte, []byte) {
	sec := randBLSTSecretKey()
	pub := new(blstPublicKey).From(sec).Compress()
	hash := crypto.Keccak256([]byte(pub))
	return sec.Serialize(), pub, hash[12:]
}

func CreateAggregateSignFromTransactions(transactions []*pb.Transaction) []byte {
	log.Debugf("CreateAggregateSignFromTransactions total transaction %v", len(transactions))
	aggregatedSignature := new(blst.P2Aggregate)
	signatures := make([][]byte, len(transactions))
	for i, v := range transactions {
		signatures[i] = v.Sign
	}
	aggregatedSignature.AggregateCompressed(signatures, false)
	log.Debugf("aggreagtesign %v", aggregatedSignature.ToAffine().Compress())

	return aggregatedSignature.ToAffine().Compress()
}
