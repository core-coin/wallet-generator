package main

import (
	"crypto/rand"
	"fmt"
	rand2 "math/rand"
	"time"

	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/common/hexutil"
	"github.com/core-coin/go-core/crypto"
	"github.com/core-coin/go-goldilocks"
)

func main() {
	rand2.Seed(time.Now().Unix())
	common.DefaultNetworkID = common.Mainnet

	PrivateKey, err := crypto.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	PublicKey := goldilocks.Ed448DerivePublicKey(*PrivateKey)

	Address := crypto.PubkeyToAddress(PublicKey)

	fmt.Println("Private Key:", hexutil.Encode(PrivateKey[:]))
	fmt.Println("Public Key:", hexutil.Encode(PublicKey[:]))
	fmt.Println("Address:", Address.Hex())
}
