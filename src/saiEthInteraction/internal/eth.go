package internal

import (
	"context"
	"crypto/ecdsa"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iamthe1whoknocks/saiEthInteraction/models"
	"go.uber.org/zap"
)

func (s *InternalService) RawTransaction(client *ethclient.Client, value *big.Int, data []byte, contract *models.Contract) (string, error) {
	d := time.Now().Add(5000 * time.Millisecond)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	privateKey, err := crypto.HexToECDSA(contract.Private)
	if err != nil {
		s.Logger.Error("handlers - api - RawTransaction - HexToECDSA", zap.Error(err))
		return "", err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		s.Logger.Error("handlers - api - RawTransaction - cast publicKey to ecdsa", zap.Error(err))
		return "", err
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Println(err)
	}

	toAddress := common.HexToAddress(contract.Address)

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		s.Logger.Error("handlers - api - RawTransaction - get suggested gas price", zap.Error(err))
		return "", err
	}

	s.Logger.Sugar().Debugf("GAS PRICE : %v", gasPrice)

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
		s.Logger.Error("handlers - api - RawTransaction - get networkID", zap.Error(err))
		return "", err
	}

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		s.Logger.Error("handlers - api - RawTransaction - signTx", zap.Error(err))
		return "", err
	}

	err = client.SendTransaction(ctx, signedTx)
	if err != nil {
		s.Logger.Error("handlers - api - RawTransaction - sendTx", zap.Error(err))
		return "", err
	}

	return signedTx.Hash().String(), nil
}
