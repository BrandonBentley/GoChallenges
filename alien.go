package main

import (
	"fmt"
	"errors"
	"math"
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
)
var wg sync.WaitGroup
func main() {
	/* TIME (Need to import "time" to use)
	start := time.Now()
	*/
	testCases, input := ParseTestCase()
	output := "Input:\n" + input + "\nOutput:\n"
	output += FindGreatestDistancesFast(testCases)
	fmt.Print(output[:len(output)-2])
	/* TIME
	elapsed := time.Now().Sub(start)
	fmt.Println("\nTime Elapsed: ")
	fmt.Println(elapsed.String())
	*/
}

type LifeForm struct {
	X float64
	Y float64
}

type BeamData struct {
	X float64
	Y float64
	Radius float64
}

type Pair struct {
	LF1, LF2 *LifeForm
	Distance, Slope float64
}

type TestCase struct {
	Data []LifeForm
	Beam BeamData
}

func (t* TestCase) Result() string{
	return fmt.Sprintf("%.2f\n%.2f %.2f", t.Beam.Radius, t.Beam.X, t.Beam.Y)
}

func (b* BeamData) Result() string{
	return fmt.Sprintf("%.2f\n%.2f %.2f", b.Radius, b.X, b.Y)
}

func (t* TestCase) DisplayDataSet() string {
	return fmt.Sprintf("%+v", t.Data)
}

func FindGreatestDistancesFast(testCases []TestCase) string {
	length := len(testCases)
	wg.Add(length)
	beams := make([]BeamData, length)
	for i, v := range testCases {
		go FindGreatestDistanceRoutine(&v, &beams[i])
	}
	results := ""
	wg.Wait()
	for _, v := range beams {
		results += v.Result() + "\n"
	}
	return results
}

func FindGreatestDistanceRoutine(testCase *TestCase, beam *BeamData) {
	*beam, _ = FindGreatestDistance(testCase.Data)
	(*testCase).Beam = *beam
	wg.Done()
}

func FindGreatestDistance(lifeForms []LifeForm) (BeamData, error) {
	UniqueValues := CheckForDuplicates(MinMax(lifeForms))
	return OptimalCoverage(UniqueValues)
}

func MinMax(lifeForms []LifeForm) (MinX LifeForm, MaxX LifeForm, MinY LifeForm, MaxY LifeForm) {
	xMin, xMax, yMin, yMax := singlePass(lifeForms)
	//fmt.Printf("\nxMin: %+v\txMax: %+v\nyMin: %+v\tyMax: %+v\n\n", xMin, xMax, yMin, yMax)
	return xMin, xMax, yMin, yMax
}

func singlePass(lifeForms []LifeForm) (MinX LifeForm, MaxX LifeForm, MinY LifeForm, MaxY LifeForm) {
	xMin, xMax, yMin, yMax := lifeForms[0], lifeForms[0], lifeForms[0], lifeForms[0]
	var wg2 sync.WaitGroup
	wg2.Add(4)
	go func() {
		QXMin(&lifeForms, &xMin)
		wg2.Done()
	}()
	go func() {
		QXMax(&lifeForms, &xMax)
		wg2.Done()
	}()
	go func() {
		QYMin(&lifeForms, &yMin)
		wg2.Done()
	}()
	go func() {
		QYMax(&lifeForms, &yMax)
		wg2.Done()
	}()
	wg2.Wait()
	return xMin, xMax, yMin, yMax
}

func QXMin(lf *[]LifeForm, value *LifeForm) {
	lifeForms := *lf
	xMin := lifeForms[0]
	for i := 1; i<len(lifeForms); i++ {
		if lifeForms[i].X < xMin.X {
			xMin = lifeForms[i]
		}
	}
	*value = xMin

}
func QXMax(lf *[]LifeForm, value *LifeForm) {
	lifeForms := *lf
	xMax := lifeForms[0]
	for i := 1; i<len(lifeForms); i++ {
		if lifeForms[i].X > xMax.X {
			xMax = lifeForms[i]
		}
	}
	*value = xMax
}
func QYMin(lf *[]LifeForm, value *LifeForm) {
	lifeForms := *lf
	yMin := lifeForms[0]
	for i := 1; i<len(lifeForms); i++ {
		if lifeForms[i].Y < yMin.Y {
			yMin = lifeForms[i]
		}
	}
	*value = yMin
}
func QYMax(lf *[]LifeForm, value *LifeForm) {
	lifeForms := *lf
	yMax := lifeForms[0]
	for i := 1; i<len(lifeForms); i++ {
		if lifeForms[i].Y > yMax.Y {
			yMax = lifeForms[i]
		}
	}
	*value = yMax
}

