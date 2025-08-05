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

func main() {

	m := []float64{} //T of states
	n := []string{}  //states label

	//Parameter
	expectedValueSession := 500.0
	expectedValueT0 := 60.0
	expectedValueT1 := 40.0
	// alpha := 5.0
	// xm := 100.0

	Ts := utils.GenerateTs(expectedValueSession) //generate Ts
	fmt.Printf("Session Time: %f second\n", Ts)

	fs := 10000.00 //File size in MB
	fmt.Printf("File Size: %f MB\n", fs)
	fmt.Printf("Total File Size: %f Mb\n", fs*8)
	speed := 150.00 //Download speed in Mbps

	state := utils.InitState(expectedValueT0, expectedValueT1) //Initial state
	fs = fs * 8                                                //Convert to Mb

	//loop to generate states
	for {
		if state == "disconnect" {
			//if file is downloaded exit loop
			if fs <= 0 {
				break
			}

			T0 := utils.GenerateT0(expectedValueT0) //Generate T0

			n = append(n, "disconnect")
			m = append(m, T0)

			//if sum of T in state exceed Ts, exit loop
			if utils.SumFloat64Array(m) > Ts {

				break
			}

		} else {
			if fs <= 0 {
				break
			}

			T1 := utils.GenerateT1(expectedValueT1)

			//if sum of T in state exceed Ts
			if (utils.SumFloat64Array(m) + T1) > Ts {

				downloadable := Ts - utils.SumFloat64Array(m) //Only download in range of Ts
				m = append(m, downloadable)
				n = append(n, "connect")

				//File download
				fs = DownloadFile(speed, downloadable, fs)

				break
			}

			n = append(n, "connect")
			m = append(m, T1)

			//File download
			fs = DownloadFile(speed, T1, fs)

		}

		state = utils.NextState(state) //generate next state
	}

	//Output
	fmt.Println(m)
	fmt.Println(n)

	//If file download is finished print done, else print remaining size to download
	if fs <= 0 {
		fmt.Println("Done")
	} else {
		fmt.Printf("Remaining File: %f Mb\n", fs)
		fmt.Printf("Remaining File: %f MB\n", fs/8)
	}

	fmt.Println(utils.SumFloat64Array(m)) //sum of T
}
