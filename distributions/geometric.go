package distributions

import (
	"fmt"
	"math/rand"
)

type GeometricWithPositiveSupport struct{}

func (_ *GeometricWithPositiveSupport) Sample(desiredMean float64) (int, error) {
	if desiredMean < 1 {
		return -1, fmt.Errorf("desiredMean must be >= 1")
	}
	probSuccess := 1.0 / desiredMean
	return countTrialsBeforeSuccess(probSuccess)
}

func countTrialsBeforeSuccess(probSuccess float64) (int, error) {
	const MAX_TRIALS = 1 << 16
	for i := 1; i < MAX_TRIALS; i++ {
		if rand.Float64() < probSuccess {
			return i, nil
		}
	}
	return -1, fmt.Errorf("exceeded max trials: %d", MAX_TRIALS)
}
