package pablo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
)

type TokenContractService interface {
	IsSymbolSupportedByBlockchain(blockchain string, symbol string) bool

	IsNative(blockchain string, symbol string) bool

	GetContractAddress(blockchain string, symbol string) string
}

type addresses struct {
	Addresses []blockchainTokens `json:"addresses"`
}

type blockchainTokens struct {
	Blockchain string      `json:"blockchain"`
	Tokens     []TokenInfo `json:"tokens"`
}

type TokenInfo struct {
	Symbol  string `json:"symbol"`
	Native  bool   `json:"native"`
	Address string `json:"address"`
}

type TokenContractServiceImpl struct {
	tokenInfoMap map[string]*map[string]*TokenInfo
}

func NewTokenContractService() TokenContractService {
	file, err := ioutil.ReadFile("../smart-contract-addresses.json")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(file)

	data := addresses{}

	err = json.Unmarshal([]byte(file), &data)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(data)

	result := make(map[string]*map[string]*TokenInfo)

	for _, bcTokens := range data.Addresses {
		innerMap := make(map[string]*TokenInfo)

		for _, token := range bcTokens.Tokens {
			tokenCopy := token
			innerMap[token.Symbol] = &tokenCopy
		}

		innerMapCopy := innerMap
		result[bcTokens.Blockchain] = &innerMapCopy
	}

	return &TokenContractServiceImpl{
		tokenInfoMap: result,
	}
}

func (s *TokenContractServiceImpl) IsSymbolSupportedByBlockchain(blockchain string, symbol string) bool {
	innerMapPtr, ok := s.tokenInfoMap[blockchain]
	if !ok {
		return false
	}
	innerMap := *innerMapPtr
	token, ok := innerMap[symbol]
	if !ok {
		return false
	}

	fmt.Println(token)
	return true
}

func (s *TokenContractServiceImpl) IsNative(blockchain string, symbol string) bool {
	innerMapPtr, ok := s.tokenInfoMap[blockchain]
	if !ok {
		log.Fatal("")
	}
	innerMap := *innerMapPtr
	token, ok := innerMap[symbol]
	if !ok {
		log.Fatal("")
	}

	return token.Native
}

func (s *TokenContractServiceImpl) GetContractAddress(blockchain string, symbol string) string {
	innerMapPtr := s.tokenInfoMap[blockchain]
	if innerMapPtr == nil {
		log.Fatal("sadasdasd")
	}
	innerMap := *innerMapPtr
	token := innerMap[symbol]
	if token == nil {
		log.Fatal("gggggg")
	}

	return token.Address
}
