package main

import (
	"bitbucket.org/dhaliwalprince/funlang/context"
	"bitbucket.org/dhaliwalprince/funlang/parse"
	"bitbucket.org/dhaliwalprince/funlang/sema"
	"bitbucket.org/dhaliwalprince/funlang/ssa"
	_ "bitbucket.org/dhaliwalprince/funlang/ssa/analysis"
	_ "bitbucket.org/dhaliwalprince/funlang/ssa/passes"
	"flag"
	"fmt"
	"os"
	"runtime/pprof"
)

var cpuProfile = flag.Bool("cpuprofile", false, "collect cpu profile information")
var filename = flag.String("input", "", "input file")
var help = flag.Bool("help", false, "print help")

func main() {
	flag.Parse()
	ctx := &context.Context{}
	if *filename == "" {
		fmt.Println("please provide input file using -input option")
		os.Exit(1)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	var f *os.File
	var err error
	if *cpuProfile {
		f, err = os.Create("cpu.profile")
		if err != nil {
			fmt.Println("Unable to create file:  ", err.Error())
			os.Exit(2)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	p := parse.NewParserFromFile(ctx, *filename)
	a, err := p.Parse()
	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("::== resolving variables")
	errs := sema.ResolveProgram(a)
	if len(errs) > 0 {
		fmt.Print(errs)
	}
	fmt.Println(a)

	program := ssa.Emit(a, ctx)
	fmt.Println("::== ssa generation done")
	fmt.Print(program)

	fmt.Println("::== trying to optimize")
	passRunner := ssa.NewPassRunner(program)
	passRunner.AddNext(ssa.GetPass("dominators"))
	passRunner.Add(ssa.GetPass("mem2reg"))
	passRunner.Add(ssa.GetPass("dce"))
	passRunner.Add(ssa.GetPass("verifier"))
	passRunner.RunPasses()
	fmt.Println("::== done")
	fmt.Print(program)
}