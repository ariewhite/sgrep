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

func parseFile(filePath string, pattern string){
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

		fmt.Printf("%d  %s%s%s%s%s\n", line, before, Yellow, found, Reset, other);
		line++;
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error");
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
	index := strings.Index(line, pattern)
	if index == -1 {
		return;
	}

	sample := Yellow + pattern + Reset

	res := strings.ReplaceAll(line, pattern, sample)

	fmt.Printf("%d  %s\n", lineNum, res);
}


// .\sgrep some -f=test.txt
func main(){
	// check for flags
	filePtr := flag.String("f", "", "parse file")
	flag.Parse()

	// check for pattern
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: no pattern");
		os.Exit(1)
	}
	pattern := args[0]

	fmt.Printf("%s\n", *filePtr);
	// fmt.Printf("pattern: %s\n\n", pattern);

	lineNumber := 1
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		handleBasicPattern(line, pattern, lineNumber)
		lineNumber++;
	}

	os.Exit(0)
}
