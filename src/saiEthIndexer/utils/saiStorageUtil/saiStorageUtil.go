package saiStorageUtil

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/saiset-co/saiEthIndexer/utils/saiUtil"
	"go.mongodb.org/mongo-driver/bson"
)

type Database struct {
	url      string
	email    string
	password string
}

func Storage(Url string, Email string, Password string) Database {
	return Database{
		url:      Url,
		email:    Email,
		password: Password,
	}
}

type StorageRequest struct {
	token      string
	collection string
	options    interface{}
	criteria   interface{}
	data       interface{}
}

func (s StorageRequest) toJson() ([]byte, error) {
	jsonObj := bson.M{"collection": s.collection}

	if s.data != nil {
		jsonObj["data"] = s.data
	}

	if s.criteria != nil {
		jsonObj["select"] = s.criteria
	}

	if s.options != nil {
		jsonObj["options"] = s.options
	}

	return json.Marshal(jsonObj)
}

func (db Database) Get(collectionName string, criteria interface{}, options interface{}, token string) (error, []byte) {
	request := StorageRequest{collection: collectionName, criteria: criteria, options: options}
	return db.makeRequest("get", request, token)
}

func (db Database) Put(collectionName string, data interface{}, token string) (error, []byte) {
	request := StorageRequest{collection: collectionName, data: data}
	return db.makeRequest("save", request, token)
}

func (db Database) Update(collectionName string, criteria interface{}, data interface{}, token string) (error, []byte) {
	request := StorageRequest{collection: collectionName, criteria: criteria, data: data}
	return db.makeRequest("update", request, token)
}

func (db Database) makeRequest(method string, request StorageRequest, token string) (error, []byte) {
	jsonStr, jsonErr := request.toJson()

	if jsonErr != nil {
		fmt.Println("Database request error: ", jsonErr)
		return jsonErr, []byte("")
	}

	return saiUtil.Send(db.url+"/"+method, bytes.NewBuffer(jsonStr), token)
}
