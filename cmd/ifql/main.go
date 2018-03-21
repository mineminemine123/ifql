package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"runtime"

	"github.com/influxdata/ifql"
	"github.com/influxdata/ifql/repl"
)

var verbose = flag.Bool("v", false, "print verbose output")

var hosts = make(hostList, 0)

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

func main() {
	flag.Usage = usage
	flag.Parse()

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
	repl := repl.New(c)

	args := flag.Args()
	switch len(args) {
	case 0:
		repl.Run()
	case 1:
		q, err := loadQuery(args[0])
		if err != nil {
			log.Fatal(err)
		}
		err = repl.Input(q)
		if err != nil {
			fmt.Println(err)
		}
	default:
		flag.Usage()
		os.Exit(1)
	}
}

func loadQuery(q string) (string, error) {
	if len(q) > 0 && q[0] == '@' {
		f, err := os.Open(q[1:])
		if err != nil {
			return "", err
		}
		defer f.Close()

		data, err := ioutil.ReadAll(f)
		if err != nil {
			return "", err
		}
		q = string(data)
	}
	return q, nil
}
