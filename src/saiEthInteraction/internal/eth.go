package internal

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iamthe1whoknocks/saiEthInteraction/models"
	"go.uber.org/zap"
)

var mux sync.Mutex
var nonceList = map[string]uint64{}

func (is *InternalService) getNonce(client *ethclient.Client, contract *models.Contract, fromAddress common.Address) (uint64, error) {
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return 0, err
	}

	mux.Lock()
	if prevNonce, ok := nonceList[contract.Address]; ok {
		if prevNonce == nonce {
			nonce += 10
		}
	}

	nonceList[contract.Address] = nonce
	mux.Unlock()

	return nonce, nil
}

func (is *InternalService) RawTransaction(client *ethclient.Client, value *big.Int, data []byte, contract *models.Contract) (string, error) {
	d := time.Now().Add(5000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	privateKey, err := crypto.HexToECDSA(contract.Private)
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - HexToECDSA", zap.Error(err))
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		is.Logger.Error("handlers - api - RawTransaction - cast publicKey to ecdsa", zap.Error(err))
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := is.getNonce(client, contract, fromAddress)
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - get nonce", zap.Error(err))
		return "", err
	}

	toAddress := common.HexToAddress(contract.Address)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - get suggested gas price", zap.Error(err))
		return "", err
	}

	is.Logger.Sugar().Debugf("GAS PRICE : %v", gasPrice)

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      contract.GasLimit,
		GasPrice: gasPrice,
		Data:     data,
	})

	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - get networkID", zap.Error(err))
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - signTx", zap.Error(err))
		return "", err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		is.Logger.Error("handlers - api - RawTransaction - sendTx", zap.Error(err))
		return "", err
	}

	return signedTx.Hash().String(), nil
}
