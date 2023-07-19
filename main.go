package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var currentPath string
var currentInternetPath string
var lstStaticStates []string

var lstObjStateLya []ObjStateLya

var tmp_geo []string

var internet_geo []ObjInternetType

var root = flag.String("root", ".", "file system path")

func main() {
	fmt.Println("Welcome to lya")

	lstStaticStates = append(lstStaticStates, "ACT")
	lstStaticStates = append(lstStaticStates, "NSW")
	lstStaticStates = append(lstStaticStates, "NT")
	lstStaticStates = append(lstStaticStates, "OT")
	lstStaticStates = append(lstStaticStates, "QLD")
	lstStaticStates = append(lstStaticStates, "SA")
	lstStaticStates = append(lstStaticStates, "TAS")
	lstStaticStates = append(lstStaticStates, "VIC")
	lstStaticStates = append(lstStaticStates, "WA")

	//currentPath = "C:\\Users\\plane\\OneDrive\\Documents\\development\\data\\addy\\MAY21_GNAF_PipeSeparatedValue\\"
	//currentPath = "/mnt/c/Users/plane/OneDrive/Documents/development/data/addy/MAY21_GNAF_PipeSeparatedValue/"
	cp := readLocal(".currentPath")

	//exportPath := os.Getenv("PATH")
	fmt.Println("loading from" + cp[0])
	importLya(cp[0])
	/*if folderExists(exportPath) {
		importLya(exportPath)
	} else {
		fmt.Println("processing and exporting?")
		loadStates()

		loadPostcodes()

		loadStreets()

		checkPreprocessedInternet()

		setupLocalStreetAccess()
	}*/
	//fmt.Println(cp[0])

	//internetToStreetTesting()
	//loadStreetGeo()
	//loadStreetGeo()
	//loadInternet()
	//loadInternetSpeeds()
	//loadSuburbs()

	/*findFrankston := findSuburb("FRANKSTON")
	fmt.Println(findFrankston)

	fmt.Println(lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0])
	fmt.Println(lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[1])
	//lat1 string, long1 string, lat2 string, long2 string
	getDistance(
		lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0].LATITUDE,
		lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0].LONGITUDE,
		lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[1].LATITUDE,
		lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[1].LONGITUDE,
	)*/

	/*getsubs := getSuburbDistance(
	lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0].LATITUDE,
	lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0].LONGITUDE)*/

	//get suburbs
	/*getsubs := getSuburbDistance(
		findFrankston[1].LATITUDE,
		findFrankston[1].LONGITUDE)
	sorted := sortClosestSuburb(getsubs)
	fmt.Println("Print top 14 closes to frankston south")
	for i := 0; i < 50; i++ {
		fmt.Println(sorted[i])
	}*/
	internet_geo = nil

	//exportLya()

	//runtime.GC()
	fmt.Println("Lotyouraddress running")
	http.HandleFunc("/", handlerFunc)
	http.Handle("/web/", http.FileServer(http.Dir(*root)))
	http.ListenAndServe(":3001", nil)
}

func exportLya() {
	tmp_folder_path := currentPath + "G-NAF/G-NAF MAY 2021/Standard/"
	tmp_full_local_export := tmp_folder_path + "export/"
	createFolder(tmp_full_local_export)

	tmp_full_postcodes_export := tmp_folder_path + "export/postcodes/"
	createFolder(tmp_full_postcodes_export)

	tmp_full_suburbs_export := tmp_folder_path + "export/suburbs/"
	createFolder(tmp_full_suburbs_export)

	states_local := tmp_full_local_export + "STATES.csv"
	if !fileExists(states_local) {
		exportStates(states_local)
	}

	exportPostcodes(tmp_full_postcodes_export)
	exportSuburbs(tmp_full_suburbs_export)
}

/*STREET_NAME         string
STREET_TYPE_CODE    string
LONGITUDE           string
LATITUDE            string
STREET_LOCALITY_PID string
Data                []string
internetType        int
local_pid           string
Selected_Internet   ObjApiInternetType*/

//internet

/*Fixed_wireless bool
FTTB           bool
FTTDP_FTTC     bool
FTTN           bool
FTTP           bool
HFC            bool
Satellite      bool*/

func exportStreets(path string) {
	fmt.Println("exporting suburbs")
	for s, _ := range lstObjStateLya {
		state_folder := path + lstObjStateLya[s].State_Abbr + "/"
		createFolder(state_folder)
		for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {
			postcode_folder := state_folder + lstObjStateLya[s].LstObjPostcodeLya[p].Postcode + "/"
			createFolder(postcode_folder)

			for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {
				save_path := postcode_folder + lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID + "_"
				save_path += strings.ReplaceAll(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID+".csv", " ", "")
				if !fileExists(save_path) {
					var data []string
					title := "STREET_LOCALITY_PID,STREET_NAME,STREET_TYPE_CODE,LONGITUDE,LATITUDE,"
					title += "Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite"
					data = append(data, title)
					for r, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya {
						row := lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_LOCALITY_PID + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_NAME + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_TYPE_CODE + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].LONGITUDE + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].LATITUDE + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.Fixed_wireless) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTB) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTDP_FTTC) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTN) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTP) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.HFC) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.Satellite)
						data = append(data, row)
					}
					saveFile(data, save_path)
				}
			}

		}
	}
}

func exportSuburbs(path string) {
	fmt.Println("exporting suburbs")
	for s, _ := range lstObjStateLya {
		state_folder := path + lstObjStateLya[s].State_Abbr + "/"
		createFolder(state_folder)
		for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {
			postcode_folder := state_folder + lstObjStateLya[s].LstObjPostcodeLya[p].Postcode + "/"
			createFolder(postcode_folder)

			for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {
				save_path := postcode_folder + lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID + "_"
				save_path += strings.ReplaceAll(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID+".csv", " ", "")
				if !fileExists(save_path) {
					var data []string
					title := "STREET_LOCALITY_PID,STREET_NAME,STREET_TYPE_CODE,LONGITUDE,LATITUDE,"
					title += "Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite"
					data = append(data, title)
					for r, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya {
						row := lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_LOCALITY_PID + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_NAME + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].STREET_TYPE_CODE + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].LONGITUDE + ","
						row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].LATITUDE + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.Fixed_wireless) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTB) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTDP_FTTC) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTN) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.FTTP) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.HFC) + ","
						row += strconv.FormatBool(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya[r].Selected_Internet.Satellite)
						data = append(data, row)
					}
					saveFile(data, save_path)
				}
			}

		}
	}
}

func exportPostcodes(path string) {
	fmt.Println("exporting postcodes")
	for s, _ := range lstObjStateLya {
		state_folder := path + lstObjStateLya[s].State_Abbr + "/"
		createFolder(state_folder)
		for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {
			save_path := state_folder + lstObjStateLya[s].LstObjPostcodeLya[p].Postcode + ".csv"
			if !fileExists(save_path) {
				var data []string
				data = append(data, "LOCALITY_PID,Suburb,LONGITUDE,LATITUDE,Street_Count")
				for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {
					street_count := len(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya)
					row := lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LOCALITY_PID + ","
					row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].Suburb + ","
					row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LONGITUDE + ","
					row += lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LATITUDE + ","
					row += strconv.Itoa(street_count)
					data = append(data, row)
				}

				saveFile(data, save_path)
			}

		}
	}
}

