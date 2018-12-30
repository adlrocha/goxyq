package main

import "testing"

// func checkInvoke(t *testing.T, stub *shim.MockStub, args [][]byte) {
// 	fmt.Println("[*] Checking invoke function:", string(args[0]))
// 	res := stub.MockInvoke("1", args)
// 	if res.Status != shim.OK {
// 		fmt.Println("Invoke", args, "failed", string(res.Message))
// 		t.FailNow()
// 	}
// }
func TestTrust(t *testing.T) {
	// 	scc := new(SmartContract)
	// 	stub := shim.NewMockStub("trust", scc)

	// 	checkInit(t, stub, nil)

	// 	asset := []byte(`{"id":"as1","owner":"admin","data":{"metric1":"a","metric2":"b"}}`)
	// 	metadata := []byte(`{"metric1":"c","metric2":"b", "metric3":"c"}`)
	// 	checkInvoke(t, stub, [][]byte{[]byte("create"), asset})
	// 	checkGetAsset(t, stub, "as1")
	// 	checkInvoke(t, stub, [][]byte{[]byte("update"), []byte("as1"), metadata})
	// 	checkGetAsset(t, stub, "as1")
	// 	// checkHistory(t, stub, "as1")
}
