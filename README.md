# Core Blockchain ICAN Address Generator

## This is simple ICAN address generator for Core Blockchain network

### Command-line flags

1. `-n` - Core coin network. Values could be: `1` - Mainnet, `3` - Devin (Testnet), `5` - Private Network
2. `-k` - Keystore path where to store wallet file. Must be rooted (absolute) path like `/home/keystore`
3. `-t` - Boolean which enables printing values in terminal. With this flag json keyfile will not be generated, instead of this all values will be printed in terminal`

### How to use

1. Grant permissions `chmod 100 address-generator`
2. Generate data via 2 ways:
   1. Decrypted - store the private key, public key and an address into wallet.txt file by running the tool with command `./address-generator -n {n} -t >> wallet.txt`
   2. Encrypted - generate json wallet file `/address-generator -n {n} -k /home/keystore`
3. Clean history (optional, recommended) `clear && history -c`


### Returned values

1. After running of generator in 1 way you will receive such data:

      `Private Key: 0x8f7…949` - private key in go-core
      
      `Public Key: 0xe23…280` - public key in go-core
      
      `Address: cb723…c61` - address in go-core

2. After running of generator in 1 way you will receive json file which can be imported in go-core


### How to run

* Build from source
  1. `go build .`
  2. Grant permissions `chmod 100 address-generator`
  3. Generate data in one of two ways
     - Write wallet.txt file `./address-generator -n {n} -t >> wallet.txt`
     - Generate json wallet file `/address-generator -n {n} -k /home/keystore`
  4. Clean history (optional, recommended) `clear && history -c`
* Use prebuilt binaries
  1. Download binary from the [Release page](https://github.com/core-coin/address-generator/releases)
  2. Grant permissions `chmod 100 address-generator`
  3. Generate data in one of two ways
     - Write wallet.txt file `./address-generator -n {n} -t >> wallet.txt`
     - Generate json wallet file `/address-generator -n {n} -k /home/keystore`
  4. Clean history (optional, recommended) `clear && history -c`

### License

[CORE License](LICENSE)