package main

import (
	"fmt"
	"io"
	"strconv"
	"bytes"
	"log"
	"encoding/json"
	"net/http"
)


type measurements struct {
	Measurements []int `json:"measurements"`
}

func crop(w http.ResponseWriter, r *http.Request){
	if r.URL.Path != "/"{
		http.Error(w,"404 NOT FOUND!!!", http.StatusNotFound)
		return
	}

switch r.Method {
case "GET":
	http.ServeFile(w,r,"predicts.html")
case "POST":
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w,"ParseForm() err",err)
		return 
	}
	fmt.Fprintf(w,"Post request values r.postfrom = %v\n",r.PostForm)
	p, err1 := strconv.Atoi(r.FormValue("p"))
	if err1 != nil {
		panic(err1) 
	}

	n, err2 := strconv.Atoi(r.FormValue("n"))
	if err2 != nil {
		panic(err2) 
	}

	k, err3 := strconv.Atoi(r.FormValue("k"))

	if err3 != nil {
		panic(err3) 
	}

	temperature, err4 := strconv.Atoi(r.FormValue("temperature"))
	if err4 != nil {
		panic(err4) 
	}

	humidity, err5 := strconv.Atoi(r.FormValue("humidity"))
	if err5 != nil {
		panic(err5) 
	}

	ph, err6 := strconv.Atoi(r.FormValue("ph"))

	if err6 != nil {
		panic(err6) 
	}

	rainfall, err7 := strconv.Atoi(r.FormValue("rainfall"))

	if err7 != nil {
		panic(err7) 
	}

	fmt.Fprintf(w,"Phosphorus = %s\n",p)
	fmt.Fprintf(w,"Pottasium = %s\n",k)
	fmt.Fprintf(w,"Temperature = %s\n",temperature)
	fmt.Fprintf(w,"Humidity = %s\n",humidity)
	fmt.Fprintf(w,"PH = %s\n",ph)
	fmt.Fprintf(w,"Rainfall = %s\n",rainfall)

	sendData(p,k,n,temperature,humidity,ph,rainfall,w)
	
	
	

default:
	fmt.Fprintf(w,"Only get and post allowed")
	return
}
}

func sendData(P int, N int, K int, Temperature int, Humidity int, PH int, Rainfall int,w http.ResponseWriter,){
	const url = "http://localhost:8000/add/"
	measurement := []measurements{
		{Measurements:[]int{P,N,K,Temperature,Humidity,PH,Rainfall}},
	}

	

	requests, err := json.Marshal(measurement)
	if err != nil {
		panic(err) 
	}
	
	res, err1 := http.NewRequest("POST",url,bytes.NewBuffer(requests))
	res.Header.Add("Content-Type","application/json")
	if err1 != nil {
		panic(err1) 
	}

	

	client := &http.Client{}
	res1, err := client.Do(res)
	if err != nil {
		panic(err)
	}

	defer res1.Body.Close()
	
	b, err3 := io.ReadAll(res1.Body)

	if err3 != nil {
		log.Fatalln(err3)
	}

	fmt.Fprintf(w,"Crop is:%s\n",string(b))
	

}

func main(){
	http.HandleFunc("/", crop)
	fmt.Printf("Starting server!!!\n")
	if err := http.ListenAndServe(":9000",nil);  err != nil{
		log.Fatal(err)
	}
}