func exportStates(path string) {
	fmt.Println("exporting states")
	var data []string
	data = append(data, "State,Abbr,Postcode_count,Suburb_count,Street_Count")

	for s, _ := range lstObjStateLya {
		postcode_count := 0
		suburb_count := 0
		street_count := 0
		postcode_count += len(lstObjStateLya[s].LstObjPostcodeLya)
		for p, _ := range lstObjStateLya[s].LstObjPostcodeLya {

			suburb_count += len(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya)

			for u, _ := range lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya {
				street_count += len(lstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[u].LstObjStreetsLya)
			}
		}
		row := lstObjStateLya[s].State_Name + "," + lstObjStateLya[s].State_Abbr
		row += "," + strconv.Itoa(postcode_count)
		row += "," + strconv.Itoa(suburb_count)
		row += "," + strconv.Itoa(street_count)
		data = append(data, row)
	}

	saveFile(data, path)
}

func setupLocalStreetAccess() {
	tmp_folder_path := currentPath + "G-NAF/G-NAF MAY 2021/Standard/"
	tmp_full_local_street := tmp_folder_path + "StreetAccess/"
	createFolder(tmp_full_local_street)
	if fileExists(tmp_full_local_street) {
		fmt.Println("setup ready")
	} else {
		createLocalStreetFolders(tmp_full_local_street)
	}

	fmt.Println("setting up local street access")
	createLocalStreetFiles(tmp_full_local_street, tmp_folder_path)
}

func createLocalStreetFolders(tmp_full_local_street string) {
	for states, _ := range lstObjStateLya {
		state_path := tmp_full_local_street + lstObjStateLya[states].State_Abbr + "/"
		createFolder(state_path)
		sliceLength := len(lstObjStateLya[states].LstObjPostcodeLya)
		var wg sync.WaitGroup
		wg.Add(sliceLength)
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {

			go func(postcodes int) {
				defer wg.Done()
				postcode_path := state_path + lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode + "/"
				createFolder(postcode_path)
				for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
					suburb_path := postcode_path + lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID + "/"
					createFolder(suburb_path)
				}
			}(postcodes)

		}
		wg.Wait()
	}
}

func createLocalStreetFiles(tmp_full_local_street string, tmp_folder_path string) {
	for states, _ := range lstObjStateLya {
		state_path := tmp_full_local_street + lstObjStateLya[states].State_Abbr + "/"
		read_state := readLocal(tmp_folder_path + lstObjStateLya[states].State_Abbr + "_ADDRESS_DETAIL_psv.psv")
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			postcode_path := state_path + lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode + "/"

			sliceLength := len(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya)
			var wg sync.WaitGroup
			wg.Add(sliceLength)

			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				go func(states int, postcodes int, suburbs int) {
					defer wg.Done()
					suburb_path := postcode_path + lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID + "/"
					_ = suburb_path
					//fmt.Println(suburb_path)
					for i := 1; i < len(read_state); i++ {
					}
				}(states, postcodes, suburbs)

			}

			wg.Wait()

		}
		read_state = nil
		//runtime.GC()
		//debug.FreeOSMemory()
	}
}

func worker(wg *sync.WaitGroup, states int, postcodes int, suburbs int, postcode_path string, read_state []string) {
	defer wg.Done()
	suburb_path := postcode_path + lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID + "/"
	fmt.Println(suburb_path)
	for i := 1; i < len(read_state); i++ {
	}
}

func findSuburb(name string, limiter int) []ObjSuburbLya {
	var lstSuburbs []ObjSuburbLya
	var count int
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				tmp_name := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb
				if strings.Contains(tmp_name, name) {
					lstSuburbs = append(lstSuburbs, lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
					count += 1
					if count >= limiter {
						return lstSuburbs
					}
				}
			}
		}
	}
	return lstSuburbs
}

func findPIDSuburb(name string) []ObjSuburbLya {
	var lstSuburbs []ObjSuburbLya
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				/*tmp_name := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb
				if strings.Contains(tmp_name, name) {
					lstSuburbs = append(lstSuburbs, lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
				}*/
				if name == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
					//lstSuburbs = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs]

					lstSuburbs = append(lstSuburbs, lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
					return lstSuburbs
				}
			}
		}
	}
	return lstSuburbs
}

func getSuburbDistance(lat1 string, long1 string) []ObjDistance {
	lat_2, _ := strconv.ParseFloat(lat1, 32)
	long_2, _ := strconv.ParseFloat(long1, 32)
	var lstContents []ObjDistance
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				tmp_lat, _ := strconv.ParseFloat(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LATITUDE, 32)
				tmp_long, _ := strconv.ParseFloat(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LONGITUDE, 32)
				km := distance(lat_2, long_2, tmp_lat, tmp_long, "K")
				tmp_ObjDistance := ObjDistance{
					Km:           km,
					Suburb:       lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb,
					LOCALITY_PID: lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID,
					sorted:       false,
					Postcode:     lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode,
					State:        lstObjStateLya[states].State_Abbr,
				}
				lstContents = append(lstContents, tmp_ObjDistance)
			}
		}
	}
	//fmt.Println(len(lstContents))
	return lstContents
}

//examples

//address
//vic_address_detail_psv.psv
//ADDRESS_DETAIL_PID|DATE_CREATED|DATE_LAST_MODIFIED|DATE_RETIRED|BUILDING_NAME|LOT_NUMBER_PREFIX|LOT_NUMBER|LOT_NUMBER_SUFFIX|FLAT_TYPE_CODE|FLAT_NUMBER_PREFIX|FLAT_NUMBER|FLAT_NUMBER_SUFFIX|LEVEL_TYPE_CODE|LEVEL_NUMBER_PREFIX|LEVEL_NUMBER|LEVEL_NUMBER_SUFFIX|NUMBER_FIRST_PREFIX|NUMBER_FIRST|NUMBER_FIRST_SUFFIX|NUMBER_LAST_PREFIX|NUMBER_LAST|NUMBER_LAST_SUFFIX|STREET_LOCALITY_PID|LOCATION_DESCRIPTION|LOCALITY_PID|ALIAS_PRINCIPAL|POSTCODE|PRIVATE_STREET|LEGAL_PARCEL_ID|CONFIDENCE|ADDRESS_SITE_PID|LEVEL_GEOCODED_CODE|PROPERTY_PID|GNAF_PROPERTY_PID|PRIMARY_SECONDARY
//GAVIC420168306|2004-04-29|2021-05-20|||||||||||||||61|||||VIC2021100||VIC941|P|3199||43\LP56908|2|420304753|7||1141578|

//geo
//vic_address_detail_geocode_psv.psv
//ADDRESS_DEFAULT_GEOCODE_PID|DATE_CREATED|DATE_RETIRED|ADDRESS_DETAIL_PID|GEOCODE_TYPE_CODE|LONGITUDE|LATITUDE
//2163875|2012-11-01||GAVIC420168306|FCS|145.124738|-38.162404

//street
//VIC_ADDRESS_DETAIL_psv.psv
//STREET_LOCALITY_PID|DATE_CREATED|DATE_RETIRED|STREET_CLASS_CODE|STREET_NAME|STREET_TYPE_CODE|STREET_SUFFIX_CODE|LOCALITY_PID|GNAF_STREET_PID|GNAF_STREET_CONFIDENCE|GNAF_RELIABILITY_CODE
//VIC1982231|2017-11-01||C|LAWSON|AVENUE||VIC941|253044441|2|4

//suburb
//LOCALITY_PID|DATE_CREATED|DATE_RETIRED|LOCALITY_NAME|PRIMARY_POSTCODE|LOCALITY_CLASS_CODE|STATE_PID|GNAF_LOCALITY_PID|GNAF_RELIABILITY_CODE
//VIC941|2012-04-27||FRANKSTON SOUTH||G|2|250184905|5

