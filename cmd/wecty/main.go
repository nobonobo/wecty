package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nobonobo/wecty/cmd/wecty/commands"
	_ "github.com/nobonobo/wecty/cmd/wecty/commands/generate"
	_ "github.com/nobonobo/wecty/cmd/wecty/commands/server"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of wecty:\n")
		fmt.Fprintf(os.Stderr, "  commands:\n")
		fmt.Fprintf(os.Stderr, "    generate\thtml to go for wecty code generator\n")
		fmt.Fprintf(os.Stderr, "    server\tdevelopment http server\n")
	}
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		flag.Usage()
		os.Exit(1)
	}
	if err := commands.Execute(args[0], args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
