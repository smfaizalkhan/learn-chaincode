/*
Copyright 2016 Chinasystems
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
Licensed Materials - Property of Chinasystems
© Copyright Chinasystems Corp. 2016
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// GoodsInspectionChaincode example simple Chaincode implementation
type GoodsInspectionChaincode struct {
}

// GoodsInspection Data
type GoodsInspection struct {
	Reference string `json:"reference"`
	Bank      string `json:"bank"`
	Inspector string `json:"inspector"`
	Status    string `json:"status"`
	Document  string `json:"document"`
}

// Init function
func (t *GoodsInspectionChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
          fmt.Println("Initialising")
	if len(args) != 0 {
		return nil, errors.New("Incorrect number of arguments. Expecting 0")
	}

	return nil, nil
}

// Invoke runs callback representing the invocation of a chaincode
func (t *GoodsInspectionChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "delete" {
		// Deletes an entity from its state
		return t.delete(stub, args)
	} else if function == "register" {
		// Register goods
		return t.register(stub, args)
	} else if function == "inspect" {
		// Inspect goods
		return t.inspect(stub, args)
	}

	return nil, errors.New("Incorrect function name")

}

func (t *GoodsInspectionChaincode) register(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Registering Goods...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var goodsInspection GoodsInspection
	var err error

	goodsInspection.Reference = args[0]
	goodsInspection.Bank = args[1]
	goodsInspection.Inspector = args[2]
	goodsInspection.Status = "pending"
	goodsInspection.Document = ""

	goodsInspectionBytes, err := json.Marshal(&goodsInspection)
	if err != nil {
		fmt.Println("Error marshalling goodsInspection")
		return nil, errors.New("Error registering goods inspector")
	}
	err = stub.PutState(goodsInspection.Reference, goodsInspectionBytes)
	if err != nil {
		fmt.Println("Error registering goods inspector")
		return nil, errors.New("Error registering goods")
	}

	return nil, nil
}

// Inspect goods
func (t *GoodsInspectionChaincode) inspect(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inspecting Goods...")

	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	var goodsInspection GoodsInspection
	var goodsInspectionOld GoodsInspection
	var err error

	goodsInspection.Reference = args[0]
	goodsInspection.Status = args[1]
	goodsInspection.Document = args[2]

	fmt.Println("Getting State on reference number " + goodsInspection.Reference)
	goodsInspectionBytes, err := stub.GetState(goodsInspection.Reference)
	if goodsInspectionBytes == nil {
		return nil, errors.New("Failed retrieving goods on reference number " + goodsInspection.Reference)
	}

	err = json.Unmarshal(goodsInspectionBytes, &goodsInspectionOld)
	if err != nil {
		fmt.Println("Error Unmarshalling goodsInspectionBytes")
		return nil, errors.New("Error retrieving goods on reference number " + goodsInspection.Reference)
	}

	goodsInspection.Bank = goodsInspectionOld.Bank
	goodsInspection.Inspector = goodsInspectionOld.Inspector

	goodsInspectionBytesNew, err := json.Marshal(&goodsInspection)
	if err != nil {
		fmt.Println("Error marshalling goodsInspection")
		return nil, errors.New("Error marshalling goodsInspection")
	}

	err = stub.PutState(goodsInspection.Reference, goodsInspectionBytesNew)
	if err != nil {
		fmt.Println("Error Inspecting goods")
		return nil, errors.New("Error Inspecting goods")
	}

	return nil, err
}

// Deletes an entity from state
func (t *GoodsInspectionChaincode) delete(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	Ref := args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(Ref)
	if err != nil {
		return nil, errors.New("Failed to delete state")
	}

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *GoodsInspectionChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function != "query" {
		return nil, errors.New("Invalid query function name. Expecting \"query\"")
	}

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	// ref number
	Ref := args[0]

	fmt.Println("Getting State on reference number " + Ref)
	GoodsInspectionBytes, err := stub.GetState(Ref)
	if err != nil {
		return nil, errors.New("Failed to query state")
	}
	if GoodsInspectionBytes == nil {
		return nil, errors.New("Failed retrieving goods on reference number " + Ref)
	}
	//fmt.Println("Result: " + string(GoodsInspectionBytes))

	return GoodsInspectionBytes, nil
}

func main() {
	err := shim.Start(new(GoodsInspectionChaincode))
	if err != nil {
		fmt.Printf("Error starting Goods Inspection chaincode: %s", err)
	}
}