func loadStates() {

	//C:\Users\plane\OneDrive\Documents\development\data\addy\MAY21_GNAF_PipeSeparatedValue\G-NAF\G-NAF MAY 2021\Standard
	for index, _ := range lstStaticStates {
		fmt.Println("adding: " + lstStaticStates[index])
		item := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstStaticStates[index] + "_STATE_psv.psv")
		tmp_line := strings.Split(item[1], "|")
		//fmt.Println(tmp_line[3])
		tmp_ObjStateLya := ObjStateLya{
			State_Name: strings.TrimSpace(tmp_line[3]),
			State_Abbr: strings.TrimSpace(tmp_line[4]),
		}
		lstObjStateLya = append(lstObjStateLya, tmp_ObjStateLya)
		item = nil
		//runtime.GC()
	}
}

func loadPostcodes() {

	for index, _ := range lstObjStateLya {
		fmt.Println("loading suburbs in : " + lstObjStateLya[index].State_Name)
		lst_tmp_ObjSuburbLya := loadSuburbs(lstObjStateLya[index].State_Abbr)
		tmp_saved_postcodes := currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_POSTCODES.csv"
		check_postcodes := fileExists(tmp_saved_postcodes)
		if !check_postcodes {
			postcodes := readFindAllPostcodes(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_ADDRESS_DETAIL_psv.psv")
			saveFile(postcodes, tmp_saved_postcodes)
			//fmt.Println(postcodes)
		}

		load_postcodes := readLocal(tmp_saved_postcodes)
		var unique_postcodes []string
		for lp, _ := range load_postcodes {
			if strings.Contains(load_postcodes[lp], ",") {
				tmp_line := strings.Split(load_postcodes[lp], ",")
				found := contains(unique_postcodes, tmp_line[0])
				if !found {
					unique_postcodes = append(unique_postcodes, tmp_line[0])
				}
			}
		}

		for up, _ := range unique_postcodes {
			var tmp_ObjSuburbLya []ObjSuburbLya
			for tlp, _ := range load_postcodes {
				tmp_line2 := strings.Split(load_postcodes[tlp], ",")
				if tmp_line2[0] == unique_postcodes[up] {
					for ltosub, _ := range lst_tmp_ObjSuburbLya {
						if lst_tmp_ObjSuburbLya[ltosub].LOCALITY_PID == tmp_line2[1] {
							tmp_ObjSuburbLya = append(tmp_ObjSuburbLya, lst_tmp_ObjSuburbLya[ltosub])
							break
						}
					}
				}
			}
			tmp_ObjPostcodeLya := ObjPostcodeLya{
				Postcode:        unique_postcodes[up],
				LstObjSuburbLya: tmp_ObjSuburbLya,
			}
			lstObjStateLya[index].LstObjPostcodeLya = append(lstObjStateLya[index].LstObjPostcodeLya, tmp_ObjPostcodeLya)
		}

		load_postcodes = nil
		//runtime.GC()
	}
}

func loadSuburbs(state string) []ObjSuburbLya {
	var lst_tmp_ObjSuburbLya []ObjSuburbLya
	tmp_sub_names := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + state + "_LOCALITY_psv.psv")
	tmp_sub_geos := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + state + "_LOCALITY_POINT_psv.psv")
	for tsn, _ := range tmp_sub_names {
		tmp_line1 := strings.Split(tmp_sub_names[tsn], "|")
		for tsg, _ := range tmp_sub_geos {
			tmp_line2 := strings.Split(tmp_sub_geos[tsg], "|")
			if tmp_line1[0] == tmp_line2[3] {
				tmp_ObjSuburbLya := ObjSuburbLya{
					LOCALITY_PID: tmp_line1[0],
					LONGITUDE:    tmp_line2[5],
					LATITUDE:     tmp_line2[6],
					Suburb:       tmp_line1[3],
				}
				lst_tmp_ObjSuburbLya = append(lst_tmp_ObjSuburbLya, tmp_ObjSuburbLya)
				break
			}
		}
	}
	fmt.Println("load suburbs done")
	tmp_sub_geos = nil
	tmp_sub_names = nil
	runtime.GC()
	return lst_tmp_ObjSuburbLya
}

//street
//STREET_LOCALITY_PID|DATE_CREATED|DATE_RETIRED|STREET_CLASS_CODE|STREET_NAME|STREET_TYPE_CODE|STREET_SUFFIX_CODE|LOCALITY_PID|GNAF_STREET_PID|GNAF_STREET_CONFIDENCE|GNAF_RELIABILITY_CODE
//TAS3349313|2018-05-05||C|EMMETT|STREET||TAS487|502237149|2|4
/*func loadStreets() {

	for index, _ := range lstObjStateLya {
		tmp_saved_path := currentPath + "G-NAF/G-NAF MAY 2021/Standard/geo_internet_" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_psv.psv"
		check_geo_internet_streets := fileExists(tmp_saved_path)
		if !check_geo_internet_streets {
			tmp_streets := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_psv.psv")
			tmp_streets_geo := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_POINT_psv.psv")
			fmt.Println("Creating tmp lst for state " + lstObjStateLya[index].State_Abbr)
			var lst_tmp_ObjStreetsLya []ObjStreetsLya
			for i := 1; i < len(tmp_streets); i++ {
				if strings.Contains(tmp_streets[i], "|") {
					tmp_line := strings.Split(tmp_streets[i], "|")
					STREET_LOCALITY_PID := tmp_line[0]
					STREET_NAME := tmp_line[4]
					STREET_TYPE_CODE := tmp_line[5]
					//INTERNET_TYPE := ""
					//LONGITUDE := ""
					//LATITUDE := ""
					tmp_ObjStreetsLya := ObjStreetsLya{
						STREET_LOCALITY_PID: STREET_LOCALITY_PID,
						STREET_NAME:         STREET_NAME,
						STREET_TYPE_CODE:    STREET_TYPE_CODE,
						LONGITUDE:           "",
						LATITUDE:            "",
						InternetType:        "",
					}
					lst_tmp_ObjStreetsLya = append(lst_tmp_ObjStreetsLya, tmp_ObjStreetsLya)
				}
			}

			fmt.Println("Appending geo locations to streets")
			oldTime := time.Now()
			sliceLength := len(tmp_streets_geo)
			var wg sync.WaitGroup
			wg.Add(sliceLength - 1)
			for i := 1; i < len(tmp_streets_geo); i++ {
				if strings.Contains(tmp_streets_geo[i], "|") {
					go func(i int) {
						defer wg.Done()
						tmp_line := strings.Split(tmp_streets_geo[i], "|")
						STREET_LOCALITY_PID := tmp_line[3]
						LONGITUDE := tmp_line[6]
						LATITUDE := tmp_line[7]
						for tmp_lst, _ := range lst_tmp_ObjStreetsLya {
							if STREET_LOCALITY_PID == lst_tmp_ObjStreetsLya[tmp_lst].STREET_LOCALITY_PID {
								lst_tmp_ObjStreetsLya[tmp_lst].LATITUDE = LATITUDE
								lst_tmp_ObjStreetsLya[tmp_lst].LONGITUDE = LONGITUDE
								break
							}
						}
						//updateGeoOnStreet(index, STREET_LOCALITY_PID, LATITUDE, LONGITUDE)
					}(i)

				}
			}
			wg.Wait()

			fmt.Println("setup internet types on streets")

			var internet_wg sync.WaitGroup
			internetSliceLength := len(lst_tmp_ObjStreetsLya)
			internet_wg.Add(internetSliceLength)
			for i := 0; i < len(lst_tmp_ObjStreetsLya); i++ {
				go func(i int) {
					defer internet_wg.Done()
					//lon string, lat string
					tmp_internet_type := getInternetType(lst_tmp_ObjStreetsLya[i].LONGITUDE, lst_tmp_ObjStreetsLya[i].LATITUDE)
					lst_tmp_ObjStreetsLya[i].InternetType = tmp_internet_type.InternetType
					//updateGeoOnStreet(index, STREET_LOCALITY_PID, LATITUDE, LONGITUDE)
				}(i)

			}
			internet_wg.Wait()

			currentTime := time.Now()
			diff := currentTime.Sub(oldTime)
			//In hours
			fmt.Printf("Hours: %f\n", diff.Hours())

			//In minutes
			fmt.Printf("Minutes: %f\n", diff.Minutes())

			//In seconds
			fmt.Printf("Seconds: %f\n", diff.Seconds())

			//In nanoseconds
			fmt.Printf("Nanoseconds: %d\n", diff.Nanoseconds())

			*var tmp_save_street_interet_geo_data []string
			tmp_save_street_interet_geo_data = append(tmp_save_street_interet_geo_data, "STREET_LOCALITY_PID,STREET_NAME,STREET_TYPE_CODE,INTERNET_TYPE,LONGITUDE,LATITUDE")
			for i := 1; i < len(tmp_streets); i++ {
				if strings.Contains(tmp_streets[i], "|") {
					tmp_line := strings.Split(tmp_streets[i], "|")
					STREET_LOCALITY_PID := tmp_line[0]
					STREET_NAME := tmp_line[4]
					STREET_TYPE_CODE := tmp_line[5]
					INTERNET_TYPE := ""
					LONGITUDE := ""
					LATITUDE := ""
					for i_geo := 1; i_geo < len(tmp_streets_geo); i_geo++ {
						tmp_line_geo := strings.Split(tmp_streets_geo[i_geo], "|")
						if tmp_line_geo[3] == tmp_line[0] {
							LATITUDE = tmp_line_geo[6]
							LONGITUDE = tmp_line_geo[7]
							tmp_streets_geo = remove(tmp_streets_geo, i_geo)
							break
						}
					}
					INTERNET_TYPE = getInternetType(LONGITUDE, LATITUDE).InternetType
					tmp_save_street_interet_geo_data = append(tmp_save_street_interet_geo_data, STREET_LOCALITY_PID+"|"+STREET_NAME+"|"+STREET_TYPE_CODE+"|"+INTERNET_TYPE+"|"+LONGITUDE+"|"+LATITUDE)
				}
			}
			saveFile(tmp_save_street_interet_geo_data, tmp_saved_path)*
		}

		break
	}
}*/

