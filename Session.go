package dccli

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type Session struct {
	Account    Account
	isLoggedin bool
	NoLogID    string
	NoLogPW    string
	Appid      string
	Apptoken   string
	NowGallID  string
	NowPostNo  int
}

type AppCheckstruct struct {
	Result        bool   `json:"result"`
	Ver           string `json:"ver"`
	Notice        bool   `json:"notice"`
	Notice_update bool   `json:"notice_update"`
	Date          string `json:"date"`
}

type Account struct {
	Result           bool   `json:"result"`
	User_id          string `json:"user_id"`
	User_no          string `json:"user_no"`
	Name             string `json:"name"`
	Is_adult         string `json:"is_adult"`
	Is_dormancy      int    `json:"is_dormancy"`
	Otp_token        string `json:"otp_token"`
	Is_gonick        int    `json:"is_gonick"`
	Is_security_code string `json:"is_security_code"`
	Auth_change      int    `json:"auth_change"`
	Stype            string `json:"stype"`
	Pw_campaign      int    `json:"pw_campaign"`
}

func HashedURLmake(gallid string, appid string) string {
	input := []byte(
		fmt.Sprintf("https://app.dcinside.com/api/gall_list_new.php?id=%s&page=1&app_id=%s",
			gallid,
			appid,
		),
	)
	return fmt.Sprintf("https://app.dcinside.com/api/redirect.php?hash=%s", base64.StdEncoding.EncodeToString(input))
}

func Base64EncodeLink(input string) string {
	return fmt.Sprintf("https://app.dcinside.com/api/redirect.php?hash=%s", base64.StdEncoding.EncodeToString([]byte(input)))
}

func (s *Session) GetAppID() error {
	//{
	//	"fid": "fT-9GN8ASwOa9ihWpuokdn",
	//	"appId": "1:477369754343:android:d2ffdd960120a207727842",
	//	"authVersion": "FIS_v2",
	//	"sdkVersion": "a:17.0.2"
	//}
	//{
	//	"name": "projects/477369754343/installations/fT-9GN8ASwOa9ihWpuokdn",
	//	"fid": "fT-9GN8ASwOa9ihWpuokdn",
	//	"refreshToken": "3_AS3qfwKZ1zsz4C0dvSZdg9CBYSKG4MBEoYrNuKiGg-908_yTBGRkxTD1qeI_vuzCOGb5GSj8O8cxdCwWFXT0fEBlPNEmAkbPV5ZFOVRg-yQKojU",
	//	"authToken": {
	//	"token": "eyJhbGciOiJFUzI1NiIsInR5cCI6IkpXVCJ9.eyJhcHBJZCI6IjE6NDc3MzY5NzU0MzQzOmFuZHJvaWQ6ZDJmZmRkOTYwMTIwYTIwNzcyNzg0MiIsImV4cCI6MTY3NjIwOTYxOCwiZmlkIjoiZlQtOUdOOEFTd09hOWloV3B1b2tkbiIsInByb2plY3ROdW1iZXIiOjQ3NzM2OTc1NDM0M30.AB2LPV8wRgIhALHo8OYiKb41UxwuCyjPLJ21qQQM2Ofme63jdbQc0YzHAiEAvbRYIf13I0NqMmHBe5iRz7-Hglcx0-RfCf0sOi8XWnw",
	//		"expiresIn": "604800s"
	//}
	//token=fT-9GN8ASwOa9ihWpuokdn:APA91bHW2DbvpDTeJxUA_ACwoLzPkCfJpWqj5N2Eb9H7gYz9D28e1jJH_RRXZoDDMKClZSlXXVosI10BlHGcFgOg1dkkJRm8qCaU9Fci7V2q9ZSRSefw0tA7xW1A_3jl8UU5GG3_uLNL
	res, err := http.Get("http://json2.dcinside.com/json0/app_check_A_rina.php")
	if err != nil {
		log.Fatal(err)
	}

	bod, _ := io.ReadAll(res.Body)
	//fmt.Println(string(bod))
	var Appc []AppCheckstruct
	err = json.Unmarshal(bod, &Appc)
	if err != nil {
		log.Fatal(err)
	}
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("dcArdchk_%s", Appc[0].Date))) // value token calculated
	res, err = http.PostForm(
		"https://msign.dcinside.com/auth/mobile_app_verification",
		url.Values{
			"value_token":  {fmt.Sprintf("%x", h.Sum(nil))},
			"signature":    {"ReOo4u96nnv8Njd7707KpYiIVYQ3FlcKHDJE046Pg6s="},
			"client_token": {"fT-9GN8ASwOa9ihWpuokdn:APA91bHW2DbvpDTeJxUA_ACwoLzPkCfJpWqj5N2Eb9H7gYz9D28e1jJH_RRXZoDDMKClZSlXXVosI10BlHGcFgOg1dkkJRm8qCaU9Fci7V2q9ZSRSefw0tA7xW1A_3jl8UU5GG3_uLNL"},
		},
	)
	if err != nil {
		return errors.New("Error GetAppID function")
	}
	bod, _ = io.ReadAll(res.Body)
	s.Appid = gjson.Get(string(bod), "app_id").String()
	return nil
}

