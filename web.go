package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	//var page_loaded = false
	//keys, ok := r.URL.Query()["key"]
	kstate, okst := r.URL.Query()["state"]
	kpostcode, okpo := r.URL.Query()["postcode"]
	ksuburb, oksu := r.URL.Query()["suburb"]
	kstreet, okstr := r.URL.Query()["street"]

	search_kstreet, search_okstr := r.URL.Query()["search_street"]
	search_ksuburb, search_oksu := r.URL.Query()["search_suburb"]

	knearsuburbs, oknearsuburbs := r.URL.Query()["near_suburb"]
	kinternet, okinternet := r.URL.Query()["internet"]

	/*if !ok || len(keys[0]) < 1 {
	}else
	*/
	found := false

	if okst && !okpo && !oksu && !okstr {
		kstate := kstate[0]
		lstApiStates := apiStates(kstate)
		tmp_json := convertStatesJson(lstApiStates)
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if okst && okpo && !oksu && !okstr {
		kstate := kstate[0]
		kpostcode := kpostcode[0]
		showJustPostcodes := apiPostcodes(kstate, kpostcode)
		tmp_json := convertPostcodesJson(showJustPostcodes)
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if okst && okpo && oksu && !okstr {
		ksuburb := ksuburb[0]
		kpostcode := kpostcode[0]
		kstate := kstate[0]
		showJustSuburbs := apiSuburbs(kstate, kpostcode, ksuburb)
		//fmt.Println(showJustSuburbs)
		tmp_json := convertSuburbsJson(showJustSuburbs)
		//fmt.Fprintf(os.Stdout, "%s", tmp_json)
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if okst && okpo && oksu && okstr {
		kstreet := kstreet[0]
		ksuburb := ksuburb[0]
		kpostcode := kpostcode[0]
		kstate := kstate[0]
		showJustStreets := apiStreet(kstate, kpostcode, ksuburb, kstreet)
		tmp_json := convertStreetsJson(showJustStreets)
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if !okst && !okpo && !oksu && !okstr && search_okstr {
		showStreetSearch := streetSearch(search_kstreet[0], 20)
		tmp_json, err := json.Marshal(showStreetSearch)
		if err != nil {
			log.Fatal("Cannot encode street search to JSON ", err)
		}
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if !okst && !okpo && !oksu && !okstr && search_oksu {
		showSuburbSearch := suburbSearch(search_ksuburb[0], 20)
		tmp_json, err := json.Marshal(showSuburbSearch)
		if err != nil {
			log.Fatal("Cannot encode suburb search to JSON ", err)
		}
		fmt.Fprintf(w, "%s", tmp_json)
		found = true
	} else if !okst && !okpo && !oksu && !okstr && !search_oksu && !search_okstr && oknearsuburbs {
		knearsuburbs := knearsuburbs[0]
		findNear := findPIDSuburb(knearsuburbs)
		//fmt.Println(findNear)
		if len(findNear) == 1 {
			getsubs := getSuburbDistance(
				findNear[0].LATITUDE,
				findNear[0].LONGITUDE)
			sorted := sortClosestSuburb(getsubs)
			//var top20 []ObjDistance
			//fmt.Println(sorted[0])
			tmp_json, err := json.Marshal(sorted)
			if err != nil {
				log.Fatal("Cannot encode near suburb to JSON ", err)
			}

			fmt.Fprintf(w, "%s", tmp_json)
			//manualRebuild := manualJsonSuburbSearch(getsubs)
			found = true
			//fmt.Fprint(w, manualRebuild)
		}
	} else if okinternet {
		_ = kinternet
		/*street_pid := kinternet[0]
		tmp_internet := internetSearch(street_pid)
		fmt.Fprintf(w, "%s", tmp_internet.InternetType)*/
		fmt.Fprint(w, "Currently not setup")
		found = true
	} else {
		fmt.Fprint(w, "Welcome to lotyouraddress api")
		//found = true
	}

	if !found {
		fmt.Fprint(w, "Welcome to lotyouraddress api")
	}
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func writeImage(w http.ResponseWriter, img *image.Image) {

	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		//log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		//log.Println("unable to write image.")
	}
}

func convertStatesJson(tmp_states []ObjApiStateLya) []byte {
	pagesJson, err := json.Marshal(tmp_states)
	if err != nil {
		log.Fatal("Cannot encode state to JSON ", err)
	}
	return pagesJson
}

func convertPostcodesJson(tmp_postcode ObjApiStateLyaPostcodes) []byte {
	pagesJson, err := json.Marshal(tmp_postcode)
	if err != nil {
		log.Fatal("Cannot encode postcode to JSON ", err)
	}
	return pagesJson
}

func convertSuburbsJson(tmp_suburbs []ObjApiSuburbLya) []byte {
	pagesJson, err := json.Marshal(tmp_suburbs)
	if err != nil {
		log.Fatal("Cannot encode suburbs to JSON ", err)
	}
	return pagesJson
}

func convertStreetsJson(tmp_streets ObjApiSuburbStreetLya) []byte {
	pagesJson, err := json.Marshal(tmp_streets)
	if err != nil {
		log.Fatal("Cannot encode streets to JSON ", err)
	}
	return pagesJson
}

func apiStates(search string) []ObjApiStateLya {
	var lstApiStates []ObjApiStateLya
	if search == "" {
		for states, _ := range lstObjStateLya {
			tmp_ObjStateLya := convertStateToApi(lstObjStateLya[states])
			lstApiStates = append(lstApiStates, tmp_ObjStateLya)
		}
	} else {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				tmp_ObjStateLya := convertStateToApi(lstObjStateLya[states])
				lstApiStates = append(lstApiStates, tmp_ObjStateLya)
			}
		}
	}
	return lstApiStates
}

func apiPostcodes(search string, search_postcode string) ObjApiStateLyaPostcodes {
	var tmp_state_postcodes ObjApiStateLyaPostcodes
	if search_postcode == "" {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				var tmp_ObjApiPostcodeLya []ObjApiPostcodeLya
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					tmp_api_postcode := convertPostcodeToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes])
					tmp_ObjApiPostcodeLya = append(tmp_ObjApiPostcodeLya, tmp_api_postcode)
				}
				tmp_state_postcodes = convertStatePostcodeToApi(lstObjStateLya[states], tmp_ObjApiPostcodeLya)
				break
			}
		}
	} else {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				var tmp_ObjApiPostcodeLya []ObjApiPostcodeLya
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					if search_postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
						tmp_api_postcode := convertPostcodeToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes])
						tmp_ObjApiPostcodeLya = append(tmp_ObjApiPostcodeLya, tmp_api_postcode)
						break
					}
				}
				tmp_state_postcodes = convertStatePostcodeToApi(lstObjStateLya[states], tmp_ObjApiPostcodeLya)
				break
			}
		}
	}
	return tmp_state_postcodes
}

