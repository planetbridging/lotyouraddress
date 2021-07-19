package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var currentPath string
var lstStaticStates []string

var lstObjStateLya []ObjStateLya

var tmp_geo []string

var root = flag.String("root", ".", "file system path")

func main() {
	fmt.Println("Welcome to lya")
	//currentPath = "C:\\Users\\plane\\OneDrive\\Documents\\development\\data\\addy\\MAY21_GNAF_PipeSeparatedValue\\"
	//currentPath = "/mnt/c/Users/plane/OneDrive/Documents/development/data/addy/MAY21_GNAF_PipeSeparatedValue/"
	cp := readLocal(".currentPath")
	currentPath = cp[0]
	//fmt.Println(cp[0])
	loadStates()
	loadPostcodes()
	loadStreets()
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

	fmt.Println("Lotyouraddress running")
	http.HandleFunc("/", handlerFunc)
	http.Handle("/web/", http.FileServer(http.Dir(*root)))
	http.ListenAndServe(":3000", nil)
}

func findSuburb(name string) []ObjSuburbLya {
	var lstSuburbs []ObjSuburbLya
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				tmp_name := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb
				if strings.Contains(tmp_name, name) {
					lstSuburbs = append(lstSuburbs, lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
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
	lstStaticStates = append(lstStaticStates, "ACT")
	lstStaticStates = append(lstStaticStates, "NSW")
	lstStaticStates = append(lstStaticStates, "NT")
	lstStaticStates = append(lstStaticStates, "OT")
	lstStaticStates = append(lstStaticStates, "QLD")
	lstStaticStates = append(lstStaticStates, "SA")
	lstStaticStates = append(lstStaticStates, "TAS")
	lstStaticStates = append(lstStaticStates, "VIC")
	lstStaticStates = append(lstStaticStates, "WA")

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
	return lst_tmp_ObjSuburbLya
}

//street
//STREET_LOCALITY_PID|DATE_CREATED|DATE_RETIRED|STREET_CLASS_CODE|STREET_NAME|STREET_TYPE_CODE|STREET_SUFFIX_CODE|LOCALITY_PID|GNAF_STREET_PID|GNAF_STREET_CONFIDENCE|GNAF_RELIABILITY_CODE
//TAS3349313|2018-05-05||C|EMMETT|STREET||TAS487|502237149|2|4

func loadStreets() {

	for index, _ := range lstObjStateLya {
		fmt.Println("loading streets in : " + lstObjStateLya[index].State_Name)
		tmp_streets := readLocal(currentPath + "G-NAF/G-NAF MAY 2021/Standard/" + lstObjStateLya[index].State_Abbr + "_STREET_LOCALITY_psv.psv")
		for i := 1; i < len(tmp_streets); i++ {
			if strings.Contains(tmp_streets[i], "|") {
				tmp_line := strings.Split(tmp_streets[i], "|")
				tmp_ObjStreetsLya := ObjStreetsLya{
					STREET_LOCALITY_PID: tmp_line[0],
					STREET_NAME:         tmp_line[4],
					STREET_TYPE_CODE:    tmp_line[5],
				}
				found := false
				for pc, _ := range lstObjStateLya[index].LstObjPostcodeLya {
					for sub, _ := range lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya {
						if lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LOCALITY_PID == tmp_line[7] {
							lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya = append(lstObjStateLya[index].LstObjPostcodeLya[pc].LstObjSuburbLya[sub].LstObjStreetsLya, tmp_ObjStreetsLya)
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
	}
	//fmt.Println(lstObjStateLya[0].LstObjPostcodeLya[0].LstObjSuburbLya[0])
}

/*STREET_NAME         string
STREET_TYPE_CODE    string
LONGITUDE           string
LATITUDE            string
STREET_LOCALITY_PID string*/
