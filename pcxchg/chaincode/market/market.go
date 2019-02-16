package main

import (
    "github.com/hyperledger/fabric/core/chaincode/shim"
    pb "github.com/hyperledger/fabric/protos/peer"
    "fmt"
)

type Market struct {
}

func (c *Market) Init(stub shim.ChaincodeStubInterface) pb.Response { 
    return shim.Success(nil) 
}

func (c *Market) Invoke(stub shim.ChaincodeStubInterface) pb.Response { 

    function, args := stub.GetFunctionAndParameters()

    if(function=="buyPC") {
        if(len(args)!=2) {
            return shim.Error("This function needs the serial number and the chaincode to be invoked")
        }

        callArgs:= make([][]byte, 3)

        callArgs[0]= []byte("updateStatus")
        callArgs[1]= []byte(args[0])
        callArgs[2]= []byte("bought")

        return stub.InvokeChaincode(args[1],  callArgs, "")

        } else {
            return shim.Error("You can buyPC(serialnumber, chaincode)")
        }

        return shim.Success(nil)
}

func main() {
    err := shim.Start(new(Market))
    if err != nil {
        fmt.Printf("Error starting chaincode sample: %s", err)
    }
}
