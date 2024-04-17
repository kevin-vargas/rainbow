package main

import (
	"rainbow/joy"
	"rainbow/wiz"

	"github.com/alexflint/go-arg"
)

var args struct {
	Playlist     []string `arg:"required"`
	SpotyAddress string   `arg:"required"`
	Address      string   `arg:"required"`
}

func main() {
	arg.MustParse(&args)
	messages := make(chan []wiz.Option)
	m := wiz.New(args.Address, "", messages)
	j := joy.New(messages, args.SpotyAddress, args.Playlist)
	go m.Start()
	j.Start()
}
