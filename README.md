# Core coin Address Generator

This is simple address generator for core coin networks


Command-line flags:
1. `network` - core coin network. Values could be: `1` - Mainnet, `3` - Devin (Testnet), `5` - Private Network


How to use:
1. Generate the private key, public key and an address by running the tool with command `address-generator -network {n}`
2. Clear screen with command `clear`
3. Clear terminal history for current session with command `history -c`


Returned values:
After running of generator you will receive such data:

`Private Key: 0x8f7...949` - private key in go-core

`Public Key: 0xe23...280` - public key in go-core

`Address: ab723...c61` - address in go-core

