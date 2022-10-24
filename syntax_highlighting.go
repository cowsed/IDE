package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

// literally just a nano syntax highlighter
var definition_paths = []string{"Highlighters/go.nanorc"}

// file ending to Highlighter
var definitions []Highlighter

type HighlightedExpression struct {
	reg    *regexp.Regexp
	fg_col string
	bg_col string
}
type Highlighter struct {
	name        string
	file_ending *regexp.Regexp
	comment     string
	//string regex to string color
	expressions []HighlightedExpression
}

func ParseHighlighter(source string) (Highlighter, error) {
	hl := Highlighter{
		file_ending: &regexp.Regexp{},
		comment:     "",
		expressions: []HighlightedExpression{},
	}
	for _, line := range strings.Split(source, "\n") {
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "syntax":
			lang := parts[1]
			regex_string := parts[2][1 : len(parts[2])-1]
			regex, err := regexp.Compile(regex_string)
			if err != nil {
				//we don't know what files this applies to, it's basically useless
				log.Println(err)
				return hl, fmt.Errorf("error parsing syntax definition %v", err)
			}
			hl.name = lang
			hl.file_ending = regex
		case "comment":
			comment_with_quotes := parts[1]
			comment_wout_quotes := comment_with_quotes[1 : len(comment_with_quotes)-1]
			hl.comment = comment_wout_quotes
			fmt.Println("comment is ", parts[1])
		case "color":
			colordef := parts[1]
			color_parts := strings.Split(colordef, ",")

			var fg_col, bg_col string
			if len(color_parts) == 1 {
				fg_col = color_parts[0]
			} else {
				fg_col = color_parts[0]
				bg_col = color_parts[1]
			}

			//regex for the thing to be colored
			regex_index := len(parts) - 1
			regex_with_quotes := parts[regex_index]
			regex_wout_quotes := regex_with_quotes[1 : len(regex_with_quotes)-1]
			regex, err := regexp.Compile(regex_wout_quotes)
			if err != nil {
				log.Println("err parsing regex", err)
			}
			he := HighlightedExpression{
				reg:    regex,
				fg_col: fg_col,
				bg_col: bg_col,
			}
			hl.expressions = append(hl.expressions, he)
		}
	}
	return hl, nil
}

func ParseSyntaxHighlightingDefinitions() {
	for _, filename := range definition_paths {
		f, err := os.Open(filename)
		if err != nil {
			log.Println(err)
			continue
		}
		bs, err := io.ReadAll(f)
		if err != nil {
			log.Println(err)
			continue
		}
		s := string(bs)
		hl, err := ParseHighlighter(s)
		if err != nil {
			log.Printf("error parsing syntax highlighting file %s: %v\n", filename, err)
		}
		fmt.Println(hl)
		definitions = append(definitions, hl)
	}
}
