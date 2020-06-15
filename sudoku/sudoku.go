package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (s Sudoku) print() {
	fmt.Printf("%v\n", strings.Repeat("-", 29))
	for i := 0; i < 81; i++ {
		if s[i] == 0 {
			fmt.Print(" _ ")
		} else {
			fmt.Printf(" %v ", s[i])
		}
		if (i+1)%9 == 0 {
			fmt.Printf("\n")
		} else if (i+1)%3 == 0 {
			fmt.Print("|")
		}
		if (i+1)%27 == 0 {
			fmt.Printf("%v\n", strings.Repeat("-", 29))
		}
	}
}

func (s *Sudoku) readFromString(sudokuString string) error {
	splitted := strings.Fields(sudokuString)
	for i, cell := range splitted {
		val, err := strconv.Atoi(cell)
		if err != nil {
			return err
		}
		s[i] = val
	}
	return nil
}

// Sudoku is a struct to hold sudoku data
type Sudoku [81]int

var boxStartIndex = map[int]int{0: 0, 1: 3, 2: 6, 3: 27, 4: 30, 5: 33, 6: 54, 7: 57, 8: 60}
var numSet = map[int]struct{}{1: {}, 2: {}, 3: {}, 4: {}, 5: {}, 6: {}, 7: {}, 8: {}, 9: {}}

// IsResolved checks if s is correctly resolved
func (s Sudoku) IsResolved() bool {
	// for i := 0; i < 81; i++ {
	// 	if s[i] == 0 {
	// 		return false
	// 	}
	// }
	for y := 0; y < 9; y++ {
		if !(s.isValid(y, s.getRowMap) && s.isValid(y, s.getColMap) && s.isValid(y, s.getBoxMap)) {
			return false
		}
	}
	return true
}

func (s Sudoku) isValid(no int, mapFn func(int) map[int]struct{}) bool {
	dataMap := mapFn(no)
	if len(intersection(numSet, dataMap)) != 9 {
		return false
	}
	return true
}

func getRowNoFromIndex(index int) int {
	return index / 9
}

func getColNoFromIndex(index int) int {
	return index % 9
}

func getBoxNoFromIndex(index int) int {
	xNo := getColNoFromIndex(index) / 3
	yNo := getRowNoFromIndex(index) / 3
	return xNo + 3*yNo
}

func intersection(firstMap map[int]struct{}, secondMap map[int]struct{}) map[int]struct{} {
	result := map[int]struct{}{}
	for key := range firstMap {
		if _, ok := secondMap[key]; ok {
			result[key] = struct{}{}
		}
	}
	return result
}

func subtraction(map1, map2 map[int]struct{}) map[int]struct{} {
	result := map[int]struct{}{}
	for key := range map1 {
		if _, ok := map2[key]; !ok {
			result[key] = struct{}{}
		}
	}
	return result
}

func (s Sudoku) getRow(rowNo int) [9]int {
	var row [9]int
	copy(row[:], s[rowNo*9:(rowNo+1)*9])
	return row
}

func (s Sudoku) getCol(colNo int) [9]int {
	var col [9]int
	for i := 0; i < 9; i++ {
		col[i] = s[colNo+9*i]
	}
	return col
}

func (s Sudoku) getBox(boxNo int) [9]int {
	var box [9]int
	init := ((boxNo / 3) * 27) + ((boxNo % 3) * 3)
	for i := 0; i < 3; i++ {
		box[3*i] = s[init+i*9]
		box[3*i+1] = s[init+i*9+1]
		box[3*i+2] = s[init+i*9+2]
	}
	return box
}

func (s Sudoku) getMap(no int, getDataFn func(int) [9]int) map[int]struct{} {
	resultMap := map[int]struct{}{}
	resultValues := getDataFn(no)
	for _, val := range resultValues {
		if val != 0 {
			resultMap[val] = struct{}{}
		}
	}
	return resultMap
}

func (s Sudoku) getRowMap(rowNo int) map[int]struct{} {
	return s.getMap(rowNo, s.getRow)
}

func (s Sudoku) getColMap(colNo int) map[int]struct{} {
	return s.getMap(colNo, s.getCol)
}

func (s Sudoku) getBoxMap(boxNo int) map[int]struct{} {
	return s.getMap(boxNo, s.getBox)
}

func (s Sudoku) getLackingNums(no int, rowMapFn func(int) map[int]struct{}) map[int]struct{} {
	return subtraction(numSet, rowMapFn(no))
}
func (s Sudoku) getPossibleNumsForIndex(index int) []int {
	lackingRow := s.getLackingNums(getRowNoFromIndex(index), s.getRowMap)
	lackingCol := s.getLackingNums(getColNoFromIndex(index), s.getColMap)
	lackingBox := s.getLackingNums(getBoxNoFromIndex(index), s.getBoxMap)
	result := intersection(lackingRow, lackingCol)
	result = intersection(result, lackingBox)
	possibleNumsList := make([]int, len(result))
	i := 0
	for key := range result {
		possibleNumsList[i] = key
		i++
	}
	return possibleNumsList
}

// Solve solves sudoku
func (s *Sudoku) humanLikeSolve() {
	for !s.IsResolved() {
		iterChange := 0
		for i, val := range s {
			if val == 0 {
				result := s.getPossibleNumsForIndex(i)

				if len(result) == 1 {
					s[i] = result[0]
					iterChange++
				}
			}
		}
		if iterChange == 0 {
			return
		}
	}
}

func (s *Sudoku) backtrackingSolve() bool {
	if s.IsResolved() {
		return true
	}
	for i := 0; i < 81; i++ {
		if s[i] == 0 {
			possibleNums := s.getPossibleNumsForIndex(i)
			for _, num := range possibleNums {
				s[i] = num
				if s.backtrackingSolve() {
					return true
				}
				s[i] = 0
			}
			return false
		}
	}
	return false
}

func main() {
	sudoku := `0 4 0 0 0 0 1 7 9 
0 0 2 0 0 8 0 5 4 
0 0 6 0 0 5 0 0 8 
0 8 0 0 7 0 9 1 0 
0 5 0 0 9 0 0 3 0 
0 1 9 0 6 0 0 4 0 
3 0 0 4 0 0 7 0 0 
5 7 0 1 0 0 2 0 0 
9 2 8 0 0 0 0 6 0`
	s := Sudoku{}
	s.readFromString(sudoku)
	s.print()
	start := time.Now()
	fmt.Println(s.backtrackingSolve())
	duration := time.Since(start)
	s.print()
	fmt.Println(duration)

}
