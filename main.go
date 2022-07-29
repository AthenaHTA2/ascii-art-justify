file: main.go

***Note: The ascii-art code can not print " = double quotes***
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	Banner       string
	Alignment    string
	flagTerminal string
	ToConvert    string
	numLetters   int
)

func main() {
	flagTerminal = os.Args[3]
	n := len(os.Args)
	numArgs(n)
	checkFlag(flagTerminal)
	Banner = os.Args[2] + ".txt"
	data, _ := os.ReadFile(Banner) // take data from Banner onto data
	words := strings.Split(os.Args[1], "\\n")
	lines := strings.Split(string(data), "\n")

	for index := range words {
		if words[index] == "" { // break if there is no word
			break
		}

		letters := strings.Split(words[index], "") // slice the words into symbols
		var ascii []int
		for i := range letters {
			l := int([]rune(letters[i])[0]) // store the rune value of the letter in l
			ascii = append(ascii, l)        // append l value to ascii untill all letters are appended
			numLetters += 1
		}
		for j := 1; j < 9; j++ { // loop 8 times because a row is 8 lines long
			str := ""
			for k := range ascii { // loop through all runes in ascii
				line := (ascii[k] - 32) * 9 // this finds the row before the ascii symbol in the txt file
				str += lines[line+j]        // gives the string all the lines into the string
				// in terminal, the 'stty size' command returns current screen size, e.g. rows 25; columns 98; command stty -a shows all stats.
				// E.g.: fmt.Printf("%100s\n",str) this aligns to right
				// E.g.: fmt.Printf("%-100s\n",str) this aligns to left
			}

			align := os.Args[3][8:]
			const space = "\u0020"
			width := consoleWidth()

			switch align {
			case "left": //  My original align to left was: fmt.Printf("%-98s\n", str)
				fmt.Println(str)
			case "right": // My original hard coded 'align to right' was: fmt.Printf("%98s\n", str)
				spot := (width - utf8.RuneCountInString(str))
				fmt.Printf("%s%s\n", strings.Repeat(space, spot), str) // align to right
			case "center":
				middle := ((width - utf8.RuneCountInString(str)) / 2)
				if middle < 1 { // check if number of columns is smaller than string length
					fmt.Print(str)
				} else {
					fmt.Printf("%s%s\n", strings.Repeat(space, middle), str) // align to middle
				}
			case "justify":
				inBetweenFirst := ((width - (utf8.RuneCountInString(str))) / numLetters)
				inBetween := inBetweenFirst + (inBetweenFirst/numLetters)/2 // adding the second term to reduce space to the right
				str := ""
				for k := range ascii { // loop through all runes in ascii
					line := (ascii[k] - 32) * 9                                  // this finds the row before the ascii symbol in the txt file
					str = str + lines[line+j] + strings.Repeat(space, inBetween) // gives the string all the lines into the string, padded with 'inBetween' number of spaces to justify.
				}
				fmt.Printf("%s", str) // print across the length of the terminal
				fmt.Println()
			}

		}
	}
}

func numArgs(n int) {
	if n != 4 {
		fmt.Println()
		fmt.Print("EX: go run . something standard --align=right")
		fmt.Println()
		os.Exit(0)
	}
}

func checkFlag(flagTerminal string) {
	ToConvert = os.Args[1]
	style := os.Args[2]
	where := strings.LastIndex(os.Args[3][0:], "=")
	Alignment = os.Args[3][where+1:]
	if flagTerminal != "--align=center" && flagTerminal != "--align=left" && flagTerminal != "--align=right" && flagTerminal != "--align=justify" {
		fmt.Println()
		fmt.Printf("EX: go run . %v %v --align=%v\n", ToConvert, style, Alignment)
	}
}

func consoleWidth() int { // function provided by Steven Pearson
	cmd := exec.Command("stty", "size")
	defer exec.Command("clear")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	s := string(out)
	s = strings.TrimSpace(s)
	sArr := strings.Split(s, " ")

	width, err := strconv.Atoi(sArr[1])
	if err != nil {
		log.Fatal(err)
	}
	return width
}

