package main

type ObjStreetsLya struct {
	STREET_NAME         string
	STREET_TYPE_CODE    string
	LONGITUDE           string
	LATITUDE            string
	STREET_LOCALITY_PID string
	Data                []string
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
