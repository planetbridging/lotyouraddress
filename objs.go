package main

import (
	"encoding/json"
	"fmt"
)

type ObjProcessInternet struct {
	Street_id      string
	data           map[string]bool
	Fixed_wireless bool
	FTTB           bool
	FTTDP_FTTC     bool
	FTTN           bool
	FTTP           bool
	HFC            bool
	Satellite      bool
}

type ObjTmpInternetSuburbSort struct {
	LONGITUDE    string
	LATITUDE     string
	LOCALITY_PID string
	LstInternet  []ObjInternetType
	sorted       []ObjDistance
}

type ObjInternetType struct {
	LONGITUDE string
	LATITUDE  string
	//InternetType    string
	selected_Street string

	Fixed_wireless bool
	FTTB           bool
	FTTDP_FTTC     bool
	FTTN           bool
	FTTP           bool
	HFC            bool
	Satellite      bool
}

type ObjStreetsLya struct {
	STREET_NAME         string
	STREET_TYPE_CODE    string
	LONGITUDE           string
	LATITUDE            string
	STREET_LOCALITY_PID string
	Data                []string
	internetType        int
	local_pid           string
	Selected_Internet   ObjApiInternetType
}

type ObjSuburbLya struct {
	Suburb           string
	LOCALITY_PID     string
	LONGITUDE        string
	LATITUDE         string
	LstObjStreetsLya []ObjStreetsLya
}

type ObjPostcodeLya struct {
	Postcode        string
	LstObjSuburbLya []ObjSuburbLya
}

type ObjStateLya struct {
	State_Abbr        string
	State_Name        string
	LstObjPostcodeLya []ObjPostcodeLya
	//HtmlStateTemplate string
}

type ObjDistance struct {
	Suburb       string
	LOCALITY_PID string
	Km           float64
	sorted       bool
	Postcode     string
	State        string
}

//----------------------------------api
type ObjApiInternetType struct {
	Fixed_wireless bool
	FTTB           bool
	FTTDP_FTTC     bool
	FTTN           bool
	FTTP           bool
	HFC            bool
	Satellite      bool
}
type ObjApiStateLya struct {
	State_Abbr     string
	State_Name     string
	Postcode_Count int
	//HtmlStateTemplate string
}

type ObjApiPostcodeLya struct {
	Postcode     string
	Suburb_Count int
}

type ObjApiStateLyaPostcodes struct {
	State_Abbr     string
	State_Name     string
	Postcode_Count int
	LstPostcodes   []ObjApiPostcodeLya
	//HtmlStateTemplate string
}

type ObjApiSuburbLya struct {
	Suburb       string
	LOCALITY_PID string
	LONGITUDE    string
	LATITUDE     string
	Street_count int
}

type ObjApiStreetsLya struct {
	STREET_NAME         string
	STREET_TYPE_CODE    string
	LONGITUDE           string
	LATITUDE            string
	STREET_LOCALITY_PID string
	Selected_Internet   ObjApiInternetType
}

type ObjApiSuburbStreetLya struct {
	Suburb       string
	LOCALITY_PID string
	LONGITUDE    string
	LATITUDE     string
	LstStreets   []ObjApiStreetsLya
}

type ObjApiPostcodeSuburbsLya struct {
	Postcode   string
	lstSuburbs []ObjApiSuburbLya
}

type ObjApiStatePostcodeSuburbs struct {
	State_Abbr     string
	State_Name     string
	Postcode_Count int
	LstPostcodes   []ObjApiPostcodeSuburbsLya
	//HtmlStateTemplate string
}

//---------------------------search api

type ObjApiStreetSearch struct {
	STREET_NAME         string
	STREET_TYPE_CODE    string
	STREET_LOCALITY_PID string
	Suburb              string
	Postcode            string
	State_Abbr          string
	LOCALITY_PID        string
}

type ObjApiSuburbSearch struct {
	Street_Count int
	Suburb       string
	Postcode     string
	State_Abbr   string
	LOCALITY_PID string
}

func convertStateToApi(state ObjStateLya) ObjApiStateLya {
	tmp_ObjApiStateLya := ObjApiStateLya{
		State_Name:     state.State_Name,
		State_Abbr:     state.State_Abbr,
		Postcode_Count: len(state.LstObjPostcodeLya),
	}
	return tmp_ObjApiStateLya
}

func convertPostcodeToApi(postcode ObjPostcodeLya) ObjApiPostcodeLya {
	tmp_ObjApiPostcodeLya := ObjApiPostcodeLya{
		Postcode:     postcode.Postcode,
		Suburb_Count: len(postcode.LstObjSuburbLya),
	}
	return tmp_ObjApiPostcodeLya
}

