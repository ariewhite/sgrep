package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
	"regexp"
)


const(
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

var(
	found = false
)

// interface to matches
type Matcher interface {
	findIndex(line string) (stary, end int)
	Match(line string) bool
}

type Match struct{
	lineNumber int
	line string
	pattern string
	filePath string
	matchStart int
	matchEnd int
}

type Options struct{
	showLineNum  bool
	showFileName bool
	onlyMatched  bool
	invertMatch  bool
}

// -------------- string -------------------
type StringMatcher struct {
	pattern string
	ignoreCase bool
}

func (s *StringMatcher) Match(line string) bool {
	if s.ignoreCase {
		return strings.Contains(strings.ToLower(line), strings.ToLower(s.pattern))
	}
	return strings.Contains(line, s.pattern)
}

func (s *StringMatcher) findIndex(line string) (int, int) {
	index := strings.Index(line, s.pattern)
	if index == -1 {
		return -1, -1
	}
	return index, index + len(s.pattern)
}

// -------------- regex ------------------
type RegexMatcher struct {
	pattern string 
	re *regexp.Regexp
}

func (s *RegexMatcher) Match(line string) bool{
	return s.re.MatchString(line)
}

func (s *RegexMatcher) findIndex(line string) (int, int) {
	l := s.re.FindStringIndex(line)
	if l == nil {
		return -1, -1
	} 
	return l[0], l[1]
}


// -------------- matcher select -----------------
func newMatcher(pattern string, reg bool, ignoreCase bool) (Matcher, error) {
	if reg {
		if ignoreCase {
			pattern = "(?i)" + pattern;
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("Invalid regex %w", err)
		}
		return &RegexMatcher{re: re}, nil
	}
	return &StringMatcher{pattern: pattern, ignoreCase: ignoreCase}, nil
}

func (o *Options) print(m Match) string {
	var sb strings.Builder
	if o.showFileName {
		sb.WriteString(colorize(m.filePath, Green) + colorize(":", Blue))
	}
	if o.showLineNum {
		sb.WriteString(Green + fmt.Sprintf("%d", m.lineNumber) + Reset + ":");
	}
	if o.onlyMatched {
		sb.WriteString(colorize(m.pattern, Yellow))
	} else {
		sb.WriteString(highlight(m, Yellow))
	}

	return sb.String()
}

func highlight(m Match, color string) string {
	if m.matchStart == -1{
		return m.line
	}
	return m.line[:m.matchStart] + color + m.line[m.matchStart:m.matchEnd] + Reset + m.line[:m.matchEnd];
}

func colorize(original string, color string) string {
	return color + original + Reset;
}

//main searcher
func search(io *os.File, matcher Matcher, filePath string, invert bool, match func(m Match)) error {
	scanner := bufio.NewScanner(io)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNum := 1 

	for scanner.Scan() {
		line := scanner.Text()
		matched := matcher.Match(line)
		matchStart, matchEnd := matcher.findIndex(line)

		if matched != invert {
			match(Match{
				lineNumber: lineNum,
				line: line,
				filePath: filePath,
				matchStart: matchStart,
				matchEnd: matchEnd,
			})
		}
		lineNum++
	}
	return scanner.Err(); 
}

// BASIC USAGE
// sgrep <pattern> -f=<path/to/file>
// sgrep <pattern> file1.txt file2.txt
// <stdin> | sgrep <pattern>
func main(){
 	// check for single file mode
	filePtr := flag.String("f", "", "Parse file instead stdin");
	lineNumberOut 	:= flag.Bool("n", false, "Print line number");
	fileNameOnly  	:= flag.Bool("l", false, "Print only file name");
	invertOut       := flag.Bool("v", false, "Invert output");
	onlyMatchedOut  := flag.Bool("m", false, "Print only matched out");
	regularExp 		:= flag.Bool("e", false, "Regular Expressions");
	ignoreCase 	    := flag.Bool("i", false, "Ignore case");

	
	flag.Parse()

	// check for pattern exist
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "no pattern")
		os.Exit(2)
	}
	// get pattern
	pattern := flag.Arg(0)
	files := flag.Args()[1:]

	fmt.Printf("pattern: %s\nfile: %s\n", pattern, *filePtr)

	options := &Options{
		showLineNum: *lineNumberOut,
		showFileName: flag.NArg() > 2,	
		onlyMatched: *onlyMatchedOut,
		invertMatch: *invertOut,
	}

	if *fileNameOnly {
		fmt.Fprintf(os.Stderr, "add feature in future\n")
		os.Exit(0)
	}

	onMatch := func(x Match) {
		found = true
		fmt.Printf(options.print(x) + "\n")
	}

	matcher, err := newMatcher(pattern, *regularExp, *ignoreCase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(2)
	}

	// selecting mode
	if len(files) > 0 {
		for _, path := range files {
			// open file
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s: %v\n", path, err)
				continue
			}

			if err := search(f, matcher, path, *invertOut, onMatch); err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s: %v\n", path, err)
			}
			f.Close()
		}
	} else if *filePtr != "" { // -f singe file mode
		f, err := os.Open(*filePtr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s: %v\n", err, *filePtr)
			os.Exit(2)
		} 
		defer f.Close()

		if err := search(f, matcher, *filePtr, *invertOut, onMatch); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)	
		}
	} else {                  // stdin
		if err := search(os.Stdin, matcher, *filePtr, *invertOut, onMatch); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	if !found {
		os.Exit(1)
	}

	os.Exit(0)
}
