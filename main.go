package main

import (
	"context"
	"crypto/rand"
	"fmt"

	crpc "github.com/celestiaorg/celestia-node/api/rpc/client"
	"github.com/celestiaorg/celestia-node/share"
	"github.com/rollkit/celestia-da/celestia"
	goDA "github.com/rollkit/go-da"
)

var daRpcAddr string = "http://localhost:26658"
var gasPrice float64 = 0.1

func main() {
	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJBbGxvdyI6WyJwdWJsaWMiLCJyZWFkIiwid3JpdGUiLCJhZG1pbiJdfQ.Un4ghSc25h0IN6BK7Sg7sKZL-GVDIx5aEfFPZ-pf7H8"
	rpcClient, err := crpc.NewClient(context.TODO(), daRpcAddr, token)
	if err != nil {
		panic(err)
	}

	balance, err := rpcClient.State.Balance(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(balance)

	id := make([]byte, 3)
	id[0] = 1
	id[1] = 3
	id[2] = 5

	ns, err := share.NewBlobNamespaceV0(id)
	if err != nil {
		panic(err)
	}

	b, err := rpcClient.State.Balance(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	da := celestia.NewCelestiaDA(rpcClient, ns, gasPrice, context.TODO())

	maxBlobSize, err := da.MaxBlobSize(context.TODO())
	if err != nil {
		panic(err)
	}
	fmt.Printf("maxBlobSize: %d\n", maxBlobSize)

	// random block bytes
	data := make([]byte, 100)
	_, err = rand.Read(data)
	if err != nil {
		panic(err)
	}

	// kinda misleading, since [da.Submit] wrap up a blob, the blobs here are actually datas
	blobs := make([]goDA.Blob, 0, 1)
	blobs = append(blobs, data)

	// submit takes a long time, 15+s
	// set gasPrice to -1, where [da] will use the global configured gasPrice
	ids, err := da.Submit(context.TODO(), blobs, -1, nil)
	if err != nil {
		panic(err)
	}

	proofs, err := da.GetProofs(context.TODO(), ids, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(proofs[0])

	verified, err := da.Validate(context.TODO(), ids, proofs, nil)
	if err != nil {
		panic(err)
	}

	for _, v := range verified {
		if v != true {
			panic(fmt.Errorf("not verified"))
		}
	}
}