func apiSuburbs(search string, search_postcode string, search_suburb string) []ObjApiSuburbLya {
	var lst_tmp_ObjApiSuburbLya []ObjApiSuburbLya
	if search_suburb == "" {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					if search_postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
						var lst_tmp_ObjApiSuburbLya []ObjApiSuburbLya
						for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
							tmp_suburb := convertSuburbToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
							lst_tmp_ObjApiSuburbLya = append(lst_tmp_ObjApiSuburbLya, tmp_suburb)
						}
						/*tmp_convert_suburb := convertStatePostcodeSuburb(
							lstObjStateLya[states],
							lstObjStateLya[states].LstObjPostcodeLya[postcodes],
							lst_tmp_ObjApiSuburbLya)
						return tmp_convert_suburb*/
						return lst_tmp_ObjApiSuburbLya
					}
				}
			}

		}
	} else {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					if search_postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
						var lst_tmp_ObjApiSuburbLya []ObjApiSuburbLya
						for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
							if search_suburb == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
								tmp_suburb := convertSuburbToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
								lst_tmp_ObjApiSuburbLya = append(lst_tmp_ObjApiSuburbLya, tmp_suburb)
								break
							}

						}
						/*tmp_convert_suburb := convertStatePostcodeSuburb(
							lstObjStateLya[states],
							lstObjStateLya[states].LstObjPostcodeLya[postcodes],
							lst_tmp_ObjApiSuburbLya)
						return tmp_convert_suburb*/
						return lst_tmp_ObjApiSuburbLya
					}
				}
			}

		}
	}
	return lst_tmp_ObjApiSuburbLya
}

