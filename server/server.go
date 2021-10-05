package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"

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

func paymentHandler(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/orderconfirmed", http.StatusSeeOther)
}

func orderConfirmed(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "orderPlaced.html", nil)
}

func main() {
	server := http.Server{
		Addr: ":" + os.Getenv("PORT"),
		// Addr: "localhost:3000",
	}

	http.Handle("/stylesheet/", http.StripPrefix("/stylesheet", http.FileServer(http.Dir("templates/stylesheet/"))))
	http.HandleFunc("/", cartHandler)
	http.HandleFunc("/pay", paymentHandler)
	http.HandleFunc("/orderconfirmed", orderConfirmed)

	fmt.Println("Server listening @PORT: 3000")
	server.ListenAndServe()
}
