package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"
)

type QRCodeRequest struct {
	DevToken     string `json:"token"`
	UserID       string `json:"uid"`
	UserName     string `json:"uname"`
	UserPassword string `json:"utoken"`
	APIVersion   int    `json:"v"`
}

type QRCodeResponse struct {
	ErrorCode int                 `json:"code"`
	Message   string              `json:"message"`
	Result    bool                `json:"result"`
	Data      *QRCodeResponseData `json:"data"`
}

type QRCodeResponseData struct {
	QR   string `json:"qr"`
	Code string `json:"code"`
}

func main() {
	token := flag.String("token", "", "Developer token from the developer dashboard")
	flag.Parse()

	if *token == "" {
		panic("No token provided")
	}

	http.HandleFunc("/qr", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now().Local(), r.RemoteAddr, r.URL.RequestURI())

		q := r.URL.Query()
		uname := q.Get("username")
		utoken := q.Get("password")
		if uname == "" || utoken == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		qrUrl, err := getQRCodeURL(token, uname, utoken)
		if err != nil {
			panic(err)
		}

		w.Write([]byte(qrUrl))

	})

	if err := http.ListenAndServe(":3001", nil); err != nil {
		panic(err)
	}
}

func getQRCodeURL(token *string, uname string, utoken string) (string, error) {
	qrReq := QRCodeRequest{DevToken: *token, UserName: uname, UserID: "test", UserPassword: utoken, APIVersion: 2}

	req, err := json.Marshal(qrReq)
	if err != nil {
		return "", errors.New("cannot marshal request")
	}

	resp, err := http.Post("https://api.lovense-api.com/api/lan/getQrCode", "application/json", bytes.NewReader(req))
	if err != nil {
		return "", errors.New("request failed")
	}

	lRes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("cannot read response")
	}

	var lovenseRes QRCodeResponse
	if err := json.Unmarshal(lRes, &lovenseRes); err != nil {
		return "", errors.New("cannot unmarshal response")
	}

	return lovenseRes.Data.QR, nil
}
