package main

import (
	"fmt"
	"os"
	"transactional-kvs/kvs"
)

func main() {
	inputs := os.Args[1:]

	// validate args
	err := kvs.ValidateCmdInputs(inputs)
	if err != nil {
		fmt.Println(err, inputs)
	}

	kvsInst := kvs.NewSimpleKvs()
	kvsInst.Set("asd", "dsa")
	fmt.Println(kvsInst.Get("asd"))

}