func loadStreets() {

	for index, _ := range lstObjStateLya {
		//oldTime := time.Now()
		fmt.Println("loading streets in : " + lstObjStateLya[index].State_Name)
		tmp_saved_path := currentPath + "G-NAF/G-NAF MAY 2021/Standard/geo_" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_psv.psv"
		check_geo_streets := fileExists(tmp_saved_path)
		if !check_geo_streets {
			generateGeoStreets(index, tmp_saved_path)
		}

		tmp_streets_geo := readLocal(tmp_saved_path)

		//"STREET_LOCALITY_PID,LONGITUDE,LATITUDE,STREET_NAME,STREET_TYPE_CODE,LOCALITY_PID"
		for i := 1; i < len(tmp_streets_geo); i++ {
			if strings.Contains(tmp_streets_geo[i], ",") {
				tmp_line := strings.Split(tmp_streets_geo[i], ",")
				tmp_ObjStreetsLya := ObjStreetsLya{
					STREET_LOCALITY_PID: tmp_line[0],
					STREET_NAME:         tmp_line[3],
					STREET_TYPE_CODE:    tmp_line[4],
					local_pid:           tmp_line[5],
					internetType:        -1,
					LONGITUDE:           tmp_line[1],
					LATITUDE:            tmp_line[2],
				}
				found := false
				for p, _ := range lstObjStateLya[index].LstObjPostcodeLya {
					for s, _ := range lstObjStateLya[index].LstObjPostcodeLya[p].LstObjSuburbLya {
						//lstObjStateLya[index].LstObjPostcodeLya[p].LstObjSuburbLya
						if tmp_ObjStreetsLya.local_pid == lstObjStateLya[index].LstObjPostcodeLya[p].LstObjSuburbLya[s].LOCALITY_PID {
							lstObjStateLya[index].LstObjPostcodeLya[p].LstObjSuburbLya[s].LstObjStreetsLya = append(lstObjStateLya[index].LstObjPostcodeLya[p].LstObjSuburbLya[s].LstObjStreetsLya, tmp_ObjStreetsLya)
							found = true
							break
						}
					}
					if found {
						break
					}
				}
			}

		}

		tmp_streets_geo = nil
		//runtime.GC()

		/*currentTime := time.Now()
		diff := currentTime.Sub(oldTime)
		//In hours
		fmt.Printf("Hours: %f\n", diff.Hours())

		//In minutes
		fmt.Printf("Minutes: %f\n", diff.Minutes())

		//In seconds
		fmt.Printf("Seconds: %f\n", diff.Seconds())

		//In nanoseconds
		fmt.Printf("Nanoseconds: %d\n", diff.Nanoseconds())*/
	}
	//fmt.Println(lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0])
}
func generateGeoStreets(index int, save string) {
	tmp_streets := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_psv.psv")
	var lst_tmp_streets_noGeo []ObjStreetsLya
	for i := 1; i < len(tmp_streets); i++ {
		if strings.Contains(tmp_streets[i], "|") {
			tmp_line := strings.Split(tmp_streets[i], "|")
			STREET_LOCALITY_PID := tmp_line[0]
			STREET_NAME := tmp_line[4]
			STREET_TYPE_CODE := tmp_line[5]
			LOCALITY_PID := tmp_line[7]

			tmp_ObjStreetsLya := ObjStreetsLya{
				STREET_LOCALITY_PID: STREET_LOCALITY_PID,
				STREET_NAME:         STREET_NAME,
				STREET_TYPE_CODE:    STREET_TYPE_CODE,
				local_pid:           LOCALITY_PID,
				internetType:        -1,
				LONGITUDE:           "",
				LATITUDE:            "",
			}
			lst_tmp_streets_noGeo = append(lst_tmp_streets_noGeo, tmp_ObjStreetsLya)
		}
	}

	fmt.Println("updating geo streets in state: " + lstObjStateLya[index].State_Abbr)
	tmp_streets_geo := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_POINT_psv.psv")
	sliceLength := len(tmp_streets_geo)
	var wg sync.WaitGroup
	wg.Add(sliceLength - 1)
	for i := 1; i < len(tmp_streets_geo); i++ {
		if strings.Contains(tmp_streets_geo[i], "|") {
			go func(i int) {
				defer wg.Done()
				tmp_line := strings.Split(tmp_streets_geo[i], "|")
				STREET_LOCALITY_PID := tmp_line[3]
				LONGITUDE := tmp_line[6]
				LATITUDE := tmp_line[7]

				for s, _ := range lst_tmp_streets_noGeo {
					if lst_tmp_streets_noGeo[s].STREET_LOCALITY_PID == STREET_LOCALITY_PID {
						lst_tmp_streets_noGeo[s].LATITUDE = LATITUDE
						lst_tmp_streets_noGeo[s].LONGITUDE = LONGITUDE
						break
					}
				}
				//updateGeoOnStreet(index, STREET_LOCALITY_PID, LATITUDE, LONGITUDE)
			}(i)

		}
	}
	wg.Wait()
	tmp_streets = nil
	tmp_streets_geo = nil
	//runtime.GC()
	convertObjStreetToGeoSaving(lst_tmp_streets_noGeo, save)
}

