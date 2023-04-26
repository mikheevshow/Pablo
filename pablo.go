package pablo

import (
	"log"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mikheevshow/Pablo/cex"
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
}

func CreatePablo() *Pablo {
	blockchainClientProvider := NewBlockchainClientProvider()
	return &Pablo{
		blockchainClientProvider: blockchainClientProvider,
	}
}

func (p *Pablo) TransferFromCex(name string, credentials cex.Creds, symbol string, to Address, blockchain string, amount string) *Pablo {
	return p
}

func (p *Pablo) Transfer(from PrivateKey, to Address, amount string, symbol string, blockchain string) *Pablo {
	fromAddress := from.GetAddress().ToString()
	log.Println("Start transfering from address " + fromAddress + " to address " + to.ToString() + " amount " + amount + " symbol " + symbol + " on blockchain " + blockchain)
	client := p.blockchainClientProvider.GetClient(blockchain)
	err := transfer(client, from, to, amount, symbol)
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

func transfer(blockchainClient *ethclient.Client, from PrivateKey, to Address, amount string, symbol string) error {
	return nil
}

func swapDex(dexClient *ethclient.Client, amount string, fromSymbol string, toSymbol string, wallet PrivateKey, blockchain string) error {

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
