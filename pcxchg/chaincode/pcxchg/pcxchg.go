
package main

import (
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim" // import for Chaincode Interface
    pb "github.com/hyperledger/fabric/protos/peer"      // import for peer response
    "encoding/json"
    "github.com/hyperledger/fabric/core/chaincode/lib/cid" // import CID
)

type PcXchg struct {
}

// Define our struct to store PCs in Blockchain, start fields upper case for JSON
type PC struct {
    Snumber string  // This one will be our key
    Serie string
    Others string
    Status string   // this will contain its status on the exchange
}

func (c *PcXchg) Init(stub shim.ChaincodeStubInterface) pb.Response { 
    mspid, err := cid.GetMSPID(stub)
    if(err != nil) {
        return shim.Error("GetMSPID err")
    }
    stub.PutState("producer", []byte(mspid))
    return shim.Success(nil) 
}


// Implement Invoke
func (c *PcXchg) Invoke(stub shim.ChaincodeStubInterface) pb.Response {

    // Get function name and args
    function, args := stub.GetFunctionAndParameters()

    // Check if the producer calls the chaincode
    if c.producer(stub) {
        if function == "createPC" {
            // A computer is produced and available
            return c.createPC(stub, args)
        }
    } else {
        if function == "buyPC" {
            // A market bought a computer
            return c.updateStatus(stub, args, "bought")
        } else if function == "handBackPC" {
            // A market handed back a computer
            return c.updateStatus(stub, args, "returned")
        }
    }

    if function == "queryStock" {
        // Stock query
        return c.queryStock(stub, args)
    } else if function == "queryDetail" {
        // Get details of a computer
        return c.queryDetail(stub, args)
    }

    return shim.Error("You can call createPc as producer and buyPC/handBackPC as market. Producer and market can call queryStock, queryDetail")
}

func (c *PcXchg) producer(stub shim.ChaincodeStubInterface) bool {
    mspid, _ := cid.GetMSPID(stub)

    v, err := stub.GetState("producer")
    if (err != nil) {
        return false
    }
    return string(v) == mspid
}


// createPC puts an available PC in the Blockchain
func (c *PcXchg) createPC(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 3 {
        return shim.Error("createPC arguments usage: Serialnumber, Serie, Others")
    }

    // A newly created computer is available
    pc := PC{args[0], args[1], args[2], "available"} 

    // Use JSON to store in the Blockchain
    pcAsBytes, err := json.Marshal(pc)

    if err != nil {
        return shim.Error(err.Error())
    }

    // Use serial number as key
    err = stub.PutState(pc.Snumber, pcAsBytes)

    if err != nil {
        return shim.Error(err.Error())
    }
    return shim.Success(nil)
}

// updateStatus handles sell and hand back
func (c *PcXchg) updateStatus(stub shim.ChaincodeStubInterface, args []string, status string) pb.Response {
    if len(args) != 1 {
        return shim.Error("This function needs the serial number as argument")
    }

    // Look for the serial number
    v, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("Serialnumber " + args[0] + " not found ")
    }

    // Get Information from Blockchain
    var pc PC
    // Decode JSON data
    json.Unmarshal(v, &pc)

    // Change the status
    pc.Status = status 
    // Encode JSON data
    pcAsBytes, err := json.Marshal(pc)

    // Store in the Blockchain
    err = stub.PutState(pc.Snumber, pcAsBytes)
    if err != nil {
        return shim.Error(err.Error())
    }

    return shim.Success(nil)
}

// queryDetail gives all fields of stored data and wants to have the serial number
func (c *PcXchg) queryDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // Look for the serial number
    value, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("Serial number " + args[0] + " not found")
    }

    var pc PC
    // Decode value
    json.Unmarshal(value, &pc)

    fmt.Print(pc)
    // Response info
    return shim.Success([]byte(" SNMBR: " + pc.Snumber + " Serie: " + pc.Serie + " Others: " + pc.Others + " Status: " + pc.Status))
}

// queryStock give all stored keys in the database
func (c *PcXchg) queryStock(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    // See stub.GetStateByRange in interfaces.go
    start, end := "",""

    if len(args) == 2 {
        start, end = args[0], args[1]
    } 

    // resultIterator is a StateQueryIteratorInterface
    resultsIterator, err := stub.GetStateByRange(start, end)
    if err != nil {
        return shim.Error(err.Error())
    }
    defer resultsIterator.Close()

    keys := " \n"
    // This interface includes HasNext,Close and Next
    for resultsIterator.HasNext() {
        queryResponse, err := resultsIterator.Next()
        if err != nil {
            return shim.Error(err.Error())
        }
        keys+=queryResponse.Key + " \n"
    }

    fmt.Println(keys)

    return shim.Success([]byte(keys))
}

func main() {
    err := shim.Start(new(PcXchg))
    if err != nil {
        fmt.Printf("Error starting chaincode sample: %s", err)
    }
}
