package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
)

func importLya(path string) {
	fmt.Println("starting importing states")
	importStates(path)
	fmt.Println("finished importing states")
	fmt.Println("starting importing postcodes")
	importPostcodes(path)
	fmt.Println("finished importing postcodes")
	fmt.Println("starting importing suburbs")
	importSuburbs(path)
	fmt.Println("finished importing suburbs")
	fmt.Println("starting importing addresses")
	importAddresses()
	fmt.Println("finished importing addresses")
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

func importAddresses() {
	addressesPath := currentPath + "G-NAF/G-NAF MAY 2021/Standard/"
	streets := addressesPath + "export/streets/"
	if !folderExists(streets) {
		//createFolder(streets)
		for s, _ := range lstObjStateLya {
			fmt.Println("importing addresses state ", lstObjStateLya[s].State_Abbr)
			//state_path := streets + lstObjStateLya[s].State_Abbr + "/"
			//createFolder(state_path)
			tmp := readLocal(addressesPath + lstObjStateLya[s].State_Abbr + "_ADDRESS_DETAIL_psv.psv")
			for i := 1; i < len(tmp); i++ {
				//cleaned := strings.ReplaceAll(tmp[i], " ", "")
				row := strings.Split(tmp[i], "|")
				tmp_ObjAddressesLya := convertStreetRowToObj(row)
				for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {
					if row[26] == lstObjStateLya[s].LstObjPostcodeLya[p].Postcode {
						for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {
							if row[24] == lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID {
								for r, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya {
									if row[22] == lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_LOCALITY_PID {
										lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Addresses = append(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Addresses, tmp_ObjAddressesLya)
										break
									}
								}
								break
							}
						}
						break
					}

				}

			}
			fmt.Println("finished importing addresses state ", lstObjStateLya[s].State_Abbr)
		}
	}
}

func convertStreetRowToObj(row []string) ObjAddressesLya {
	/*tmp_ObjAddressesLya := ObjAddressesLya{
		ADDRESS_DETAIL_PID: row[0],
		//DATE_CREATED|DATE_LAST_MODIFIED|DATE_RETIRED|
		BUILDING_NAME:        row[4],
		LOT_NUMBER_PREFIX:    row[5],
		LOT_NUMBER:           row[6],
		LOT_NUMBER_SUFFIX:    row[7],
		FLAT_TYPE_CODE:       row[8],
		FLAT_NUMBER_PREFIX:   row[9],
		FLAT_NUMBER:          row[10],
		FLAT_NUMBER_SUFFIX:   row[11],
		LEVEL_TYPE_CODE:      row[12],
		LEVEL_NUMBER_PREFIX:  row[13],
		LEVEL_NUMBER:         row[14],
		LEVEL_NUMBER_SUFFIX:  row[15],
		NUMBER_FIRST_PREFIX:  row[16],
		NUMBER_FIRST:         row[17],
		NUMBER_FIRST_SUFFIX:  row[18],
		NUMBER_LAST_PREFIX:   row[19],
		NUMBER_LAST:          row[20],
		NUMBER_LAST_SUFFIX:   row[21],
		STREET_LOCALITY_PID:  row[22],
		LOCATION_DESCRIPTION: row[23],
		LOCALITY_PID:         row[24],
		ALIAS_PRINCIPAL:      row[25],
		//POSTCODE|
		PRIVATE_STREET:      row[27],
		LEGAL_PARCEL_ID:     row[28],
		CONFIDENCE:          row[29],
		ADDRESS_SITE_PID:    row[30],
		LEVEL_GEOCODED_CODE: row[31],
		PROPERTY_PID:        row[32],
		GNAF_PROPERTY_PID:   row[33],
		PRIMARY_SECONDARY:   row[34],
	}*/
	tmp_ObjAddressesLya := ObjAddressesLya{
		ADDRESS_DETAIL_PID: row[0],
		//DATE_CREATED|DATE_LAST_MODIFIED|DATE_RETIRED|
		BUILDING_NAME: row[4],
		//LOT_NUMBER_PREFIX:    row[5],
		LOT_NUMBER: row[6],
		//LOT_NUMBER_SUFFIX:    row[7],
		//FLAT_TYPE_CODE:       row[8],
		//FLAT_NUMBER_PREFIX:   row[9],
		FLAT_NUMBER: row[10],
		//FLAT_NUMBER_SUFFIX:   row[11],
		LEVEL_TYPE_CODE: row[12],
		//LEVEL_NUMBER_PREFIX:  row[13],
		LEVEL_NUMBER: row[14],
		//LEVEL_NUMBER_SUFFIX:  row[15],
		//NUMBER_FIRST_PREFIX:  row[16],
		NUMBER_FIRST: row[17],
		//NUMBER_FIRST_SUFFIX:  row[18],
		//NUMBER_LAST_PREFIX:   row[19],
		NUMBER_LAST: row[20],
		//NUMBER_LAST_SUFFIX:   row[21],
		//STREET_LOCALITY_PID:  row[22],
		//LOCATION_DESCRIPTION: row[23],
		//LOCALITY_PID:         row[24],
		//ALIAS_PRINCIPAL:      row[25],
		//POSTCODE|
		//PRIVATE_STREET:      row[27],
		//LEGAL_PARCEL_ID:     row[28],
		//CONFIDENCE:          row[29],
		//ADDRESS_SITE_PID:    row[30],
		//LEVEL_GEOCODED_CODE: row[31],
		//PROPERTY_PID:        row[32],
		//GNAF_PROPERTY_PID:   row[33],
		//PRIMARY_SECONDARY:   row[34],
	}
	return tmp_ObjAddressesLya
}
