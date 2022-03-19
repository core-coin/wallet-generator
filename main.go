package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	rand2 "math/rand"
	"time"

	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/common/hexutil"
	"github.com/core-coin/go-core/crypto"
	"github.com/core-coin/go-goldilocks"
)

func main() {
	var network int64
	rand2.Seed(time.Now().Unix())

	flag.Int64Var(&network, "network", 1, "Core Coin network ID (1-Mainnet, 3-Devin(Testnet), 5-Private Network)")
	flag.Int64Var(&network, "n", 1, "Core Coin network ID (1-Mainnet, 3-Devin(Testnet), 5-Private Network) (shorthand)")

	flag.Parse()
	if network != 1 && network != 3 && network != 5 {
		panic("Wrong network ID")
	}
	common.DefaultNetworkID = common.NetworkID(network)

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
