package math

import "math"

func Delta(a, b []float32) []float32 {
	delta := make([]float32, len(a))
	for i := range delta {
		delta[i] = b[i] - a[i]
	}
	return delta
}

func Norm2(vec []float32) float32 {
	norm2 := float32(0)
	for i := 0; i < len(vec); i++ {
		norm2 += vec[i] * vec[i]
	}
	return norm2
}

func Norm(vec []float32) float32 {
	return float32(math.Sqrt(float64(Norm2(vec))))
}

func Normalise(vec []float32) []float32 {
	norm2 := Norm2(vec)
	norm := float32(math.Sqrt(float64(norm2)))
	for i := 0; i < len(vec); i++ {
		vec[i] /= norm
	}
	return vec
}

func Dot(a, b []float32) float32 {
	sum := float32(0)
	for i := 0; i < len(a); i++ {
		sum += a[i] * b[i]
	}
	return sum
}

func Ip(a, b []float32) float32 {
	return 1 - Dot(a, b)
}

func Cosine(a, b []float32) float32 {
	return Ip(a, b) / (Norm(a) * Norm(b))
}

func L2(a, b []float32) float32 {
	return Norm2(Delta(a, b))
}
