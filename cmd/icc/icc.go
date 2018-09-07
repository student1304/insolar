/*
 *    Copyright 2018 INS Ecosystem
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/insolar/insolar/logicrunner/goplugin/preprocessor"
	"github.com/insolar/insolar/logicrunner/goplugin/testutil"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

func printUsage() {
	fmt.Println("usage: icc <command> [<args>]")
	fmt.Println("Commands: ")
	fmt.Println(" wrapper   generate contract's wrapper")
	fmt.Println(" proxy     generate contract's proxy")
	fmt.Println(" compile   compile contract")
	fmt.Println(" imports   rewrite imports")
}

type outputFlag struct {
	path   string
	writer io.Writer
}

func newOutputFlag() *outputFlag {
	return &outputFlag{path: "-", writer: os.Stdout}
}

func (r *outputFlag) String() string {
	return r.path
}
func (r *outputFlag) Set(arg string) error {
	var res io.Writer
	if arg == "-" {
		res = os.Stdout
	} else {
		var err error
		res, err = os.OpenFile(arg, os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return errors.Wrap(err, "couldn't open file for writing")
		}
	}
	r.path = arg
	r.writer = res
	return nil
}
func (r *outputFlag) Type() string {
	return "file"
}

func main() {

	if len(os.Args) == 1 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "wrapper":
		fs := flag.NewFlagSet("wrapper", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		for _, fn := range fs.Args() {
			err := preprocessor.GenerateContractWrapper(fn, output.writer)
			if err != nil {
				panic(err)
			}
		}
	case "proxy":
		fs := flag.NewFlagSet("proxy", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		var reference string
		fs.StringVarP(&reference, "code-reference", "r", "testRef", "reference to code of")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("proxy command should be followed by exactly one file name to process"))
		}
		err = preprocessor.GenerateContractProxy(fs.Arg(0), reference, output.writer)
		if err != nil {
			panic(err)
		}
	case "imports":
		fs := flag.NewFlagSet("imports", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output file (use - for STDOUT)")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		if fs.NArg() != 1 {
			panic(errors.New("imports command should be followed by exactly one file name to process"))
		}

		err = preprocessor.CmdRewriteImports(fs.Arg(0), output.writer)
		if err != nil {
			panic(err)
		}
	case "compile":
		fs := flag.NewFlagSet("compile", flag.ExitOnError)
		output := newOutputFlag()
		fs.VarP(output, "output", "o", "output directory")
		name := newOutputFlag()
		fs.VarP(name, "name", "n", "contract's file name")
		err := fs.Parse(os.Args[2:])
		if err != nil {
			panic(err)
		}

		dstDir := output.String() + "/plugins/"
		err = os.MkdirAll(dstDir, 0777)
		if err != nil {
			panic(err)
		}

		origGoPath, err := testutil.ChangeGoPath(output.String())
		if err != nil {
			panic(err)
		}
		defer os.Setenv("GOPATH", origGoPath) // nolint: errcheck

		//contractPath := root + "/src/contract/" + name + "/main.go"

		out, err := exec.Command("go", "build", "-buildmode=plugin", "-o", dstDir+"/"+name.String()+".so", "contract/"+name.String()).CombinedOutput()
		if err != nil {
			panic(errors.Wrap(err, "can't build contract: "+string(out)))
		}
	default:
		printUsage()
		fmt.Printf("\n\n%q is not valid command.\n", os.Args[1])
		os.Exit(2)
	}
}
