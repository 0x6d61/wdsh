/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/spf13/cobra"
)

var pattern string

func searchWordCmd(c *cobra.Command, args []string) {
	if len(args) == 0 {
		c.Help()
		os.Exit(1)
	}
	if pattern == "" {
		fmt.Println("Please specify the regular expression to search in the file")
		c.Help()
		os.Exit(1)
	}

	lineNumberFlag, err := c.PersistentFlags().GetBool("number")
	if err != nil {
		c.Help()
		os.Exit(1)
	}

	f, err := os.Open(args[0])
	if err != nil {
		fmt.Println("Unable to open the file.")
		os.Exit(1)
	}

	defer f.Close()

	scanner := bufio.NewScanner(f)

	lineNumber := 1
	for scanner.Scan() {
		line := scanner.Text()
		regexp := regexp.MustCompile(pattern)
		if regexp.MatchString(line) {
			findString := regexp.ReplaceAllString(line, "\x1b[31m$0\x1b[0m")
			if lineNumberFlag {
				fmt.Printf("%d,%s\n", lineNumber, findString)
			} else {
				fmt.Println(findString)
			}
		}
		lineNumber += 1
	}

	if err = scanner.Err(); err != nil {
		fmt.Println("Unable to open the file.")
		os.Exit(1)
	}

}

func main() {

	rootCmd := &cobra.Command{
		Use: "wdsh",
		Run: searchWordCmd,
	}

	rootCmd.PersistentFlags().StringVarP(&pattern, "pattern", "p", ``, "Regular expression to search from file")
	rootCmd.PersistentFlags().BoolP("number", "n", true, "Output the number of lines")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
