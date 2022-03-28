package main

import (
	"bytes"
	"crypto/rand"
	"flag"
	"fmt"
	"github.com/core-coin/go-core/accounts/keystore"
	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/common/hexutil"
	"github.com/core-coin/go-core/crypto"
	"github.com/core-coin/go-goldilocks"
	"github.com/pkg/errors"
	"golang.org/x/term"
	rand2 "math/rand"
	"os"
	"path"
	"strings"
	"syscall"
	"time"
)

var (
	network           int64
	displayInTerminal bool
	keydir            string
)

func main() {
	rand2.Seed(time.Now().Unix())

	if err := parseFlags(); err != nil {
		fmt.Println(err)
		return
	}

	if displayInTerminal {
		displayDataInTerminal()
		return
	}

	if err := storeDataInJSON(); err != nil {
		fmt.Println(err)
	}
}

func parseFlags() error {
	wd, err := os.Getwd()
	if err != nil {
		return errors.New("Cannot get current directory: " + err.Error())
	}

	flag.Int64Var(&network, "n", 1, "Core Coin network ID (1-Mainnet, 3-Devin(Testnet), 5-Private Network) (shorthand)")
	flag.Int64Var(&network, "network", 1, "Core Coin network ID (1-Mainnet, 3-Devin(Testnet), 5-Private Network)")

	flag.BoolVar(&displayInTerminal, "t", false, "Display decrypted data without password in terminal (optional) (shorthand)")
	flag.BoolVar(&displayInTerminal, "terminal", false, "Display decrypted data without password in terminal (optional)")

	flag.StringVar(&keydir, "k", wd,
		`A rooted (absolute) path where to store encrypted json keyfile. Example : "/root/core-coin". Default value is current directory`)
	flag.StringVar(&keydir, "keydir", wd,
		`A rooted (absolute) path where to store encrypted json keyfile. Example : "/root/core-coin". Default value is current directory (shorthand)`)

	flag.Parse()

	if !path.IsAbs(keydir) {
		return errors.New("Key directory is not rooted (absolute)")
	}

	if network == 2 {
		return errors.New("There is not network with id = 2")
	}
	common.DefaultNetworkID = common.NetworkID(network)

	return nil
}

func displayDataInTerminal() {
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

func readPasswordFromStdin() (string, error) {
	fmt.Println("Enter password for wallet: ")
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	fmt.Println("Repeat password for wallet: ")
	repeatedBytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}

	if bytes.Compare(bytePassword, repeatedBytePassword) != 0 {
		return "", errors.New("Passwords does not match!")
	}
	password := string(bytePassword)
	return strings.TrimSpace(password), nil
}

func storeDataInJSON() error {
	password, err := readPasswordFromStdin()
	if err != nil {
		return err
	}
	account, err := keystore.StoreKey(keydir, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		return err
	}
	fmt.Printf("\nYour new key was generated\n\n")
	fmt.Printf("Public address of the key:   %s\n", account.Address.Hex())
	fmt.Printf("Path of the secret key file: %s\n\n", account.URL.Path)
	fmt.Printf("- You can share your public address with anyone. Others need it to interact with you.\n")
	fmt.Printf("- You must NEVER share the secret key with anyone! The key controls access to your funds!\n")
	fmt.Printf("- You must BACKUP your key file! Without the key, it's impossible to access account funds!\n")
	fmt.Printf("- You must REMEMBER your password! Without the password, it's impossible to decrypt the key!\n\n")
	return nil
}
