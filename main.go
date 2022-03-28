package main

import (
	"crypto/rand"
	"fmt"
	"github.com/core-coin/go-core/accounts/keystore"
	"github.com/core-coin/go-core/common"
	"github.com/core-coin/go-core/crypto"
	"github.com/core-coin/go-goldilocks"
	"html/template"
	rand2 "math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"
)

type DecryptedWalletData struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

type EncryptedWalletData struct {
	FilePath string
	Address  string
}

func main() {
	rand2.Seed(time.Now().Unix())

	http.HandleFunc("/index", index)
	http.HandleFunc("/generate_raw", generateRaw)
	http.HandleFunc("/generate_encrypted", generateEncrypted)
	http.HandleFunc("/exit", exit)

	go open("http://localhost:8080/index")

	panic(http.ListenAndServe(":8080", nil))

}

func index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("").ParseFiles(
		path.Join("templates", "base.html"),
		path.Join("templates", "index.html"))
	if err != nil {
		fmt.Fprintf(w, "Cannot parse template files: %v ", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func generateRaw(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	network := r.FormValue("network_id")
	if network == "" {
		fmt.Fprint(w, "You need to define network ID, 1=Mainnet, 3=Devin, everything bigger then 4=private networks")
		return
	}
	networkId, err := strconv.Atoi(network)
	if err != nil {
		fmt.Fprintf(w, "Wrong network id: %v ", err)
		return
	}
	if networkId == 2 {
		fmt.Fprintln(w, "There is not network with id = 2")
		return
	}
	common.DefaultNetworkID = common.NetworkID(networkId)

	PrivateKey, err := crypto.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(w, "Cannot generate private key: %v ", err)
		return
	}
	PublicKey := goldilocks.Ed448DerivePublicKey(*PrivateKey)

	Address := crypto.PubkeyToAddress(PublicKey)

	tmpl, err := template.New("").ParseFiles(
		path.Join("templates", "base.html"),
		path.Join("templates", "decrypted.html"))
	if err != nil {
		fmt.Fprintf(w, "Cannot parse template files: %v ", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", DecryptedWalletData{
		PrivateKey: common.Bytes2Hex(PrivateKey[:]),
		PublicKey:  common.Bytes2Hex(PublicKey[:]),
		Address:    Address.Hex(),
	})
	if err != nil {
		fmt.Println(err)
	}
}

func exit(w http.ResponseWriter, req *http.Request) {
	os.Exit(1)
}

func generateEncrypted(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	password := r.FormValue("pass")
	passwordRepeat := r.FormValue("pass_repeat")
	if password != passwordRepeat {
		fmt.Fprintln(w, "Passwords does not match")
		return
	}

	network := r.FormValue("network_id")
	if network == "" {
		fmt.Fprint(w, "You need to define network ID, 1=Mainnet, 3=Devin, everything bigger then 4=private networks")
		return
	}
	networkId, err := strconv.Atoi(network)
	if err != nil {
		fmt.Fprintf(w, "Wrong network id: %v ", err)
		return
	}

	if networkId == 2 {
		fmt.Fprintln(w, "There is not network with id = 2")
		return
	}
	common.DefaultNetworkID = common.NetworkID(networkId)

	keyPath := r.FormValue("path")
	if keyPath == "" {
		keyPath, err = os.Getwd()
		if err != nil {
			fmt.Fprintf(w, "Cannot get working directory: %v ", err)
			return
		}
	}
	if !path.IsAbs(keyPath) {
		fmt.Fprintln(w, "Path for keyfile is not absolute")
		return
	}
	account, err := keystore.StoreKey(keyPath, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		fmt.Fprintf(w, "Cannot create a keyfile: %v", err)
		return
	}

	tmpl, err := template.New("").ParseFiles(
		path.Join("templates", "base.html"),
		path.Join("templates", "encrypted.html"))
	if err != nil {
		fmt.Fprintf(w, "Cannot parse template files: %v ", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", EncryptedWalletData{
		FilePath: account.URL.Path,
		Address:  account.Address.Hex(),
	})
	if err != nil {
		fmt.Println(err)
	}
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
