package main

import (
    "fmt"
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
    "encoding/json"
)

type Producer struct {
}

type PC struct {
    Snumber string
    Serie string
    Others string
    Status string
}

func (c *Producer) Init(stub shim.ChaincodeStubInterface) pb.Response { 
    return shim.Success(nil) 
}

func (c *Producer) Invoke(stub shim.ChaincodeStubInterface) pb.Response { 

    function, args := stub.GetFunctionAndParameters()

    switch function {
    case "createPC":
        return c.createPC(stub, args)
    case "updateStatus":
        return c.updateStatus(stub, args)
    case "queryDetail":
        return c.queryDetail(stub, args)
    default:
        return shim.Error("Functions: createPC, updateStatus, queryDetail")
    }
}

// createPC puts an aviable PC in the Blockchain
func (c *Producer) createPC(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    if len(args) != 3 {
        return shim.Error("createPC needs Serialnumber, Serie, Others as arguments")
    }

    pc := PC{args[0], args[1], args[2], "available"}

    pcAsBytes, err := json.Marshal(pc)
    if err != nil {
        return shim.Error(err.Error())
    }

    stub.PutState(pc.Snumber, pcAsBytes)
    return shim.Success(nil)
}

// updateStatus handles sell and hand back
func (c *Producer) updateStatus(stub shim.ChaincodeStubInterface, args []string) pb.Response {
    if len(args) != 2 {
        return shim.Error("This function needs the serial number and the new status as arguments")
    }

    v, err := stub.GetState(args[0])
    if err != nil {
        return shim.Error("Serialnumber " + args[0] + " not found ")
    }

    var pc PC
    json.Unmarshal(v, &pc)

    pc.Status = args[1]
    pcAsBytes, err := json.Marshal(pc)

    stub.PutState(pc.Snumber, pcAsBytes)
    return shim.Success(nil)
}

// queryDetail gives all fields of stored data and want to have the serial number
func (c *Producer) queryDetail(stub shim.ChaincodeStubInterface, args []string) pb.Response {

    value, err := stub.GetState(args[0]) // look for the serial number
    if err != nil {
        return shim.Error("Serial number " + args[0] + " not found")
    }

    var pc PC
    json.Unmarshal(value, &pc) // decode value

    fmt.Print(pc)
    return shim.Success([]byte(" SNMBR: " + pc.Snumber + " Serie: " + pc.Serie + " Others: " + pc.Others + " Status: " + pc.Status)) // response info
}

func main() {
    err := shim.Start(new(Producer))
    if err != nil {
        fmt.Printf("Error starting chaincode sample: %s", err)
    }
}