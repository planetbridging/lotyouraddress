package main

import (
  "fmt"
  "net/http"
  "strings"
  "io/ioutil"
  "image/jpeg"
  "bytes"
  "image"
  "strconv"
  "flag"
  "os"
  "sync"
  "time"
  "pack/simpleinterest"
  "pack/what"
)

type ObjStreetsLya struct {
    STREET_NAME string
    STREET_TYPE_CODE string
    LONGITUDE string
    LATITUDE string
    STREET_LOCALITY_PID string
    Data[] string
}

type ObjSuburbLya struct {
    Suburb string
    LOCALITY_PID string
    LstObjStreetsLya[] ObjStreetsLya
}

type ObjPostcodeLya struct {
    Postcode string
    LstObjSuburbLya[] ObjSuburbLya
}

type ObjStateLya struct {
    State string
    LstObjPostcodeLya[] ObjPostcodeLya
    //HtmlStateTemplate string
}

type ObjGeo struct{
  STREET_LOCALITY_PID string
  Data[] string
}

var LstObjStateLya[] ObjStateLya
var LstLya[] string
var root = flag.String("root", ".", "file system path")
var HomePageTemplate = ""

var ImageTemplate string = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
        <title>Clair</title>
		<link rel="stylesheet" href="/web/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
        <script src="/web/js/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
		<script src="/web/js/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
		<script src="/web/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
		<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
		<link href="/web/css/master.css" rel="stylesheet">
		<script src="/web/js/jquery.min.js"></script>
		<script src="/web/js/master.js"></script>