func convertObjStreetToGeoSaving(lst []ObjStreetsLya, path string) {
	fmt.Println("saving: " + path)
	var data []string
	data = append(data, "STREET_LOCALITY_PID,LONGITUDE,LATITUDE,STREET_NAME,STREET_TYPE_CODE,LOCALITY_PID")
	for s, _ := range lst {
		line := lst[s].STREET_LOCALITY_PID + "," + lst[s].LONGITUDE + "," + lst[s].LATITUDE + "," + lst[s].STREET_NAME + "," + lst[s].STREET_TYPE_CODE + "," + lst[s].local_pid
		data = append(data, line)
	}
	saveFile(data, path)
}

func remove(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func loadStreetGeo() {
	//STREET_LOCALITY_POINT_PID|DATE_CREATED|DATE_RETIRED|STREET_LOCALITY_PID|BOUNDARY_EXTENT|PLANIMETRIC_ACCURACY|LONGITUDE|LATITUDE
	//L3163461|2018-08-03||VIC1982231|469||145.13349036|-38.1699638
	for index, _ := range lstObjStateLya {
		oldTime := time.Now()
		fmt.Println("updating geo streets in state: " + lstObjStateLya[index].State_Abbr)
		tmp_streets_geo := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_POINT_psv.psv")
		sliceLength := len(tmp_streets_geo)
		var wg sync.WaitGroup
		wg.Add(sliceLength - 1)
		for i := 1; i < len(tmp_streets_geo); i++ {
			if strings.Contains(tmp_streets_geo[i], "|") {
				go func(i int) {
					defer wg.Done()
					tmp_line := strings.Split(tmp_streets_geo[i], "|")
					STREET_LOCALITY_PID := tmp_line[3]
					LONGITUDE := tmp_line[6]
					LATITUDE := tmp_line[7]
					updateGeoOnStreet(index, STREET_LOCALITY_PID, LATITUDE, LONGITUDE)
				}(i)

			}
		}
		wg.Wait()
		currentTime := time.Now()
		diff := currentTime.Sub(oldTime)
		//In seconds
		fmt.Printf("Seconds: %f\n", diff.Seconds())

		//In nanoseconds
		fmt.Printf("Nanoseconds: %d\n", diff.Nanoseconds())
		tmp_streets_geo = nil
		//runtime.GC()
	}
}

func updateGeoOnStreet(index int, street_id string, lat string, lon string) bool {
	for pc, _ := range lstObjStateLya[index].LstObjPostcodeLya {
		for sub, _ := range lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya {
			for stre, _ := range lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya {
				if street_id == lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya[stre].STREET_LOCALITY_PID {
					lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya[stre].LATITUDE = lat
					lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya[stre].LONGITUDE = lon
					return true
				}
			}

		}

	}
	return false
}

/*STREET_NAME         string
STREET_TYPE_CODE    string
LONGITUDE           string
LATITUDE            string
STREET_LOCALITY_PID string*/

func checkPreprocessedInternet() {
	//.currentInternetPathProcessed

	cppi := readLocal(".currentInternetPathProcessed")
	currentInternetPath = cppi[0]
	internet_path := currentInternetPath + "internet_done.csv"
	internet_path_processed := currentInternetPath + "internet_processed.csv"
	check_preprocessed := fileExists(internet_path)
	check_postprocessed := fileExists(internet_path_processed)
	if !check_postprocessed {
		if !check_preprocessed {
			loadInternet()
			internetToSuburbToStreetTesting(internet_path)
		} else {
			//loading and testing processed data
			loadInternetIntoStreets(internet_path)
			saveNewInternet(internet_path_processed)
		}
	} else {
		finishedInternetProcessingLoadingFile(internet_path_processed)
	}

}

func finishedInternetProcessingLoadingFile(path string) {
	fmt.Println("Loading into lya pool")
	data := readLocal(path)
	var wg sync.WaitGroup
	wg.Add(len(data) - 1)
	for i := 1; i < len(data); i++ {
		go func(i int) {
			defer wg.Done()
			if strings.Contains(data[i], ",") {
				//STATE,Postcode,Suburb,Street,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite
				row := strings.Split(data[i], ",")

				tmp_internet_type := ObjApiInternetType{
					Fixed_wireless: false,
					FTTB:           false,
					FTTDP_FTTC:     false,
					FTTN:           false,
					FTTP:           false,
					HFC:            false,
					Satellite:      false,
				}

				if strings.Contains(row[4], "true") {
					tmp_internet_type.Fixed_wireless = true
				}
				if strings.Contains(row[5], "true") {
					tmp_internet_type.FTTB = true
				}
				if strings.Contains(row[6], "true") {
					tmp_internet_type.FTTDP_FTTC = true
				}
				if strings.Contains(row[7], "true") {
					tmp_internet_type.FTTN = true
				}
				if strings.Contains(row[8], "true") {
					tmp_internet_type.FTTP = true
				}
				if strings.Contains(row[9], "true") {
					tmp_internet_type.HFC = true
				}
				if strings.Contains(row[10], "true") {
					tmp_internet_type.Satellite = true
				}
				setInternetToStatePostSubStreet(row[0], row[1], row[2], row[3], tmp_internet_type)
			}
		}(i)
	}

}

func loadInternetIntoStreets(path string) {
	fmt.Println("loadInternetIntoStreets")
	//LONGITUDE,LATITUDE,selected_Street,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite
	occured := map[string]bool{}
	occured_street_id := map[string]bool{}
	data := readLocal(path)
	var uni_street_data []string
	var uni_street_codes []string
	var lst_uni_objs []ObjProcessInternet
	for i := 1; i < len(data); i++ {
		if strings.Contains(data[i], ",") {
			row := strings.Split(data[i], ",")
			//selected_Street,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite
			row_without_geo := row[2] + "," + row[3] + "," + row[4] + "," + row[5] + "," + row[6] + "," + row[7] + "," + row[8] + "," + row[9]
			if occured[row_without_geo] != true {
				occured[row_without_geo] = true
				uni_street_data = append(uni_street_data, row_without_geo)
			}

			if occured_street_id[row[2]] != true {
				occured_street_id[row[2]] = true
				tmp_uni_data := map[string]bool{}
				uni_street_codes = append(uni_street_codes, row[2])
				tmp := ObjProcessInternet{
					Street_id:      row[2],
					Fixed_wireless: false,
					FTTB:           false,
					FTTDP_FTTC:     false,
					FTTN:           false,
					FTTP:           false,
					HFC:            false,
					Satellite:      false,
					data:           tmp_uni_data,
				}
				lst_uni_objs = append(lst_uni_objs, tmp)
			}
		}
	}

	fmt.Println("Uni Street data: ", len(uni_street_data))
	fmt.Println("Uni Street id: ", len(occured_street_id))
	fmt.Println("pre processing start")

	//way to slow even with multi threading
	/*var wg sync.WaitGroup
	wg.Add(len(lst_uni_objs))
	for r := 0; r < 10000; r++ {
		//for r, _ := range lst_uni_objs {
		go func(r int) {
			defer wg.Done()
			for d, _ := range uni_street_data {
				tmp_row := strings.Split(uni_street_data[d], ",")
				if tmp_row[0] == lst_uni_objs[r].Street_id {
					if lst_uni_objs[r].data[uni_street_data[r]] != true {
						lst_uni_objs[r].data[uni_street_data[r]] = true
					}
				}
			}
		}(r)
	}
	wg.Wait()*/

	//for r := 0; r < 10000; r++ {
	for r, _ := range uni_street_data {
		tmp_row := strings.Split(uni_street_data[r], ",")

		for u, _ := range lst_uni_objs {
			if lst_uni_objs[u].Street_id == tmp_row[0] {
				if lst_uni_objs[u].data[uni_street_data[r]] != true {
					lst_uni_objs[u].data[uni_street_data[r]] = true
				}
				break
			}
		}
		fmt.Println(r, "/", len(uni_street_data))
	}
	fmt.Println("pre processing complete")

	var wg sync.WaitGroup
	wg.Add(len(lst_uni_objs))
	for luo, _ := range lst_uni_objs {
		go func(luo int) {
			defer wg.Done()
			tmp_internet_type := ObjApiInternetType{
				Fixed_wireless: false,
				FTTB:           false,
				FTTDP_FTTC:     false,
				FTTN:           false,
				FTTP:           false,
				HFC:            false,
				Satellite:      false,
			}

			for element, _ := range lst_uni_objs[luo].data {
				tmp_row := strings.Split(element, ",")
				if strings.Contains(tmp_row[1], "true") {
					tmp_internet_type.Fixed_wireless = true
				}
				if strings.Contains(tmp_row[2], "true") {
					tmp_internet_type.FTTB = true
				}
				if strings.Contains(tmp_row[3], "true") {
					tmp_internet_type.FTTDP_FTTC = true
				}
				if strings.Contains(tmp_row[4], "true") {
					tmp_internet_type.FTTN = true
				}
				if strings.Contains(tmp_row[5], "true") {
					tmp_internet_type.FTTP = true
				}
				if strings.Contains(tmp_row[6], "true") {
					tmp_internet_type.HFC = true
				}
				if strings.Contains(tmp_row[7], "true") {
					tmp_internet_type.Satellite = true
				}
				//fmt.Println("Key:", key, "=>", "Element:", element)
			}

			setInternetToStreet(lst_uni_objs[luo].Street_id, tmp_internet_type)
		}(luo)
	}
	wg.Wait()
}

func setInternetToStreet(street_id string, tmp ObjApiInternetType) bool {
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
					if street_id == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID {
						lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet = tmp
						return true
					}
				}
			}
		}
	}
	return false
}