func CheckForDuplicates(lfs...LifeForm) []LifeForm {
	removeList := make([]int, 0)
	maxLength := len(lfs) - 2
	for i, _ := range lfs {
		for j := i+1; j < len(lfs); j++ {
			if lfs[i].X == lfs[j].X && lfs[i].Y == lfs[j].Y  {
				if len(removeList) < maxLength {
					removeList = append(removeList, j)
				} else {
					fmt.Println("Error Occured in checking for duplicates", removeList, lfs)
				}
			}
		}
	}
	if len(removeList) <= 0 {
		return lfs
	} else {
		return RemoveDuplicate(lfs, removeList)
	}
}

func RemoveDuplicate(lfs []LifeForm, indexes []int) []LifeForm {
	indexes = SortIntSlice(indexes)
	for _, index := range indexes {
		lfs[index] = lfs[len(lfs)-1]
		lfs = lfs[:len(lfs)-1]
	}
	return lfs
}

func OptimalCoverage(lfs []LifeForm) (BeamData, error) {
	pair := Pair{nil, nil, -1, 0}
	for i, _ := range lfs {
		for j := i+1; j < len(lfs); j++ {
			distance, slope := CalculateDistanceBetween(lfs[i], lfs[j])
			if distance > pair.Distance {
				pair.Distance = distance
				pair.Slope = slope
				pair.LF1 = &lfs[i]
				pair.LF2 = &lfs[j]
			}
		}
	}
	if pair.LF1 == nil || pair.LF2 == nil {
		return BeamData{}, errors.New("invalid LifeForms error")
	} else {
		return CalculateBeamData(pair), nil
	}
}

func CalculateDistanceBetween(lf1 LifeForm, lf2 LifeForm) (distance float64, slope float64) {
	X := lf1.X - lf2.X
	Y := lf1.Y - lf2.Y
	return math.Sqrt(math.Pow(X,2) + math.Pow(Y,2)), X/Y
}

func CalculateBeamData(pair Pair) BeamData {
	X, Y := CalculateCenter(pair)
	return BeamData{X,Y,pair.Distance/2}
}

func CalculateCenter(pair Pair) (X float64, Y float64) {
	lf := *pair.LF2
	if pair.LF1.X < pair.LF2.X {
		lf = *pair.LF1
	}
	m, y, x := pair.Slope, lf.Y, lf.X

	b := y - m*x
	//y = mx + b
	x = x + pair.Distance/2
	y = m*x + b
	return x, y

}

func ParseTestCase() ([]TestCase, string) {
	reader := bufio.NewReader(os.Stdin)
	textNumTests, _, _ := reader.ReadLine()
	inputString :=  string(textNumTests) + "\n"
	numTests, _ := strconv.ParseInt(string(textNumTests), 10, 64)
	testCases := make ([]TestCase, numTests)
	for i := 0; i < int(numTests); i++ {
		textNumRows, _, _ := reader.ReadLine()
		inputString += string(textNumRows) + "\n"
		numRows, _ := strconv.Atoi(string(textNumRows))
		lifeForms := make([]LifeForm, numRows)
		for j := 0; j < numRows; j++ {
			textValue, _, _ := reader.ReadLine()
			inputString += string(textValue) + "\n"
			value := strings.Split(string(textValue), " ")
			value1, _ := strconv.ParseFloat(value[0], 64)
			value2, _ := strconv.ParseFloat(value[1], 64)
			lifeForms[j] = LifeForm{value1, value2}
		}
		testCases[i].Data = lifeForms
	}
	return testCases, inputString
}

func SortIntSlice(nums []int) []int {
	for i := 0; i < len(nums); i++ {
		for j := i+1; j < len(nums); j++ {
			if nums[i] < nums[j] {
				temp := nums[i]
				nums[i] = nums[j]
				nums[j] = temp

			}
		}
	}
	return nums
}
