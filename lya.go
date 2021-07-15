package main

import (
	"fmt"
	"strings"
)

var currentPath string
var lstStaticStates []string

var lstObjStateLya []ObjStateLya

var tmp_geo []string

func main() {
	fmt.Println("Welcome to lya")
	//currentPath = "C:\\Users\\plane\\OneDrive\\Documents\\development\\data\\addy\\MAY21_GNAF_PipeSeparatedValue\\"
	currentPath = "/mnt/c/Users/plane/OneDrive/Documents/development/data/addy/MAY21_GNAF_PipeSeparatedValue/"
	loadStates()
	loadPostcodes()
	//loadSuburbs()
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
	fmt.Println(lstObjStateLya[0].LstObjPostcodeLya[0])
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
