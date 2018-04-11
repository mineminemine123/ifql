package influxql

import (
	"encoding/json"
	"io"
	"log"
	"sort"
	"strconv"
	"time"

	"github.com/influxdata/ifql/query/execute"
)

type Writer struct{}

func (Writer) WriteTo(w io.Writer, results map[string]execute.Result) error {
	// TODO: This code was copy pasted from the hackathon branch that enabled Chronograf to talk to IFQL.
	// It may need a refactor and it will tests to ensure that it accurately transforms the data to the correct format.
	resp := response{
		Results: make([]result, len(results)),
	}

	// This process needs to be deterministic.
	// Process results in lexigraphical order.
	order := make([]string, 0, len(results))
	for name := range results {
		order = append(order, name)
	}
	sort.Strings(order)

	for i, name := range order {
		r := results[name]
		blocks := r.Blocks()
		seriesID := 0
		result := result{}
		err := blocks.Do(func(b execute.Block) error {
			s := series{Tags: b.Tags()}
			for _, c := range b.Cols() {
				if !c.Common {
					s.Columns = append(s.Columns, c.Label)
				}
			}
			seriesID++
			s.Name = strconv.Itoa(seriesID)
			times := b.Times()
			times.DoTime(func(ts []execute.Time, rr execute.RowReader) {
				for i := range ts {
					var v []interface{}
					for j, c := range rr.Cols() {
						if c.Common {
							continue
						}
						switch c.Type {
						case execute.TFloat:
							v = append(v, rr.AtFloat(i, j))
						case execute.TInt:
							v = append(v, rr.AtInt(i, j))
						case execute.TString:
							v = append(v, rr.AtString(i, j))
						case execute.TUInt:
							v = append(v, rr.AtUInt(i, j))
						case execute.TBool:
							v = append(v, rr.AtBool(i, j))
						case execute.TTime:
							v = append(v, rr.AtTime(i, j)/execute.Time(time.Second))
						default:
							v = append(v, "unknown")
						}
					}
					s.Values = append(s.Values, v)
				}
			})
			result.Series = append(result.Series, s)
			return nil
		})
		if err != nil {
			log.Println("Error iterating through results:", err)
		}
		resp.Results[i] = result
	}

	return json.NewEncoder(w).Encode(resp)
}

type response struct {
	Results []result `json:"results"`
	Err     string   `json:"err"`
}

type result struct {
	Series []series `json:"series"`
}

type series struct {
	Name    string            `json:"name"`
	Columns []string          `json:"columns"`
	Tags    map[string]string `json:"tags"`
	Values  [][]interface{}   `json:"values"`
}
