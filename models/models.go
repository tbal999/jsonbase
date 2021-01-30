package models

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"time"
	"math/rand"

	"github.com/cheggaaa/pb"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

//BRUTE FORCE KNN ALGORITHM

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


//3 LAYER NEURAL NETWORK MODEL - initial code courtesy of ---> github.com/sausheong/gonn
//I have added multiple layers of functionality to make it avaialble to all sorts of custom datasets.
type Network struct {
	inputs        int
	hiddens       int
	outputs       int
	hiddenWeights *mat.Dense
	outputWeights *mat.Dense
	learningRate  float64
	targetmap map[int]string
}

func CreateNN(input, hidden, output int, rate float64) (net Network) {
	net = Network{
		inputs:       input,
		hiddens:      hidden,
		outputs:      output,
		learningRate: rate,
	}
	net.hiddenWeights = mat.NewDense(net.hiddens, net.inputs, randomArray(net.inputs*net.hiddens, float64(net.inputs)))
	net.outputWeights = mat.NewDense(net.outputs, net.hiddens, randomArray(net.hiddens*net.outputs, float64(net.hiddens)))
	return
}

func (net *Network) train(inputData []float64, targetData []float64) {
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.hiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.outputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)

	targets := mat.NewDense(len(targetData), 1, targetData)
	outputErrors := subtract(targets, finalOutputs)
	hiddenErrors := dot(net.outputWeights.T(), outputErrors)

	net.outputWeights = add(net.outputWeights,
		scale(net.learningRate,
			dot(multiply(outputErrors, sigmoidPrime(finalOutputs)),
				hiddenOutputs.T()))).(*mat.Dense)

	net.hiddenWeights = add(net.hiddenWeights,
		scale(net.learningRate,
			dot(multiply(hiddenErrors, sigmoidPrime(hiddenOutputs)),
				inputs.T()))).(*mat.Dense)
}

func (net Network) predict(inputData []float64) mat.Matrix {
	inputs := mat.NewDense(len(inputData), 1, inputData)
	hiddenInputs := dot(net.hiddenWeights, inputs)
	hiddenOutputs := apply(sigmoid, hiddenInputs)
	finalInputs := dot(net.outputWeights, hiddenOutputs)
	finalOutputs := apply(sigmoid, finalInputs)
	return finalOutputs
}

func sigmoid(r, c int, z float64) float64 {
	return 1.0 / (1 + math.Exp(-1*z))
}

func sigmoidPrime(m mat.Matrix) mat.Matrix {
	rows, _ := m.Dims()
	o := make([]float64, rows)
	for i := range o {
		o[i] = 1
	}
	ones := mat.NewDense(rows, 1, o)
	return multiply(m, subtract(ones, m))
}

func matrixPrint(X mat.Matrix) {
	fa := mat.Formatted(X, mat.Prefix(""), mat.Squeeze())
	fmt.Printf("%v\n", fa)
}

func dot(m, n mat.Matrix) mat.Matrix {
	r, _ := m.Dims()
	_, c := n.Dims()
	o := mat.NewDense(r, c, nil)
	o.Product(m, n)
	return o
}

func apply(fn func(i, j int, v float64) float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Apply(fn, m)
	return o
}

func scale(s float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Scale(s, m)
	return o
}

func multiply(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.MulElem(m, n)
	return o
}

func add(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Add(m, n)
	return o
}

func addScalar(i float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	a := make([]float64, r*c)
	for x := 0; x < r*c; x++ {
		a[x] = i
	}
	n := mat.NewDense(r, c, a)
	return add(m, n)
}

func subtract(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Sub(m, n)
	return o
}

func randomArray(size int, v float64) (data []float64) {
	dist := distuv.Uniform{
		Min: -1 / math.Sqrt(v),
		Max: 1 / math.Sqrt(v),
	}

	data = make([]float64, size)
	for i := 0; i < size; i++ {
		data[i] = dist.Rand()
	}
	return
}

//Not used in original NN - might come in handy
func addBiasNodeTo(m mat.Matrix, b float64) mat.Matrix {
	r, _ := m.Dims()
	a := mat.NewDense(r+1, 1, nil)
	a.Set(0, 0, b)
	for i := 0; i < r; i++ {
		a.Set(i+1, 0, m.At(i, 0))
	}
	return a
}

func (net *Network) targets(names []string) map[string]int {
	trainmap := make(map[string]int)
	net.targetmap = make(map[int]string)
	iter := 0
	for index := range names {
		_, ok := trainmap[names[index]]
		if ok {
			continue
		} else {
			trainmap[names[index]] = iter
			net.targetmap[iter] = names[index]
			iter++
		}
	}
	return trainmap
}

func (net *Network) Train(training [][]float64, trainingname []string, epochcount int) {
	tmap := net.targets(trainingname)
	rand.Seed(time.Now().UTC().UnixNano())
	t1 := time.Now()
	bar := pb.StartNew(epochcount)
	for epochs := 0; epochs < epochcount; epochs++ {
		bar.Increment()
			for index := range training {
				traindata := training[index]
				targets := make([]float64, net.outputs)
				inputss := make([]float64, net.inputs)
				for i := range targets {
						targets[i] = 0.001
				}
				for i := range inputss {
					inputss[i] = (traindata[i] / 255.0 * 0.999) + 0.001
				}
				targets[tmap[trainingname[index]]] = 0.999
				net.train(inputss, targets)
			}
		}
	bar.Finish()
	elapsed := time.Since(t1)
	fmt.Printf("\nTime taken to train: %s\n", elapsed)
}

func (net *Network) Predict(training [][]float64, trainingname []string) [][]string {
	output := [][]string{}
	Headers := []string{"INDEX", "PREDICTION"}
	output = append(output, Headers)
	for index := range training {
		inputss := make([]float64, net.inputs)
		for i := range inputss {
			inputss[i] = (training[index][i] / 255.0 * 0.999) + 0.001
		}
		outputs := net.predict(inputss)
		best := 0
		highest := 0.0
		for i := 0; i < net.outputs; i++ {
			if outputs.At(i, 0) > highest {
				best = i
				highest = outputs.At(i, 0)
			}
		}
		row := []string{strconv.Itoa(index),net.targetmap[best]} 
		output = append(output, row)
	}
	return output
}

