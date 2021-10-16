package main

import (
	"crypto/sha512"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/disharjayanth/golangCart/data"
	"github.com/joho/godotenv"
)

var temp *template.Template
var err error

func init() {
	godotenv.Load()
	temp, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Error parsing glob of template files:", err)
	}
}

func cartHandler(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "cart.html", nil)
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	order := data.Order{}
	order.Items = make(map[string]string)

	total := 0
	for phone, cost := range r.URL.Query() {
		fmt.Println(phone, cost[0])
		n, err := strconv.Atoi(cost[0])
		if err != nil {
			fmt.Println("cannot parse string price: ", err)
		}
		total += n

		order.OrderId = "1"
		order.Items[phone] = cost[0]
		order.TotalAmount = total

		fmt.Println("Order:", order)
	}

	order.Store()
	fmt.Println("Toal:", total)
	http.Redirect(w, r, "/user", http.StatusSeeOther)
}

func orderConfirmed(w http.ResponseWriter, r *http.Request) {
	orderDetails := data.Order{}
	orderDetails.Get()
	fmt.Println("Order details:", orderDetails)
	temp.ExecuteTemplate(w, "orderPlaced.html", orderDetails)
}

func userDetails(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		temp.ExecuteTemplate(w, "userDetails.html", nil)
	case "POST":
		phone := r.FormValue("phone")
		cardNumber := r.FormValue("cardNumber")
		cardExpiryMonth := r.FormValue("cardExpiryMonth")
		cardExpiryYear := r.FormValue("cardExpiryYear")
		cardCVV := r.FormValue("cardCVV")

		order := data.Order{}
		order.Get()

		merchant_key := os.Getenv("MERCHANT_KEY")
		salt := os.Getenv("SALT")
		txnid := "s7hhDQVWvbhBdN"
		amount := "240"
		productInfo := "phones"
		firstname := r.FormValue("name")
		email := r.FormValue("email")

		hashString := merchant_key + "|" + txnid + "|" + amount + "|" + productInfo + "|" + firstname + "|" + email + "|||||||||||" + salt

		sha_512 := sha512.New()
		sha_512.Write([]byte(hashString))
		final_key := fmt.Sprintf("%x", sha_512.Sum(nil))

		params := url.Values{}
		params.Add("key", merchant_key)
		params.Add("amount", amount)
		params.Add("txnid", txnid)
		params.Add("firstname", firstname)
		params.Add("email", email)
		params.Add("phone", phone)
		params.Add("productinfo", productInfo)
		params.Add("surl", `https://apiplayground-response.herokuapp.com/`)
		params.Add("furl", `https://apiplayground-response.herokuapp.com/`)
		params.Add("pg", `DC`)
		params.Add("bankcode", `VISA`)
		params.Add("ccname", "demo")
		params.Add("ccnum", cardNumber)
		params.Add("ccexpmon", cardExpiryMonth)
		params.Add("ccexpyr", cardExpiryYear)
		params.Add("ccvv", cardCVV)
		params.Add("txn_s2s_flow", ``)
		params.Add("hash", final_key)

		body := strings.NewReader(params.Encode())
		req, err := http.NewRequest("POST", "https://secure.payu.in/_payment", body)
		if err != nil {
			// handle err
			fmt.Println("could create make request!", err)
			return
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request to payu server: ", err)
			return
		}
		defer resp.Body.Close()

		http.Redirect(w, r, resp.Request.URL.String(), http.StatusSeeOther)
	}
}

func main() {
	var server http.Server
	if os.Getenv("ENV") == "DEV" {
		server = http.Server{
			Addr: "localhost:3000",
		}
	} else {
		server = http.Server{
			Addr: ":" + os.Getenv("PORT"),
		}
	}

	http.Handle("/stylesheet/", http.StripPrefix("/stylesheet", http.FileServer(http.Dir("templates/stylesheet/"))))
	http.HandleFunc("/", cartHandler)
	http.HandleFunc("/user", userDetails)
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/orderconfirmed", orderConfirmed)

	fmt.Println("Server listening @PORT: 3000")
	server.ListenAndServe()
}
