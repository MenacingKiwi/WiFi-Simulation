package main

import (
	"fmt"
	"math"
	"math/rand"
	"wifiProject/utils"
)

// State represents a state and a T value.
type State struct {
	state string
	T     float64
}

var FileSize float64
var remainingSize float64

func SimulateFileSize(xmin float64, alpha float64, num int) []float64 {
	files := []float64{}
	for range num {
		u := rand.Float64()
		files = append(files, xmin/math.Pow(1-u, 1/alpha))
	}
	return files
}

// DownloadFile simulates a download and updates the remaining FileSize.
func DownloadFile(structNames []string, value float64, speeds map[string]float64) float64 {
	if len(structNames) > 0 {
		fmt.Printf("Downloading files for the following structs with T value: %.2f\n", value)
		fmt.Printf("Remainin file to download: %f\n", remainingSize)
		for _, name := range structNames {
			speed := speeds[name]
			if remainingSize > 0 {

				fmt.Printf(" - %s at speed %.2f Mbps\n", name, speed)
				remainingSize -= (speed * value)

			}
		}
	}
	if remainingSize > 0 {
		return remainingSize
	} else {
		return 0
	}
}

// SumOfTinState sums the T values in a slice of State structs.
func SumOfTinState(states []State) float64 {
	sum := 0.0
	for _, s := range states {
		sum += s.T
	}
	return sum
}

// GenerateBand generates a series of states and T values for a given band.
func GenerateBand(initialState string, expectedValueT0, expectedValueT1, Ts float64) ([]State, int) {
	states := []State{}
	currentState := initialState
	connectCount := 0
	for {
		if currentState == "disconnect" {
			t0 := InverseCDFExponential(rand.Float64(), expectedValueT0)
			states = append(states, State{state: "Disconnect", T: t0})
			if SumOfTinState(states) >= Ts {
				break
			}
		} else { // "connect"
			t1 := InverseCDFExponential(rand.Float64(), expectedValueT1)
			if SumOfTinState(states)+t1 >= Ts {
				downloadable := Ts - SumOfTinState(states)
				states = append(states, State{state: "Connect", T: downloadable})
				connectCount++
				break
			}
			states = append(states, State{state: "Connect", T: t1})
			connectCount++
		}
		currentState = NextState(currentState)
	}
	return states, connectCount
}

// FindMinMaxPerStates contains the core logic for the iterative comparison.
func FindMinMaxPerStates(structData map[string][]State, speeds map[string]float64) {
	// A slice of all the struct names for consistent ordering
	structNames := []string{"2.4G", "5G", "6G"}

	// Initialize slices to hold the current elements and indices
	currentElements := make([]State, len(structNames))
	indices := make([]int, len(structNames))

	// Populate initial currentElements with the first element of each struct
	for i, name := range structNames {
		if len(structData[name]) > 0 {
			currentElements[i] = structData[name][0]
		}
	}

	fmt.Println("=====================================================")

	for iteration := 1; ; iteration++ {

		if remainingSize <= 0 {
			fmt.Printf("\n------------------\nFile Download Done\n------------------\n")
			break
		}

		minT := math.MaxFloat64
		minIndex := -1

		fmt.Printf("\n---------------------\n")
		fmt.Printf("\nIteration %d:\n\n", iteration)

		for i, s := range currentElements {
			fmt.Printf("Struct: %s, State: %s, T: %.2f\n", structNames[i], s.state, s.T)
		}

		// Find the minimum T value and its index
		for i, s := range currentElements {
			if s.T < minT {
				minT = s.T
				minIndex = i
			}
		}

		// Check for "Connect" states and call DownloadFile
		fmt.Printf("\nFound minimum T value: %.2f in state '%s' at struct '%s'\n", minT, currentElements[minIndex].state, structNames[minIndex])

		var connectStructs []string
		for i, s := range currentElements {
			if s.state == "Connect" {
				connectStructs = append(connectStructs, structNames[i])
			}
		}

		remainingSize := DownloadFile(connectStructs, minT, speeds)
		fmt.Printf("After download remaining: %.2f Mb\n", remainingSize)

		// Replace the min element with the next one from its struct
		indices[minIndex]++
		if indices[minIndex] >= len(structData[structNames[minIndex]]) {
			fmt.Printf("\nStruct '%s' has run out of elements. Stopping.\n", structNames[minIndex])
			break
		}

		currentElements[minIndex] = structData[structNames[minIndex]][indices[minIndex]]
		fmt.Printf("Replaced minimum element from '%s' with its next value: State: %s, T: %.2f\n", structNames[minIndex], currentElements[minIndex].state, currentElements[minIndex].T)

		// Subtract min T from all other values
		for i := 0; i < len(currentElements); i++ {
			if i != minIndex {
				currentElements[i].T -= minT
			}
		}
	}

	fmt.Println("\n=====================================================")
	fmt.Println("Final elements in the comparison set:")
	for i, s := range currentElements {
		fmt.Printf("Struct: %s, State: %s, T: %.2f\n", structNames[i], s.state, s.T)
	}
}

// Helper functions (formerly in generator.go)
func InverseCDFExponential(u, val float64) float64 {
	return (-val) * math.Log(1-u)
}

