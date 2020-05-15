package server

import (
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/nobonobo/wecty/cmd/wecty/commands"
)

// Server ...
type Server struct {
	*flag.FlagSet
	addr     string
	isTinyGo bool
}

// Usage ...
func (s *Server) Usage() {
	s.FlagSet.Usage()
}

// Execute ...
func (s *Server) Execute(args []string) error {
	if err := s.FlagSet.Parse(args); err != nil {
		return err
	}
	log.Println("listen and serve:", s.addr)
	tempDir, err := ioutil.TempDir("", "gobuild")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		return err
	}
	d, err := os.Getwd()
	if err != nil {
		return err
	}
	dh := &DevHandler{
		goCmd:        "go",
		wasmExecPath: filepath.Join(runtime.GOROOT(), "misc", "wasm", "wasm_exec.js"),
		workDir:      d,
		tempDir:      tempDir,
	}
	if s.isTinyGo {
		dh.goCmd = "tinygo"
		output, err := commands.RunCmd(d, nil, "tinygo", "env", "TINYGOROOT")
		if err != nil {
			return err
		}
		dh.wasmExecPath = filepath.Join(strings.TrimSpace(output), "targets", "wasm_exec.js")
	}
	return http.Serve(l, dh)
}

func init() {
	s := &Server{
		FlagSet: flag.NewFlagSet("server", flag.ContinueOnError),
	}
	s.FlagSet.StringVar(&s.addr, "addr", ":8080", "listen address")
	s.FlagSet.BoolVar(&s.isTinyGo, "tinygo", false, "use tinygo tool chain")
	commands.Register(s)
}
