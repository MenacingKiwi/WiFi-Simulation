package main

import (
	"fmt"
	"wifiProject/utils"
	"math"
)

type State struct {
	state string
	T     float64
}

func GenerateBand24G(state string, E_T0 float64, E_T1 float64, T_24G []float64, st_24G []string, Ts float64 ) ([]float64, []string) {
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
	type Result struct {
		Min float64
		Max float64
	}

	// A slice of all the struct slices
	allStructs := [][]State{struct1, struct2, struct3}

	// Create a new slice to hold the current elements being compared
	// Initialize it with the first elements of each struct
	currentElements := make([]State, len(allStructs))
	for i := range allStructs {
		if len(allStructs[i]) > 0 {
			currentElements[i] = allStructs[i][0]
		}
	}

	// Keep track of the current index for each struct
	indices := make([]int, len(allStructs))

	fmt.Println("Starting repetitive logic until one struct runs out...")
	fmt.Println("=====================================================")

	// Main loop: continue as long as we can pull an element from each struct
	for {
		// --- Step 1: Find the minimum T value and its index from the current elements ---
		minT := math.MaxFloat64
		minIndex := -1
		
		fmt.Printf("\nIteration %d:\n", indices[0] + 1)
		fmt.Println("Current elements being compared:")
		for _, s := range currentElements {
			fmt.Printf("State: %s, T: %.2f\n", s.state, s.T)
		}

		for i, s := range currentElements {
			if s.T < minT {
				minT = s.T
				minIndex = i
			}
		}

		// --- Step 2: Apply the logic based on the state of the minimum value ---
		minState := currentElements[minIndex].state
		fmt.Printf("\nFound minimum T value: %.2f in state '%s' at struct index %d\n", minT, minState, minIndex)

		if minState == "Disconnect" {
			fmt.Println("Condition: min state is 'Disconnect'. Applying logic...")
			// Subtract the min T from all other values
			for i := 0; i < len(currentElements); i++ {
				if i != minIndex {
					currentElements[i].T -= minT
				}
			}

		} else if minState == "Connect" {
			fmt.Println("Condition: min state is 'Connect'. Applying logic...")
			DownloadFile(minT)
			
			// Subtract the min T from all other values
			for i := 0; i < len(currentElements); i++ {
				if i != minIndex {
					currentElements[i].T -= minT
				}
			}
		}

		// --- Step 3: Replace the min element with the next one from its struct ---
		// Increment the index for the struct that had the minimum value
		indices[minIndex]++
		
		// Check if the struct has run out of elements. If so, break the loop.
		if indices[minIndex] >= len(allStructs[minIndex]) {
			fmt.Printf("\nStruct %d has run out of elements. Stopping.\n", minIndex)
			break
		}
		
		// Update the element in the comparison slice with the next one
		currentElements[minIndex] = allStructs[minIndex][indices[minIndex]]
		
		fmt.Printf("Replaced element with: State: %s, T: %.2f\n", currentElements[minIndex].state, currentElements[minIndex].T)
	}

	fmt.Println("\n=====================================================")
	fmt.Println("Final elements in the comparison set:")
	for _, s := range currentElements {
		fmt.Printf("State: %s, T: %.2f\n", s.state, s.T)
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
	expectedValueSession := 500.0
	expectedValueT0 := 60.0
	expectedValueT1 := 40.0
	// alpha := 5.0
	// xm := 100.0

	Ts := utils.GenerateTs(expectedValueSession) //generate Ts
	fmt.Println("=======================================================================")
	fmt.Printf("Session Time: %f second\n", Ts)

	fs := 1000.00 //File size in MB
	fmt.Printf("File Size: %f MB\n", fs)
	fmt.Printf("Total File Size: %f Mb\n", fs*8)
	fmt.Println("=======================================================================")
	//speed := 150.00 //Download speed in Mbps
	fs = fs * 8 //Convert to Mb

	state24G := utils.InitState(60.0, 40.0) //Initial state
	state5G := utils.InitState(120.0, 80.0)
	state6G := utils.InitState(180.0, 120.0)

	//loop to generate states
	T_24G = GenerateBand(state24G, expectedValueT0, expectedValueT1, Ts)
	T_5G = GenerateBand(state5G, 120.0, 80.0, Ts)
	T_6G = GenerateBand(state6G, 180.0, 120.0, Ts)

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

	//If file download is finished print done, else print remaining size to download
	if fs <= 0 {
		fmt.Println("Done")
	} else {
		fmt.Printf("Remaining File: %f Mb\n", fs)
		fmt.Printf("Remaining File: %f MB\n", fs/8)
	}
	fmt.Println("=======================================================================")
	FindMinMaxPerStates(T_24G, T_5G, T_6G)
}
