package main

import (
	"log"

	prompt "github.com/c-bata/go-prompt"
	_ "github.com/influxdata/ifql/functions"
	"github.com/influxdata/ifql/query"
)

func init() {
	query.FinalizeRegistration()
}

var repl *query.REPL

func completer(d prompt.Document) []prompt.Suggest {
	names := repl.Scope.Names()
	log.Println(names)
	s := make([]prompt.Suggest, len(names))
	for i, n := range repl.Scope.Names() {
		s[i] = prompt.Suggest{
			Text: n,
		}
	}
	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
}

func main() {
	repl = query.NewREPL()
	p := prompt.New(
		repl.Input,
		completer,
		prompt.OptionPrefix("> "),
		prompt.OptionTitle("ifql"),
	)
	p.Run()
}
