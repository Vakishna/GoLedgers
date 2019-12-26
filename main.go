/*
VAKISHNA THAYALAN
1 Ursula Close,
Wheelers Hill, VIC 3150,
Australia
DATE: Wednesday 25 December 2019
*/

package main

import (
	"cloud.google.com/go/civil"
	"encoding/json"
	"fmt"
	"github.com/auth0-community/auth0"
	_ "github.com/auth0-community/auth0"
	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/dgrijalva/jwt-go/request"
	handlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	jose "gopkg.in/square/go-jose.v2"
	"log"
	"net/http"
	"os"
	"time"
)

/* Global String Secret */
var mySigningKey = []byte("ZjzMgq%nvU1zW#5C$b64za&*WqH0vS")

type AccountClass int

const (
	Equity    AccountClass = 0
	Asset     AccountClass = 1
	Liability AccountClass = 2
	Expense   AccountClass = 3
	Revenue   AccountClass = 4
)

// CODE REFERENCE - https://github.com/auth0-community/auth0-go/blob/master/example/main.go
type Product struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

// CODE REFERENCE - https://github.com/auth0-community/auth0-go/blob/master/example/main.go
var products = []Product{
	Product{Id: 1, Name: "Hover Shooters", Slug: "hover-shooters", Description: "Shoot your way to the top on 14 different hoverboards"},
	Product{Id: 2, Name: "Ocean Explorer", Slug: "ocean-explorer", Description: "Explore the depths of the sea in this one of a kind underwater experience"},
	Product{Id: 3, Name: "Dinosaur Park", Slug: "dinosaur-park", Description: "Go back 65 million years in the past and ride a T-Rex"},
	Product{Id: 4, Name: "Cars VR", Slug: "cars-vr", Description: "Get behind the wheel of the fastest cars in the world."},
	Product{Id: 5, Name: "Robin Hood", Slug: "robin-hood", Description: "Pick up the bow and arrow and master the art of archery"},
	Product{Id: 6, Name: "Real World VR", Slug: "real-world-vr", Description: "Explore the seven wonders of the world in VR"},
}

type EmployeeType int

const (
	Permanent  EmployeeType = 0
	Casual     EmployeeType = 1
	Apprentice EmployeeType = 2
	Agency     EmployeeType = 3
	Contract   EmployeeType = 4
)

type Employee struct {
	Id           int            `json:"id"`
	GivenName    string         `json:"first_name"`
	Surname      string         `json:"last_name"`
	EmployeeType *EmployeeType  `json:"employee_type"`
	TFN          int16          `json:"tfn"`
	DateOfBirth  civil.DateTime `json:"date_of_birth"`
	Address      *Address       `json:"address"`
}

type Account struct {
	Id           int           `json:"id"`
	FriendlyName string        `json:"friendly_name"`
	BalanceDR    float64       `json:"bal_dr"`
	BalanceCR    float64       `json:"bal_cr"`
	Description  string        `json:"description"`
	Nature       *AccountClass `json:"nature"`
}

type Address struct {
	Id               int    `json:"id"`
	PrimaryRecipient string `json:"primary_recipient"`
	UnitOrPrefix     string `json:"unit_or_prefix"`
	StreetNumber     int    `json:"street_number"`
	StreetName       string `json:"street_name"`
	Suburb           string `json:"suburb"`
	State            string `json:"state"`
	PostCode         int    `json:"post_code"`
}

type PurchaseItems struct {
	Quantity     int            `json:"quantity"`
	Description  string         `json:"description"`
	UnitPrice    float64        `json:"unit_price"`
	Gst          float64        `json:"gst"`
	ItemTotal    float64        `json:"item_total"`
	PurchaseDate civil.DateTime `json:"purchased_date"`
}

type Purchase struct {
	Id        int              `json:"id"`
	Merchant  string           `json:"merchant"`
	Abn       int64            `json:"abn"`
	Acn       int64            `json:"acn"`
	Items     []*PurchaseItems `json:"purchase_item"`
	CRAccount *Account         `json:"credit_account"`
	DRAccount *Account         `json:"debit_account"`
}

type Supplier struct {
	Id             int              `json:"id"`
	Abn            int64            `json:"abn"`
	Acn            int64            `json:"acn"`
	InventoryItems []*InventoryItem `json:"inventory_items"`
	Address        *Address         `json:"address"`
}

type InventoryItem struct {
	Id            int      `json:"id"`
	ItemName      string   `json:"item_name"`
	Description   string   `json:"description"`
	Brand         string   `json:"item_brand"`
	Barcode       string   `json:"barcode"`
	Supplier      Supplier `json:"supplier"`
	SupplierPrice float64  `json:"supplier_price"`
	CurrentStock  int      `json:"current_stock"`
	LifeTimeSales int      `json:"lifetime_sales"`
}

type Inventory struct {
	Id             int              `json:"id"`
	InventoryItems []*InventoryItem `json:"inventory_item"`
	Description    string           `json:"description"`
	AssocBranch    *Address         `json:"assoc_branch"`
}

type SalesItems struct {
	Id           int            `json:"id"`
	UnitPrice    float64        `json:"unit_price"`
	Gst          float64        `json:"gst"`
	ItemTotal    float64        `json:"item_total"`
	PurchaseDate civil.DateTime `json:"purchased_date"`
}

// TODO: When a sale is processed, decrement CurrentStock
type Sales struct {
	Id            int            `json:"id"`
	CustomerABN   int64          `json:"customer_abn"`
	CustomerACN   int64          `json:"customer_acn"`
	InventoryItem *InventoryItem `json:"inventory_item"`
}