func setInternetToStatePostSubStreet(state string, postcode string, suburb string, street_id string, tmp ObjApiInternetType) bool {
	for states, _ := range lstObjStateLya {
		if state == lstObjStateLya[states].State_Abbr {
			for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
				if postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
					for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
						if suburb == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
							for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
								if street_id == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID {
									lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet = tmp
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}

func saveNewInternet(save_path string) {
	var save_complete_processing []string
	save_complete_processing = append(save_complete_processing, "STATE,Postcode,Suburb,Street,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite")
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
					save := lstObjStateLya[states].State_Abbr + ","
					save += lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode + ","
					save += lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID + ","
					save += lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.Fixed_wireless) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.FTTB) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.FTTDP_FTTC) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.FTTN) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.FTTP) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.HFC) + ","
					save += strconv.FormatBool(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].Selected_Internet.Satellite)
					save_complete_processing = append(save_complete_processing, save)
				}
			}
		}
	}
	saveFile(save_complete_processing, save_path)
}

func loadInternet() {
	//C:\Users\plane\OneDrive\Documents\development\data\internetnrenting\kmlconvert
	//lawson frankston example
	//tmp_lat := "-38.17599781"
	//tmp_lon := "145.12900074"
	//145.12900074,-38.17599781

	//lamington mango hill example
	//tmp_lat := "-27.24106463"
	//tmp_lon := "153.04441987"
	//closest internet node
	//fttp -27.241274 153.044381
	//153.044381,-27.241274

	//NSW2927371 example
	//tmp_lat := "-34.57924181"
	//tmp_lon := "150.74007459"

	cp := readLocal(".currentInternetPath")
	currentInternetPath = cp[0]

	//current location :D files, err := ioutil.ReadDir("./")
	files, err := ioutil.ReadDir(currentInternetPath)
	if err != nil {
		log.Fatal(err)
	}
	var tmp_lst_internet_get []string

	for _, f := range files {
		internet_type := strings.Replace(f.Name(), ".csv", "", -1)
		//if strings.Contains(internet_type, "wireless") || strings.Contains(internet_type, "satellite") {

		//} else {
		fmt.Println("reading: " + internet_type)

		tmp_data := readLocal(currentInternetPath + f.Name())

		for i := 1; i < len(tmp_data); i++ {
			//fmt.Println(tmp_data[i])
			cleaned := strings.ReplaceAll(tmp_data[i], " ", "")
			row := strings.Split(cleaned, ",")
			value := row[0] + "," + row[1]
			tmp_lst_internet_get = append(tmp_lst_internet_get, value)
			//value := row[0] + "," + row[1] + "," + internet_type
			//tmp_lst_internet_get = append(tmp_lst_internet_get, value)
			/*Fixed_wireless := false
			FTTB := false
			FTTDP_FTTC := false
			FTTN := false
			FTTP := false
			HFC := false
			Satellite := false
			tmp_ObjInternetType := ObjInternetType{
				LONGITUDE:      row[0],
				LATITUDE:       row[1],
				Fixed_wireless: Fixed_wireless,
				FTTB:           FTTB,
				FTTDP_FTTC:     FTTDP_FTTC,
				FTTN:           FTTN,
				FTTP:           FTTP,
				HFC:            HFC,
				Satellite:      Satellite,
				//InternetType: internet_type,
			}*/

			//internet_geo = append(internet_geo, tmp_ObjInternetType)

		}
		//}

	}

	cleaned_data := unique(tmp_lst_internet_get)

	//selected_internet := getInternetType(tmp_lon, tmp_lat)
	//fmt.Println("Tmp: ", len(tmp_lst_internet_get))
	fmt.Println("Internet: ", len(cleaned_data))
	fmt.Println("setupInterGeo")
	setupInterGeo(cleaned_data)
	fmt.Println("setupInternetType")
	setupInternetType()
	//fmt.Println(selected_internet.InternetType, selected_internet.LATITUDE, selected_internet.LONGITUDE)
}

func setupInterGeo(cleaned_data []string) {
	for cd := range cleaned_data {
		Fixed_wireless := false
		FTTB := false
		FTTDP_FTTC := false
		FTTN := false
		FTTP := false
		HFC := false
		Satellite := false
		row := strings.Split(cleaned_data[cd], ",")
		tmp_ObjInternetType := ObjInternetType{
			LONGITUDE:      row[0],
			LATITUDE:       row[1],
			Fixed_wireless: Fixed_wireless,
			FTTB:           FTTB,
			FTTDP_FTTC:     FTTDP_FTTC,
			FTTN:           FTTN,
			FTTP:           FTTP,
			HFC:            HFC,
			Satellite:      Satellite,
			//InternetType: internet_type,
		}
		internet_geo = append(internet_geo, tmp_ObjInternetType)
	}
}

