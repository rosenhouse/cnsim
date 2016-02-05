package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"sort"
)

// PoissonDist encodes an Poisson probability distribution
type PoissonDist struct {
	Mean float32
}

func (dist PoissonDist) Sample() int {
	λ := float64(dist.Mean)
	x := 0
	p := math.Exp(-λ)
	s := p
	u := rand.Float64()

	for u > s {
		x = x + 1
		p = p * λ / float64(x)
		s = s + p
	}
	return x
}

// EmpiricalDist encodes an empircal probability distribution on the integers
type EmpiricalDist map[int]float32

func (e EmpiricalDist) MarshalJSON() ([]byte, error) {
	var keys = make([]int, len(e))
	i := 0
	for k, _ := range e {
		keys[i] = k
		i++
	}
	sort.Ints(keys)

	var rows = make([][]float32, len(e))
	for i, k := range keys {
		rows[i] = []float32{float32(k), e[k]}
	}

	return json.Marshal(rows)
}

// Inputs encodes the parameters of the simulation
type Inputs struct {
	NumHosts      int
	NumApps       int
	DistAppSize   PoissonDist
	ProbReflexive float32
	DistAppDegree PoissonDist
}

// Outputs encodes the results of the simulation
type Outputs struct {
	MedianInstancesPerHost int
	DistHostDegree         EmpiricalDist
	DistAppSize            EmpiricalDist
}

func main() {
	stdinBytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	var inputs Inputs
	err = json.Unmarshal(stdinBytes, &inputs)
	if err != nil {
		panic(err)
	}

	outputs, err := SteadyState(inputs)
	if err != nil {
		panic(err)
	}

	outputBytes, err := json.MarshalIndent(outputs, "", "  ")
	if err != nil {
		panic(err)
	}

	_, err = os.Stdout.Write(outputBytes)
	if err != nil {
		panic(err)
	}
	fmt.Println()
}