type GrossPaymentType string

const (
	P GrossPaymentType = "P"
	H GrossPaymentType = "H"
	N GrossPaymentType = "N"
)

type GrossPayment struct {
	Type  *GrossPaymentType `json:"gross_pay"`
	LumpA float64           `json:"lump_a"`
	LumpB float64           `json:"lump_b"`
	LumpC float64           `json:"lump_c"`
	LumpD float64           `json:"lump_d"`
}

// TODO: COMPLETE PAYG PAYMENT STRUCT
type IPayGPayment struct {
	Id               int64         `json:"id"`
	Employee         *Employee     `json:"employee"`
	GrossPay         int64         `json:"gross_pay"`
	CDEPPay          int64         `json:"cdep_pay"`
	FringePay        int64         `json:"fringe_pay"`
	PayPeriodStart   civil.Date    `json:"pay_period_start"`
	PayPeriodEnd     civil.Date    `json:"pay_period_end"`
	TotalTaxWithheld int64         `json:"total_tax"`
	FBTAAExempt      bool          `json:"fbtaa_exempt"`
	LessAnnuity      float64       `json:"less_annuity"`
	GrossPayment     *GrossPayment `json:"gross_pay"`
	//
}

type Transactions struct {
}

/* ********************************* INTERFACES **************************************** */

type Payroll interface {
	PayEmployee(employee Employee)
}

/* ********************************** HANDLERS ****************************************** */
// Current TO-DO Methods substitute
var ToDoHandle = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var EmployeeHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the employee section"))
})

var ProductHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(products)
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(payload))
})

var AccountHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("The list of accounts will be displayed here"))
})

var SupplierHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the supplier handler"))
})

var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["name"] = true
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	tokenString, _ := token.SignedString(mySigningKey)
	w.Write([]byte(tokenString))
})

// END OF HANDLERS

/* ********************************** MIDDLEWARE *************************************** */

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return mySigningKey, nil
	},
	SigningMethod: jwt.SigningMethodHS256,
})

/* ************************************************************************************** */
func main() {

	// Instantiate a gorilla/mux router
	r := mux.NewRouter()
	// Serve a simple static index page
	r.Handle("/", http.FileServer(http.Dir("./views/")))
	// We will setup our server so we can serve static assest like images, css from the /static/{file} route
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	/*		API ROUTES		*/
	// This is the API Routes
	// /status - this route will be called to verify the API is up and running
	// /accounts - this will be called to retrieve a list of the customer accounts
	// /addresses - which I will use to retrieve a list of business locations and/or branches
	// /banking - this route will be called and display a list of bank related transactions
	// /products - this will be called to retrieve a list of items from the inventory
	// /transactions - this will be called and a list of all transactions will be displayed
	/* The transactions class does not have a post method, it only displays sales, purchases, income and expenses API CALLS */
	// /sales - this will be called to display all the sales as a list
	// /purchases - this will be called to display a list of all the purchases
	// /contacts - this will display all the contacts of type supplier and customer
	// /customers - when this is called this will display all the customers
	// /suppliers - this will be called to display all the suppliers
	// /employees - this will display all inactive and active employees

	// This is a list of all the GET Methods
	r.Handle("/status", ToDoHandle).Methods("GET")
	r.Handle("/accounts", jwtMiddleware.Handler(AccountHandler)).Methods("GET")
	r.Handle("/addresses", ToDoHandle).Methods("GET")
	r.Handle("/products", jwtMiddleware.Handler(ProductHandler)).Methods("Get")
	r.Handle("/banking", ToDoHandle).Methods("GET")
	r.Handle("/transactions", ToDoHandle).Methods("GET")
	r.Handle("/sales", ToDoHandle).Methods("GET")
	r.Handle("/purchases", ToDoHandle).Methods("GET")
	r.Handle("/contacts", ToDoHandle).Methods("GET")
	r.Handle("/suppliers", SupplierHandler).Methods("Get")
	r.Handle("/get-token", GetTokenHandler).Methods("GET")
	r.Handle("/employees", EmployeeHandler).Methods("GET")

	// This is a list of all the POST Methods
	r.Handle("/account", ToDoHandle).Methods("POST")
	r.Handle("/address", ToDoHandle).Methods("POST")
	r.Handle("/product", ToDoHandle).Methods("POST")
	r.Handle("/banking", ToDoHandle).Methods("POST")
	r.Handle("/sale", ToDoHandle).Methods("POST")
	r.Handle("/purchase", ToDoHandle).Methods("POST")
	r.Handle("/contact", ToDoHandle).Methods("POST")
	r.Handle("/supplier", ToDoHandle).Methods("POST")

	// Token Handler
	// Run the program on port: 8090 and pass in the mux router
	log.Fatal(http.ListenAndServe(":8040", handlers.LoggingHandler(os.Stdout, r)))
}

/* **************************************---MONGODB DATABASE DATA PERSISTENCE---********************************** */

// DATE: FRIDAY 27th DECEMBER
// AUTHOR: VAKISHNA THAYALAN

func ConnectMongo() {

}

/************************************** AUTHENTICATION MIDDLEWARE **********************************************/
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := []byte("")
		secretProvider := auth0.NewKeyProvider(secret)
		audience := []string{"{YOUR-AUTH0-API-AUDIENCE}"}

		configuration := auth0.NewConfiguration(secretProvider, audience, "", jose.HS256)
		validator := auth0.NewValidator(configuration, nil)
		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

/**********************************************************************************************************************/
/*
 * Copyright (c) 2020 - Vakishna Thayalan
 * LICENCE: MPL-2.0
 */