func setupInternetType() {

	cp := readLocal(".currentInternetPath")
	currentInternetPath = cp[0]

	files, err := ioutil.ReadDir(currentInternetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		internet_type := strings.Replace(f.Name(), ".csv", "", -1)
		fmt.Println("loading up " + internet_type)
		//var tmp_lst_internet_get []string
		tmp_data := readLocal(currentInternetPath + f.Name())

		occured := map[string]bool{}
		//result := []string{}

		for i := 1; i < len(tmp_data); i++ {
			cleaned := strings.ReplaceAll(tmp_data[i], " ", "")
			row := strings.Split(cleaned, ",")
			value := row[0] + "," + row[1]
			if occured[value] != true {
				occured[value] = true

				// Append to result slice.
				//result = append(result, arr[e])
			}
		}

		//Fixed_wireless := false //0
		//FTTB := false           //1
		//FTTDP_FTTC := false     //2
		//FTTN := false           //3
		//FTTP := false           //4
		//HFC := false            //5
		//Satellite := false      //6

		if strings.Contains(internet_type, "fixed-wireless") {
			//Fixed_wireless = true
			fmt.Println("setting " + internet_type)
			changingInternetType(0, occured)
			fmt.Println("finished " + internet_type)
		}

		if strings.Contains(internet_type, "fttb") {
			//FTTB = true
			changingInternetType(1, occured)
		}

		if strings.Contains(internet_type, "fttdp_fttc") {
			//FTTDP_FTTC = true
			changingInternetType(2, occured)
		}

		if strings.Contains(internet_type, "fttn") {
			//FTTN = true
			changingInternetType(3, occured)
		}

		if strings.Contains(internet_type, "fttp") {
			//FTTP = true
			changingInternetType(4, occured)
		}
		if strings.Contains(internet_type, "hfc") {
			//HFC = true
			changingInternetType(5, occured)
		}
		if strings.Contains(internet_type, "satellite") {
			//Satellite = true
			changingInternetType(6, occured)
		}

	}
}

func changingInternetType(it int, occured map[string]bool) {
	for i, _ := range internet_geo {
		value := internet_geo[i].LONGITUDE + "," + internet_geo[i].LATITUDE
		if occured[value] == true {

			switch it {
			case 0:
				internet_geo[i].Fixed_wireless = true
			case 1:
				internet_geo[i].FTTB = true
			case 2:
				internet_geo[i].FTTDP_FTTC = true
			case 3:
				internet_geo[i].FTTN = true
			case 4:
				internet_geo[i].FTTP = true
			case 5:
				internet_geo[i].HFC = true
			case 6:
				internet_geo[i].Satellite = true
			}

		}
	}
}

// C:\Users\plane\OneDrive\Documents\development\data\internetnrenting\cleaned
func cleanInternetDups() {
	cp := readLocal(".currentInternetPath")
	currentInternetPath = cp[0]
	saveTo := "/mnt/c/Users/plane/OneDrive/Documents/development/data/internetnrenting/cleaned/"
	//current location :D files, err := ioutil.ReadDir("./")
	files, err := ioutil.ReadDir(currentInternetPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		fmt.Println("cleaning: " + f.Name())
		var newFile []string
		newFile = append(newFile, "longitude,latitude")
		tmp_data := readLocal(currentInternetPath + f.Name())
		for i := 1; i < len(tmp_data); i++ {
			cleaned := strings.ReplaceAll(tmp_data[i], " ", "")
			row := strings.Split(cleaned, ",")
			d := row[0] + "," + row[1]
			newFile = append(newFile, d)
		}

		lst_cleaned := unique(newFile)

		//d []string, location string
		saveFile(lst_cleaned, saveTo+f.Name())
	}
}

func unique(arr []string) []string {
	occured := map[string]bool{}
	result := []string{}

	for e := range arr {

		// check if already the mapped
		// variable is set to true or not
		if occured[arr[e]] != true {
			occured[arr[e]] = true

			// Append to result slice.
			result = append(result, arr[e])
		}
	}

	return result
}

func appendCategory(a []string, b []string) []string {

	check := make(map[string]int)
	d := append(a, b...)
	res := make([]string, 0)
	for _, val := range d {
		check[val] = 1
	}

	for letter, _ := range check {
		res = append(res, letter)
	}

	return res
}

func getInternetType(lon string, lat string) ObjInternetType {
	smallest := 9999999999999999.99999
	var selected ObjInternetType
	for index, _ := range internet_geo {
		//lat1 string, long1 string, lat2 string, long2 string
		tmp_dis := getDistance(lat, lon, internet_geo[index].LATITUDE, internet_geo[index].LONGITUDE)
		if tmp_dis <= smallest {
			smallest = tmp_dis
			selected = internet_geo[index]
		}
	}
	return selected
}

/*func loadInternetSpeeds() {
	fmt.Println("loading internet speeds to streets")
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
					tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LATITUDE
					tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LONGITUDE
					selected_internet := getInternetType(tmp_lon, tmp_lat)
					lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].InternetType = selected_internet.InternetType
				}
			}
		}
		fmt.Println("states done: " + lstObjStateLya[states].State_Abbr)
	}
}*/

func internetToStreetTesting() {
	fmt.Println("internetToStreetTesting")

	//testing 20000
	oldTime := time.Now()
	//sliceLength := len(internet_geo)
	var wg sync.WaitGroup
	//wg.Add(sliceLength)
	wg.Add(5000)
	for ig := 0; ig < 5000; ig++ {
		//for ig, _ := range internet_geo {
		go func(ig int) {
			defer wg.Done()
			smallest := 9999999999999999.99999
			var selected string
			for states, _ := range lstObjStateLya {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
						for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
							tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LATITUDE
							tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LONGITUDE
							dis := getDistance(tmp_lat, tmp_lon, internet_geo[ig].LATITUDE, internet_geo[ig].LONGITUDE)
							if dis <= smallest {
								smallest = dis
								selected = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID
							}
						}
					}
				}
			}
			internet_geo[ig].selected_Street = selected
		}(ig)
	}

	wg.Wait()
	currentTime := time.Now()
	diff := currentTime.Sub(oldTime)
	//In seconds
	fmt.Printf("Seconds: %f\n", diff.Seconds())
}

func internetToSuburbTesting() {
	fmt.Println("internetToSuburbTesting")

	//testing 20000
	oldTime := time.Now()
	//sliceLength := len(internet_geo)
	var wg sync.WaitGroup
	//wg.Add(sliceLength)
	wg.Add(5000)
	for ig := 0; ig < 5000; ig++ {
		//for ig, _ := range internet_geo {
		go func(ig int) {
			defer wg.Done()
			smallest := 9999999999999999.99999
			var selected string
			for states, _ := range lstObjStateLya {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
						tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LATITUDE
						tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LONGITUDE
						dis := getDistance(tmp_lat, tmp_lon, internet_geo[ig].LATITUDE, internet_geo[ig].LONGITUDE)
						if dis <= smallest {
							smallest = dis
							selected = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID
						}
					}
				}
			}
			internet_geo[ig].selected_Street = selected
		}(ig)
	}

	wg.Wait()
	currentTime := time.Now()
	diff := currentTime.Sub(oldTime)
	//In seconds
	fmt.Printf("Seconds: %f\n", diff.Seconds())
}

