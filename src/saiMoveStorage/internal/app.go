package internal

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type AddressesStorageResult struct {
	Result []struct {
		Address string `json:"address"`
	}
}

type StorageResult struct {
	Result []interface{}
}

func (is InternalService) MoveWhiteList() error {
	errGet, data := is.StorageFrom.Get("whitelist", bson.M{}, bson.M{})

	if errGet != nil {
		is.Logger.Error("MoveWhiteList", zap.Error(errGet))
		return errGet
	}

	var _data StorageResult
	errMarshal := json.Unmarshal(data, &_data)

	if errMarshal != nil {
		is.Logger.Error("MoveWhiteList", zap.Error(errMarshal))
		return errMarshal
	}

	for _, document := range _data.Result {
		errPut, _ := is.StorageTo.Put("whitelist", document)

		if errPut != nil {
			is.Logger.Error("MoveWhiteList", zap.Error(errPut))
			return errPut
		}
	}

	return nil
}

func (is InternalService) MoveWhiteListChart() error {
	errGet, data := is.StorageFrom.Get("whitelistChart", bson.M{}, bson.M{})

	if errGet != nil {
		is.Logger.Error("whitelistChart", zap.Error(errGet))
		return errGet
	}

	var _data StorageResult
	errMarshal := json.Unmarshal(data, &_data)

	if errMarshal != nil {
		is.Logger.Error("whitelistChart", zap.Error(errMarshal))
		return errMarshal
	}

	//is.Logger.Debug("whitelistChart", zap.ByteString("result", data))

	for _, document := range _data.Result {
		errPut, _ := is.StorageTo.Put("whitelistChart", document)

		if errPut != nil {
			is.Logger.Error("whitelistChart", zap.Error(errPut))
			return errPut
		}
	}

	return nil
}

func (is InternalService) MoveWhiteListedAddresses() error {
	errGet, data := is.StorageFrom.Get("whiteListedAddresses", bson.M{}, bson.M{})

	if errGet != nil {
		is.Logger.Error("whiteListedAddresses", zap.Error(errGet))
		return errGet
	}

	var _data StorageResult
	errMarshal := json.Unmarshal(data, &_data)

	if errMarshal != nil {
		is.Logger.Error("whiteListedAddresses", zap.Error(errMarshal))
		return errMarshal
	}

	//is.Logger.Debug("MoveWhiteListedAddresses", zap.ByteString("result", data))

	for _, document := range _data.Result {
		errPut, _ := is.StorageTo.Put("whiteListedAddresses", document)

		if errPut != nil {
			is.Logger.Error("whiteListedAddresses", zap.Error(errPut))
			return errPut
		}
	}

	return nil
}

func (is InternalService) MoveUsers() error {
	errGet1, data1 := is.StorageFrom.Get("whiteListedAddresses", bson.M{}, bson.M{})

	if errGet1 != nil {
		is.Logger.Error("users", zap.Error(errGet1))
		return errGet1
	}

	var _data1 AddressesStorageResult
	errMarshal := json.Unmarshal(data1, &_data1)

	if errMarshal != nil {
		is.Logger.Error("users", zap.Error(errMarshal))
		return errMarshal
	}

	list := CollectAddresses(_data1)

	errGet2, data2 := is.StorageFrom.Get("users", bson.M{"walletAddresses": bson.M{"$in": list}}, bson.M{})

	if errGet2 != nil {
		is.Logger.Error("users", zap.Error(errGet2))
		return errGet2
	}

	var _data2 StorageResult
	errMarshal2 := json.Unmarshal(data2, &_data2)

	if errMarshal2 != nil {
		is.Logger.Error("users", zap.Error(errMarshal2))
		return errMarshal2
	}

	for _, document := range _data2.Result {
		errPut, _ := is.StorageTo.Put("users", document)

		if errPut != nil {
			is.Logger.Error("users", zap.Error(errPut))
			return errPut
		}
	}

	return nil
}

func (is InternalService) Move() (interface{}, error) {
	if errMoveWhiteList := is.MoveWhiteList(); errMoveWhiteList != nil {
		return "NOK", errMoveWhiteList
	}

	if errMoveWhiteListChart := is.MoveWhiteListChart(); errMoveWhiteListChart != nil {
		return "NOK", errMoveWhiteListChart
	}

	if errMoveWhiteListedAddresses := is.MoveWhiteListedAddresses(); errMoveWhiteListedAddresses != nil {
		return "NOK", errMoveWhiteListedAddresses
	}

	if errMoveUsers := is.MoveUsers(); errMoveUsers != nil {
		return "NOK", errMoveUsers
	}

	return "Ok", nil
}

func CollectAddresses(data AddressesStorageResult) []string {
	var list []string
	for _, addrress := range data.Result {
		list = append(list, addrress.Address)
	}
	return list
}
