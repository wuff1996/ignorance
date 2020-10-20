# A peek of geth

##### 		Recently, I have read over Mastering Ethereum book, and I want to see how does it work.So I create  a private blockchain at ubuntu18.04.

#### What We need:

Golang	: $ sudo snap install go

Git 	       : $ sudo install git

geth	     : $ git clone [https://github.com/ethereum/go-ethereum.git](https://github.com/ethereum/go-ethereum.git;cd)

Solc	      : $ sudo snap install solc

Now let's do it!

--------------------------------

1:

```shell
mkdir privatechain

```

I create a directory named "privatechain"  as the collection of my project.

```shell
 cd priveteChain;mkdir data code
```

I just create a data and a code directory.data is the geth datadir,which means the geth will work at the context of data,the code directory save the solidity code.

```shell
touch genesis.json
```

I create a new genesis.json file which can initialize the genesisBlock. To understand what it really does, let’s have a look of the source code of the geth init funciton:

```go
// initGenesis will initialise the given JSON format genesis file and writes it as
// the zero'd block (i.e. genesis) or will fail hard if it can't succeed.
func initGenesis(ctx *cli.Context) error {
   // Make sure we have a valid genesis JSON
   genesisPath := ctx.Args().First()
   if len(genesisPath) == 0 {
      utils.Fatalf("Must supply path to genesis JSON file")
   }
   file, err := os.Open(genesisPath)
   if err != nil {
      utils.Fatalf("Failed to read genesis file: %v", err)
   }
   defer file.Close()

   genesis := new(core.Genesis)
   if err := json.NewDecoder(file).Decode(genesis); err != nil {
      utils.Fatalf("invalid genesis file: %v", err)
   }
   // Open and initialise both full and light databases
   stack, _ := makeConfigNode(ctx)
   defer stack.Close()

   for _, name := range []string{"chaindata", "lightchaindata"} {
      chaindb, err := stack.OpenDatabase(name, 0, 0, "")
      if err != nil {
         utils.Fatalf("Failed to open database: %v", err)
      }
      _, hash, err := core.SetupGenesisBlock(chaindb, genesis)
      if err != nil {
         utils.Fatalf("Failed to write genesis block: %v", err)
      }
      chaindb.Close()
      log.Info("Successfully wrote genesis state", "database", name, "hash", hash)
   }
   return nil
}
```

So the the function above just open the file given from the path, and decode what inside the genesis.json to core.Genesis type.So let’s have a look about core.Genesis struct:

```go
// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
   Config     *params.ChainConfig `json:"config"`
   Nonce      uint64              `json:"nonce"`
   Timestamp  uint64              `json:"timestamp"`
   ExtraData  []byte              `json:"extraData"`
   GasLimit   uint64              `json:"gasLimit"   gencodec:"required"`
   Difficulty *big.Int            `json:"difficulty" gencodec:"required"`
   Mixhash    common.Hash         `json:"mixHash"`
   Coinbase   common.Address      `json:"coinbase"`
   Alloc      GenesisAlloc        `json:"alloc"      gencodec:"required"`

   // These fields are used for consensus tests. Please don't use them
   // in actual genesis blocks.
   Number     uint64      `json:"number"`
   GasUsed    uint64      `json:"gasUsed"`
   ParentHash common.Hash `json:"parentHash"`
}
```

So it’s more clear that what we will write in genesis.json is what we saw above.so the struct is just like this:

```json
{
“config”:{},
“nonce”: uint64,
“timestamp”:uint64,
“extraData”:[]byte,
“gaslimit”:uint64,
“difficulty”: int,
“mixHash”:hash,
“conbase”:address,
“alloc”:{}
}
```

And the config struct is this:

```go
// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
   ChainID *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection

   HomesteadBlock *big.Int `json:"homesteadBlock,omitempty"` // Homestead switch block (nil = no fork, 0 = already homestead)

   DAOForkBlock   *big.Int `json:"daoForkBlock,omitempty"`   // TheDAO hard-fork switch block (nil = no fork)
   DAOForkSupport bool     `json:"daoForkSupport,omitempty"` // Whether the nodes supports or opposes the DAO hard-fork

   // EIP150 implements the Gas price changes (https://github.com/ethereum/EIPs/issues/150)
   EIP150Block *big.Int    `json:"eip150Block,omitempty"` // EIP150 HF block (nil = no fork)
   EIP150Hash  common.Hash `json:"eip150Hash,omitempty"`  // EIP150 HF hash (needed for header only clients as only gas pricing changed)

   EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
   EIP158Block *big.Int `json:"eip158Block,omitempty"` // EIP158 HF block

   ByzantiumBlock      *big.Int `json:"byzantiumBlock,omitempty"`      // Byzantium switch block (nil = no fork, 0 = already on byzantium)
   ConstantinopleBlock *big.Int `json:"constantinopleBlock,omitempty"` // Constantinople switch block (nil = no fork, 0 = already activated)
   PetersburgBlock     *big.Int `json:"petersburgBlock,omitempty"`     // Petersburg switch block (nil = same as Constantinople)
   IstanbulBlock       *big.Int `json:"istanbulBlock,omitempty"`       // Istanbul switch block (nil = no fork, 0 = already on istanbul)
   MuirGlacierBlock    *big.Int `json:"muirGlacierBlock,omitempty"`    // Eip-2384 (bomb delay) switch block (nil = no fork, 0 = already activated)

   YoloV1Block *big.Int `json:"yoloV1Block,omitempty"` // YOLO v1: https://github.com/ethereum/EIPs/pull/2657 (Ephemeral testnet)
   EWASMBlock  *big.Int `json:"ewasmBlock,omitempty"`  // EWASM switch block (nil = no fork, 0 = already activated)

   // Various consensus engines
   Ethash *EthashConfig `json:"ethash,omitempty"`
   Clique *CliqueConfig `json:"clique,omitempty"`
}
```

And alloc struct is this:

```go
// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount
And genesisAccount struct is this:
// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
   Code       []byte                      `json:"code,omitempty"`
   Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
   Balance    *big.Int                    `json:"balance" gencodec:"required"`
   Nonce      uint64                      `json:"nonce,omitempty"`
   PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}
```

So I will write a template:

```json
{
"config": {
        "chainId": 15,
        "homesteadBlock": 0,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0
    },

  "alloc"      : {
  "0x0000000000000000000000000000000000000001": {"balance": "111111111"},
  "0x0000000000000000000000000000000000000002": {"balance": "222222222"}
    },

  "coinbase"   : "0x0000000000000000000000000000000000000000",
  "difficulty" : "0x00001",
  "extraData"  : "",
  "gasLimit"   : "0x2fefd8",
  "nonce"      : "0x0000000000000107",
  "mixhash"    : "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp"  : "0x00"
}
```

As we saw in the source code of the config struct,the config is the core config which determines the blockchain settings. And it stored in the database on a per block basis, which means that any network, identified by its genesis block, can have it’s own set of configuration options.

The chainId identifies the current chain, and some specific number is used :

| mainnet | ETH  | 1    | 1    | Production  |
| ------- | ---- | ---- | ---- | ----------- |
| ropsten | ETH  | 3    | 3    | Test        |
| rinkeby | ETH  | 4    | 4    | Test        |
| goerli  | ETH  | 5    | 5    | Test        |
| dev     | ETH  | 2018 | 2018 | Development |
| classic | ETC  | 61   | 1    | Production  |
| mordor  | ETC  | 63   | 7    | Test        |
| kotti   | ETC  | 6    | 6    | Test        |

So we have to specific a different chainId so that the geth will not misunderstand that I want to connect to the mainnet.And it is used for relay protection at the transaction signature as “v”.

homesteadBlock switch block (nil = no fork, 0 = already homestead)

EIP150 HF block (nil = no fork):EIP150 change the IO-heavy to mitigate transaction spam attacks

EIP155Block:Simple replay attack protection,add “v”.

EIP158Block:State clearing.

And alloc specifies the initial state that is part of the genesis block.It’s just a key-value map, the key is the 20bytes account address,and the value is an genesisAccount which the code,storage,balance,nonce,privateKey.

Coinbase set the receiver who will receive the miner-reward.

Difficulty specific the difficulty of the total network to mine this block,note that the subBlock’s difficulty is relevant to parent block, so I set it to 0x0001 for more easier to mine the block out.

Gaslimit limit the total gas used in this block , it set the computation of this blockchain.

Nonce and mixhash is to mine out the correct nonce number to make sure that the hash of this block to POW.

-------------------------------------

2:

After that,we just save the file,and $ geth --datadir init genesis.json

The geth will decode the genesis,json and create a genesis block in the database.

Next.

```
geth --datadir data --networkid 666 console
```

To run the node.

Now we can interactive with node with console,or we can open another tty.And put

```
geth attach data/geth.ipc
```

 to interactive with it.

-----------------------------

3:

Now that we are interactive with the node,let's do something within current blockchain:

```js
eth.accounts
```

![img](C:\Users\Administrator\Desktop\love\ignorance\Blockchain\wps1.jpg)

As I have created some accounts so It looks like this,Now let’s create a new account with man-readable password.

```js
personal.newAccount("wuff")
```

And you will get your address.And the key configuration of you accounts is stored at data/keystore/UTC--xxx

So What happend when I input the command?

When I input the personal.newAccount(“wuff”),it mostly just invokes the PrivateAccountAPI method:

```go
// NewAccount will create a new account and returns the address for the new account.
func (s *PrivateAccountAPI) NewAccount(password string) (common.Address, error) {
   ks, err := fetchKeystore(s.am)
   if err != nil {
      return common.Address{}, err
   }
   acc, err := ks.NewAccount(password)
   if err == nil {
      log.Info("Your new key was generated", "address", acc.Address)
      log.Warn("Please backup your key file!", "path", acc.URL.Path)
      log.Warn("Please remember your password!")
      return acc.Address, nil
   }
   return common.Address{}, err
}
```

So, the program will fetch the keystore which is used to store the key configuration on the disk.

And ks.NewAccount(“wuff”),to generate a new key and store it into the key directory, encrypting it with the passphrase.Let’s have a deeper sight in it:

```go
// NewAccount generates a new key and stores it into the key directory,
// encrypting it with the passphrase.
func (ks *KeyStore) NewAccount(passphrase string) (accounts.Account, error) {
   _, account, err := storeNewKey(ks.storage, crand.Reader, passphrase)
   if err != nil {
      return accounts.Account{}, err
   }
   // Add the account to the cache immediately rather
   // than waiting for file system notifications to pick it up.
   ks.cache.add(account)
   ks.refreshWallets()
   return account, nil
}
```

So NewAccount function just using passphrase what we passed to it, a rand number, and the keystore file.Then add the account to the cache,then refreshWallets retrieves the current account list and based on that does any necessary wallet refreshes.

But how stroreNewKey function works?

```go
func storeNewKey(ks keyStore, rand io.Reader, auth string) (*Key, accounts.Account, error) {
   key, err := newKey(rand)
   if err != nil {
      return nil, accounts.Account{}, err
   }
   a := accounts.Account{
      Address: key.Address,
      URL:     accounts.URL{Scheme: KeyStoreScheme, Path: ks.JoinPath(keyFileName(key.Address))},
   }
   if err := ks.StoreKey(a.URL.Path, key, auth); err != nil {
      zeroKey(key.PrivateKey)
      return nil, a, err
   }
   return key, a, err
}
```

It’s more clear now,newKey function receive the rand number and generate a key who warps the uuid ,the address,and the private key ,and then create An account object at the set of the key address & keystore-inside file with name: UTC--<created_at UTC ISO8601>-<address hex>.After that, ks writes the key and auth->”wuff ”what we write before to the account file and then clean the private key in memory.

At this layer, It’s obvious that what we input is just an auth in the keystore.

So, how the key object created by newKey function?

Let's check:

```go
func newKey(rand io.Reader) (*Key, error) {
   privateKeyECDSA, err := ecdsa.GenerateKey(crypto.S256(), rand)
   if err != nil {
      return nil, err
   }
   return newKeyFromECDSA(privateKeyECDSA), nil
}
```

The crypto.S256() returns an instance of secp2561k1 curve.As NIST defines:
$$
Y^2 mod p=(x^3+7)modp
$$
It creates a curve and p is the prime number which limit the curve can’t escape from it.
$$
P = 2^ (256) - 2^(32)- 2^9 -2^8 -2^7 - 2^6 - 2^4 -1
$$
So  the ecdsa.GenerateKey method recevice the rand number and use the secp2561k1 to create numbe:

```go
// randFieldElement returns a random element of the field underlying the given
// curve using the procedure given in [NSA] A.2.1.
func randFieldElement(c elliptic.Curve, rand io.Reader) (k *big.Int, err error) {
   params := c.Params()
   b := make([]byte, params.BitSize/8+8)
   _, err = io.ReadFull(rand, b)
   if err != nil {
      return
   }

   k = new(big.Int).SetBytes(b)
   n := new(big.Int).Sub(params.N, one)
   k.Mod(k, n)
   k.Add(k, one)
   return
}
```

The one = 1,so this is sure that the private key is bigger than the order of the basic point,so how big the params.N(order point) is ?
$$
N = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141
$$
Now that we get the correct privateKey, use the algorithm:

K = k * G to get the publicKey point (x,y) at the secp256k1 curve. The k is the time G + G,G is the unchanged point in the secp265k1 curve called creating point.

The algorithm G + G is defined: the intersection that Tangent of point G with secp256k1 curve.

##### So now we clear out what happened when we input the personal.newAccount(“wuff”):

##### 1:retrieves the encrypted keystore from the account manager.

##### 2: get a safe rand number and create a private key,using private key to generate the public key(point),and use the public key to generate address(last 20bytes of keccak256(pk)).

##### 3:create an account with key.Address and create a file whose name platform is : 

```
UTC--<created_at UTC ISO8601>-<address hex>
```

#####   inside the keystore directory.And write the key and what we write(here is “wuff”) in it,then clean the private key in memory.

##### 4:add the account created before to the cache and refresh the wallet. Then returns the account’s address and it’s path to us in console.

![img](C:\Users\Administrator\Desktop\love\ignorance\Blockchain\wps2.jpg)

So let’s put

```js
 eth.accounts 
```

And console shows that:

![img](C:\Users\Administrator\Desktop\love\ignorance\Blockchain\wps3.jpg)

Great! Let’s take a look of what happened this time:

```go
// listAccounts will return a list of addresses for accounts this node manages.
func (s *PrivateAccountAPI) ListAccounts() []common.Address {
   return s.am.Accounts()
}
```

The web3 just invokes this function to return a list of accounts’ addresses in this node.

:), fine , let’s see the balance of the account

Let’s see the accounts list:

```go

 /*Accounts retrieves the list of signing accounts the wallet is currently aware of. For hierarchical deterministic wallets, the list will not be exhaustive, rather only contain the accounts explicitly pinned during account derivation.*/
```

So let’s see the balance of an account: 

```js
web3.fromWei(eth.getBalance(eth.accounts[0]),"ether")
```

<img src="C:\Users\Administrator\Desktop\love\ignorance\Blockchain\wps4.jpg" alt="img"  />

Because I have mined before , so it shows that first account of this node have 766 ether.

Let’s figure out what works this out:

```go
// GetBalance returns the amount of wei for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *PublicBlockChainAPI) GetBalance(ctx context.Context, address common.Address, blockNrOrHash rpc.BlockNumberOrHash) (*hexutil.Big, error) {
   state, _, err := s.b.StateAndHeaderByNumberOrHash(ctx, blockNrOrHash)
   if state == nil || err != nil {
      return nil, err
   }
   return (*hexutil.Big)(state.GetBalance(address)), state.Error()
}
```

So it will retrieve the the stateDB given blockNumber or blockHash ... too late :)

End-->