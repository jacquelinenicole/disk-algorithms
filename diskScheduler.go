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
	
	algorithm := getValue(s, "Use")
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
	fmt.Printf("Lower cylinder: %3d\n", lowerCyl)
	fmt.Printf("Upper cylinder: %3d\n", upperCyl)
	fmt.Println("Cylinder requests:")
	for _, req := range cylReqs {
		fmt.Printf("Cylinder %3d\n", req)
	}

    traversalCount := 0
	if (algorithm == "fcfs") {
    	traversalCount = fcfs(lowerCyl, upperCyl, initCyl, cylReqs)
    } else if (algorithm == "sstf") {
    	traversalCount = sstf(lowerCyl, upperCyl, initCyl, cylReqs)
    } 


    fmt.Printf("%s traversal count = %d", strings.ToUpper(algorithm), traversalCount)
}



/* fcfs */
func fcfs(lowerCyl int, upperCyl int, lastCyl int, cylReqs []int) int {
	traversalCount := 0

	for _, req := range cylReqs {
		if cylError(req, lowerCyl, upperCyl) {
			continue
		}

		fmt.Printf("Servicing %3d\n", req)
		traversalCount += calcTraversal(lastCyl, req, lowerCyl, upperCyl)
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
		currCylIndex := findShortestSeekIndex(lastCyl, lowerCyl, upperCyl, unserviced)
		fmt.Printf("Servicing %3d\n", unserviced[currCylIndex])

		traversalCount += calcTraversal(lastCyl, unserviced[currCylIndex], lowerCyl, upperCyl)
		lastCyl = unserviced[currCylIndex]
		unserviced = remove(unserviced, currCylIndex)
	}

	return traversalCount
}

func findShortestSeekIndex(lastCyl int, lowerCyl int, upperCyl int, unserviced []int) int {
	index := 0
	for i := 0 ; i < len(unserviced) ; i++ {
		if calcTraversal(lastCyl, unserviced[i], lowerCyl, upperCyl) < unserviced[index] {
			index = i
		}
	}

	return index
}

func remove(s []int, i int) []int {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}

func calcTraversal(lastCyl int, currCyl int, start int, max int) int {
	higher := 0
	lower := 0

	if lastCyl > currCyl {
		higher = lastCyl
		lower = currCyl
	} else {
		higher = currCyl
		lower = lastCyl
	}

	circular := lower - start + max - higher
	regular := higher - lower

	if circular < regular {
		return circular
	}

	return regular
}



/* scan */
func scan(lowerCyl int, upperCyl int, initCyl int, cylReqs []int, direction int) int {
	// bubble sort cylinder requests
	for i := 0 ; i < len(cylReqs) ; i++ {
		for j := 0 ; j < len(cylReqs) - i - 1 ; j++ {
			if cylReqs[j] > cylReqs[j+1] {
				cylReqs[j], cylReqs[j+1] = cylReqs[j+1], cylReqs[j]
			}
		}
	}

	direction = scanTraversal(initCyl, cylReqs, direction)
	direction = scanTraversal(initCyl, cylReqs, direction)

	// started moving up
	if direction > 0 {
		return upperCyl - (initCyl + lowerCyl) + upperCyl - lowerCyl
	} else {
		// started moving down
		return initCyl - lowerCyl + upperCyl - lowerCyl
	}
}

// TODO: error message for out of bounds cylinder req
func scanTraversal(initCyl int, cylReqs []int, direction int) int {
	if direction < 0 {
		for i := len(cylReqs) - 1 ; i >= 0 ; i-- {
			if cylReqs[i] < initCyl {
				fmt.Printf("Servicing %3d\n", cylReqs[i])
			}
		}

		direction = 1
	} else if direction > 0 {
		for i := 0 ; i < len(cylReqs) ; i++ {
			if cylReqs[i] > initCyl {
				fmt.Printf("Servicing %3d\n", cylReqs[i])
			}
		}

		direction = -1
	}

	return direction
}



