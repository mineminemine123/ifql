package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/influxdata/ifql"
	"github.com/influxdata/ifql/query"
	"github.com/influxdata/ifql/query/crossexecute/influxql"
	"github.com/influxdata/ifql/repl"
)

var (
	verbose   = flag.Bool("v", false, "print verbose output")
	hosts     = make(hostList, 0)
	transpile = flag.Bool("transpile", false, "transpiles an influxql query and outputs the spec (does not execute the query)")
)

func init() {
	flag.Var(&hosts, "host", "An InfluxDB host to connect to. Can be provided multiple times.")
}

type hostList []string

func (l *hostList) String() string {
	return "<host>..."
}

func (l *hostList) Set(s string) error {
	*l = append(*l, s)
	return nil
}

var defaultStorageHosts = []string{"localhost:8082"}

func usage() {
	fmt.Println("Usage: ifql [OPTIONS] [query]")
	fmt.Println()
	fmt.Println("Runs queries using the IFQL engine.")
	fmt.Println()
	fmt.Println("If no query is provided an interactive REPL will be run.")
	fmt.Println()
	fmt.Println("The query argument is either a string query or a path to a file prefixed with an '@'.")
	fmt.Println()
	fmt.Println("Options:")

	flag.PrintDefaults()
}

func runTranspile(args []string) {
	q, err := repl.LoadQuery(args[0])
	if err != nil {
		log.Fatal(err)
	}

	transpiler := &influxql.Transpiler{}
	spec, err := transpiler.Transpile(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("query", query.Formatted(spec, query.FmtJSON))
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *transpile {
		if flag.NArg() != 1 {
			log.Fatal("Transpile must have exactly one argument")
		}
		runTranspile(flag.Args())
		return
	}

	if len(hosts) == 0 {
		hosts = defaultStorageHosts
	}

	c, err := ifql.NewController(ifql.Config{
		Hosts:            hosts,
		ConcurrencyQuota: runtime.NumCPU() * 2,
		MemoryBytesQuota: math.MaxInt64,
		Verbose:          *verbose,
	})
	if err != nil {
		log.Fatal(err)
	}
	replCmd := repl.New(c)

	args := flag.Args()
	switch len(args) {
	case 0:
		replCmd.Run()
	case 1:
		q, err := repl.LoadQuery(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = replCmd.Input(q)
		if err != nil {
			fmt.Println(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}
