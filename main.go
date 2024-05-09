package main

import (
	"flag"
	"log"
	"os/exec"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/sanity-io/litter"
	"github.com/reinhrst/fzf-lib"
)

const scriptsource = `tell application "Safari"
	set urls to URL of tabs of windows
	set u to first item of urls
	set names to name of tabs of windows
	set n to first item of names
	set result to {u, n}
	-- return count of n
	return result
end tell
`

const command = "/usr/bin/osascript"

type Value struct {
	String *string  `@String`
}

type Tabs struct  {
	Props []*PropArray `"{" @@  "," @@ "}"`
}

type  PropArray struct {
	Values []*Value `"{" ( @@ ( "," @@ )* )? "}"`
}

func main() {
	flag.Parse()

	log.Println("goody")

	// Run the scriptsource command
	cmd := exec.Command(command, "-s", "s", "-e", scriptsource)
	output, err := cmd.Output()
	
	if err != nil {
		log.Println("had an error running the command", err)
		// TODO(rjk): Print stderr.
		log.Println(string(output))
		os.Exit(-1)
	}

	// It's not valid JSON. I need a parser. A simple parser.
	log.Println("output", string(output))

	parser, err := participle.Build[Tabs]()
	if err != nil {
		log.Println("parser creation failed", err)
		os.Exit(-2)
	}
	
	ast, err := parser.ParseBytes("", output)
	if err != nil {
		log.Println("parsing failed", err)
		os.Exit(-3)
	}

	// Example: how to dump the AST
	log.Println(*ast.Props[0].Values[0].String)

//	litter.Dump(ast)

	// TODO(rjk): Consider caching the results during the execution of the filter so that
	// re-executions don't have to run the AppleScript.

	if fa := flag.Args(); len(fa) < 1 {
		// Do I need to exit here?
		log.Println("Do I need to exit here?")
		os.Exit(-3)
	}

	// Prepare input for the fzf library
	entries := make([]string, 0)
	for i, _ := range(ast.Props[0].Values) {
		entries = append(entries, *ast.Props[1].Values[i].String, *ast.Props[0].Values[i].String)
	}

	// TODO(rjk): Explore the options that I want.
	options := fzf.DefaultOptions()
	myFzf := fzf.New(entries, options)
	
	// needz to fix this up.
	myFzf.Search(flag.Args()[0])

	// This API is (apparently) designed to permit running multiple searches and
	// getting results back.
	// The "right" way is probaby to have a slave process that's feeding me the content
	// to the slave.
	result := <- myFzf.GetResultChannel()


	// search the ast with it.
	// accumulate the matchines
	litter.Dump(result)


	
}

type MatchTab struct {
	Title string
	Url string
	Relevance int
}



func makeRegexp(arg string) (*regexp.Regexp, error) {
	// TODO(rjk): It's possible to do this more efficiently.
	ss := strings.Split(arg, "")
	res := ".*" + strings.Join(ss, ".*",) + ".*"
	return regexp.Compile(res)
}