</head>
<body><div></div><img src="/web/imgs/clair.png">`

func handlerFunc(w http.ResponseWriter, r *http.Request){
  //var page_loaded = false
  //keys, ok := r.URL.Query()["key"]
  kstate, okst := r.URL.Query()["state"]
  kpostcode, okpo := r.URL.Query()["postcode"]
  ksuburb, oksu := r.URL.Query()["suburb"]
  kstreet, okstr := r.URL.Query()["street"]

  /*if okst{kstate := kstate[0]}
  if okpo{kpostcode := kpostcode[0]}
  if oksu{ksuburb := ksuburb[0]}
  if okstr{kstreet := kstreet[0]}*/

  if okst && okpo && oksu && okstr{
    fmt.Print(kpostcode[0] + ksuburb[0] + kstreet[0])
  }else if okst && okpo && oksu{

  }else if okst && okpo{

  }else if okst{
    //fmt.Println("State: " + kstate)
    HtmlStateOutput(kstate[0],w)
  }

  /*if !ok || len(keys[0]) < 1 {

  }else
  */
/*  if ok{
        key := keys[0]

        fmt.Println(w,"Url Param 'key' is: " + string(key))
  }

  if okpo{
    kpostcode := kpostcode[0]
    fmt.Println("Loading: " + kpostcode)
  }*/


  if okst{



    /*if okpo{

      if oksu{

        if okst{
          kstreet := kstreet[0]
          fmt.Println("Loading: " + kstreet)
        }else{
          ksuburb := ksuburb[0]
          fmt.Println("Loading: " + ksuburb)
        }

      }else{

        kpostcode := kpostcode[0]
        fmt.Println("Loading: " + kpostcode)
      }

    }else{

    }
    kstate := kstate[0]
    fmt.Println("Loading: " + kstate)
    HtmlStateOutput(kstate,w)*/
  }
  /*else if r.URL.Path == "/"{
      fmt.Fprint(w,HomePageTemplate)
      page_loaded = true
  }

    if !page_loaded{
      for s := 1; s < len(LstObjStateLya); s++ {
        var tmpnewpathstate string = "/" + strings.TrimSpace(LstObjStateLya[s].State)
        if r.URL.Path == tmpnewpathstate{
          fmt.Fprint(w,ImageTemplate + GetSelectionPostcode(s) + "</body></html>")
          page_loaded = true
          break
        }else if r.URL.Path == tmpnewpathstate + "/"{
          fmt.Fprint(w,ImageTemplate + SetupStateHtmlTemplates(s) + "</body></html>")
          page_loaded = true
          break
        }
      }
    }*/
}

func HtmlStateOutput(kstate string, w http.ResponseWriter){
  for s := 1; s < len(LstObjStateLya); s++ {
    if kstate == strings.TrimSpace(LstObjStateLya[s].State) {
      fmt.Fprint(w,ImageTemplate + GetSelectionPostcode(s) + "</body></html>")
      break
    }
  }
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

func GetStateTitle(num int) string{
  var statename = strings.TrimSpace(LstObjStateLya[num].State)
  var postcode_count = len(LstObjStateLya[num].LstObjPostcodeLya)
  var suburb_count = 0
  var street_count = 0

  for p:= range LstObjStateLya[num].LstObjPostcodeLya{
    suburb_count += len(LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya)
    for s:= range LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya{
      street_count += len(LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya[s].LstObjStreetsLya)
    }
  }
  var title = "State: " + statename + " Postcodes: " + strconv.Itoa(postcode_count)
  title += " Suburbs: " + strconv.Itoa(suburb_count)
  title += " Streets: " + strconv.Itoa(street_count)
  return title
}

func GetHtmlAccordionCard(acctype string,title string,name string,content string) string{
  name = "Lst" + strings.TrimSpace(name)
  title = strings.TrimSpace(title)
  var CardTemplate string = `<div class="card">
    <div class="card-header" id="heading`+name+`">
      <h5 class="mb-0">
        <button class="btn btn-link collapsed" data-toggle="collapse" data-target="#`+name+`" aria-expanded="false" aria-controls="`+name+`">
          `+title+`
        </button>
      </h5>
    </div>
    <div id="`+name+`" class="collapse" aria-labelledby="heading`+name+`" data-parent="#`+acctype+`accordion">
      <div class="card-body">
        `+content+`
    </div>
  </div>`
  return CardTemplate
}

func GetSelectionPostcode(num int) string{
    var postcodetbl = `<table class="table table-dark table-striped">
    <thead>
      <tr>
        <th scope="col">Postcode</th>
        <th scope="col">Suburbs</th>
        <th scope="col">Streets</th>
      </tr>
    </thead>
    <tbody>`
    for p:= range LstObjStateLya[num].LstObjPostcodeLya {
      postcodetbl+= "<tr><td>" + LstObjStateLya[num].LstObjPostcodeLya[p].Postcode + "</td>"
      postcodetbl+= "<td>" + strconv.Itoa(len(LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya)) + "</td>"
      var streetcount = 0
      for s:= range LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya{
        streetcount+=len(LstObjStateLya[num].LstObjPostcodeLya[p].LstObjSuburbLya[s].LstObjStreetsLya)
      }
      postcodetbl+= "<td>" + strconv.Itoa(streetcount) + "</td></tr>"
    }
    postcodetbl += "</tbody></table>"
    return postcodetbl
}

func GetStatesHtml() string{
  var homepage = `<div id="Statesaccordion">`
  //for s:= range LstObjStateLya[0:len(LstObjStateLya)]  {
  for s := 1; s < len(LstObjStateLya); s++ {
  //  if s > 0{
      var statetitle = GetStateTitle(s)
      homepage += GetHtmlAccordionCard("States",statetitle,LstObjStateLya[s].State,GetSelectionPostcode(s))
    //}
  }
  homepage += `</div>`
  return homepage
}

func Contains(arr []string, str string) bool {
   for _, a := range arr {
      if a == str {
         return true
      }
   }
   return false
}

func Unique(elements []string) []string {
    encountered := map[string]bool{}

    // Create a map of all unique elements.
    for v:= range elements {
        encountered[elements[v]] = true
    }

    // Place all keys from the map into a slice.
    result := []string{}
    for key, _ := range encountered {
        result = append(result, key)
    }
    return result
}

func ReadFile(filename string) []string{
  data, err := ioutil.ReadFile(filename)
  var LstFile[]string
  if err != nil {
      fmt.Println("File reading error", err)
      return LstFile
  }

  for _, line := range strings.Split(strings.TrimSuffix(string(data), "\n"), "\n") {
    LstFile = append(LstFile, line)
  }
  return LstFile
}

func SortingHat(r string){
  //STREET_LOCALITY_PID,STREET_NAME,STREET_TYPE_CODE,LOCALITY_PID,Suburb,Postcode,LONGITUDE,LATITUDE,State
  var r_split = strings.Split(r,",")

  tmp_ObjStreetsLya := ObjStreetsLya{
    STREET_LOCALITY_PID: r_split[0],
    STREET_NAME: r_split[1],
    STREET_TYPE_CODE: r_split[2],
    LONGITUDE: r_split[6],
    LATITUDE: r_split[7],
  }

  found := false
  var s_num int = -1
  //for s:= range LstObjStateLya{
    for s := 1; s < len(LstObjStateLya); s++ {
  //  StatesSortingHat(LstLya[v])
    //LstStates = append(LstStates, r_split[8])
    result := LstObjStateLya[s].State == strings.TrimSpace(r_split[8])
    if result{
      found = true
      s_num = s
      break
    }
  }

  var tmplya [] ObjStreetsLya
  tmplya = append(tmplya, tmp_ObjStreetsLya)
  tmp_ObjSuburbLya := ObjSuburbLya{
    Suburb: r_split[4],
    LOCALITY_PID: r_split[3],
    LstObjStreetsLya: tmplya,
  }
  tmp_ObjPostcodeLya := ObjPostcodeLya{
    Postcode: r_split[5],
  }
  tmp_ObjPostcodeLya.LstObjSuburbLya = append(tmp_ObjPostcodeLya.LstObjSuburbLya, tmp_ObjSuburbLya)
  tmp_ObjStateLya := ObjStateLya{
    State: strings.TrimSpace(r_split[8]),
  }
  tmp_ObjStateLya.LstObjPostcodeLya = append(tmp_ObjStateLya.LstObjPostcodeLya, tmp_ObjPostcodeLya)

  if !found{
    LstObjStateLya = append(LstObjStateLya, tmp_ObjStateLya)
  }else{
    p_found := false
    var p_num int = -1
    for p:= range LstObjStateLya[s_num].LstObjPostcodeLya{
      result := LstObjStateLya[s_num].LstObjPostcodeLya[p].Postcode == r_split[5]
      if result{
        p_found = true
        p_num = p
        break
      }
    }

    if p_found{
      sub_found := false
      var sub_num int = -1
      for sub:= range LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya{
        result := LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya[sub].LOCALITY_PID == r_split[3]
        if result{
          sub_found = true
          sub_num = sub
          break
        }
      }

      if sub_found{
        LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya[sub_num].LstObjStreetsLya = append(
          LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya[sub_num].LstObjStreetsLya, tmp_ObjStreetsLya)
      }else{
        LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya = append(
          LstObjStateLya[s_num].LstObjPostcodeLya[p_num].LstObjSuburbLya, tmp_ObjSuburbLya)
      }

    }else{
      LstObjStateLya[s_num].LstObjPostcodeLya = append(LstObjStateLya[s_num].LstObjPostcodeLya, tmp_ObjPostcodeLya)
    }

  }


}

func SetupStreetSearch(){
  fmt.Println("Lotyouraddress street setup")
  LstLya = ReadFile("slotyouraddress_streets_geo.csv")
  for v:= range LstLya{
    SortingHat(LstLya[v])
  }

  HomePageTemplate = ImageTemplate +GetStatesHtml()+"</body></html>"
}

func SetupStateHtmlTemplates(v int) string{
  //for v:= range LstObjStateLya{
    var statehtml = `<table class="table table-dark table-striped">
    <thead>
      <tr>
        <th scope="col">Postcode</th>
        <th scope="col">LOCALITY_PID</th>
        <th scope="col">Suburb</th>
        <th scope="col">Streets</th>
      </tr>
    </thead>
    <tbody>`
    for p:= range LstObjStateLya[v].LstObjPostcodeLya{
      var tmppost = LstObjStateLya[v].LstObjPostcodeLya[p].Postcode
      for s:= range LstObjStateLya[v].LstObjPostcodeLya[p].LstObjSuburbLya{
        statehtml+= "<tr><td>" + tmppost + "</td>"
        statehtml+= "<td>" + LstObjStateLya[v].LstObjPostcodeLya[p].LstObjSuburbLya[s].LOCALITY_PID + "</td>"
        statehtml+= "<td>" + LstObjStateLya[v].LstObjPostcodeLya[p].LstObjSuburbLya[s].Suburb + "</td>"
        statehtml+= "<td>" + strconv.Itoa(len(LstObjStateLya[v].LstObjPostcodeLya[p].LstObjSuburbLya[s].LstObjStreetsLya)) + "</td></tr>"
      }
    }
    statehtml += "</tbody></table>"
    //LstObjStateLya[v].HtmlStateTemplate = statehtml
  //}
  return statehtml
}

func GeoSortingHat(data[] string, s int, p int , sub int){
  //ADDRESS_DETAIL_PID|BUILDING_NAME|LOT_NUMBER_PREFIX|LOT_NUMBER|LOT_NUMBER_SUFFIX|FLAT_TYPE_CODE|FLAT_NUMBER_PREFIX|FLAT_NUMBER|FLAT_NUMBER_SUFFIX|LEVEL_TYPE_CODE|LEVEL_NUMBER_PREFIX|LEVEL_NUMBER|LEVEL_NUMBER_SUFFIX|NUMBER_FIRST_PREFIX|NUMBER_FIRST|NUMBER_FIRST_SUFFIX|NUMBER_LAST_PREFIX|NUMBER_LAST|NUMBER_LAST_SUFFIX|STREET_LOCALITY_PID|LEGAL_PARCEL_ID|LEVEL_GEOCODED_CODE|LONGITUDE|LATITUDE
//STREET_LOCALITY_PID = 19
  for d:= range data{
    var d_split[] string = strings.Split(data[d],"|")
    if len(d_split) > 19{
      for street:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya{
        var street_pid = LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[street].STREET_LOCALITY_PID
        if street_pid == d_split[19]{
          //fmt.Println("yay")
        //  LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[street].Data
          LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[street].Data =
          append(LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[street].Data, data[d])
          break
        }
      }
    }else{
      fmt.Println(LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LOCALITY_PID)
      fmt.Println(data[d])
    }
  }
}

func SetupGeoLocation(){
  dir, _ := os.Getwd()
  for s := 1; s < len(LstObjStateLya); s++ {
    for p:= range LstObjStateLya[s].LstObjPostcodeLya{
      for sub:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya{
        var file_geo_path = dir + "/"+LstObjStateLya[s].State+"/" + LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LOCALITY_PID + ".txt"
        //fmt.Println(file_geo_path)
        var data[] string = ReadFile(file_geo_path)
        if len(data) > 0{
          GeoSortingHat(data,s,p,sub)
        }
      }
    }
  }

  //fmt.Println(dir)
}

func SetupGeoBulk(s int ){
  var LstObjGeo[] ObjGeo
  start := time.Now()
  //var wg sync.WaitGroup
  fmt.Println(start)
  fmt.Println(LstObjStateLya[s].State + "_ADDRESS_DETAIL_Extracted.txt")
  os.MkdirAll(LstObjStateLya[s].State, os.ModePerm)
  var fn_state string = LstObjStateLya[s].State + "_ADDRESS_DETAIL_Extracted.txt"
  var data[] string = ReadFile(fn_state)
  for fnr := 1; fnr < len(data); fnr++ {
    var r_split = strings.Split(data[fnr],"|")
    if len(r_split) < 22{
    }else{
      var found = false
      for ogeo:= range LstObjGeo{
        if LstObjGeo[ogeo].STREET_LOCALITY_PID == r_split[19]{
          LstObjGeo[ogeo].Data = append(LstObjGeo[ogeo].Data,data[fnr])
          found = true
          break
        }
      }

      if !found{
        tmp_ObjGeo := ObjGeo{
          STREET_LOCALITY_PID: r_split[19],
        }
        tmp_ObjGeo.Data = append(tmp_ObjGeo.Data,data[fnr])
        LstObjGeo = append(LstObjGeo,tmp_ObjGeo)
      }
    }
    //wg.Add(1)
    //go AppendGeoData(s,data[fnr],&wg)
    //fmt.Println(fnr)
  }
  t := time.Now()
  elapsed := t.Sub(start)
  fmt.Println(elapsed)

  for g:= range LstObjGeo{
    var found = false
    for p:= range LstObjStateLya[s].LstObjPostcodeLya{
      for sub:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya{
        for str:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya{
          var tmplip = LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].STREET_LOCALITY_PID
          if tmplip == LstObjGeo[g].STREET_LOCALITY_PID{
            LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].Data = LstObjGeo[g].Data
            found = true
            break
          }
        }
        if found{
          break
        }
      }
      if found{
        break
      }
    }
  }
  SaveStateSuburbID(s)
  fmt.Println("done")
}

func SaveStateSuburbID(s int){
  var template_title = "ADDRESS_DETAIL_PID|BUILDING_NAME|LOT_NUMBER_PREFIX|LOT_NUMBER|LOT_NUMBER_SUFFIX|FLAT_TYPE_CODE|FLAT_NUMBER_PREFIX|FLAT_NUMBER|FLAT_NUMBER_SUFFIX|LEVEL_TYPE_CODE|LEVEL_NUMBER_PREFIX|LEVEL_NUMBER|LEVEL_NUMBER_SUFFIX|NUMBER_FIRST_PREFIX|NUMBER_FIRST|NUMBER_FIRST_SUFFIX|NUMBER_LAST_PREFIX|NUMBER_LAST|NUMBER_LAST_SUFFIX|STREET_LOCALITY_PID|LEGAL_PARCEL_ID|LEVEL_GEOCODED_CODE|LONGITUDE|LATITUDE"
  for p:= range LstObjStateLya[s].LstObjPostcodeLya{
    for sub:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya{
      var data[] string
      data = append(data,template_title)
      for str:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya{
        for dstreet:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].Data{
          data = append(data,LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].Data[dstreet])
        }
      }
      SaveFile(data,LstObjStateLya[s].State + "/" + LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LOCALITY_PID + ".txt")
    }
  }
}

func SaveFile(d[] string, name string){
  f, err := os.Create(name)
    if err != nil {
        fmt.Println(err)
                f.Close()
        return
    }
    //d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}
    //fmt.Fprintln(f, "YESSS")
    for _, v := range d {

        fmt.Fprintln(f, v)
        if err != nil {
            fmt.Println(err)
            return
        }
    }
    err = f.Close()
    if err != nil {
        fmt.Println(err)
        return
    }
    fmt.Println("file written successfully")
}

func AppendGeoData(s int,data string, wg *sync.WaitGroup) bool{
  defer wg.Done()
  data = strings.TrimSpace(data)
  var r_split = strings.Split(data,"|")
  if len(r_split) < 22{
    return false
  }else{
    var found = false
    for p:= range LstObjStateLya[s].LstObjPostcodeLya{
      for sub:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya{
        for str:= range LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya{
          var str_loc string = LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].STREET_LOCALITY_PID
          if str_loc == r_split[19]{
            LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].Data =
            append(
              LstObjStateLya[s].LstObjPostcodeLya[p].LstObjSuburbLya[sub].LstObjStreetsLya[str].Data,
              data)
            found = true
            return true
            break
          }
          if found{
            break
          }
        }
        if found{
          break
        }
      }
      if found{
        break
      }
    }
  }
  return true
}

func main(){
  SetupStreetSearch()
  SetupGeoLocation()
  fmt.Println("Simple interest calculation")
  p := 5000.0
  r := 10.0
  t := 1.0

  si := simpleinterest.Calculate(p, r, t)
  siwhat := what.AgainCalculate(p, r, t)
  //simpleinterest.what()
  //silol := simpleinterest.lolCalculate(p, r, t)
  fmt.Println("Simple interest is", si)
  fmt.Println("Simple interest is", siwhat)
  //fmt.Println("Simple interest is", silol)
  //HtmlOutput.Hello()
  //d := []string{"Welcome to the world of Go1.", "Go is a compiled language.", "It is easy to learn Go."}
  //SaveFile(d,"TAS/id.txt")
  //SetupStateHtmlTemplates()

  //SetupGeoBulk()
  /*for s := 1; s < len(LstObjStateLya); s++ {
    SetupGeoBulk(s)
  }*/
    f := what.Newfoo()
    fmt.Println(f)
    //f.Bar()

  //  what.NewFoo().Bar()
  /*tmpobj := what.Newfoo(){
    dog:"poo",
  }*/
  fmt.Println("Lotyouraddress running")
  http.HandleFunc("/",handlerFunc)
  http.Handle("/web/", http.FileServer(http.Dir(*root)))
  http.ListenAndServe(":8120",nil)
}