func NextState(state string) string {
	if state == "disconnect" {
		return "connect"
	}
	return "disconnect"
}

func InitState(expectedValueT0, expectedValueT1 float64) string {
	p0 := expectedValueT0 / (expectedValueT1 + expectedValueT0)
	u := rand.Float64()
	if u <= p0 {
		return "disconnect"
	}
	return "connect"
}

func main() {

	// T of states
	// T_24G := []State{}
	// T_5G := []State{}
	// T_6G := []State{}

	round := 100
	result := []float64{}
	miss := 0.0
	guarateeBandwidth := 100
	satisfy := 0
	f := SimulateFileSize(1000.0, 1.5, round)

	for i := range round {
		fmt.Printf("Round %d\n", i)
		// Parameters
		expectedValueSession := 200.0
		expectedValueT0 := 60.0
		expectedValueT1 := 40.0

		Ts := InverseCDFExponential(rand.Float64(), expectedValueSession)
		fmt.Println("=======================================================================")
		fmt.Printf("Session Time: %.2f second\n", Ts)

		FileSize := f[i] // File size in MB
		fmt.Printf("File Size: %.2f MB\n", FileSize)
		totalFileSizeMb := FileSize * 8
		fmt.Printf("Total File Size: %.2f Mb\n", totalFileSizeMb)
		fmt.Println("=======================================================================")

		speeds := map[string]float64{
			"2.4G": 150.00, // Download speed in Mbps
			"5G":   500.00,
			"6G":   500.00,
		}

		state24G := InitState(expectedValueT0, expectedValueT1)
		state5G := InitState(expectedValueT0, expectedValueT1)
		state6G := InitState(expectedValueT0, expectedValueT1)

		T_24G, connectCount24G := GenerateBand(state24G, expectedValueT0, expectedValueT1, Ts)
		T_5G, connectCount5G := GenerateBand(state5G, expectedValueT0, expectedValueT1, Ts)
		T_6G, connectCount6G := GenerateBand(state6G, expectedValueT0, expectedValueT1, Ts)

		// Combine generated structs into a single map for easier handling
		allStructData := map[string][]State{
			"2.4G": T_24G,
			"5G":   T_5G,
			"6G":   T_6G,
		}

		// Initial output
		fmt.Println("=======================================================================")
		fmt.Println("Generated States:")
		fmt.Println("2.4G Band")
		for _, s := range T_24G {
			fmt.Printf("(%s %.2f)\n", s.state, s.T)
		}
		fmt.Printf("Connect Count: %d\n", connectCount24G)
		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println("5G Band")
		for _, s := range T_5G {
			fmt.Printf("(%s %.2f)\n", s.state, s.T)
		}
		fmt.Printf("Connect Count: %d\n", connectCount5G)
		fmt.Println("-----------------------------------------------------------------------")
		fmt.Println("6G Band")
		for _, s := range T_6G {
			fmt.Printf("(%s %.2f)\n", s.state, s.T)
		}
		fmt.Printf("Connect Count: %d\n", connectCount6G)
		fmt.Println("=======================================================================")

		remainingSize = totalFileSizeMb
		FindMinMaxPerStates(allStructData, speeds)

		fmt.Println("=======================================================================")
		// Final download status
		if remainingSize <= 0 {
			fmt.Println("Done")
			result = append(result, 0)

		} else {
			result = append(result, remainingSize/8)
			fmt.Printf("Remaining File: %.2f Mb\n", remainingSize)
			fmt.Printf("Remaining File: %.2f MB\n", remainingSize/8)
			miss++
		}
		fmt.Printf("\nBandwidth Satisfation\n")
		fmt.Printf("%d\n", (connectCount24G*150)/len(T_24G))
		if (connectCount24G*150)/len(T_24G) >= guarateeBandwidth {
			fmt.Println("2.4G Bandwidth Satisfied")
		} else {
			fmt.Println("2.4G Bandwidth Not Satisfied")
		}
		if (connectCount5G*500)/len(T_5G) >= guarateeBandwidth {
			fmt.Println("5G Bandwidth Satisfied")
		} else {
			fmt.Println("5G Bandwidth Not Satisfied")
		}
		if (connectCount6G*500)/len(T_6G) >= guarateeBandwidth {
			fmt.Println("6G Bandwidth Satisfied")
		} else {
			fmt.Println("6G Bandwidth Not Satisfied")
		}
		if ((connectCount24G*150)/len(T_24G)+(connectCount5G*500)/len(T_5G)+(connectCount6G*500)/len(T_6G))/3 >= guarateeBandwidth {
			fmt.Println("Overall Bandwidth Satisfied")
			satisfy++
		}
		fmt.Println("=======================================================================")
	}
	fmt.Println("Statistics")
	fmt.Println(result)
	fmt.Printf("Average Remaining File Size: %f MB\n", utils.SumFloat64Array(result)/float64(len(result)))
	fmt.Printf("Deadline Miss Rate: %.2f%%\n", (miss/float64(round))*100)
	fmt.Printf("Bandwidth Satisfy: %d out of %d \n", satisfy, round)
	fmt.Printf("Bandwidth Satisfy Rate: %.2f%%\n", (float64(satisfy)/float64(round))*100)
	fmt.Println("=======================================================================")

}