/* c-scan */
func cscan(lowerCyl int, upperCyl int, initCyl int, cylReqs []int, direction int) int {
	// bubble sort cylinder requests
	for i := 0 ; i < len(cylReqs) ; i++ {
		for j := 0 ; j < len(cylReqs) - i - 1 ; j++ {
			if cylReqs[j] > cylReqs[j+1] {
				cylReqs[j], cylReqs[j+1] = cylReqs[j+1], cylReqs[j]
			}
		}
	}

	direction = cscanTraversal(initCyl, cylReqs, direction)
	cscanTraversal(initCyl, cylReqs, direction)

	return upperCyl - lowerCyl + upperCyl - (initCyl + lowerCyl) + initCyl - lowerCyl
}

// TODO: error message for out of bounds cylinder req
func cscanTraversal(initCyl int, cylReqs []int, direction int) int {
	if direction == 1 { // start moving up
		for i := 0 ; i < len(cylReqs) ; i++ {
			if cylReqs[i] > initCyl {
				fmt.Printf("Servicing %3d\n", cylReqs[i])
			}
		}
		direction = 2
	} else if  direction == -1 { // start moving down
		for i := len(cylReqs) - 1 ; i >= 0 ; i-- {
			if cylReqs[i] < initCyl {
				fmt.Printf("Servicing %3d\n", cylReqs[i])
			}
		}
		direction = -2
	} else if direction == 2 { // finish from the front
		for i := 0 ; i < len(cylReqs) ; i++ {
			if cylReqs[i] >= initCyl {
				break
			}
			fmt.Printf("Servicing %3d\n", cylReqs[i])
		}

	} else if direction == -2 { // finish from the back
		for i := len(cylReqs) - 1 ; i >= 0 ; i-- {
			if cylReqs[i] <= initCyl {
				break
			}
			fmt.Printf("Servicing %3d\n", cylReqs[i])
		}
	}


	return direction
}



/* look */
func look(lowerCyl int, upperCyl int, initCyl int, cylReqs []int, direction int) int {
	// bubble sort cylinder requests
	for i := 0 ; i < len(cylReqs) ; i++ {
		for j := 0 ; j < len(cylReqs) - i - 1 ; j++ {
			if cylReqs[j] > cylReqs[j+1] {
				cylReqs[j], cylReqs[j+1] = cylReqs[j+1], cylReqs[j]
			}
		}
	}

	direction = scanTraversal(initCyl, cylReqs, direction)
	direction = scanTraversal(initCyl, cylReqs, direction)

	// started moving up
	if direction > 0 {
		return cylReqs[len(cylReqs) - 1] - (initCyl + cylReqs[0]) + cylReqs[len(cylReqs) - 1] - cylReqs[0]
	} else {
		// started moving down
		return initCyl - cylReqs[0] + cylReqs[len(cylReqs) - 1] - cylReqs[0]
	}
}



/* c-look */
func clook(lowerCyl int, upperCyl int, initCyl int, cylReqs []int, direction int) int {
	// bubble sort cylinder requests
	for i := 0 ; i < len(cylReqs) ; i++ {
		for j := 0 ; j < len(cylReqs) - i - 1 ; j++ {
			if cylReqs[j] > cylReqs[j+1] {
				cylReqs[j], cylReqs[j+1] = cylReqs[j+1], cylReqs[j]
			}
		}
	}

	direction = cscanTraversal(initCyl, cylReqs, direction)
	cscanTraversal(initCyl, cylReqs, direction)

	return cylReqs[len(cylReqs) - 1] - cylReqs[0] + cylReqs[len(cylReqs) - 1] - (initCyl + cylReqs[0]) + initCyl - cylReqs[0]
}

func cylError(req int, lowerCyl int, upperCyl int) bool {
	if req > upperCyl || req < lowerCyl {
		fmt.Printf("Cylinder request %3d outside of bounds.\n", req)
		return true
	}

	return false
}

