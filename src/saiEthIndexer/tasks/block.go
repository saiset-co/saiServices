package tasks

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/adam-lavrik/go-imath/ix"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/onrik/ethrpc"
	"github.com/saiset-co/saiEthIndexer/config"
	"github.com/saiset-co/saiEthIndexer/utils/saiStorageUtil"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

var startBlock int

type BlockManager struct {
	config    *config.Configuration
	storage   saiStorageUtil.Database
	logger    *zap.Logger
	websocket *WebsocketManager
}

type Block struct {
	ID int `json:"id"`
}

type LogTransfer struct {
	Type   string
	From   common.Address
	To     common.Address
	Tokens big.Int
}

func NewBlockManager(c config.Configuration, logger *zap.Logger) *BlockManager {

	manager := &BlockManager{
		config:    &c,
		storage:   saiStorageUtil.Storage(c.Specific.Storage.URL, c.Storage.Email, c.Storage.Password),
		logger:    logger,
		websocket: NewWebSocketManager(c),
	}

	return manager
}

func (bm *BlockManager) GetLastBlock(id int) (*Block, error) {
	block := Block{ID: id}
	pwd, err := os.Getwd()
	if err != nil {
		bm.logger.Error("tasks - BlockManager - get currect directory", zap.Error(err))
		return &block, nil
	}

	data, err := ioutil.ReadFile(pwd + "/block.data")
	if err != nil {
		bm.logger.Error("tasks - BlockManager - read file", zap.Error(err))
		return &block, nil
	}

	lastDataBlock, strErr := strconv.Atoi(string(data))

	if strErr != nil {
		log.Println("Data from file can't be converted to int:", err)
		return &block, nil
	}

	var lastBlocks []int
	for i := 0; i < len(bm.config.EthContracts.Contracts); i++ {
		lastBlocks = append(lastBlocks, bm.config.EthContracts.Contracts[i].StartBlock)
	}

	if lastDataBlock > 0 {
		startBlock = lastDataBlock
	} else if len(lastBlocks) > 0 {
		startBlock = ix.MinSlice(lastBlocks)
	} else if bm.config.StartBlock > 0 {
		startBlock = bm.config.StartBlock
	} else {
		startBlock = id
	}

	return &Block{ID: startBlock}, nil
}

func (bm *BlockManager) SetLastBlock(blk *Block) error {
	pwd, err := os.Getwd()

	if err != nil {
		bm.logger.Error("tasks - BlockManager - set last block - read currect directory", zap.Error(err))
		return err
	}

	lastBlock := strconv.Itoa(blk.ID)
	err = ioutil.WriteFile(pwd+"/block.data", []byte(lastBlock), 0777)
	if err != nil {
		bm.logger.Error("tasks - BlockManager - set last block - write to file", zap.Error(err))
	}

	bm.logger.Sugar().Debugf("block %d was saved to the temp file", blk.ID)
	return nil
}

func (bm *BlockManager) HandleReceipts(receipt *ethrpc.TransactionReceipt, _abi abi.ABI) ([]map[string]interface{}, error) {
	var events []map[string]interface{}

	for _, l := range receipt.Logs {
		id := common.HexToHash(l.Topics[0])
		_event, eventErr := _abi.EventByID(id)
		if eventErr != nil {
			continue
		}

		data := map[string]interface{}{}
		event := map[string]interface{}{}

		d, _ := hex.DecodeString(l.Data[2:])
		unpackErr := _event.Inputs.UnpackIntoMap(data, d)

		if unpackErr != nil {
			fmt.Println("can't unpack event:", unpackErr)
			continue
		}

		for eventId, eventData := range _event.Inputs {
			if eventData.Indexed {
				data[eventData.Name] = l.Topics[eventId+1]
			}
		}

		data["name"] = _event.Name

		event["Data"] = data
		event["Log"] = l

		events = append(events, event)
	}

	return events, nil
}