func internetToSuburbToStreetTesting(saveto string) {
	fmt.Println("internetToSuburbToStreetTesting")

	//testing 20000
	oldTime := time.Now()
	//sliceLength := len(internet_geo)
	//var wg sync.WaitGroup
	//wg.Add(sliceLength)
	//wg.Add(1000)

	var tmpObjInternetSuburb []ObjTmpInternetSuburbSort

	//testing for ig := 0; ig < 10000; ig++ {
	for ig, _ := range internet_geo {
		//go func(ig int) {
		//	defer wg.Done()
		smallest := 9999999999999999.99999
		var selected string
		var lat string
		var lon string
		//_ = lon
		//_ = lat
		//_ = selected
		for states, _ := range lstObjStateLya {
			for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
				for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
					tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LATITUDE
					tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LONGITUDE
					dis := getDistance(tmp_lat, tmp_lon, internet_geo[ig].LATITUDE, internet_geo[ig].LONGITUDE)
					if dis <= smallest {
						smallest = dis
						selected = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID
						lat = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LATITUDE
						lon = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LONGITUDE
					}
				}
			}
		}

		found := false
		for tmpobj, _ := range tmpObjInternetSuburb {
			if tmpObjInternetSuburb[tmpobj].LOCALITY_PID == selected {
				tmpObjInternetSuburb[tmpobj].LstInternet = append(tmpObjInternetSuburb[tmpobj].LstInternet, internet_geo[ig])
				found = true
				break
			}
		}

		if !found {
			var tmplstinternet []ObjInternetType
			tmplstinternet = append(tmplstinternet, internet_geo[ig])

			getsubs := getSuburbDistance(
				lat,
				lon)
			sorted := sortClosestSuburb(getsubs)

			tmpAdd := ObjTmpInternetSuburbSort{
				LONGITUDE:    lon,
				LATITUDE:     lat,
				LstInternet:  tmplstinternet,
				LOCALITY_PID: selected,
				sorted:       sorted,
			}
			tmpObjInternetSuburb = append(tmpObjInternetSuburb, tmpAdd)
		}

		//internet_geo[ig].selected_Street = selected
		//	}(ig)
	}

	fmt.Println("Completed adding suburbs")

	done := sortTmpInternetSuburbs(tmpObjInternetSuburb)
	saveInternetProcessed(done, saveto)
	//wg.Wait()
	fmt.Println("SelectedSuburbs: ", len(tmpObjInternetSuburb))
	currentTime := time.Now()
	diff := currentTime.Sub(oldTime)
	//In seconds
	fmt.Printf("Seconds: %f\n", diff.Seconds())
}

func saveInternetProcessed(done []ObjTmpInternetSuburbSort, path string) {
	var save []string
	titles := "LONGITUDE,LATITUDE,selected_Street,Fixed_wireless,FTTB,FTTDP_FTTC,FTTN,FTTP,HFC,Satellite"
	save = append(save, titles)
	/*LONGITUDE string
	LATITUDE  string
	//InternetType    string
	selected_Street string

	Fixed_wireless bool
	FTTB           bool
	FTTDP_FTTC     bool
	FTTN           bool
	FTTP           bool
	HFC            bool
	Satellite      bool*/
	for checkl1, _ := range done {
		for checkl2, _ := range done[checkl1].LstInternet {
			row := done[checkl1].LstInternet[checkl2].LONGITUDE + ","
			row += done[checkl1].LstInternet[checkl2].LATITUDE + ","
			row += done[checkl1].LstInternet[checkl2].selected_Street + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].Fixed_wireless) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].FTTB) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].FTTDP_FTTC) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].FTTN) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].FTTP) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].HFC) + ","
			row += strconv.FormatBool(done[checkl1].LstInternet[checkl2].Satellite)
			save = append(save, row)
			//fmt.Println(done[checkl1].LstInternet[checkl2], done[checkl1].LstInternet[checkl2].selected_Street)
		}
	}

	saveFile(save, path)
}

func sortTmpInternetSuburbs(tmpObjInternetSuburb []ObjTmpInternetSuburbSort) []ObjTmpInternetSuburbSort {
	var wg sync.WaitGroup
	wg.Add(len(tmpObjInternetSuburb))
	for t, _ := range tmpObjInternetSuburb {
		go func(t int) {
			defer wg.Done()
			var tmpstreets []ObjStreetsLya
			limit := 1000
			for ts, _ := range tmpObjInternetSuburb[t].sorted {
				tmpgetstreets := getStreets(tmpObjInternetSuburb[t].sorted[ts].LOCALITY_PID)

				for tmp_str, _ := range tmpgetstreets {
					tmpstreets = append(tmpstreets, tmpgetstreets[tmp_str])
				}

				if ts <= limit {
					break
				}
			}

			for internet, _ := range tmpObjInternetSuburb[t].LstInternet {
				smallest := 9999999999999999.99999
				var selected string
				for f, _ := range tmpstreets {
					tmp_lat := tmpstreets[f].LATITUDE
					tmp_lon := tmpstreets[f].LONGITUDE
					dis := getDistance(tmp_lat, tmp_lon, tmpObjInternetSuburb[t].LstInternet[internet].LATITUDE, tmpObjInternetSuburb[t].LstInternet[internet].LONGITUDE)
					if dis <= smallest {
						smallest = dis
						selected = tmpstreets[f].STREET_LOCALITY_PID
					}
				}
				tmpObjInternetSuburb[t].LstInternet[internet].selected_Street = selected
			}

		}(t)
	}

	wg.Wait()
	return tmpObjInternetSuburb
}

func getStreets(local_pid string) []ObjStreetsLya {
	var tmpstreets []ObjStreetsLya
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				if local_pid == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {

					return lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya
				}
			}
		}
	}
	return tmpstreets
}

/*getsubs := getSuburbDistance(
	lat,
	lon)
sorted := sortClosestSuburb(getsubs)
smallest_stret := 9999999999999999.99999
_ = smallest_stret
var selected_street string
subcount := 50
for sort, _ := range sorted {
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				if sorted[sort].LOCALITY_PID == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
					tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LATITUDE
					tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LONGITUDE
					dis := getDistance(tmp_lat, tmp_lon, internet_geo[ig].LATITUDE, internet_geo[ig].LONGITUDE)
					if dis <= smallest {
						smallest_stret = dis
						selected_street = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID
					}
				}

			}
		}
	}
	subcount += 1
	if subcount >= 50 {
		break
	}
}*/
//---------------------------------------------------
/*oldTime := time.Now()
fmt.Println("updating geo streets in state: " + lstObjStateLya[index].State_Abbr)
tmp_streets_geo := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_POINT_psv.psv")
sliceLength := len(tmp_streets_geo)
var wg sync.WaitGroup
wg.Add(sliceLength - 1)
for i := 1; i < len(tmp_streets_geo); i++ {
	if strings.Contains(tmp_streets_geo[i], "|") {
		go func(i int) {
			defer wg.Done()
			tmp_line := strings.Split(tmp_streets_geo[i], "|")
			STREET_LOCALITY_PID := tmp_line[3]
			LONGITUDE := tmp_line[6]
			LATITUDE := tmp_line[7]
			updateGeoOnStreet(index, STREET_LOCALITY_PID, LATITUDE, LONGITUDE)
		}(i)

	}
}
wg.Wait()
currentTime := time.Now()
diff := currentTime.Sub(oldTime)
//In seconds
fmt.Printf("Seconds: %f\n", diff.Seconds())*/

//single thread
/*for ig, _ := range internet_geo {
	smallest := 9999999999999999.99999
	var selected ObjStreetsLya
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
					tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LATITUDE
					tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LONGITUDE
					dis := getDistance(tmp_lat, tmp_lon, internet_geo[ig].LATITUDE, internet_geo[ig].LONGITUDE)
					if dis <= smallest {
						smallest = dis
						selected = lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets]
					}
				}
			}
		}
	}
	internet_geo[ig].selected_Street = selected
	fmt.Println(ig, "/", len(internet_geo))
}*/
