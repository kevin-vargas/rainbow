package main

import (
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
	messages := make(chan []wiz.Option)
	m := wiz.New(args.Address, args.MacAddress, messages)
	j := joy.New(messages)
	go m.Start()
	j.Start()
}
