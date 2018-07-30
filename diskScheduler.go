/*
Jacqueline van der Meulen | 08/02/2018 | COP 4600

*/
package main
import (
	"fmt"
	"bufio"
	"os"
	"strconv"
	"strings"
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
func parse(fileName string) (string, int, int, int, []int) {
    file, err := os.Open(fileName)
    check(err)

    s := bufio.NewScanner(file)
	s.Split(bufio.ScanWords)
	
	algorithm := getValue(s, "use")
	lowerCyl, _ := strconv.Atoi(getValue(s, "lowerCYL"))
	upperCyl, _ := strconv.Atoi(getValue(s, "upperCYL"))
	initCyl, _ := strconv.Atoi(getValue(s, "initCYL"))

	var cylReqs []int

	for i := 0 ; ; i++ {
		if foundCyl(s) {
			s.Scan()
			nextCyl, _ := strconv.Atoi(string(s.Bytes()))
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
	for string(s.Bytes()) != word {
		s.Scan()
	}
	s.Scan()

	return string(s.Bytes())
}

// found "cylreq" --> true
// found "end" --> false
func foundCyl(s *bufio.Scanner) bool {
	for string(s.Bytes()) != "cylreq" && string(s.Bytes()) != "end" {
		s.Scan()
	}

	if string(s.Bytes()) == "end" {
		return false
	}

	return true
}

// prints data parsed from input file and calls algorithm requested
func runAlgorithm(algorithm string, lowerCyl int, upperCyl int, initCyl int, cylReqs []int) {
	fmt.Printf("Seek algorithm: %s\n", strings.ToUpper(algorithm))
	fmt.Printf("\tLower cylinder: %5d\n", lowerCyl)
	fmt.Printf("\tUpper cylinder: %5d\n", upperCyl)
	fmt.Printf("\tInit cylinder: %5d\n", initCyl)
	fmt.Println("\tCylinder requests:")
	for _, req := range cylReqs {
		fmt.Printf("\t\tCylinder %5d\n", req)
	}

    traversalCount := 0
	if (algorithm == "fcfs") {
    	traversalCount = fcfs(lowerCyl, upperCyl, initCyl, cylReqs)
    } else if (algorithm == "sstf") {
    	traversalCount = sstf(lowerCyl, upperCyl, initCyl, cylReqs)
    } else if (algorithm == "scan" || algorithm == "c-scan" || algorithm == "look" || algorithm == "c-look") {
    	traversalCount = scanlook(algorithm, lowerCyl, upperCyl, initCyl, cylReqs)
    } else {
    	fmt.Println("Invalid algorithm requested")
    }

    fmt.Printf("%s traversal count = %d\n", strings.ToUpper(algorithm), traversalCount)
}


/* fcfs */
func fcfs(lowerCyl int, upperCyl int, lastCyl int, cylReqs []int) int {
	traversalCount := 0

	for _, req := range cylReqs {
		if cylError(req, lowerCyl, upperCyl) {
			continue
		}

		fmt.Printf("Servicing %5d\n", req)
		traversalCount += calcTraversal(lastCyl, req)
		lastCyl = req
	}

	return traversalCount
}


/* sstf */
func sstf(lowerCyl int, upperCyl int, lastCyl int, cylReqs []int) int {
	traversalCount := 0

	unserviced := cylReqs

	// if out of bounds, print error message and remove
	for i := 0 ; i < len(unserviced) ; i++ {
		if cylError(unserviced[i], lowerCyl, upperCyl) {
			unserviced = remove(unserviced, i)
		}
	}

	for len(unserviced) > 0 {
		currCylIndex := findShortestSeekIndex(lastCyl, lowerCyl, unserviced)
		fmt.Printf("Servicing %d\n", unserviced[currCylIndex])

		traversalCount += calcTraversal(lastCyl, unserviced[currCylIndex])
		lastCyl = unserviced[currCylIndex]
		unserviced = remove(unserviced, currCylIndex)
	}

	return traversalCount
}

func findShortestSeekIndex(lastCyl int, lowerCyl int, unserviced []int) int {
	index := 0
	for i := 0 ; i < len(unserviced) ; i++ {
		if calcTraversal(lastCyl, unserviced[i]) < calcTraversal(lastCyl, unserviced[index]) {
			index = i
		}
	}

	return index
}

func remove(s []int, i int) []int {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func calcTraversal(lastCyl int, currCyl int) int {
	higher := currCyl
	lower := lastCyl

	if lastCyl > currCyl {
		higher = lastCyl
		lower = currCyl
	}

	return higher - lower
}


/* scan, c-scan, look, c-look */
func scanlook(algorithm string, lowerCyl int, upperCyl int, initCyl int, cylReqs []int) int {
	// bubble sort cylinder requests
	for i := 0 ; i < len(cylReqs) ; i++ {
		for j := 0 ; j < len(cylReqs) - i - 1 ; j++ {
			if cylReqs[j] > cylReqs[j+1] {
				cylReqs[j], cylReqs[j+1] = cylReqs[j+1], cylReqs[j]
			}
		}
	}

	// servicing cylinders moving up
	for i := 0 ; i < len(cylReqs) ; i++ {
		if cylError(cylReqs[i], lowerCyl, upperCyl) {
			continue
		}

		if cylReqs[i] >= initCyl {
			fmt.Printf("Servicing %d\n", cylReqs[i])
		}
	}

	// servicing remaining cylinders
	lastCyl := 0
	onePass := true
	if algorithm[0:1] == "c" {
		lastCyl, onePass = frontTraversal(initCyl, lowerCyl, upperCyl, onePass, cylReqs)
	} else {
		lastCyl, onePass = reverseTraversal(initCyl, lowerCyl, upperCyl, onePass, cylReqs)
	}

	// calculating traversal count
	if onePass {
		return cylReqs[len(cylReqs) - 1] - initCyl
	} else if algorithm == "scan" {
		return upperCyl - initCyl + upperCyl - lowerCyl - (lastCyl - lowerCyl)
	} else if algorithm == "c-scan" {
		return upperCyl - initCyl + upperCyl - lowerCyl + lastCyl - lowerCyl
	} else if algorithm == "look" {
		return cylReqs[len(cylReqs) - 1] - initCyl + cylReqs[len(cylReqs) - 1] - lowerCyl - (lastCyl - lowerCyl)
	} else { // c-look
		return cylReqs[len(cylReqs) - 1] - initCyl + cylReqs[len(cylReqs) - 1] - cylReqs[0] + lastCyl - cylReqs[0]
	}
}

// going back down
func reverseTraversal(initCyl int, lowerCyl int, upperCyl int, onePass bool, cylReqs []int) (int, bool) {
	lastCyl := initCyl

	for i := len(cylReqs) - 1 ; i >= 0 ; i-- {
		if cylReqs[i] < initCyl && cylReqs[i] >= lowerCyl {
			fmt.Printf("Servicing %d\n", cylReqs[i])
			lastCyl = cylReqs[i]	
			onePass = false		
		}
	}	

	return lastCyl, onePass
}

// servicing remaining cylinders from front
func frontTraversal(initCyl int, lowerCyl int, upperCyl int, onePass bool, cylReqs []int) (int, bool) {
	lastCyl := initCyl
	i := 0

	for cylReqs[i] < initCyl && cylReqs[i] > lowerCyl {
		fmt.Printf("Servicing %d\n", cylReqs[i])
		lastCyl = cylReqs[i]
		i++
		onePass = false
	}

	return lastCyl, onePass
}

// checks bounds, prints error message
func cylError(req int, lowerCyl int, upperCyl int) bool {
	if req > upperCyl || req < lowerCyl {
		fmt.Printf("Cylinder request %5d outside of bounds.\n", req)
		return true
	}

	return false
}