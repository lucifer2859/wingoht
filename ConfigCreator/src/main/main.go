package main

import (
	"github.com/dchest/captcha"
	"log"
	"net/http"
	"encoding/json" 
)

type CaptchaInfo struct {
	Base64			string		`json:"base64"`
	Digit 			string		`json:"digit"`
}

func showFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	
	var result CaptchaInfo
	
	result.Base64, result.Digit = captcha.New()
	bytes, _ := json.Marshal(result)
	
	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func main() {
	http.HandleFunc("/", showFormHandler)
	if err := http.ListenAndServe("192.168.40.47:8666", nil); err != nil {
		log.Fatal(err)
	}
}