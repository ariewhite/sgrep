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
	found bool = false
)

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

func parseMultiFiles(filesCount int, pattern string) {
	for i := 2; i <= filesCount+1; i++  {
		// get file path
		parseFile(os.Args[i], pattern, true)
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
	// lineNumberPtr := flag.Bool("n", false, "Print line number");
	
	flag.Parse()
	
	nFlags := flag.NArg()

	// check for pattern exist
	if nFlags < 1 {
		fmt.Fprintln(os.Stderr, "no pattern")
		os.Exit(2)
	}

	// get pattern
	pattern := flag.Arg(0)

	fmt.Printf("pattern: %s\nfile: %s\n", pattern, *filePtr)

	var input *os.File
	if *filePtr != "" {
		parseFile(*filePtr, pattern, false)
	} else if nFlags > 2 {
		parseMultiFiles(nFlags-1, pattern)
	} else {
		input = os.Stdin
	}
	

	// stdin (nofiles) mode

	lineNumber := 1
	found = false

	scanner := bufio.NewScanner(input)
	scanner.Buffer(make([]byte, 1024), 1024*1024)
	
	for scanner.Scan() {
		line := scanner.Text()

		if (strings.Contains(line, pattern)) {
			found = true
			highlighted := strings.ReplaceAll(line, pattern, Yellow + pattern + Reset)
			fmt.Printf("%d  %s\n", lineNumber, highlighted)
		}
		lineNumber++;
	}

	// exit code
	exitCode := 0
	if !found {
		exitCode = 1
	}

	os.Exit(exitCode)
}
