package embedding

import (
	"github.com/hypermodeinc/modus/sdk/go/pkg/vectors"
)

type VectorOperations struct{}

func NewVectorOperations() *VectorOperations {
	return &VectorOperations{}
}

func (vo *VectorOperations) CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) {
		return 0
	}

	aFloat := make([]float64, len(a))
	bFloat := make([]float64, len(b))

	for i := range a {
		aFloat[i] = float64(a[i])
		bFloat[i] = float64(b[i])
	}

	dotProduct := vectors.Dot(aFloat, bFloat)
	magnitudeA := vectors.Magnitude(aFloat)
	magnitudeB := vectors.Magnitude(bFloat)

	if magnitudeA == 0 || magnitudeB == 0 {
		return 0
	}

	return dotProduct / (magnitudeA * magnitudeB)
}
