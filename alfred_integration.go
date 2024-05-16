package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/reinhrst/fzf-lib"
)

// tpf needs to be the path to the file containing the contexts (a
// complex JSON format. args is the argument to the command.
// TODO(rjk): It would be possible to make this a generic function
// so that one function would be used in many circumstances.
func genAlfredResult(entries []string, fzfies []fzf.MatchResult) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "	")
	var result Result
	dedup := make(map[string]struct{})

	// Create a result based on the hash
	for _, v := range fzfies {
		// Exclude the Uid field to make sure that the items aren't re-ordered.
		var url, title string

		if v.HayIndex % 2 == 0 {
			url = entries[v.HayIndex + 1]
			title = v.Key
		} else {
			title = entries[v.HayIndex - 1]
			url = v.Key
		}
		// Unique the objects.
		if _, ok := dedup[url]; ok {
			continue
		}
		dedup[url] = struct{}{}

		result.Items = append(result.Items, &Item{
			Title:        title,
			Arg:          url,
			Subtitle: url,
			Autocomplete: title,
			relevance:    v.Score,
			Valid:        true,
		})
	}
	// I think that they're sorted already.
	// sort.Sort(result.Items)

	if err := encoder.Encode(result); err != nil {
		log.Fatalf("can't write json %v", err)
	}
}

type Item struct {
	Uid          string `json:"uid,omitempty"`
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle,omitempty"`
	Arg          string `json:"arg"`
	Autocomplete string `json:"autocomplete"`
	Valid        bool `json:"valid"`
	relevance int
}

type Result struct {
	Items ItemCollection `json:"items"`
}

type ItemCollection []*Item

func (c ItemCollection) Len() int {
	return len(c)
}

func (c ItemCollection) Less(i, j int) bool {
	return c[i].relevance > c[j].relevance
}

func (c ItemCollection) Swap(i, j int) {
	tmp := c[i]
	c[i] = c[j]
	c[j] = tmp
}

var _ = ItemCollection(nil)
