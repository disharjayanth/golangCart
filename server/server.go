package main

import (
	"crypto/sha512"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/disharjayanth/golangCart/data"
)

var temp *template.Template
var err error

func init() {
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

	// order.Store()
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
		// fmt.Println(name, email, phone, cardNumber, cardExpiryMonth, cardExpiryYear, cardCVV)

		order := data.Order{}
		order.Get()

		merchant_key := "u2mqR7"
		salt := "wp7LclP7ie2H1G23cELOHkwqh0GIjFKJ"
		txnid := "s7hhDQVWvbhBdN"
		amount := "240"
		productInfo := "phones"
		firstname := r.FormValue("name")
		email := r.FormValue("email")

		hashString := merchant_key + "|" + txnid + "|" + amount + "|" + productInfo + "|" + firstname + "|" + email + "|||||||||||" + salt
		fmt.Println("hashString:", hashString)

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
		params.Add("productinfo", `Phones`)
		params.Add("surl", `https://apiplayground-response.herokuapp.com/`)
		params.Add("furl", `https://apiplayground-response.herokuapp.com/`)
		params.Add("pg", `cc`)
		params.Add("bankcode", `cc`)
		params.Add("ccnum", cardNumber)
		params.Add("ccexpmon", cardExpiryMonth)
		params.Add("ccexpyr", cardExpiryYear)
		params.Add("ccvv", cardCVV)
		params.Add("ccname", `undefined`)
		params.Add("txn_s2s_flow", ``)
		params.Add("hash", final_key)
		body := strings.NewReader(params.Encode())
		req, err := http.NewRequest("POST", "https://secure.payu.in/_payment", body)
		if err != nil {
			// handle err
			fmt.Println("could create make request!", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request to payu server: ", err)
			return
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		fmt.Println("Response:", string(b))
		w.Write(b)
	}
}

func main() {
	server := http.Server{
		// Addr: ":" + os.Getenv("PORT"),
		Addr: "localhost:3000",
	}

	http.Handle("/stylesheet/", http.StripPrefix("/stylesheet", http.FileServer(http.Dir("templates/stylesheet/"))))
	http.HandleFunc("/", cartHandler)
	http.HandleFunc("/user", userDetails)
	http.HandleFunc("/order", orderHandler)
	http.HandleFunc("/orderconfirmed", orderConfirmed)

	fmt.Println("Server listening @PORT: 3000")
	server.ListenAndServe()
}

// params value for payU
// params := url.Values{}
// params.Add("key", `JP***g`)
// params.Add("amount", strconv.Itoa(order.TotalAmount))
// params.Add("txnid", `C3nrapLxcTty9R`)
// params.Add("firstname", name)
// params.Add("email", email)
// params.Add("phone", phone)
// params.Add("productinfo", fmt.Sprintf("%#v", order))
// params.Add("surl", `https://apiplayground-response.herokuapp.com/`)
// params.Add("furl", `https://apiplayground-response.herokuapp.com/`)
// params.Add("pg", `cc`)
// params.Add("bankcode", `cc`)
// params.Add("ccnum", cardNumber)
// params.Add("ccexpmon", cardExpiryMonth)
// params.Add("ccexpyr", cardExpiryYear)
// params.Add("ccvv", cardCVV)
// params.Add("ccname", `undefined`)
// params.Add("txn_s2s_flow", ``)
// params.Add("hash", `5acbdf29517ba345f40b38bbea1241c79a8721f33d6b3ee704972095440ec959b7e5b19dd12e106a24ea95773eca484138d096dcac95424c8abc250131eba9f3`)
// body := strings.NewReader(params.Encode())
// req, err := http.NewRequest("POST", "https://test.payu.in/_payment -H", body)
// if err != nil {
// 	fmt.Println("Error making creating request to payu server: ", err)
// 	return
// }

// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// resp, err := http.DefaultClient.Do(req)
// if err != nil {
// 	fmt.Println("Error sending request to payu server: ", err)
// 	return
// }
// defer resp.Body.Close()

// // #2
// params := url.Values{}
// 		params.Add("key", merchant_key)
// 		params.Add("amount", amount)
// 		params.Add("txnid", txnid)
// 		params.Add("firstname", firstname)
// 		params.Add("email", email)
// 		params.Add("phone", phone)
// 		params.Add("productinfo", `Phones`)
// 		params.Add("surl", `https://apiplayground-response.herokuapp.com/`)
// 		params.Add("furl", `https://apiplayground-response.herokuapp.com/`)
// 		params.Add("pg", `cc`)
// 		params.Add("bankcode", `cc`)
// 		params.Add("ccnum", ``)
// 		params.Add("ccexpmon", ``)
// 		params.Add("ccexpyr", ``)
// 		params.Add("ccvv", ``)
// 		params.Add("ccname", ``)
// 		params.Add("txn_s2s_flow", ``)
// 		params.Add("hash", final_key)
// 		body := strings.NewReader(params.Encode())
// 		req, err := http.NewRequest("POST", "https://test.payu.in/_payment -H", body)
// 		if err != nil {
// 			// handle err
// 		}
// 		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
// 		resp, err := http.DefaultClient.Do(req)
// 		if err != nil {
// 			// handle err
// 		}
// 		defer resp.Body.Close()

// 		sb, _ := ioutil.ReadAll(resp.Body)
// 		w.Write(sb)

// url := fmt.Sprintf("https://test.payu.in/_payment/?key=%s&&txnid=%s&&amount=%s&&productinfo=%s&&firstname=%s&&email=%s&&phone=%s&&lastname=%s&&surl=%s&&furl=%s&&hash=%s&&pg=%s&&bankcode=%s",
// 	merchant_key, `s7hhDQVWvbhBdN`, strconv.Itoa(order.TotalAmount), "phones", firstname, email, phone, "mike",
// 	`https://apiplayground-response.herokuapp.com`, `https://apiplayground-response.herokuapp.com`,
// 	final_key, "cc", "cc")

// fmt.Println("Redirect;", url)

// http.Redirect(w, req, url, http.StatusSeeOther)
