package main

import (
	"log"
	"os/exec"
	"os"

	"github.com/alecthomas/participle/v2"
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
	String *string  `@String`
}

type Tabs struct  {
	Props []*PropArray `"{" @@  "," @@ "}"`
}

type  PropArray struct {
	Values []*Value `"{" ( @@ ( "," @@ )* )? "}"`
}

func main() {
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

	litter.Dump(ast)

	// TODO(rjk): 
}
