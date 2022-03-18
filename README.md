# Core Blockchain ICAN Address Generator

## This is simple ICAN address generator for Core Blockchain network

### Command-line flags

1. `-n` - Core coin network. Values could be: `1` - Mainnet, `3` - Devin (Testnet), `5` - Private Network


### How to use

1. Grant permissions `chmod 100 address-generator`
2. Store the private key, public key and an address into wallet.txt file by running the tool with command `./address-generator -n {n} >> wallet.txt`
3. Clean history (optional, recommended) `clear && history -c`


### Returned values

After running of generator you will receive such data:

`Private Key: 0x8f7…949` - private key in go-core

`Public Key: 0xe23…280` - public key in go-core

`Address: cb723…c61` - address in go-core


### How to run

* Build from source
  1. `go build .`
  2. Grant permissions `chmod 100 address-generator`
  3. Write wallet.txt file `./address-generator -n {n} >> wallet.txt`
  4. Clean history (optional, recommended) `clear && history -c`
* Use prebuilt binaries
  1. Download binary from the [Release page](https://github.com/core-coin/address-generator/releases)
  2. Grant permissions `chmod 100 address-generator`
  3. Write wallet.txt file `./address-generator -n {n} >> wallet.txt`
  4. Clean history (optional, recommended) `clear && history -c`

### License

[CORE License](LICENSE)
