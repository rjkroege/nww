package main

import (
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/reinhrst/fzf-lib"
	"github.com/sanity-io/litter"
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
	String *string `@String`
}

type Tabs struct {
	Props []*PropArray `"{" @@  "," @@ "}"`
}

type PropArray struct {
	Values []*Value `"{" ( @@ ( "," @@ )* )? "}"`
}

var debugging = flag.Bool("debug", false, "Set for extra logging")

func main() {
	flag.Parse()

	if !*debugging {
		log.SetOutput(io.Discard)
	}

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

	if *debugging {
		litter.Dump(ast)
	}

	// TODO(rjk): Consider caching the results during the execution of the filter so that
	// re-executions don't have to run the AppleScript.

	if fa := flag.Args(); len(fa) < 1 {
		// Do I need to exit here?
		log.Println("Do I need to exit here?")
		os.Exit(-3)
	}

	// Prepare input for the fzf library
	entries := make([]string, 0)
	for i := range ast.Props[0].Values {
		// Even entries are the title.
		// Odd entries are the URLs.
		entries = append(entries,
			strings.Trim(*ast.Props[1].Values[i].String, "\""),
			strings.Trim(*ast.Props[0].Values[i].String, "\""),
		)
	}

	// TODO(rjk): Explore the options that I want.
	options := fzf.DefaultOptions()
	myFzf := fzf.New(entries, options)

	// needz to fix this up.
	myFzf.Search(flag.Args()[0])

	// This API is (apparently) designed to permit running multiple searches
	// and getting results back. The "right" way is probaby to have a slave
	// process that's feeding me the content to the slave. so... how should
	// that work? Should I do that?
	result := <-myFzf.GetResultChannel()

	// search the ast with it.
	// accumulate the matchiness
	if *debugging {
		litter.Dump(result)
	}

	genAlfredResult(entries, result.Matches)

}

type MatchTab struct {
	Title     string
	Url       string
	Relevance int
}

func makeRegexp(arg string) (*regexp.Regexp, error) {
	// TODO(rjk): It's possible to do this more efficiently.
	ss := strings.Split(arg, "")
	res := ".*" + strings.Join(ss, ".*") + ".*"
	return regexp.Compile(res)
}
