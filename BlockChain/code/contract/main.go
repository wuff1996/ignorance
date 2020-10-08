package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"log"
	"math/big"
)

func main() {
	//client ,err := ethclient.Dial("http://localhost:30303")
	//if err != nil {
	//	log.Fatal(errors.Wrap(err,"main"))
	//}
	privateKey,err:=crypto.GenerateKey()
	if err!= nil {
		log.Fatal(err)
	}
	alloc:=make(core.GenesisAlloc)
	//publicKey:=privateKey.Public()
	//publicKeyECDSA,ok:=publicKey.(*ecdsa.PublicKey)
	//if !ok {
	//	log.Fatal(errors.Wrap(err,"main"))
	//}
	//fromAddress:=crypto.PubkeyToAddress(*publicKeyECDSA)
	//nonce,err:=client.PendingNonceAt(context.Background(),fromAddress)
	//if err!= nil {
	//	log.Fatal(errors.Wrap(err,"main"))
	//}
	//gasPrice,err:=client.SuggestGasPrice(context.Background())
	//if err!= nil {
	//	log.Fatal(errors.Wrap(err, "main"))
	//}
	auth :=bind.NewKeyedTransactor(privateKey)
	//auth.Nonce=big.NewInt(int64((nonce)))
	auth.Value=big.NewInt(0)
	auth.GasLimit=uint64(300000)
	//auth.GasPrice=gasPrice
	alloc[auth.From]=core.GenesisAccount{Balance: big.NewInt(133700000)}
	sim:=backends.NewSimulatedBackend(alloc,133700000)
	address,tx,instance,err:=DeployTest(auth,sim)
	if err != nil {
		log.Fatal(errors.Wrap(err,"main"))
	}
	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())
	fmt.Println("instance: ",instance)
	bodyHash:=sim.Blockchain().CurrentHeader().Hash()
	fmt.Println("blockChain: ",sim.Blockchain().CurrentHeader())
	sim.Commit()
	fmt.Println("blockChain: ",sim.Blockchain().GetBody(bodyHash))

}
