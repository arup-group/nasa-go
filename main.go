package main

import (
    "fmt"
    "net/http"
    "log"
    "io/ioutil"
)

const NASA_API_ROOT="https://api.nasa.gov/"
const API_KEY="WlA0h7verog2bZbd7yHhWxmxOJsfCHfy5CCaRhwx"
var client = http.Client{}

func hello(w http.ResponseWriter, req *http.Request) {
    fmt.Fprintf(w, "Hello World\n")
}

func make_request(url string)string{
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

  return string(bytes)
}

func insight(w http.ResponseWriter, req *http.Request){
    const url=NASA_API_ROOT+"insight_weather/?api_key="+API_KEY+"&feedtype=json&ver=1.0"
    res:=make_request(url)
    fmt.Fprintf(w, res)
}

func earth_imagery(w http.ResponseWriter, req *http.Request){
    query := req.URL.Query()
    valid := true
    lat := query.Get("lat")
    date := query.Get("date")
    lon := query.Get("lon")
    dim := query.Get("dim")

    if lat == "" {
        fmt.Println("lat not present")
        valid = false
    }
    if lon == "" {
        fmt.Println("lon not present")
        valid = false
    }
    if date == ""{
        fmt.Println("date not present")
        valid = false
    }

    fmt.Println(lat)
    fmt.Println(lon)
    fmt.Println(date)

    if valid {
      url := NASA_API_ROOT+"planetary/earth/imagery/?api_key="+API_KEY+"&lat="+lat+"&lon="+lon+"&date="+date
      fmt.Println(url)
      if dim == "" {
          fmt.Println("dim not present")
      }else{
        url += ("&dim=" + dim)
      }

      reqImg, err := client.Get(url)
      if err != nil {
        fmt.Fprintf(w, "Error %d", err)
        return
      }
      buffer := make([]byte, reqImg.ContentLength)
      // ioutil.ReadFull(reqImg.Body, buffer)
      reqImg.Body.Close()

      w.Header().Set("Content-Length", fmt.Sprint(reqImg.ContentLength))
      w.Header().Set("Content-Type", reqImg.Header.Get("Content-Type"))
      w.Write(buffer)

      // fmt.Fprintf(w, buffer)
    }
}

func headers(w http.ResponseWriter, req *http.Request) {

    for name, headers := range req.Header {
        for _, h := range headers {
            fmt.Fprintf(w, "%v: %v\n", name, h)
        }
    }
}

func main() {
    fmt.Println("Running on 8090\n")

    http.HandleFunc("/", hello)
    http.HandleFunc("/hello", hello)
    http.HandleFunc("/headers", headers)
    http.HandleFunc("/insight", insight)
    http.HandleFunc("/earth-imagery", earth_imagery)


    http.ListenAndServe(":8090", nil)
}