func (s *Session) Login(id string, pw string) error {
	if s.isLoggedin == true {
		return errors.New("Already logged in")
	}
	rr := url.Values{}
	rr.Add("user_id", id)
	rr.Add("user_pw", pw)
	rr.Add("client_token", "fT-9GN8ASwOa9ihWpuokdn:APA91bHW2DbvpDTeJxUA_ACwoLzPkCfJpWqj5N2Eb9H7gYz9D28e1jJH_RRXZoDDMKClZSlXXVosI10BlHGcFgOg1dkkJRm8qCaU9Fci7V2q9ZSRSefw0tA7xW1A_3jl8UU5GG3_uLNL")
	rr.Add("mode", "login_normal")

	req, err := http.NewRequest(
		"POST",
		"https://msign.dcinside.com/api/login",
		strings.NewReader(rr.Encode()),
	)
	if err != nil {
		return errors.New("Error Making Request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("user-agent", "dcinside.app")
	req.Header.Set("Host", "msign.dcinside.com")
	req.Header.Set("referer", "http://www.dcinside.com")
	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		return errors.New("Error Posting Request")
	}
	bod, _ := io.ReadAll(res.Body)
	// fmt.Println(string(bod))
	var account Account
	e := json.Unmarshal(bod, &account)
	if e != nil {
		return errors.New("Error while parsing json")
	}
	// fmt.Println(account)
	if account.Result == true {
		s.isLoggedin = true
		s.Account = account
		return nil
	} else {
		return errors.New("Failed to Log in")
	}
} // 객체지향 추가 완?료

func New() *Session {
	p := &Session{}
	p.isLoggedin = false
	return p
}

//func (s *Session) FetchFCMToken() {
//	r, e := http.NewRequest("POST", "https://firebaseinstallations.googleapis.com/v1/projects/dcinside-b3f40/installations", nil)
//	if e != nil {
//		panic(e)
//	}
//
//	r.Header.Set("accept", "application/json")
//	r.Header.Set("accept-encoding", "gzip")
//	r.Header.Set("cache-control", "no-cache")
//	r.Header.Set("connection", "Keep-Alive")
//	r.Header.Set("content-encoding", "gzip")
//	r.Header.Set("host", "firebaseinstallations.googleapis.com")
//	r.Header.Set("user-agent", "Dalvik/2.1.0 (Linux; U; Android 13; Pixel 5 Build/TP1A.221105.002)")
//	r.Header.Set("x-android-cert", "43BD70DFC365EC1749F0424D28174DA44EE7659D")
//	r.Header.Set("x-android-package", "com.dcinside.app.android")
//	r.Header.Set("x-firebase-client", "H4sIAAAAAAAAAKtWykhNLCpJSk0sKVayio7VUSpLLSrOzM9TslIyUqoFAFyivEQfAAAA")
//	r.Header.Set("x-goog-api-key", "AIzaSyDcbVof_4Bi2GwJ1H8NjSwSTaMPPZeCE38")
//	b := bytes.NewBuffer([]byte(`{
//  "fid": "f7RXAqYIR6iACLGVP06qb4",
//  "appId": "1:477369754343:android:d2ffdd960120a207727842",
//  "authVersion": "FIS_v2",
//  "sdkVersion": "a:17.0.2"}`))
//	r.Body = io.NopCloser(b)
//	client := &http.Client{}
//	res, err := client.Do(r)
//	if err != nil {
//		panic(err)
//	}
//	bod, _ := io.ReadAll(res.Body)
//	fmt.Println(string(bod)) // Must get fid, appid,
//}

// 위에꺼 FCM토큰관련한건데 아직 작동도안하고 수정할거많아서 일단 주석처리함