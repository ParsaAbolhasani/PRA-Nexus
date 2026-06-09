package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BalanceReader : خواننده موجودی واقعی از بلاکچین
type BalanceReader struct {
	Client     *ethclient.Client
	ContractABI abi.ABI
	ContractAddr common.Address
}

// NewBalanceReader : ایجاد نمونه جدید
func NewBalanceReader(rpcURL, contractAddress, abiPath string) (*BalanceReader, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, err
	}

	// خواندن فایل ABI
	abiFile, err := os.ReadFile(abiPath)
	if err != nil {
		return nil, err
	}

	var contractABI abi.ABI
	err = json.Unmarshal(abiFile, &contractABI)
	if err != nil {
		return nil, err
	}

	return &BalanceReader{
		Client:       client,
		ContractABI:  contractABI,
		ContractAddr: common.HexToAddress(contractAddress),
	}, nil
}

// GetBalance : دریافت موجودی واقعی یک آدرس
func (br *BalanceReader) GetBalance(address common.Address) (*big.Float, error) {
	// بسته‌بندی تابع balanceOf
	data, err := br.ContractABI.Pack("balanceOf", address)
	if err != nil {
		return nil, err
	}

	// فراخوانی قرارداد (خواندن، بدون گس)
	msg := ethereum.CallMsg{
		To:   &br.ContractAddr,
		Data: data,
	}

	result, err := br.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	// باز کردن خروجی
	var balance *big.Int
	err = br.ContractABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return nil, err
	}

	// تبدیل به عدد اعشاری با احتساب 18 رقم اعشار
	balanceFloat := new(big.Float).SetInt(balance)
	divisor := new(big.Float).SetFloat64(1e18)
	balanceFloat.Quo(balanceFloat, divisor)

	return balanceFloat, nil
}

// GetAllowance : دریافت مقدار مجاز برای خرج کردن توسط spender
func (br *BalanceReader) GetAllowance(owner, spender common.Address) (*big.Float, error) {
	data, err := br.ContractABI.Pack("allowance", owner, spender)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &br.ContractAddr,
		Data: data,
	}

	result, err := br.Client.CallContract(context.Background(), msg, nil)
	if err != nil {
		return nil, err
	}

	var allowance *big.Int
	err = br.ContractABI.UnpackIntoInterface(&allowance, "allowance", result)
	if err != nil {
		return nil, err
	}

	allowanceFloat := new(big.Float).SetInt(allowance)
	divisor := new(big.Float).SetFloat64(1e18)
	allowanceFloat.Quo(allowanceFloat, divisor)

	return allowanceFloat, nil
}

// تابع کمکی برای تبدیل float64 به big.Int با 18 رقم اعشار
func FloatToBigInt(amount float64) *big.Int {
	bigFloat := new(big.Float).SetFloat64(amount)
	multiplier := new(big.Float).SetFloat64(1e18)
	bigFloat.Mul(bigFloat, multiplier)
	result, _ := bigFloat.Int(nil)
	return result
}

// استفاده در WalletManager قبلی (به‌روزرسانی تابع GetPRABalance)
func (wm *WalletManager) GetPRABalance(address common.Address) (*big.Float, error) {
	// اگر BalanceReader را به WalletManager اضافه کرده باشی
	if wm.BalanceReader == nil {
		return nil, fmt.Errorf("balance reader not initialized")
	}
	return wm.BalanceReader.GetBalance(address)
}