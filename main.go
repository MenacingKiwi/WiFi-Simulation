package main

import (
	"fmt"
	"wifiProject/utils"
)

func DownloadFile(s float64, T float64, f float64) float64{
	fmt.Println("=============================")
	fmt.Printf("Total Downloaded: %f Mb\n", (s * T))
	f = f - (s * T)
	fmt.Printf("Remain: %f Mb\n", f)
	fmt.Println("=============================")
	return f
}

func GenerateBand24G(state string, E_T0 float64, E_T1 float64, T_24G []float64, st_24G []string, Ts float64 ) ([]float64, []string) {
	for {
		if state == "disconnect" {
			//if file is downloaded exit loop

			T0_24G := utils.GenerateT0(E_T0) //Generate T0

			st_24G = append(st_24G, "disconnect")
			T_24G = append(T_24G, T0_24G)

			//if sum of T in state exceed Ts, exit loop
			if utils.SumFloat64Array(T_24G) > Ts {

				break
			}

		} else {

			T1 := utils.GenerateT1(E_T1)

			//if sum of T in state exceed Ts
			if (utils.SumFloat64Array(T_24G) + T1) > Ts {

				downloadable := Ts - utils.SumFloat64Array(T_24G) //Only download in range of Ts
				T_24G = append(T_24G, downloadable)
				st_24G = append(st_24G, "connect")

				//File download

				break
			}

			st_24G = append(st_24G, "connect")
			T_24G = append(T_24G, T1)

			//File download

		}

		state = utils.NextState(state) //generate next state
	}
	return T_24G, st_24G
}

func main() {
	
	//T of states
	T_24G := []float64{} 
	// T_5G := []float64{}
	// T_6G := []float64{}
	//states label
	st_24G := []string{}
	// st_5G := []string{}
	// st_6G := []string{}


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
	fs = fs * 8                                                //Convert to Mb

	state24G := utils.InitState(60.0, 40.0) //Initial state
	// state5G := utils.InitState(120.0,80.0)
	// state6G := utils.InitState(180.0,120.0)

	//loop to generate states
	T_24G, st_24G = GenerateBand24G(state24G, expectedValueT0, expectedValueT1, T_24G, st_24G, Ts)

	//Output
	fmt.Println("=======================================================================")
	fmt.Println(T_24G)
	fmt.Println(st_24G)
	fmt.Println("=======================================================================")

	//If file download is finished print done, else print remaining size to download
	if fs <= 0 {
		fmt.Println("Done")
	} else {
		fmt.Printf("Remaining File: %f Mb\n", fs)
		fmt.Printf("Remaining File: %f MB\n", fs/8)
	}

	fmt.Printf("Total T: %f\n",utils.SumFloat64Array(T_24G)) //sum of T
}
