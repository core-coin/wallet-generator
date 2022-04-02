package main

import (
	"crypto/rand"
	"errors"
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

	IsError  bool
	ErrorMsg string
}

type EncryptedWalletData struct {
	FilePath string
	Address  string

	IsError  bool
	ErrorMsg string
}

func main() {
	rand2.Seed(time.Now().Unix())

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/generate_raw", rawDataHandler)
	http.HandleFunc("/generate_encrypted", encryptedDataHandler)
	http.HandleFunc("/exit", exitHandler)

	go open("http://localhost:8080/")

	panic(http.ListenAndServe(":8080", nil))

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	returnPageToClient(w, nil, baseHTML, indexHTML, cssHTML)
}

func rawDataHandler(w http.ResponseWriter, r *http.Request) {
	_, _, _, err := getFormValues(w, r) // set network ID
	if err != nil {
		returnPageToClient(
			w,
			DecryptedWalletData{IsError: true, ErrorMsg: err.Error()},
			baseHTML, indexHTML, cssHTML,
		)
		return
	}
	PrivateKey, err := crypto.GenerateKey(rand.Reader)
	if err != nil {
		returnPageToClient(
			w,
			DecryptedWalletData{IsError: true, ErrorMsg: fmt.Sprintf("Cannot generate private key: %v ", err)},
			baseHTML, indexHTML, cssHTML,
		)
		return
	}
	PublicKey := goldilocks.Ed448DerivePublicKey(*PrivateKey)

	Address := crypto.PubkeyToAddress(PublicKey)

	returnPageToClient(
		w,
		DecryptedWalletData{
			PrivateKey: common.Bytes2Hex(PrivateKey[:]),
			PublicKey:  common.Bytes2Hex(PublicKey[:]),
			Address:    Address.Hex()},
		baseHTML,
		decryptedHTML,
		cssHTML,
	)
}

func exitHandler(w http.ResponseWriter, req *http.Request) {
	os.Exit(1)
}

func encryptedDataHandler(w http.ResponseWriter, r *http.Request) {
	keyPath, password, _, err := getFormValues(w, r)
	if err != nil {
		returnPageToClient(
			w,
			DecryptedWalletData{IsError: true, ErrorMsg: err.Error()},
			baseHTML, indexHTML, cssHTML,
		)
		return
	}

	account, err := keystore.StoreKey(keyPath, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		returnPageToClient(
			w,
			DecryptedWalletData{IsError: true, ErrorMsg: "cannot create a wallet file: " + err.Error()},
			baseHTML, indexHTML, cssHTML,
		)
		return
	}

	returnPageToClient(
		w,
		EncryptedWalletData{
			FilePath: account.URL.Path,
			Address:  account.Address.Hex(),
		},
		baseHTML,
		encryptedHTML,
		cssHTML,
	)
}

func returnPageToClient(w http.ResponseWriter, data interface{}, templates ...string) {
	tmpl, err := renderTemplates(templates...)
	if err != nil {
		fmt.Fprintf(w, "Cannot render templates: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Fprintf(w, "Cannot execute template: %v", err)
	}
}

func getFormValues(w http.ResponseWriter, r *http.Request) (keyPath, password, passwordRepeat string, err error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", errors.New(fmt.Sprintf("cannot parse data from form: %v", err))
	}

	// Setup NetworkID
	network := r.FormValue("network_id")
	if network == "" {
		return "", "", "", errors.New("you need to define network ID, 1=Mainnet, 3=Devin, everything bigger then 4=private networks")
	}
	networkId, err := strconv.Atoi(network)
	if err != nil {
		return "", "", "", errors.New(fmt.Sprintf("Wrong network id: %v ", err))
	}

	if networkId == 2 {
		return "", "", "", errors.New("there is not network with id = 2")
	}
	common.DefaultNetworkID = common.NetworkID(networkId)

	password = r.FormValue("pass")
	passwordRepeat = r.FormValue("pass_repeat")
	if password != passwordRepeat {
		return "", "", "", errors.New("Passwords does not match")
	}

	keyPath = r.FormValue("path")
	if keyPath == "" {
		keyPath, err = os.Getwd()
		if err != nil {
			return "", "", "", errors.New(fmt.Sprintf("Cannot get current directory: %v ", err))
		}
	}
	if !path.IsAbs(keyPath) {
		return "", "", "", errors.New("path for keyfile is not absolute")
	}

	return
}

func renderTemplates(templates ...string) (*template.Template, error) {
	var err error

	tmpl := template.New("")

	for _, oneTemplate := range templates {
		tmpl, err = tmpl.Parse(oneTemplate)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Cannot parse template: %v ", err))
		}
	}

	return tmpl, nil
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