func convertStatePostcodeToApi(state ObjStateLya, postcodes []ObjApiPostcodeLya) ObjApiStateLyaPostcodes {
	tmp_ObjApiStateLyaPostcodes := ObjApiStateLyaPostcodes{
		State_Name:     state.State_Name,
		State_Abbr:     state.State_Abbr,
		Postcode_Count: len(state.LstObjPostcodeLya),
		LstPostcodes:   postcodes,
	}
	return tmp_ObjApiStateLyaPostcodes
}

func convertSuburbToApi(suburb ObjSuburbLya) ObjApiSuburbLya {
	tmp_ObjApiSuburbLya := ObjApiSuburbLya{
		Suburb:       suburb.Suburb,
		LOCALITY_PID: suburb.LOCALITY_PID,
		LONGITUDE:    suburb.LONGITUDE,
		LATITUDE:     suburb.LATITUDE,
		Street_count: len(suburb.LstObjStreetsLya),
	}
	return tmp_ObjApiSuburbLya
}

func convertStreetToApi(street ObjStreetsLya) ObjApiStreetsLya {
	tmp_streets := ObjApiStreetsLya{
		STREET_NAME:         street.STREET_NAME,
		STREET_TYPE_CODE:    street.STREET_TYPE_CODE,
		STREET_LOCALITY_PID: street.STREET_LOCALITY_PID,
		LONGITUDE:           street.LONGITUDE,
		LATITUDE:            street.LATITUDE,
		Selected_Internet:   street.Selected_Internet,
	}
	return tmp_streets
}

func convertStatePostcodeSuburb(state ObjStateLya, postcode ObjPostcodeLya, suburbs []ObjApiSuburbLya) ObjApiStatePostcodeSuburbs {
	tmp_api_postcode := ObjApiPostcodeSuburbsLya{
		Postcode:   postcode.Postcode,
		lstSuburbs: suburbs,
	}
	var tmp_lst_post []ObjApiPostcodeSuburbsLya
	tmp_lst_post = append(tmp_lst_post, tmp_api_postcode)
	tmp_ObjApiStatePostcodeSuburbs := ObjApiStatePostcodeSuburbs{
		State_Abbr:     state.State_Abbr,
		State_Name:     state.State_Name,
		Postcode_Count: len(state.LstObjPostcodeLya),
		LstPostcodes:   tmp_lst_post,
	}
	return tmp_ObjApiStatePostcodeSuburbs
}

func convertSuburbStreet(suburb ObjSuburbLya, streets []ObjApiStreetsLya) ObjApiSuburbStreetLya {
	tmp_suburbs := ObjApiSuburbStreetLya{
		Suburb:       suburb.Suburb,
		LOCALITY_PID: suburb.LOCALITY_PID,
		LONGITUDE:    suburb.LONGITUDE,
		LATITUDE:     suburb.LATITUDE,
		LstStreets:   streets,
	}
	return tmp_suburbs
}

func convertStreetSearch(state ObjStateLya, postcode ObjPostcodeLya, suburb ObjSuburbLya, street ObjStreetsLya) ObjApiStreetSearch {
	tmp_obj := ObjApiStreetSearch{
		STREET_NAME:         street.STREET_NAME,
		STREET_TYPE_CODE:    street.STREET_TYPE_CODE,
		STREET_LOCALITY_PID: street.STREET_LOCALITY_PID,
		Suburb:              suburb.Suburb,
		Postcode:            postcode.Postcode,
		State_Abbr:          state.State_Abbr,
		LOCALITY_PID:        suburb.LOCALITY_PID,
	}
	return tmp_obj
}

func convertSuburbSearch(state ObjStateLya, postcode ObjPostcodeLya, suburb ObjSuburbLya) ObjApiSuburbSearch {
	tmp_obj := ObjApiSuburbSearch{
		Suburb:       suburb.Suburb,
		Postcode:     postcode.Postcode,
		State_Abbr:   state.State_Abbr,
		LOCALITY_PID: suburb.LOCALITY_PID,
		Street_Count: len(suburb.LstObjStreetsLya),
	}
	return tmp_obj
}

//---------fixing json scramble
func manualJsonSuburbSearch(lst []ObjDistance) string {
	var results string
	fmt.Println(lst[0])
	for i := 0; i < len(lst)-1; i++ {
		pagesJson, _ := json.Marshal(lst[i])
		results += string(pagesJson) + ","
	}
	//pagesJson, _ := json.Marshal(lst[len(lst)])
	//results += string(pagesJson)
	return "[" + results + "]"
}