func apiStreet(search string, search_postcode string, search_suburb string, search_street string) ObjApiSuburbStreetLya {
	var tmp_streets ObjApiSuburbStreetLya
	if search_street == "" {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					if search_postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
						for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
							if search_suburb == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
								var tmp_lst_streets []ObjApiStreetsLya
								for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
									tmp_convert_streets := convertStreetToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets])
									tmp_lst_streets = append(tmp_lst_streets, tmp_convert_streets)
								}
								tmp_streets := convertSuburbStreet(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs], tmp_lst_streets)
								return tmp_streets
							}

						}
					}
				}
			}

		}
	} else {
		for states, _ := range lstObjStateLya {
			if search == lstObjStateLya[states].State_Abbr {
				for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
					if search_postcode == lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode {
						for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
							if search_suburb == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LOCALITY_PID {
								var tmp_lst_streets []ObjApiStreetsLya
								for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
									if search_street == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID {
										tmp_convert_streets := convertStreetToApi(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets])
										tmp_lst_streets = append(tmp_lst_streets, tmp_convert_streets)
										break
									}

								}
								tmp_streets := convertSuburbStreet(lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs], tmp_lst_streets)
								return tmp_streets
							}

						}
					}
				}
			}

		}
	}
	return tmp_streets
}

func streetSearch(search string, limiter int) []ObjApiStreetSearch {
	var results []ObjApiStreetSearch
	count := 0
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {
					start := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_NAME
					end := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_TYPE_CODE
					street_name := start + " " + end
					state_string := lstObjStateLya[states].State_Abbr
					postcode_string := lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode
					suburb_string := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb
					full := street_name + " " + suburb_string + " " + postcode_string + " " + state_string
					if strings.Contains(full, search) && count <= 20 {
						//tmp_full := street_name + "," + suburb_string + "," + postcode_string + "," + state_string
						tmp_full := convertStreetSearch(lstObjStateLya[states],
							lstObjStateLya[states].LstObjPostcodeLya[postcodes],
							lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs],
							lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets])
						results = append(results, tmp_full)
						count += 1

						if count >= limiter {
							return results
						}
					}

				}

			}

		}

	}
	return results
}

func suburbSearch(search string, limiter int) []ObjApiSuburbSearch {
	count := 0
	var results []ObjApiSuburbSearch
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				state_string := lstObjStateLya[states].State_Abbr
				postcode_string := lstObjStateLya[states].LstObjPostcodeLya[postcodes].Postcode
				suburb_string := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].Suburb
				full := suburb_string + " " + postcode_string + " " + state_string
				if strings.Contains(full, search) && count <= 20 {
					//tmp_full := street_name + "," + suburb_string + "," + postcode_string + "," + state_string
					tmp_full := convertSuburbSearch(lstObjStateLya[states],
						lstObjStateLya[states].LstObjPostcodeLya[postcodes],
						lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs])
					results = append(results, tmp_full)
					count += 1
					if count >= limiter {
						return results
					}
				}

			}
		}

	}
	return results
}

/*func internetSearch(search string) ObjInternetType {
	tmp_empty := ObjInternetType{
		InternetType: "empty",
	}
	for states, _ := range lstObjStateLya {
		for postcodes, _ := range lstObjStateLya[states].LstObjPostcodeLya {
			for suburbs, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya {
				for streets, _ := range lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya {

					if search == lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].STREET_LOCALITY_PID {
						tmp_lat := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LATITUDE
						tmp_lon := lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].LONGITUDE
						selected_internet := getInternetType(tmp_lon, tmp_lat)
						return selected_internet
					}
					//lstObjStateLya[states].LstObjPostcodeLya[postcodes].LstObjSuburbLya[suburbs].LstObjStreetsLya[streets].InternetType = selected_internet.InternetType
				}
			}
		}
		//fmt.Println("states done: " + lstObjStateLya[states].State_Abbr)
	}
	return tmp_empty
}*/

func getCloseStreet(lat string, long string) {

}