func (bm *BlockManager) HandleTransactions(trs []ethrpc.Transaction, receipts map[string]*ethrpc.TransactionReceipt) {
	for j := 0; j < len(trs); j++ {
		for i := 0; i < len(bm.config.EthContracts.Contracts); i++ {
			methodName := "Unknown"
			decodedInput := map[string]interface{}{}
			status, _ := strconv.ParseBool(receipts[trs[j].Hash].Status[2:])

			if bm.config.SkipFailedTransactions && !status {
				continue
			}

			if strings.ToLower(trs[j].From) != strings.ToLower(bm.config.EthContracts.Contracts[i].Address) && strings.ToLower(trs[j].To) != strings.ToLower(bm.config.EthContracts.Contracts[i].Address) {
				continue
			}

			bm.logger.Sugar().Debugf("TO %s", trs[j].To)
			bm.logger.Sugar().Debugf("From %s", trs[j].From)

			raw, err := json.Marshal(trs[j])
			if err != nil {
				bm.logger.Error("block manager - handle transaction - marshal transaction", zap.String("transaction hash", trs[j].Hash), zap.Error(err))
				continue
			}

			data := bson.M{
				"Number": trs[j].BlockNumber,
				"Hash":   trs[j].Hash,
				"From":   trs[j].From,
				"To":     trs[j].To,
				"Amount": trs[j].Value,
			}

			_abi, abiErr := abi.JSON(strings.NewReader(bm.config.EthContracts.Contracts[i].ABI))
			if abiErr != nil {
				bm.logger.Error("block manager - handle transaction - parse abi from config", zap.String("address", bm.config.EthContracts.Contracts[i].Address), zap.Error(abiErr))
				continue
			}

			events, trErr := bm.HandleReceipts(receipts[trs[j].Hash], _abi)
			if trErr != nil {
				bm.logger.Error("block manager - handle transaction events - HandleReceipts", zap.String("transaction hash", trs[j].Hash), zap.Error(trErr))
				continue
			}

			if len(trs[j].Input) >= 10 {
				decodedInput = map[string]interface{}{}
				decodedSig, decSigErr := hex.DecodeString(trs[j].Input[2:10])
				if decSigErr != nil {
					bm.logger.Error("block manager - handle transaction - decode transaction function identifier", zap.String("transaction hash", trs[j].Hash), zap.Error(decSigErr))
					continue
				}

				method, methodErr := _abi.MethodById(decodedSig)
				if methodErr != nil {
					bm.logger.Error("block manager - handle transaction - MethodById", zap.String("transaction hash", trs[j].Hash), zap.Error(methodErr))
					continue
				}

				methodName = method.Name
				decodedData, decInErr := hex.DecodeString(trs[j].Input[2:])
				if decInErr != nil {
					bm.logger.Error("block manager - handle transaction - decode input", zap.String("transaction hash", trs[j].Hash), zap.Error(decInErr))
					continue
				}

				unpackErr := method.Inputs.UnpackIntoMap(decodedInput, decodedData[4:])
				if unpackErr != nil {
					bm.logger.Error("block manager - handle transaction - UnpackIntoMap", zap.String("transaction hash", trs[j].Hash), zap.Error(unpackErr))
					continue
				}
			}

			data["Events"] = events
			data["Status"] = status
			data["Operation"] = methodName
			data["Input"] = decodedInput

			for _, operation := range bm.config.Operations {
				if operation == methodName {
					wsErr := bm.websocket.SendMessage(string(raw), bm.config.WebSocket.Token)
					if wsErr != nil {
						bm.logger.Error("block manager - handle transaction - SendWebSocketMsg", zap.String("transaction hash", trs[j].Hash), zap.Error(wsErr))
						continue
					}
				}
			}

			storageErr, _ := bm.storage.Put("transactions", data, bm.config.Storage.Token)
			if storageErr != nil {
				bm.logger.Error("block manager - handle transaction - storage.Put", zap.String("transaction hash", trs[j].Hash), zap.Error(storageErr))
				continue
			}

			bm.logger.Sugar().Infof("%d transaction from %s to %s has been updated.\n", trs[j].TransactionIndex, trs[j].From, trs[j].To)
		}
	}
}
