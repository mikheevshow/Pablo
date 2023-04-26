package pablo

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mikheevshow/Pablo/cex"
	"golang.org/x/crypto/sha3"
)

type Address string
type PrivateKey string

func (a Address) ToString() string {
	return string(a)
}

func (pk PrivateKey) ToString() string {
	return string(pk)
}

func (pk PrivateKey) GetAddress() Address {
	addr, err := getAddressFrom(pk)
	if err != nil {
		panic(err)
	}
	return addr
}

type SwapConfig struct {
	AbiFile string
}

var swapMap = map[string]*SwapConfig{
	"sushi":        &SwapConfig{""},
	"1inch":        &SwapConfig{""},
	"woo":          &SwapConfig{""},
	"stargate":     &SwapConfig{""},
	"traderjoexyz": &SwapConfig{""},
}

type BridgeConfig struct {
	AbiFile string
}

var bridgeMap = map[string]*BridgeConfig{
	"bitcoinbridge": &BridgeConfig{""},
}

type Pablo struct {
	blockchainClientProvider BlockchainClientProvider
	dexClientProvider        DexClientProvider
	tokenContractService     TokenContractService
}

func CreatePablo() *Pablo {
	blockchainClientProvider := NewBlockchainClientProvider()
	tokenContractService := NewTokenContractService()
	return &Pablo{
		blockchainClientProvider: blockchainClientProvider,
		tokenContractService: tokenContractService,
	}
}

func (p *Pablo) TransferFromCex(name string, credentials cex.Creds, symbol string, to Address, blockchain string, amount string) *Pablo {
	return p
}

func (p *Pablo) Transfer(from PrivateKey, to Address, amount string, symbol string, blockchain string) *Pablo {
	fromAddress := from.GetAddress().ToString()
	log.Println("Start transfering from address " + fromAddress + " to address " + to.ToString() + " amount " + amount + " symbol " + symbol + " on blockchain " + blockchain)

	supported := p.tokenContractService.IsSymbolSupportedByBlockchain(blockchain, symbol)
	if !supported {
		log.Fatal("Symbol " + symbol + " not supported")
	}

	client := p.blockchainClientProvider.GetClient(blockchain)
	isNative := p.tokenContractService.IsNative(blockchain, symbol)
	var err error
	if isNative {
		err = transferNative(client, from, to, amount)
	} else {
		tokenAddress := p.tokenContractService.GetContractAddress(blockchain, symbol)
		err = transferErc20(client, from, to, amount, 10, tokenAddress)
	}

	if err != nil {
		panic(err)
	}
	log.Println("Finished successfully transfer from address " + fromAddress + " to address " + to.ToString() + " amount " + amount + " symbol " + symbol + " on blockchain " + blockchain)
	return p
}

func (p *Pablo) SwapDex(name string, amount string, fromSymbol string, toSymbol string, wallet PrivateKey, blockchain string) *Pablo {
	log.Println("Startign swap ")
	return p
}

func (p *Pablo) Bridge(name string, privateKey PrivateKey, symbol string, amount string, fromBlockchain string, toBlockchain string) *Pablo {
	return p
}

func (p *Pablo) Wait(duration time.Duration) *Pablo {
	log.Panicln("Start waiting " + duration.String())
	time.Sleep(duration)
	log.Panicln("Continue pipeline execution")
	return p
}

func bridge(bridgeClient *ethclient.Client, privateKey string, symbol string, amount string, fromBlockchain string, toBlockchain string) error {
	return nil
}

func transferNative(blockchainClient *ethclient.Client, from PrivateKey, to Address, amount string) error {

	privateKey, err := crypto.HexToECDSA(from.ToString())
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := blockchainClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	value := new(big.Int)
	value.SetString(amount, 0)

	gasLimit := uint64(21000)
	gasPrice, err := blockchainClient.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	toAddress := common.HexToAddress(to.ToString())

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, nil)

	chainID, err := blockchainClient.NetworkID(context.Background())
	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}
	 
	err = blockchainClient.SendTransaction(context.Background(), signedTx)

	if err != nil {
		return err
	}

	return nil
}

func transferErc20(blockchainClient *ethclient.Client, from PrivateKey, to Address, amount string, denomindation int, symbolAddress string) error {

	value := big.NewInt(0)
	toAddress := common.HexToAddress(to.ToString())
	tokenAddress := common.HexToAddress(symbolAddress)

	transferSignature := []byte("transfer(address,uint256)")

	hash := sha3.NewLegacyKeccak256()
	hash.Write(transferSignature)

	methodID := hash.Sum(nil)[:4]

	paddedAddress := common.LeftPadBytes(toAddress.Bytes(), 32)

	amnt := new(big.Int)
	amnt.SetString(amount, denomindation)
	paddedAmnt := common.LeftPadBytes(amnt.Bytes(), 32)

	var data []byte
	data = append(data, methodID...)
	data = append(data, paddedAddress...)
	data = append(data, paddedAmnt...)

	privateKey, err := crypto.HexToECDSA(from.ToString())
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Public key from private key is not ok")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := blockchainClient.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return err
	}

	gasLimit, err := blockchainClient.EstimateGas(context.Background(), ethereum.CallMsg{
		To:   &tokenAddress,
		Data: data,
	})

	if err != nil {
		return err
	}

	gasPrice, err := blockchainClient.SuggestGasPrice(context.Background())

	if err != nil {
		return err
	}

	tx := types.NewTransaction(nonce, tokenAddress, value, gasLimit, gasPrice, data)

	chainID, err := blockchainClient.NetworkID(context.Background())

	if err != nil {
		return err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return err
	}

	err = blockchainClient.SendTransaction(context.Background(), signedTx)

	if err != nil {
		return err
	}

	return nil
}

func swapDex(dexClient *ethclient.Client, amount string, fromAddress string, toAddress string, wallet PrivateKey) error {

	// privateKey, err := crypto.HexToECDSA(wallet.ToString())
	// if err != nil {
	// 	return err
	// }

	// publicKey := privateKey.Public()
	// publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	// }

	// fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)





	return nil
}

func getAddressFrom(privateKey PrivateKey) (Address, error) {
	privKey, err := crypto.HexToECDSA(privateKey.ToString())
	if err != nil {
		return "", err
	}
	address := crypto.PubkeyToAddress(privKey.PublicKey).Hex()
	return Address(address), nil
}

// Blockchain client provider

type EthClientProvider interface {
	GetClient(name string) *ethclient.Client
}

type BlockchainClientProvider EthClientProvider

type BlockchainClientProviderImpl struct {
	clientMap map[string]*ethclient.Client
}

func NewBlockchainClientProvider() BlockchainClientProvider {
	return &BlockchainClientProviderImpl{}
}

func (p *BlockchainClientProviderImpl) GetClient(name string) *ethclient.Client {
	return p.clientMap[name]
}

// Dex client provider

type DexClientProvider EthClientProvider

type DexClientProviderImpl struct {
	clientMap map[string]*ethclient.Client
}

func NewDexClientProvider() DexClientProvider {
	return &DexClientProviderImpl{}
}

func (p *DexClientProviderImpl) GetClient(name string) *ethclient.Client {
	return p.clientMap[name]
}

// Cex withdrawal provider

type cexWithdrawalProvider interface {
}
