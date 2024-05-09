package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

// tpf needs to be the path to the file containing the contexts (a
// complex JSON format. args is the argument to the command.
func genAlfredResult(entries []MatchTab) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "	")
	var result Result

	// Create a result based on the hash
	for _, v := range entries {
		// Exclude the Uid field to make sure that the items aren't re-ordered.
		result.Items = append(result.Items, &Item{
			Title:        v.Title,
			Arg:          v.Url,
			Subtitle: v.Url,
			// TODO(rjk): Restructure this to do somethign sensible.
			Autocomplete: v.Title,
			relevance:    v.Relevance,
			Valid:        true,
		})
	}
	sort.Sort(result.Items)

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
	relevance    int
	Valid        bool `json:"valid"`
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
