package main

import "fmt"

func BracketCombinations(num int) int {

	// Zero brackets gives us zero combinations
	if num == 0 {
		return 0
	}

	combinationsCount := 0

	// Functoin that will walk all the pathes with brackets
	var walker func(int, int)

	walker = func(openBrackets, clozedBrackets int) {

		// if we reached the end of path
		if openBrackets == num && clozedBrackets == num {
			combinationsCount++ // scope captured variable
			return
		}

		if clozedBrackets > openBrackets {
			return // wrong path (cannot close bracket before openning)
		}

		if clozedBrackets > num {
			return // wrong path (cannot close brackets more than total `num` supply)
		}

		if openBrackets > num {
			return // wrong path (cannot open brackets more than total `num` supply)
		}

		// Path to open bracket
		walker(openBrackets+1, clozedBrackets)

		// Path to close bracket
		walker(openBrackets, clozedBrackets+1)
	}

	walker(0, 0)

	return combinationsCount
}

func main() {

	// do not modify below here, readline is our function
	// that properly reads in the input for you
	fmt.Println(BracketCombinations(3))

}
