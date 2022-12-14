package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/iamthe1whoknocks/saiEthInteraction/models"
	"github.com/saiset-co/saiService"
	"go.uber.org/zap"
)

const (
	contractsPath = "contracts.json"
)

func (is *InternalService) NewHandler() saiService.Handler {
	return saiService.Handler{
		"api": saiService.HandlerElement{
			Name:        "api",
			Description: "transact encoded transaction to contract by ABI",
			Function: func(data interface{}) (interface{}, error) {
				contractData, ok := data.(map[string]interface{})
				if !ok {
					Service.Logger.Sugar().Debugf("handling connect method, wrong type, current type : %+v", reflect.TypeOf(data))
					return nil, errors.New("wrong type of incoming data")
				}

				b, err := json.Marshal(contractData)
				if err != nil {
					Service.Logger.Error("handlers - api - marshal incoming data", zap.Error(err))
					return nil, err
				}

				req := models.EthRequest{}
				err = json.Unmarshal(b, &req)
				if err != nil {
					Service.Logger.Error("handlers - api - unmarshal data to struct", zap.Error(err))
					return nil, err
				}

				contract, err := Service.GetContractByName(req.Contract)
				if err != nil {
					Service.Logger.Error("handlers - api - GetContractByName", zap.Error(err))
					return nil, err
				}

				abiEl, err := abi.JSON(strings.NewReader(contract.ABI))
				if err != nil {
					log.Fatalf("Could not read ABI: %v", err)
				}

				ethURL := contract.Server
				if !ok {
					Service.Logger.Sugar().Fatalf("wrong type of eth_server value in config, type : %+v", reflect.TypeOf(Service.Context.GetConfig("specific.eth_server", "")))
				}

				ethClient, err := ethclient.Dial(ethURL)
				if err != nil {
					Service.Logger.Error("handlers - api - dial eth server", zap.Error(err))
					return nil, err
				}

				var args []interface{}
				for _, v := range req.Params {
					arg := v.Value

					if v.Type == "address" {
						arg = common.HexToAddress(v.Value.(string))
					}

					if v.Type == "uint256" {
						arg, ok = new(big.Int).SetString(v.Value.(string), 10)
						if !ok {
							Service.Logger.Error("handlers - api - can't convert to bigInt")
							return nil, errors.New("handlers - api - can't convert to bigInt")
						}
					}

					if v.Type == "[]string" {
						t := v.Value.([]interface{})
						s := make([]string, len(t))
						for i, a := range t {
							s[i] = fmt.Sprint(a)
						}

						arg = s
					}

					args = append(args, arg)
				}

				input, err := abiEl.Pack(req.Method, args...)

				if err != nil {
					Service.Logger.Error("handlers - api - pack eth server", zap.Error(err))
					return nil, err
				}

				response, err := Service.RawTransaction(ethClient, big.NewInt(0), input, contract)
				if err != nil {
					return nil, err
				}

				return response, nil
			},
		},

		"add": saiService.HandlerElement{
			Name:        "add",
			Description: "add contract to contracts",
			Function: func(data interface{}) (interface{}, error) {
				contractData, ok := data.(map[string]interface{})
				if !ok {
					Service.Logger.Sugar().Debugf("handlers - add - wrong data type, current type : %+v", data)
					return nil, errors.New("wrong type of incoming data")
				}

				b, err := json.Marshal(contractData)
				if err != nil {
					Service.Logger.Error("api - add - marshal incoming data", zap.Error(err))
					return nil, err
				}

				contracts := models.Contracts{}
				err = json.Unmarshal(b, &contracts)
				if err != nil {
					Service.Logger.Error("handlers - add - unmarshal data to struct", zap.Error(err))
					return nil, err
				}

				// validate all incoming contracts
				validatedContracts := make([]models.Contract, 0)
				for _, contract := range contracts.Contracts {
					err = contract.Validate()
					if err != nil {
						Service.Logger.Error("handlers - add - validate incoming contracts", zap.Any("contract", contract), zap.Error(err))
						continue
					}
					validatedContracts = append(validatedContracts, contract)
				}

				// check if incoming contracts already exists
				Service.Mutex.RLock()
				checkedContracts := Service.filterUniqueContracts(validatedContracts)
				Service.Mutex.RUnlock()

				Service.Mutex.Lock()
				Service.Contracts = append(Service.Contracts, checkedContracts...)
				Service.Mutex.Unlock()

				Service.Logger.Sugar().Debugf("ACTUAL CONTRACTS : %+v", Service.Contracts)

				err = Service.RewriteContractsConfig(contractsPath)
				if err != nil {
					Service.Logger.Error("handlers - add - rewrite contracts file", zap.Error(err))
					return nil, err
				}
				return "ok", nil

			},
		},

		"delete": saiService.HandlerElement{
			Name:        "delete",
			Description: "delete contract by name",
			Function: func(data interface{}) (interface{}, error) {
				deleteData, ok := data.(map[string]interface{})
				if !ok {
					Service.Logger.Sugar().Debugf("handlers - delete - wrong data type, current type : %+v", data)
					return nil, errors.New("wrong type of incoming data")
				}

				b, err := json.Marshal(deleteData)
				if err != nil {
					Service.Logger.Error("api - delete - marshal incoming data", zap.Error(err))
					return nil, err
				}

				deleteContractName := models.DeleteData{}
				err = json.Unmarshal(b, &deleteContractName)
				if err != nil {
					Service.Logger.Error("handlers - add - unmarshal data to struct", zap.Error(err))
					return nil, err
				}

				Service.Mutex.Lock()
				Service.DeleteContracts(&deleteContractName)
				Service.Mutex.Unlock()

				Service.Logger.Sugar().Debugf("CONTRACTS AFTER DELETION : %+v", Service.Contracts)

				err = Service.RewriteContractsConfig(contractsPath)
				if err != nil {
					Service.Logger.Error("handlers - delete - rewrite contracts file", zap.Error(err))
					return nil, err
				}
				return "ok", nil

			},
		},
	}

}
