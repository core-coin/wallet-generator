package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	rand2 "math/rand"
	"net/http"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strconv"
	"time"

	"github.com/core-coin/go-core/v2/accounts/keystore"
	"github.com/core-coin/go-core/v2/common"
	"github.com/core-coin/go-core/v2/crypto"
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
	_, _, _, err := getFormValues(r) // set network ID
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

	returnPageToClient(
		w,
		DecryptedWalletData{
			PrivateKey: common.Bytes2Hex(PrivateKey.PrivateKey()),
			PublicKey:  common.Bytes2Hex(PrivateKey.PublicKey()[:]),
			Address:    PrivateKey.Address().Hex(),
		},
		baseHTML,
		decryptedHTML,
		cssHTML,
	)
}

func exitHandler(w http.ResponseWriter, req *http.Request) {
	os.Exit(1)
}

func encryptedDataHandler(w http.ResponseWriter, r *http.Request) {
	keyPath, password, _, err := getFormValues(r)
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
			DecryptedWalletData{IsError: true, ErrorMsg: "Cannot create a wallet file: " + err.Error()},
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

func getFormValues(r *http.Request) (keyPath, password, passwordRepeat string, err error) {
	if err := r.ParseForm(); err != nil {
		return "", "", "", errors.New(fmt.Sprintf("Cannot parse data from form: %v", err))
	}

	// Setup NetworkID
	network := r.FormValue("network_id")
	if network == "" {
		return "", "", "", errors.New("You need to define network ID: 1=Mainnet, 3=Devin, 4+=Enterprise")
	}
	networkId, err := strconv.Atoi(network)
	if err != nil {
		return "", "", "", errors.New(fmt.Sprintf("Wrong network id: %v ", err))
	}

	if networkId == 2 {
		return "", "", "", errors.New("Network ID 2 is not existent!")
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
		return
	}
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if !path.IsAbs(keyPath) {
			return "", "", "", errors.New("Зath for keyfile is not absolute")
		}
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
