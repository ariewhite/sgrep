package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
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
		if strings.Contains(scanner.Text(), pattern) {
			fmt.Printf("%d %s\n", line, scanner.Text());
		}
		
		line++;
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error");
	}
}

// .\sgrep some -f=test.txt
func main(){
	// check for flags
	// -f
	filePtr := flag.String("f", "", "parse file")
	flag.Parse()
	
	// check for pattern
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: no pattern");
		os.Exit(1)
	}

	pattern := args[0]

	fmt.Printf("file: %s\n", *filePtr);
	fmt.Printf("pattern: %s\n\n", pattern);

	parseFile(*filePtr, pattern);
	os.Exit(0)
}
