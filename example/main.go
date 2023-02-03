package main

import (
	"fmt"
	"log"

	dccli "github.com/emptycan1010/dcgo"
	"github.com/tidwall/gjson"
)

func main() {
	r, e := dccli.GetAppID()
	if e != nil {
		log.Fatalln(e)
	}
	appid := gjson.Get(r, "app_id").String()
	// fmt.Println(appid)
	// print(dccli.AddComment("tsmanga", appid, 1, "aaa", "ㅇㅇ", "1111"))
	// res, e := dccli.GetComment("tsmanga", appid, 1, 1)
	// if e != nil {
	// 	log.Fatalln(e)
	// }
	// fmt.Println(res)
	fmt.Print(dccli.DelPost("tsmanga", appid, 9, "1111"))
}
