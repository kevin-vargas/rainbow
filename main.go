package main

import (
	"fmt"
	"rainbow/joy"
	"rainbow/wiz"

	"github.com/alexflint/go-arg"
)

var args struct {
	Address    string `arg:"required"`
	MacAddress string `arg:"required"`
}

func main() {
	arg.MustParse(&args)
	fmt.Println(args.Address)
	messages := make(chan []wiz.Option)
	m := wiz.New(args.Address, args.MacAddress, messages)
	j := joy.New(messages)
	go m.Start()
	j.Start()
}
