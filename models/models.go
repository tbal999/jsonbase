package models

import (
	"fmt"
	"math"
	"sort"
	"strconv"

	"github.com/cheggaaa/pb"
)

func squareDistance(features1, features2 []float64) float64 {
	var d float64
	for i := range features1 {
		d += (features1[i] - features2[i]) * (features1[i] - features2[i])
	}
	return math.Sqrt(d)
}

func calcaverage(items []float64) string {
	var average float64
	for index := range items {
		average += items[index]
	}
	output := average / float64(len(items))
	return strconv.FormatFloat(output, 'f', 4, 64)
}

func rangeMap(words []string) string {
	m := make(map[string]int)
	for _, word := range words {
		_, ok := m[word]
		if !ok {
			m[word] = 1
		} else {
			m[word]++
		}
	}
	min := 0
	var largest string
	for index := range m {
		if m[index] >= min {
			min = m[index]
			largest = index
		}
	}
	return largest
}

func KNN(training, testing [][]float64, trainingname, testingname []string, k int, train, regression bool) [][]string {
	output := [][]string{}
	Headers := []string{"INDEX", "PREDICTION"}
	output = append(output, Headers)
	bar := pb.StartNew(len(testing))
	for testindex := range testing {
		out := []struct {
			name   string
			number float64
		}{}
		bar.Increment()
		input := []string{}
		candidates := []string{}
		regdates := []float64{}
		if train == true {
			input = append(input, testingname[testindex])
		} else {
			input = append(input, strconv.Itoa(testindex+1))
		}
		var likely string
		for trainindex := range training {
			b := squareDistance(testing[testindex], training[trainindex])
			new := struct {
				name   string
				number float64
			}{trainingname[trainindex], b}
			out = append(out, new)
		}
		sort.SliceStable(out, func(i, j int) bool {
			return out[i].number < out[j].number
		})
		if k < len(out)-1 {
			for i := 0; i <= k; i++ {
				if regression == true {
					ax, _ := strconv.ParseFloat(out[i].name, 64)
					regdates = append(regdates, ax)
				} else {
					ix := out[i].name
					candidates = append(candidates, ix)
				}
			}
		} else {
			fmt.Printf("       The K number must be less than training sample size. Recommended starting number for this sample is %f\n", math.Sqrt(float64(len(training))))
			return nil
		}
		if regression == true {
			averagedata := calcaverage(regdates)
			input = append(input, averagedata)
			output = append(output, input)
		} else {
			likely = rangeMap(candidates)
			input = append(input, likely)
			output = append(output, input)
		}
	}
	bar.Finish()
	var counter float64
	if train == true && regression == false {
		for index := range output {
			if output[index][0] == output[index][1] {
				counter++
			}
		}
		fmt.Printf("%f success rate with K number of %d\n", (counter/float64(len(output)))*100, k)
	} else {
		return output
	}
	return nil
}
