/*
Jacqueline van der Meulen | 08/02/2018 | COP 4600

*/
package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"math"
)

func main() {
    if len(os.Args) > 3 {
    	fmt.Println("Invalid argument list. Correct usage: \ngo run diskScheduler.go [file name]")
    	os.Exit(-1)
    }

    fileName := os.Args[1]
   	algorithm, lowerCyl, upperCyl, initCyl, cylReqs := parse(fileName)

   	runAlgorithm(algorithm, lowerCyl, upperCyl, initCyl, cylReqs)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// gets all needed data from file
func parse(fileName string) (string, int, int, int, int[]) {
    file, err := os.Open(fileName)
    check(err)

    s := bufio.NewScanner(file)
	s.Split(bufio.ScanWords)
	
	algorithm := getValue(s, "Use")
	lowerCyl, _ := strconv.Atoi(getValue(s, "lowerCYL"))
	upperCyl, _ := strconv.Atoi(getValue(s, "upperCYL"))
	initCyl, _ := strconv.Atoi(getValue(s, "initCYL"))

	var cylReqs []int

	for int i := 0 ; ; i++ {
		if foundCyl(s) {
			s.Scan()
			nextCyl, _ := strconv.Atoi(s.Bytes())
			cylReqs =  append(cylReqs, nextCyl)
		} else {
			break
		}
	}

    file.Close()

	return algorithm, lowerCyl, upperCyl, initCyl, cylReqs
}

// finds keyword given then returns the word/value coming after it
func getValue(s *bufio.Scanner, word string) string {
	for string(s.Bytes()) != word) {
		s.Scan()
	}
	s.Scan()

	return string(s.Bytes())
}

// found "cylreq" --> true
// found "end" --> false
func foundCyl(s *bufio.Scanner) bool {
	for string(s.Bytes()) != "cylreq" && string(s.Bytes()) != "end") {
		s.Scan()
	}

	if string(s.Bytes()) == "end" {
		return false
	}

	return true
}

// prints data parsed from input file and calls algorithm requested
func runAlgorithm(algorithm String, lowerCyl int, upperCyl int, initCyl int, cylReqs int[]) {
	fmt.printf("Seek algorithm: %s\n", algorithm)
	fmt.printf("Lower cylinder: %3d\n", lowerCyl)
	fmt.printf("Upper cylinder: %3d\n", upperCyl)
	fmt.println("Cylinder requests:")
	for _, req := range cylReqs {
		fmt.printf("Cylinder %3d\n", req)
	}

	if (algorithm == "fcfs") {
    	traversalCount := fcfs(lowerCyl, upperCyl, initCyl, cylReqs)
    } else if (algorithm == "sstf") {
    	traversalCount := sstf(lowerCyl, upperCyl, initCyl, cylReqs)
    }

    fmt.printf("%s traversal count = %d", algorithm, traversalCount)
}


func fcfs(lowerCyl int, upperCyl int, lastCyl int, cylReqs int[]) int {
	traversalCount := 0

	for _, req := range cylReqs {
		if cylError(req, lowerCyl, upperCyl) {
			continue
		}

		fmt.printf("Servicing %3d\n", req)
		traversalCount += calcTraversal(lastCyl, req, lowerCyl, upperCyl)
		lastCyl = req
	}

	return traversalCount
}

func sstf(lowerCyl int, upperCyl int, lastCyl int, cylReqs int[]) int {
	traversalCount := 0

	unserviced := cylReqs

	// if out of bounds, print error message and remove
	for i := 0 ; i < len(unserviced) ; i++ {
		if cylError(req, lowerCyl, upperCyl) {
			remove(unserviced, i)
		}

	}

	for len(unserviced > 0) {
		currCylIndex := findShortestSeekTime(lastCyl, lowerCyl, upperCyl, unserviced)
		fmt.printf("Servicing %3d\n", unserviced[currCylIndex])

		traversalCount += calcTraversal(lastCyl, unserviced[currCylIndex], lowerCyl, upperCyl)
		lastCyl = unserviced[currCylIndex]
		remove(unserviced, currCylIndex)
	}

	return traversalCount
}

func findShortestSeekIndex(lastCyl int, lowerCyl int, upperCyl int, unserviced int[]) int {
	index := -1
	for (int i := 0 ; i < len(unserviced) ; i++)
	{
		if index == -1 || calcTraversal(lastCyl, unserviced[i], lowerCyl, upperCyl) < unserviced[index] {
			index = i
		}
	}

	return index
}

func remove(s []int, i int) []int {
    s[len(s)-1], s[i] = s[i], s[len(s)-1]
    return s[:len(s)-1]
}

func cylError(req int, lowerCyl int, upperCyl int) bool {
	if req > upperCyl || req < lowerCyl {
		fmt.printf("Cylinder request %3d outside of bounds.\n", req)
		return true
	}

	return false
}

func calcTraversal(lastCyl int, currCyl int, start int, max int) {
	if lastCyl > currCyl {
		higher := lastCyl
		lower := currCyl
	} else {
		higher := currCyl
		lower := lastCyl
	}

	return Math.Min(lower - start + max - higher, higher - lower)
}