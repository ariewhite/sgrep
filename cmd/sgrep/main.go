package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	LBlue  = "\033[38;2;68;199;206m"
	Pink   = "\033[38;2;200;100;200m"
)

var (
	found = false
	matchCount int = 0
	fileMatchCounter map[string]int
)

// --- OPTIONS ---
type Options struct {
	IgnoreCase bool
	Regex bool
	ShowLineNum bool
	Invert bool
	OnlyMatch bool
	Count bool
	FilesNames bool
	FileNamesOnlyM bool
}

func ParseOptions() (Options, string, []string ) {
	lineNum := flag.Bool("n", false, "Print line number")
	invert := flag.Bool("v", false, "Invert output")
	onlyMatch := flag.Bool("m", false, "Print only matched out")
	regex := flag.Bool("e", false, "Regular Expressions")
	ignoreCase := flag.Bool("i", false, "Ignore case")
	countOpt    := flag.Bool("c", false, "Show matches count")
	fileNames := flag.Bool("f", false, "Print files names")
	onlyFileNameMatch := flag.Bool("l", false, "Print onli files names with matches")

	flag.Parse()

	if flag.NArg() < 1 {
		log.Fatal("pattern not found")
	}

	pattern := flag.Arg(0)
	files := flag.Args()[1:]

	return Options{
		IgnoreCase: *ignoreCase,
		Regex: *regex,
		ShowLineNum: *lineNum,
		Invert: *invert,
		OnlyMatch: *onlyMatch,
		Count: *countOpt,
		FilesNames: *fileNames,
		FileNamesOnlyM: *onlyFileNameMatch,
	}, pattern, files
}


// interface to matches
type Matcher interface {
	Match(line string) bool
	FindAll(line string) [][]int
}

// --- string ---
type StringMatcher struct {
	pattern string
}

func (s *StringMatcher) Match(line string) bool {
	return strings.Contains(line, s.pattern)
}

func (s* StringMatcher) FindAll(line string) [][]int {
	var res [][]int
	start := 0 

	for {
		i := strings.Index(line[start:], s.pattern)	
		if i == -1 {
			break
		}

		i += start;
		res = append(res, []int{i, i+len(s.pattern)})
		start = i + len(s.pattern)
	}

	return res;
}
// -------------- regex ------------------
type RegexMatcher struct {
	re      *regexp.Regexp
}

func (s *RegexMatcher) Match(line string) bool {
	return s.re.MatchString(line);
}

func (s *RegexMatcher) FindAll(line string) [][]int {
	return s.re.FindAllStringIndex(line, -1);
}

// -------------- matcher select -----------------
func newMatcher(p string, opts Options) (Matcher, error) {
	if opts.Regex {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		return &RegexMatcher{re}, nil
	}
	return &StringMatcher{p}, nil
}

func print(line string, num int, matcher Matcher, opts Options, file string) {

	if opts.FilesNames {
		fmt.Printf("%s%s", colorize(file, Pink), colorize(":", LBlue))
	}
	if opts.ShowLineNum {
		fmt.Printf("%s%s", colorize(strconv.Itoa(num), Green), colorize(":", LBlue))
	}

	tempLine := line
	if opts.IgnoreCase {
		tempLine = strings.ToLower(line)
	}

	indexs := matcher.FindAll(tempLine)

	if opts.OnlyMatch {
		for _, idx := range indexs {
			// sb.WriteString(fmt.Sprintf(line[idx[0]:idx[1]]))
			fmt.Printf("%s\n", colorize(line[idx[0]:idx[1]], Yellow))
		}
	} else {
		offset := 0
		ofLen := len(Yellow) + len(Reset)

		for _, idx := range indexs {
			idx[0] += offset;
			idx[1] += offset;

			line = line[:idx[0]] + Yellow + line[idx[0]:idx[1]] + Reset + line[idx[1]:]

			offset += ofLen
		}
	    fmt.Printf("%s\n", line)
	}
}

func colorize(original string, color string) string {
	temp := color + original + Reset;
	return temp
}

func search(io io.Reader, matcher Matcher, opts Options, file string) error {
	scanner := bufio.NewScanner(io)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNum := 1

	for scanner.Scan() {
		line := scanner.Text()
		lineTemp := line
		if opts.IgnoreCase {
			lineTemp = strings.ToLower(line)
		}

		matched := matcher.Match(lineTemp)

		if matched != opts.Invert {
			matchCount++;
			if !opts.FileNamesOnlyM {
				print(line, lineNum, matcher, opts, file)
			} else {
				fileMatchCounter[file]++
			}
		}
		lineNum++
	}
	return scanner.Err()
}


func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sgrep <options> <pattern>\nTry: sgrep -h for more help\n")
}

func printHelp() {
	fmt.Println(`Usage:
	sgrep [OPTIONS] <pattern>

	Description:
	Search for PATTERN in input (stdin by default or file).

	Options:
	-n [Print line number]

	-l [Print only file name (if match found)]

	-L [Print only file name with out matches]

	-f [Print file names]

	-v [Invert match (select non-matching lines)]

	-m [Print only matched parts of a line]

	-e [Treat pattern as regular expression]

	-i [Ignore case distinctions]

	-h [Show this help message and exit]

	Examples stdin:
	cat log1.txt | sgrep error
	cat log2.txt | sgrep -n -i error 

	Examples with file:
	sgrep -i -e "^\d[A-Z]" input.txt

	Multifiles: 
	sgrep something log1.txt log2.txt
	`)
}

// BASIC USAGE
// sgrep <pattern> -f=<path/to/file>
// sgrep <pattern> file1.txt file2.txt
// <stdin> | sgrep <pattern>
func main() {
	opts, pattern, files := ParseOptions()
	helpOpt 	:= flag.Bool("h", false, "Show help")

	flag.Parse()
	
	// check for help
	if *helpOpt {
		printHelp()
		os.Exit(0)
	}

	matcher, err := newMatcher(pattern, opts)
	if err != nil {
		log.Fatal(err)
	}

	if opts.Count {
		// count logic
	} else if opts.FileNamesOnlyM {
		if len(files) == 0 {
			log.Fatal("-l required more than 0 files")
			os.Exit(2)
		} 
		
		for _, f := range files {
			file, err := os.Open(f)
			if err != nil {
				log.Println(err)
				continue
			}

			search(file, matcher, opts, f)
			file.Close()
		}

		for f,s := range fileMatchCounter {
			fmt.Printf("%s: %d\n", f, s)
		}
	} else {

		if len(files) == 0 {
			search(os.Stdin, matcher, opts, "")
			os.Exit(0)
		}

		for _, f := range files {
			file, err := os.Open(f)
			if err != nil {
				log.Println(err)
				continue
			}

			search(file, matcher, opts, f)
			file.Close()
		}

	}

	os.Exit(0)
}
