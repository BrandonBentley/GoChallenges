package CodingPractice

import (
	"sort"
	"fmt"
	"errors"
	"math"
)

func main() {
	lifeForms := []LifeForm{LifeForm{8.0, 9.0}, LifeForm{4.0, 7.5}, LifeForm{1.0, 2.0}, LifeForm{5.1, 8.7}, LifeForm{9.0, 2.0}, LifeForm{4.5, 1.0}}
	data, err := FindGreatestDistance(lifeForms)
	if err != nil {
		fmt.Println("Error has Occured")
	}
	fmt.Println(data.Radius)
	fmt.Println(data.X, " ", data.Y)
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

func FindGreatestDistance(lifeForms []LifeForm) (BeamData, error) {
	UniqueValues := CheckForDuplicates(MinMax(&lifeForms))
	return OptimalCoverage(UniqueValues)
}

func MinMax(lifeForms *[]LifeForm) (MinX *LifeForm, MaxX *LifeForm, MinY *LifeForm, MaxY *LifeForm) {
	xMin, xMax := MinMaxX(lifeForms)
	yMin, yMax := MinMaxY(lifeForms)
	return xMin, xMax, yMin, yMax
}

func MinMaxX(lifeForms *[]LifeForm) (Min *LifeForm, Max *LifeForm) {
	sort.Slice(lifeForms, func(i, j int) bool {
		return (*lifeForms)[i].X < (*lifeForms)[j].X
	})
	return &(*lifeForms)[0], &(*lifeForms)[len(*lifeForms)-1]
}

func MinMaxY(lifeForms *[]LifeForm) (Min *LifeForm, Max *LifeForm) {
	sort.Slice(lifeForms, func(i, j int) bool {
		return (*lifeForms)[i].Y < (*lifeForms)[j].Y
	})
	return &(*lifeForms)[0], &(*lifeForms)[len(*lifeForms)-1]
}

func CheckForDuplicates(lfs...*LifeForm) []*LifeForm {
	removeList := make([]int, 0)
	maxLength := len(lfs) - 2
	for i, _ := range lfs {
		for j := i+1; i < len(lfs); j++ {
			if lfs[i] == lfs[j] {
				if len(removeList) < maxLength {
					removeList = append(removeList, j)
				} else {
					fmt.Println("Error Occured")
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

func RemoveDuplicate(lfs []*LifeForm, indexes []int) []*LifeForm {
	for _, index := range indexes {
		lfs[index] = lfs[len(lfs)-1]
		lfs[len(lfs)-1] = nil
		lfs = lfs[:len(lfs)-1]
	}
	return lfs
}

func OptimalCoverage(lfs []*LifeForm) (BeamData, error) {
	pair := Pair{nil, nil, -1, 0}
	for i, _ := range lfs {
		for j := i+1; j < len(lfs); j++ {
			distance, slope := CalculateDistanceBetween(*lfs[i], *lfs[j])
			if distance > pair.Distance {
				pair.Distance = distance
				pair.Slope = slope
				pair.LF1 = lfs[i]
				pair.LF2 = lfs[j]
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
	beam := BeamData{0,0,pair.Distance/2}

	beam.X, beam.Y = CalculateCenter(pair)
	return beam
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

