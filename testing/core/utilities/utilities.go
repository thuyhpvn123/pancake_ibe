package utilities

import (
	"fmt"

	"github.com/holiman/uint256"
	log "github.com/sirupsen/logrus"
)

func CheckFatalErr(prefix string, err error) {
	if err != nil {
		log.Fatal(fmt.Errorf("%v %v", prefix, err))
	}
}

func CheckInfoErr(prefix string, err error) bool {
	if err != nil {
		log.Infof("%v %v \n", prefix, err)
		return true
	}
	return false
}

func BytesToUint256(bytes []byte) *uint256.Int {
	return uint256.NewInt(0).SetBytes(bytes)
}
