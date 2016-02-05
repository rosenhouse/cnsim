package main

func SteadyState(inputs Inputs) (Outputs, error) {
	totalInstances := 0
	for i := 0; i < inputs.NumApps; i++ {
		totalInstances += inputs.DistAppSize.Sample()
	}
	medianInstancesPerHost := totalInstances / inputs.NumHosts
	return Outputs{
		MedianInstancesPerHost: medianInstancesPerHost,
		DistAppSize:            SampleN(inputs.DistAppSize, inputs.NumApps),
		DistHostDegree: map[int]float32{
			0: 0.8,
			1: 0.1,
			2: 0.5,
			3: 0.25,
			4: 0.25,
		},
	}, nil
}

type Distribution interface {
	Sample() int
}

func SampleN(dist Distribution, nSamples int) EmpiricalDist {
	counts := make(map[int]int)
	for i := 0; i < nSamples; i++ {
		sample := dist.Sample()
		counts[sample]++
	}

	empDist := make(map[int]float32)
	for x, c := range counts {
		empDist[x] = float32(c) / float32(nSamples)
	}

	return empDist
}
