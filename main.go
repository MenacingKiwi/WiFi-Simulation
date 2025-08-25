package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"reflect"
	"wifiProject/utils"
)

// State represents a state and a T value.
type State struct {
	state string
	T     float64
}

var remainingSize float64

// SimulateFileSize generates a slice of random file sizes.
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
		fmt.Printf("Remaining file to download: %f\n", remainingSize)
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

func InitState(expectedValueT0 float64, expectedValueT1 float64) string {
	p0 := expectedValueT0 / (expectedValueT1 + expectedValueT0)
	u := rand.Float64()
	if u <= p0 {
		return "disconnect"
	}
	return "connect"
}

// GenerateBand generates a series of states and T values for a given band.
func GenerateBand(initialState string, expectedValueT0, expectedValueT1, Ts float64) ([]State, int) {
	states := []State{}
	currentState := initialState
	connectCount := 0
	for {
		if currentState == "disconnect" {
			t0 := utils.InverseCDFExponential(rand.Float64(), expectedValueT0)
			states = append(states, State{state: "Disconnect", T: t0})
			if SumOfTinState(states) >= Ts {
				break
			}
		} else { // "connect"
			t1 := utils.InverseCDFExponential(rand.Float64(), expectedValueT1)
			if SumOfTinState(states)+t1 >= Ts {
				downloadable := Ts - SumOfTinState(states)
				states = append(states, State{state: "Connect", T: downloadable})
				connectCount++
				break
			}
			states = append(states, State{state: "Connect", T: t1})
			connectCount++
		}
		currentState = utils.NextState(currentState)
	}
	return states, connectCount
}

// FindMinMaxPerStates contains the core logic for the iterative comparison.
func FindMinMaxPerStates(structData map[string][]State, speeds map[string]float64) {
	// A slice of all the struct names for consistent ordering
	var structNames []string
	for name := range structData {
		structNames = append(structNames, name)
	}

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

		// Find the minimum T value and its index
		for i, s := range currentElements {
			if s.T < minT {
				minT = s.T
				minIndex = i
			}
		}

		// Check for "Connect" states and call DownloadFile

		var connectStructs []string
		for i, s := range currentElements {
			if s.state == "Connect" {
				connectStructs = append(connectStructs, structNames[i])
			}
		}

		remainingSize = DownloadFile(connectStructs, minT, speeds)
		fmt.Printf("After download remaining: %.2f Mb\n", remainingSize)

		// Replace the min element with the next one from its struct
		indices[minIndex]++
		if indices[minIndex] >= len(structData[structNames[minIndex]]) {
			fmt.Printf("\nStruct '%s' has run out of elements. Stopping.\n", structNames[minIndex])
			break
		}

		currentElements[minIndex] = structData[structNames[minIndex]][indices[minIndex]]

		// Subtract min T from all other values
		for i := 0; i < len(currentElements); i++ {
			if i != minIndex {
				currentElements[i].T -= minT
			}
		}
	}

	fmt.Println("\n=====================================================")
}

