# Core Blockchain ICAN Address Generator

## This is simple ICAN address generator for Core Blockchain network

### How to use

1. Grant permissions 
   - Via terminal: `chmod 100 address-generator`
   - Via properties of file: Right click on file -> Properties -> Permissions -> Execute
2. Run it via 2 ways :
   1. Via terminal: `./address-generator`
   2. Via GUI just double click on it

### Returned values

1. You can generate data in decrypted type and u will be returned this values:

      `Private Key: 0x8f7…949` - private key in go-core
      
      `Public Key: 0xe23…280` - public key in go-core
      
      `Address: cb723…c61` - address in go-core

2. Also you can generate encrypted json wallet file, which will be protected with password

### How to run

* Build from source
  1. `go build .`
  2. Grant permissions `chmod 100 address-generator`
  3. Generate data in one of two ways
* Use prebuilt binaries
  1. Download binary from the [Release page](https://github.com/core-coin/address-generator/releases)
  2. Grant permissions `chmod 100 address-generator`
  3. Generate data in one of two ways

### License

[CORE License](LICENSE)