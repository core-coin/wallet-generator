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
}

type EncryptedWalletData struct {
	FilePath string
	Address  string
}

func main() {
	rand2.Seed(time.Now().Unix())

	http.HandleFunc("/index", indexHandler)
	http.HandleFunc("/generate_raw", rawDataHandler)
	http.HandleFunc("/generate_encrypted", encryptedDataHandler)
	http.HandleFunc("/exit", exitHandler)

	go open("http://localhost:8080/index")

	panic(http.ListenAndServe(":8080", nil))

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	returnPageToClient(w, nil, baseHTML, indexHTML, cssHTML)
}

func rawDataHandler(w http.ResponseWriter, r *http.Request) {
	getFormValues(w, r) // set network ID

	PrivateKey, err := crypto.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(w, "Cannot generate private key: %v ", err)
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
	keyPath, password, _ := getFormValues(w, r)

	account, err := keystore.StoreKey(keyPath, password, keystore.StandardScryptN, keystore.StandardScryptP)
	if err != nil {
		fmt.Fprintf(w, "Cannot create a keyfile: %v", err)
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
		fmt.Fprintf(w, "Cannot create a keyfile: %v", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		fmt.Println(err)
	}
}

func getFormValues(w http.ResponseWriter, r *http.Request) (keyPath, password, passwordRepeat string) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	// Setup NetworkID
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

	password = r.FormValue("pass")
	passwordRepeat = r.FormValue("pass_repeat")
	if password != passwordRepeat {
		fmt.Fprintln(w, "Passwords does not match")
		return
	}

	keyPath = r.FormValue("path")
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

const (
	baseHTML = `{{define "base"}}
<!DOCTYPE html>
<html lang="en">
<head>
{{template "css" .}}
</head>
<body>

{{template "content" .}}

</body>
</html>
{{end}}`

	encryptedHTML = `<!-- encrypted.html -->
{{define "content"}}
<div class="encrypted">
    <pre>Your new key was generated
    Public address of the key: {{.Address}}
    Path of the secret key file: {{.FilePath}}
    You can share your public address with anyone. Others need it to interact with you.
    You must NEVER share the secret key with anyone! The key controls access to your funds!
    You must BACKUP your key file! Without the key, it's impossible to access account funds!
    You must REMEMBER your password! Without the password, it's impossible to decrypt the key!
    </pre>
</div>
{{end}}`

	decryptedHTML = `<!-- decrypted.html -->
{{define "content"}}

<div class="decrypted">
    <p>Private Key: {{.PrivateKey}}</p>
    <p>Public Key: {{.PublicKey}}</p>
    <p>Address: {{.Address}}</p>
</div>

{{end}}`

	indexHTML = `<!-- index.html -->
{{define "content"}}

<form action="/generate_raw" method="post">
    <label for="network_id">Network ID:</label>
    <input type="text" id="network_id" name="network_id"><br><br>
    <input type="submit" value="Generate raw wallet values">
</form>
<p>&nbsp;</p>
OR
<p>&nbsp;</p>
<form action="/generate_encrypted" method="post">
    <label for="network_id">Network ID:</label>
    <input type="text" id="network_id_encrypted" name="network_id"><br><br>
    <label for="pass">Password:</label>
    <input type="password" id="pass" name="pass"><br><br>
    <label for="pass_repeat">Repeat password:</label>
    <input type="password" id="pass_repeat" name="pass_repeat"><br><br>
    <label for="path">Path to store keyfile (Leave empty to save in the same directory where is the program):</label>
    <input type="string" id="path" name="path"><br><br>
    <input type="submit" value="Generate json wallet file">
</form>
<p>&nbsp;</p>
OR
<p>&nbsp;</p>
<form action="/exit" method="post">
    <input type="submit" value="Exit">
</form>

{{end}}`

	cssHTML = `<!-- css.html -->
{{define "css"}}
<style>
.encrypted {
    color: yellow;
}
.decrypted {
    color: blue;
}
</style>
{{end}}`
)
