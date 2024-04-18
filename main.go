package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	erc20 "github.com/marco75116/my-go-project/gen"
	quoter "github.com/marco75116/my-go-project/genTwo"
)

const contractABI = `[{"inputs":[{"internalType":"address","name":"_factory","type":"address"},{"internalType":"address","name":"_WETH9","type":"address"}],"stateMutability":"nonpayable","type":"constructor"},{"inputs":[],"name":"WETH9","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"factory","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"bytes","name":"path","type":"bytes"},{"internalType":"uint256","name":"amountIn","type":"uint256"}],"name":"quoteExactInput","outputs":[{"internalType":"uint256","name":"amountOut","type":"uint256"},{"internalType":"uint160[]","name":"sqrtPriceX96AfterList","type":"uint160[]"},{"internalType":"uint32[]","name":"initializedTicksCrossedList","type":"uint32[]"},{"internalType":"uint256","name":"gasEstimate","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"tokenIn","type":"address"},{"internalType":"address","name":"tokenOut","type":"address"},{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint24","name":"fee","type":"uint24"},{"internalType":"uint160","name":"sqrtPriceLimitX96","type":"uint160"}],"internalType":"struct IQuoterV2.QuoteExactInputSingleParams","name":"params","type":"tuple"}],"name":"quoteExactInputSingle","outputs":[{"internalType":"uint256","name":"amountOut","type":"uint256"},{"internalType":"uint160","name":"sqrtPriceX96After","type":"uint160"},{"internalType":"uint32","name":"initializedTicksCrossed","type":"uint32"},{"internalType":"uint256","name":"gasEstimate","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"path","type":"bytes"},{"internalType":"uint256","name":"amountOut","type":"uint256"}],"name":"quoteExactOutput","outputs":[{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint160[]","name":"sqrtPriceX96AfterList","type":"uint160[]"},{"internalType":"uint32[]","name":"initializedTicksCrossedList","type":"uint32[]"},{"internalType":"uint256","name":"gasEstimate","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"components":[{"internalType":"address","name":"tokenIn","type":"address"},{"internalType":"address","name":"tokenOut","type":"address"},{"internalType":"uint256","name":"amount","type":"uint256"},{"internalType":"uint24","name":"fee","type":"uint24"},{"internalType":"uint160","name":"sqrtPriceLimitX96","type":"uint160"}],"internalType":"struct IQuoterV2.QuoteExactOutputSingleParams","name":"params","type":"tuple"}],"name":"quoteExactOutputSingle","outputs":[{"internalType":"uint256","name":"amountIn","type":"uint256"},{"internalType":"uint160","name":"sqrtPriceX96After","type":"uint160"},{"internalType":"uint32","name":"initializedTicksCrossed","type":"uint32"},{"internalType":"uint256","name":"gasEstimate","type":"uint256"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"int256","name":"amount0Delta","type":"int256"},{"internalType":"int256","name":"amount1Delta","type":"int256"},{"internalType":"bytes","name":"path","type":"bytes"}],"name":"uniswapV3SwapCallback","outputs":[],"stateMutability":"view","type":"function"}]`

func main() {
	clientArb, err := ethclient.Dial("https://arbitrum.llamarpc.com")

	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	clientMainnet, err := ethclient.Dial("https://eth.llamarpc.com")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	header, err := clientArb.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to get header: %v", err)
	}

	fmt.Println("Latest Ethereum block number:", header.Number.String())

	account := common.HexToAddress("0x320e4fBBdDEc43989d69d2bE4952e162d99C007d")
	balance, err := clientArb.BalanceAt(context.Background(), account, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	fbalance := new(big.Float).SetInt(balance)
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))

	fmt.Println(ethValue)

	addressERC20 := common.HexToAddress("0xaf88d065e77c8cC2239327C5EDb3A432268e5831")

	instance, err := erc20.NewErc20(addressERC20, clientArb)

	if err != nil {
		log.Fatal(err)
	}

	wallet := common.HexToAddress("0x320e4fBBdDEc43989d69d2bE4952e162d99C007d")
	bal, err := instance.BalanceOf(&bind.CallOpts{}, wallet)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("wei: %s\n", bal)

	name, err := instance.Name(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("name: %s\n", name)
	fmt.Printf("symbol: %s\n", symbol)
	fmt.Printf("decimals: %v\n", decimals)

	addressQuoter := common.HexToAddress("0x61fFE014bA17989E743c5F6cB21bF9697530B21e")
	// instanceQuoter, err := quoter.NewQuoter(addressQuoter, clientMainnet)

	amountIn, ok := new(big.Int).SetString("400436951100", 10)

	if !ok {
		log.Fatalf("Failed to parse big integer")
	}

	params := quoter.IQuoterV2QuoteExactInputSingleParams{
		TokenIn:           common.HexToAddress("0x8b42dd505b2b4c0b21f75eb284c7295bc71e580b"),
		TokenOut:          common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
		AmountIn:          amountIn,
		Fee:               new(big.Int).SetInt64(10000),
		SqrtPriceLimitX96: new(big.Int).SetInt64(0),
	}

	parsedABI, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	data, err := parsedABI.Pack("quoteExactInputSingle", params)

	if err != nil {
		log.Fatalf("Failed to pack data for quoteExactInputSingle: %v", err)
	}

	callMsg := ethereum.CallMsg{
		To:   &addressQuoter,
		Data: data,
	}

	output, err := clientMainnet.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Error calling contract: %v", err)
	}

	var result struct {
		AmountOut               *big.Int
		SqrtPriceX96After       *big.Int
		InitializedTicksCrossed uint32
		GasEstimate             *big.Int
	}
	err = parsedABI.UnpackIntoInterface(&result, "quoteExactInputSingle", output)
	if err != nil {
		log.Fatalf("Failed to unpack output: %v", err)
	}

	fmt.Printf("Amount Out: %s\n", result.AmountOut.String())
}
