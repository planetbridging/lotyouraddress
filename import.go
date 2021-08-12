package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func importLyua(path string) {
	importStates(path)
	fmt.Println("finished importing states")
	importPostcodes(path)
	fmt.Println("finished importing postcodes")
	importSuburbs(path)
	fmt.Println("finished importing suburbs")
}

func importStates(path string) {
	tmp := readLocal(path + "STATES.csv")
	for i := 1; i < len(tmp); i++ {
		//cleaned := strings.ReplaceAll(tmp[i], " ", "")
		row := strings.Split(tmp[i], ",")

		tmp_ObjStateLya := ObjStateLya{
			State_Name: row[0],
			State_Abbr: row[1],
		}
		lstObjStateLya = append(lstObjStateLya, tmp_ObjStateLya)
	}
}

func importPostcodes(path string) {

	for s, _ := range lstStaticStates {
		tmp_path := path + "postcodes/" + lstStaticStates[s] + "/"
		tmp := getFiles(tmp_path)
		for o, _ := range lstObjStateLya {
			if lstStaticStates[s] == lstObjStateLya[o].State_Abbr {

				for p, _ := range tmp {
					pc := strings.Replace(tmp[p], ".csv", "", -1)
					lsttmp := readLocal(tmp_path + tmp[p])
					var tmpLstObjSuburbLya []ObjSuburbLya
					for i := 1; i < len(lsttmp); i++ {
						row := strings.Split(lsttmp[i], ",")

						tmp_ObjSuburbLya := ObjSuburbLya{
							Suburb:       row[1],
							LONGITUDE:    row[2],
							LATITUDE:     row[3],
							LOCALITY_PID: row[0],
						}
						tmpLstObjSuburbLya = append(tmpLstObjSuburbLya, tmp_ObjSuburbLya)
					}

					tmp_ObjPostcodeLya := ObjPostcodeLya{
						Postcode:        pc,
						LstObjSuburbLya: tmpLstObjSuburbLya,
					}

					lstObjStateLya[o].LstObjPostcodeLya = append(lstObjStateLya[o].LstObjPostcodeLya, tmp_ObjPostcodeLya)
				}

				break
			}
		}
	}
}

func importSuburbs(path string) {
	for s, _ := range lstObjStateLya {
		state_folder := path + "suburbs/" + lstObjStateLya[s].State_Abbr + "/"
		createFolder(state_folder)
		for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {
			postcode_folder := state_folder + lstObjStateLya[s].LstObjPostcodeLya[p].Postcode + "/"
			var wg sync.WaitGroup
			sliceLength := len(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya)
			wg.Add(sliceLength)
			for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {

				go func(s int, p int, u int) {
					defer wg.Done()
					id := lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID
					suburb_file := postcode_folder + id + "_" + id + ".csv"
					tmp := readLocal(suburb_file)
					for i := 1; i < len(tmp); i++ {
						row := strings.Split(tmp[i], ",")
						//STREET_LOCALITY_PID,STREET_NAME,STREET_TYPE_CODE,LONGITUDE,LATITUDE,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite
						Fixed_wireless, _ := strconv.ParseBool(row[5])
						FTTB, _ := strconv.ParseBool(row[6])
						FTTDP_FTTC, _ := strconv.ParseBool(row[7])
						FTTN, _ := strconv.ParseBool(row[8])
						FTTP, _ := strconv.ParseBool(row[9])
						HFC, _ := strconv.ParseBool(row[10])
						Satellite, _ := strconv.ParseBool(row[11])

						tmp_ObjApiInternetType := ObjApiInternetType{
							Fixed_wireless: Fixed_wireless,
							FTTB:           FTTB,
							FTTDP_FTTC:     FTTDP_FTTC,
							FTTN:           FTTN,
							FTTP:           FTTP,
							HFC:            HFC,
							Satellite:      Satellite,
						}
						tmp_ObjStreetsLya := ObjStreetsLya{
							STREET_LOCALITY_PID: row[0],
							STREET_NAME:         row[1],
							STREET_TYPE_CODE:    row[2],
							LONGITUDE:           row[3],
							LATITUDE:            row[4],
							Selected_Internet:   tmp_ObjApiInternetType,
						}
						lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya = append(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya, tmp_ObjStreetsLya)
					}
				}(s, p, u)

			}
			wg.Wait()

		}
	}
}
