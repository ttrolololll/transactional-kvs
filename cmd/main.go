package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"transactional-kvs/kvs"
)

func main() {
	cmdReader := bufio.NewReader(os.Stdin)
	kvsInst := kvs.NewSimpleKvs()

	for {
		fmt.Print("> ")
		inputText, _ := cmdReader.ReadString('\n')
		inputs := strings.Fields(inputText)

		result, err := kvsInst.CommandExecutor(inputs)
		if err != nil {
			fmt.Println(err)
		}
		if result != "" {
			fmt.Println(result)
		}
	}
}
