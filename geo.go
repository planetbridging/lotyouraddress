package main

import (
	"math"
	"strconv"
)

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {
	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist
}

/*
func main() {
	//145.124738|-38.162404 to 143.95144633|-36.3316845
	fmt.Printf("%f Miles\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "M"))
	fmt.Printf("%f Kilometers\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "K"))
	fmt.Printf("%f test Kilometers\n", distance(145.124738, -38.162404, 143.95144633, -36.3316845, "K"))
	fmt.Printf("%f Nautical Miles\n", distance(32.9697, -96.80322, 29.46786, -98.53506, "N"))
}
*/

func getDistance(lat1 string, long1 string, lat2 string, long2 string) float64 {
	lat_1, _ := strconv.ParseFloat(lat1, 32)
	long_1, _ := strconv.ParseFloat(long1, 32)

	lat_2, _ := strconv.ParseFloat(lat2, 32)
	long_2, _ := strconv.ParseFloat(long2, 32)

	//fmt.Printf("%f test Kilometers\n", distance(lat_1, long_1, lat_2, long_2, "K"))
	return distance(lat_1, long_1, lat_2, long_2, "K")
}

func sortClosestSuburb(lstContents []ObjDistance) []ObjDistance {

	var tmplstContents []ObjDistance

	var smallest float64

	for suburbs1, _ := range lstContents {
		_ = suburbs1
		smallest = 9999999999999999999999999999
		selected := -1
		for suburbs2, _ := range lstContents {
			if lstContents[suburbs2].Km <= smallest && !lstContents[suburbs2].sorted {
				smallest = lstContents[suburbs2].Km
				selected = suburbs2
			}
		}
		if selected != -1 {
			tmplstContents = append(tmplstContents, lstContents[selected])
			lstContents[selected].sorted = true
		}
		//fmt.Println(suburbs1+1, ":", len(lstContents))
	}
	return tmplstContents
}
