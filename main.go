package main

import "github.com/0xsharma/compact-chain/cmd"

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
