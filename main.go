package main

import (
	"fmt"
	"math"
	"wifiProject/utils"
)

var FileSize float64
var Speed2G float64
var Speed5G float64
var Speed6G float64

type State struct {
	state string
	T     float64
}

func DownloadFile(structNames []string, value float64) {
	if len(structNames) > 0 {
		fmt.Printf("Detected 'Connect' states. Downloading files for the following structs with Min T value: %.2f\n", value)
		for _, name := range structNames {
			switch name{
			case "2.4G":
				fmt.Printf(" - %s\n", name)
				FileSize = FileSize - (Speed2G * value)
			case "5G":
				fmt.Printf(" - %s\n", name)
				FileSize = FileSize - (Speed5G * value)
			case "6G":
				fmt.Printf(" - %s\n", name)
				FileSize = FileSize - (Speed6G * value)
			default:
				fmt.Printf("No Connected\n")
			}
		}
		fmt.Printf("File Remaining:  %f\n", FileSize)
	}
}

func SumOfTinState(T []State) float64 {
	sum := 0.0
	for _, t := range T {
		sum += t.T
	}
	return sum
}

func PrintState(T []State) {
	for _, t := range T {
		fmt.Printf("(%s %f)\n", t.state, t.T)
	}
}
func FindMinMaxPerStates(struct1 []State, struct2 []State, struct3 []State) {

	// A slice of all the struct slices and their corresponding names for easy reference.
	allStructs := [][]State{struct1, struct2, struct3}
	structNames := []string{"2.4G", "5G", "6G"}

	// Create a new slice to hold the current elements being compared
	currentElements := make([]State, len(allStructs))
	for i := range allStructs {
		if len(allStructs[i]) > 0 {
			currentElements[i] = allStructs[i][0]
		}
	}

	// Keep track of the current index for each struct
	indices := make([]int, len(allStructs))

	fmt.Println("Starting repetitive logic until a struct runs out of elements...")
	fmt.Println("=====================================================")

	// Main loop: continue as long as we can pull an element from each struct
	for {
		// --- Step 1: Find the minimum T value and its index from the current elements ---
		minT := math.MaxFloat64
		minIndex := -1

		fmt.Printf("\n---------------------\n")
		fmt.Printf("\nIteration %d:\n", indices[0]+1)
		fmt.Println("Current elements being compared:")
		for i, s := range currentElements {
			fmt.Printf("Struct: %s, State: %s, T: %.2f\n", structNames[i], s.state, s.T)
		}

		for i, s := range currentElements {
			if s.T < minT {
				minT = s.T
				minIndex = i
			}
		}

		// --- Step 2: Check for "Connect" states and call DownloadFile ---
		fmt.Printf("\nFound minimum T value: %.2f in state '%s' at struct '%s'\n", minT, currentElements[minIndex].state, structNames[minIndex])

		var connectStructs []string
		for i, s := range currentElements {
			if s.state == "Connect" {
				connectStructs = append(connectStructs, structNames[i])
			}
		}
		
		DownloadFile(connectStructs, minT)

		// --- Step 3: Replace the min element with the next one from its struct ---
		// Increment the index for the struct that had the minimum value
		indices[minIndex]++

		// Check if the struct has run out of elements. If so, break the loop.
		if indices[minIndex] >= len(allStructs[minIndex]) {
			fmt.Printf("\nStruct '%s' has run out of elements. Stopping.\n", structNames[minIndex])
			break
		}

		// Update the element in the comparison slice with the next one
		currentElements[minIndex] = allStructs[minIndex][indices[minIndex]]

		fmt.Printf("Replaced minimum element from '%s' with its next value: State: %s, T: %.2f\n", structNames[minIndex], currentElements[minIndex].state, currentElements[minIndex].T)

		// --- Step 4: Subtract min T from all other values in the current set ---
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

func GenerateBand(state string, E_T0 float64, E_T1 float64, Ts float64) []State {
	T := []State{}
	for {

		if state == "disconnect" {
			//if file is downloaded exit loop

			T0 := utils.GenerateT0(E_T0) //Generate T0

			T = append(T, State{state: "Disconnect", T: T0})

			//if sum of T in state exceed Ts, exit loop
			if SumOfTinState(T) > Ts {

				break
			}

		} else {

			T1 := utils.GenerateT1(E_T1)

			//if sum of T in state exceed Ts
			if (SumOfTinState(T) + T1) > Ts {

				downloadable := Ts - SumOfTinState(T) //Only download in range of Ts
				T = append(T, State{state: "Connect", T: downloadable})

				//File download

				break
			}

			T = append(T, State{state: "Connect", T: T1})

			//File download

		}

		state = utils.NextState(state) //generate next state
	}
	return T
}

func main() {

	//T of states
	T_24G := []State{}
	T_5G := []State{}
	T_6G := []State{}

	//Parameter
	expectedValueSession := 200.0
	expectedValueT0 := 60.0
	expectedValueT1 := 40.0
	// alpha := 5.0
	// xm := 100.0

	Ts := utils.GenerateTs(expectedValueSession) //generate Ts
	fmt.Println("=======================================================================")
	fmt.Printf("Session Time: %f second\n", Ts)


	FileSize = 10000.00 //File size in MB
	fmt.Printf("File Size: %f MB\n", FileSize)
	FileSize = FileSize * 8
	fmt.Printf("Total File Size: %f Mb\n", FileSize)
	fmt.Println("=======================================================================")
	Speed2G = 150.00 //Download speed in Mbps
	Speed5G = 500.00
	Speed6G = 500.00
	

	state24G := utils.InitState(60.0, 40.0) //Initial state
	state5G := utils.InitState(60.0, 40.0)
	state6G := utils.InitState(60.0, 40.0)

	//loop to generate states
	T_24G = GenerateBand(state24G, expectedValueT0, expectedValueT1, Ts)
	T_5G = GenerateBand(state5G, expectedValueT0, expectedValueT1, Ts)
	T_6G = GenerateBand(state6G, expectedValueT0, expectedValueT1, Ts)

	//Output
	fmt.Println("=======================================================================")
	fmt.Println("2.4G Band")
	PrintState(T_24G)
	fmt.Println("=======================================================================")
	fmt.Println("5G Band")
	PrintState(T_5G)
	fmt.Println("=======================================================================")
	fmt.Println("6G Band")
	PrintState(T_6G)
	fmt.Println("=======================================================================")

	FindMinMaxPerStates(T_24G, T_5G, T_6G)

	//If file download is finished print done, else print remaining size to download
	if FileSize <= 0 {
		fmt.Println("Done")
	} else {
		fmt.Printf("Remaining File: %f Mb\n", FileSize)
		fmt.Printf("Remaining File: %f MB\n", FileSize/8)
	}
	fmt.Println("=======================================================================")
	
}
