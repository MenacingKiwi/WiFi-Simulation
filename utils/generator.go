package utils

import (
	"math"
	"math/rand"
)
//!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
//deprecated

func InverseTransformRandom() int {
	x := []float64{0.0, 0.4, 0.65, 0.85, 0.9, 1.0}

	u := rand.Float64()

	for i := 1; i < len(x); i++ {
		if u <= x[i] {
			return i
		}
	}
	return 0
}

func InverseCDFExponential(u float64, val float64) float64 {
	return (-val) * math.Log(1-u)
}

func GeneratePareto(alpha, xm float64) float64 {
	u := rand.Float64() // Uniform random number between 0 and 1
	return xm / math.Pow(1-u, 1/alpha)
}

func SumFloat64Array(arr []float64) float64 {
	var sum float64 = 0.0
	for _, value := range arr {
		sum += value
	}
	return sum
}

func Mean(m []float64) float64 {
    if len(m) == 0 {
        return 0.0
    }
    return SumFloat64Array(m) / float64(len(m))
}

func GenerateTs(expectedValueSession float64) float64 {
	u := rand.Float64() // Generate a uniform random number between 0 and 1
	x := InverseCDFExponential(u, expectedValueSession)
	return x
}

func GenerateT0(expectedValueT0 float64) float64 {
	u := rand.Float64() // Generate a uniform random number between 0 and 1
	x := InverseCDFExponential(u, expectedValueT0)
	return x
}

func GenerateT1(expectedValueT1 float64) float64 {
	u := rand.Float64() // Generate a uniform random number between 0 and 1
	x := InverseCDFExponential(u, expectedValueT1)
	return x
}

func InitState(expectedValueT0 float64, expectedValueT1 float64) string {
	p0 := expectedValueT0 / (expectedValueT1 + expectedValueT0)
	u := rand.Float64()
	if u <= p0 {
		return "disconnect"
	}
	return "connect"
}

func NextState(state string) string {
	if state == "disconnect" {
		return "connect"
	}
	return "disconnect"
}