func main() {

	result := []float64{}
	miss := 0.0
	guarateeBandwidth := 40
	satisfy := 0

	// User input for simulation mode
	var mode int
	fmt.Println("Select simulation mode:")
	fmt.Println("1. Single link")
	fmt.Println("2. Double links")
	fmt.Println("3. Triple links")
	fmt.Print("Enter your choice (1-3): ")
	fmt.Scanln(&mode)
	if reflect.TypeOf(mode).Kind() != reflect.Int || mode < 1 || mode > 3 {
		fmt.Println("Invalid mode selected. Please choose between 1 and 3.")
		os.Exit(1)
	}

	//
	var round int
	fmt.Print("Enter the number of rounds to simulate: ")
	fmt.Scanln(&round)
	if reflect.TypeOf(round).Kind() != reflect.Int || round <= 0 {
		fmt.Println("Invalid round")
		os.Exit(1)
	}
	f := SimulateFileSize(220, 8.5, round)

	speeds := make(map[string]float64)

	//
	var s float64
	for i := range mode {
		fmt.Printf("Enter speed for Link %d (in Mbps): ", i+1)
		fmt.Scanln(&s)
		if reflect.TypeOf(s).Kind() != reflect.Float64 || s <= 0 {
			fmt.Println("Invalid speed for Link 1")
			os.Exit(1)
		}
		speeds[fmt.Sprintf("Link%d", i+1)] = s
	}

	allStructData := make(map[string][]State)

	for i := range round {
		fmt.Printf("\nRound %d\n", i)
		//expectedValueSession := 100.0 //increase by 10 every n rounds till 100
		expectedValueT0 := 50.0
		expectedValueT1 := 50.0

		Ts := 100.0
		fmt.Println("=======================================================================")
		fmt.Printf("Session Time: %.2f second\n", Ts)

		FileSize := f[i]
		fmt.Printf("File Size: %.2f MB\n", FileSize)
		totalFileSizeMb := FileSize * 8
		fmt.Printf("Total File Size: %.2f Mb\n", totalFileSizeMb)
		fmt.Println("=======================================================================")

		switch mode {
		case 1:
			Link1, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			if len(Link1) > 0 {
				allStructData["Link1"] = Link1
			}
		case 2:
			Link1, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			Link2, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			if len(Link1) > 0 {
				allStructData["Link1"] = Link1
			}
			if len(Link2) > 0 {
				allStructData["Link2"] = Link2
			}
		case 3:
			Link1, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			Link2, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			Link3, _ := GenerateBand(InitState(expectedValueT0, expectedValueT1), expectedValueT0, expectedValueT1, Ts)
			if len(Link1) > 0 {
				allStructData["Link1"] = Link1
			}
			if len(Link2) > 0 {
				allStructData["Link2"] = Link2
			}
			if len(Link3) > 0 {
				allStructData["Link3"] = Link3
			}
		default:
			fmt.Println("Invalid mode selected. Please choose between 1 and 3.")
		}

		fmt.Println("=======================================================================")
		fmt.Println("Generated States:")
		for name, link := range allStructData {
			fmt.Printf("%s Band\n", name)
			for _, s := range link {
				fmt.Printf("(%s %.2f)\n", s.state, s.T)
			}
			fmt.Println("-----------------------------------------------------------------------")
		}
		fmt.Println("=======================================================================")

		remainingSize = totalFileSizeMb
		FindMinMaxPerStates(allStructData, speeds)

		fmt.Println("=======================================================================")
		if remainingSize <= 0 {
			fmt.Println("Done")
			result = append(result, 0)
		} else {
			result = append(result, remainingSize/8)
			fmt.Printf("Remaining File: %.2f Mb\n", remainingSize)
			fmt.Printf("Remaining File: %.2f MB\n", remainingSize/8)
			miss++
		}

		// Simplified bandwidth satisfaction check
		totalBandwidth1 := 0.0
		totalBandwidth2 := 0.0
		totalBandwidth3 := 0.0

		totalLen := 0
		if len(allStructData["Link1"]) > 0 {
			connectCount := 0
			for _, s := range allStructData["Link1"] {
				if s.state == "Connect" {
					connectCount++
				}
			}
			totalBandwidth1 = (float64(connectCount) * speeds["Link1"]) / float64(len(allStructData["Link1"]))
			fmt.Println(totalBandwidth1)
			totalLen += len(allStructData["Link1"])
		}
		if len(allStructData["Link2"]) > 0 {
			connectCount := 0
			for _, s := range allStructData["Link2"] {
				if s.state == "Connect" {
					connectCount++
				}
			}
			totalBandwidth2 = (float64(connectCount) * speeds["Link2"]) / float64(len(allStructData["Link2"]))
			fmt.Println(totalBandwidth2)
			totalLen += len(allStructData["Link2"])
		}
		if len(allStructData["Link3"]) > 0 {
			connectCount := 0
			for _, s := range allStructData["Link3"] {
				if s.state == "Connect" {
					connectCount++
				}
			}
			totalBandwidth3 = (float64(connectCount) * speeds["Link3"]) / float64(len(allStructData["Link3"]))
			fmt.Println(totalBandwidth3)
			totalLen += len(allStructData["Link3"])
		}
		if totalLen > 0 && (totalBandwidth1+totalBandwidth2+totalBandwidth3)/3 >= float64(guarateeBandwidth) {

			fmt.Println("Overall Bandwidth Satisfied")
			fmt.Println((totalBandwidth1 + totalBandwidth2 + totalBandwidth3) / 3)
			satisfy++
		} else {
			fmt.Println("Overall Bandwidth Not Satisfied")
			fmt.Println((totalBandwidth1 + totalBandwidth2 + totalBandwidth3) / 3)
		}
		fmt.Println("=======================================================================")
	}

	fmt.Println("Statistics")
	//fmt.Println(result)
	if len(result) > 0 {
		fmt.Printf("Average Remaining File Size: %f MB\n", utils.SumFloat64Array(result)/float64(len(result)))
	}
	fmt.Printf("Deadline Miss Rate: %.2f%%\n", (miss/float64(round))*100)
	fmt.Printf("Bandwidth Satisfy: %d out of %d \n", satisfy, round)
	fmt.Printf("Bandwidth Satisfy Rate: %.2f%%\n", (float64(satisfy)/float64(round))*100)
	fmt.Println("=======================================================================")
}
