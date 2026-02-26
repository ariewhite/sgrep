package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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

type Match struct{
	lineNumber int
	line string
	pattern string
	filePath string
}

type Options struct{
	showLineNum  bool
	showFileName bool
	onlyMatched  bool
	invertMatch  bool
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
		m.line = strings.ReplaceAll(m.line, m.pattern, Yellow + m.pattern + Reset)
		sb.WriteString(m.line)
	}

	return sb.String()
}


func colorize(original string, color string) string {
	return color + original + Reset;
}

//
func search(io *os.File, pattern string, filePath string, invert bool, match func(m Match)) error {
	scanner := bufio.NewScanner(io)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	lineNum := 1 

	for scanner.Scan() {
		line := scanner.Text()
		contain := strings.Contains(line, pattern)

		if invert {
			if !contain {
				match(Match{
					lineNumber: lineNum,
					line: line,
					pattern: pattern,
					filePath: filePath,
				})
			}
		} else {
			if contain {
				match(Match{
					lineNumber: lineNum,
					line: line,
					pattern: pattern,
					filePath: filePath,
				})
			}
		}

		lineNum++
	}
	return scanner.Err(); 
}

// TODO: remove fileNaming param
func parseFile(filePath string, pattern string, fileNaming bool){
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err);
	}
	defer file.Close();

	line := 1;
	scanner := bufio.NewScanner(file);

	for scanner.Scan() {
		index := strings.Index(scanner.Text(), pattern)
		if index == -1 {
			line++;
			continue
		}
		
		before := scanner.Text()[:index];
		found  := scanner.Text()[index : index+len(pattern)]
		other  := scanner.Text()[(index + len(pattern)):];
		if fileNaming {
			
			fmt.Printf("%s%s%s:%d%s  %s%s%s%s%s\n", Green, filePath, Red, line, Reset, before, Yellow, found, Reset, other);
		} else {
			fmt.Printf("%d  %s%s%s%s%s\n", line, before, Yellow, found, Reset, other);
		}

		line++;
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error");
	}
}

func parseMultiFiles(filesCount int, pattern string, printLN bool) {
	for i := 2; i <= filesCount+1; i++  {
		// get file path
		parseFile(os.Args[i], pattern, printLN)
	}
}

func findAllSubs(str, substr string) []int {
	var indices []int;
	for i := 0; ; {
		idx := strings.Index(str[i:], substr)
		if idx == -1{
			break;
		}
		indices = append(indices, i+idx);

		i += idx + len(substr);
	}

	return indices;
}

func handleBasicPattern(line string, pattern string, lineNum int){
	if strings.Contains(line, pattern) {
		return;
	}

	found = true

	sample := Yellow + pattern + Reset
	res := strings.ReplaceAll(line, pattern, sample)
	fmt.Printf("%d  %s\n", lineNum, res);
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

	// selecting mode
	if len(files) > 0 {
		for _, path := range files {
			// open file
			f, err := os.Open(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error: %s: %v\n", path, err)
				continue
			}

			if err := search(f, pattern, path, *invertOut, onMatch); err != nil {
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

		if err := search(f, pattern, *filePtr, *invertOut, onMatch); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)	
		}
	} else {                  // stdin
		if err := search(os.Stdin, pattern, *filePtr, *invertOut, onMatch); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
	}

	if !found {
		os.Exit(1)
	}

	os.Exit(0)
}
