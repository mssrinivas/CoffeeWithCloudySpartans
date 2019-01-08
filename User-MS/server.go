package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"github.com/codegangsta/negroni"
	// "github.com/streadway/amqp"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	// "github.com/satori/go.uuid"
	// "gopkg.in/mgo.v2"
	// "gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"time"
	"strings"
	"regexp"
	"errors"
	// "strconv"
)

var debug = true
var cluster1_elb = "http://internal-Cluster1-ELB-1830172783.us-west-1.elb.amazonaws.com:80"
var cluster2_elb = "http://internal-Test-ELB-1557580805.us-west-2.elb.amazonaws.com:80"

type Client struct {
	Endpoint string
	*http.Client
}

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

// Create a new client
func NewClient(server string) *Client {
	return &Client{
		Endpoint: server,
		Client:   &http.Client{Transport: tr},
	}
}


// NewServer configuration 
func NewServer() *negroni.Negroni {
	formatter := render.New(render.Options{
		IndentJSON: true,
	})
	n := negroni.Classic()
	mx := mux.NewRouter()
	initRoutes(mx, formatter)
	n.UseHandler(mx)
	return n
}

func getELB(UserId string)(c *Client){
	var validID = regexp.MustCompile(`^[a-n]|^[A-N]`)
	if validID.MatchString(UserId){
		c := NewClient(cluster1_elb)
		return c 
	} else {
		c := NewClient(cluster2_elb)
		return c
	}

}

func init() {

	// Riak KV Setup	

	elb2 := NewClient(cluster2_elb)
	msg, err := elb2.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server4: ", msg)		
	}

	elb1 := NewClient(cluster1_elb)
	msg, err = elb1.Ping( )
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Riak Ping Server5: ", msg)		
	}
	
}

func (c *Client) Ping() (string, error) {
	resp, err := c.Get(c.Endpoint + "/ping" )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return "Ping Error!", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/ping => " + string(body)) }
	return string(body), nil
}

// func (c *Client) Ping() (string, error) {
// 	fmt.Println("inside ping")
// 	resp, err := c.Get(c.Endpoint + "/ping" )
	
// 	if err != nil {
// 		fmt.Println("[RIAK DEBUG] " + err.Error())
// 		return "Ping Error!", err
// 	}
// 	defer resp.Body.Close()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	if debug { 
// 		fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint  + "/ping => " + string(body)) 
// 	}
// 	return string(body), nil
// }

// API Routes
func initRoutes(mx *mux.Router, formatter *render.Render) {
	mx.HandleFunc("/ping", pingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/signup", signupHandler(formatter)).Methods("POST")
	mx.HandleFunc("/login", loginHandler(formatter)).Methods("POST")
	mx.HandleFunc("/allusers", allusersHandler(formatter)).Methods("GET")
}

// API Ping Handler
func pingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"API version 1.0 alive!"})
	}
}

func signupHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var ord user
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&ord)

		if error != nil{
			ErrorWithJSON(w, "incorrect body", http.StatusBadRequest)
			fmt.Println("[HANDLER DEBUG]", error.Error())
			return 
		}

		requestbody, _ := json.Marshal(ord)
		
		elb2 := getELB(ord.UserId)

		chk_user, _ := elb2.GetUser(ord.UserId)

		if (ord.UserId == chk_user.UserId){
			// stat_ok := "ok"
			formatter.JSON(w, http.StatusOK, "user exists")
			fmt.Println("user exists")
		} else  {
			value_res, error := elb2.RegisterUser(ord.UserId,string(requestbody))
		
			if error != nil {
				log.Fatal(error)
				formatter.JSON(w, http.StatusBadRequest, error)
			} else {
				formatter.JSON(w, http.StatusOK, value_res)
			}
		}

	}
}

func loginHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var ord user
		decoder := json.NewDecoder(req.Body)
		error := decoder.Decode(&ord)

		if error != nil{
			ErrorWithJSON(w, "incorrect body", http.StatusBadRequest)
			fmt.Println("[HANDLER DEBUG]", error.Error())
			return 
		}

		elb2 := getELB(ord.UserId)

		user_details, error := elb2.GetUser(ord.UserId)

		// fmt.Println("error inside login ", error.Error())

		if error != nil {
			// log.Fatal(error)
			formatter.JSON(w, http.StatusBadRequest, error)
		} else {
			if (ord.Password == user_details.Password){
				// stat_ok := "ok"
				formatter.JSON(w, http.StatusOK, "SUCCESS")
				fmt.Println("Login Successful")
			} else {
				formatter.JSON(w, http.StatusBadRequest, "INVALID CREDENTIALS")
				fmt.Println("Invalid credentials")
			}

		}		
	}
}

func allusersHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		elb2 := NewClient(cluster1_elb)
		list_users, error := elb2.GetAllUsers()

		if error != nil {
			log.Fatal(error)
			formatter.JSON(w, http.StatusBadRequest, error)
		} else {
			formatter.JSON(w, http.StatusOK, list_users)	
		}
	}
}

func (c *Client) RegisterUser(key string, reqbody string) (user, error) {
	var ord_nil = user {}

	resp, err := c.Post(c.Endpoint + "/buckets/users/keys/"+key+"?returnbody=true", 
	"application/json", strings.NewReader(reqbody) )
	
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return ord_nil, err
	}	
 	defer resp.Body.Close()
 	body, err := ioutil.ReadAll(resp.Body)
 	if debug { 
		 fmt.Println("[RIAK DEBUG] POST: " + c.Endpoint + "/buckets/users/keys/"+key+"?returnbody=true => "  + string(body)) 
		}
 	var place user
 	msg1 := json.Unmarshal(body, &place); 
 	if msg1 != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", msg1)
		return ord_nil, msg1
	}	
	fmt.Println("place", place)
 	return place, nil
}

type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	return &errorString{text}
}

func (c *Client) GetUser(key string) (user, error) {
	var ord_nil = user {}
	
	resp, err := c.Get(c.Endpoint + "/buckets/users/keys/"+key)

	fmt.Println(resp.StatusCode)
	

	if err != nil{
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return ord_nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	// var usrerr error

	if resp.StatusCode== 404{
		return ord_nil, errors.New("user not found")
	}

	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/maps/keys/"+key +" => " + string(body)) }
	var ord = user { }

	if err := json.Unmarshal(body, &ord); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return ord_nil, err
	}
	fmt.Println("ord is",ord)
	return ord, nil
}


func (c *Client) GetAllUsers() ([]string, error) {
	var all_keys []string
	
	resp, error := c.Get(c.Endpoint + "/buckets/users/keys?keys=true")
	
	if error != nil {
		fmt.Println("[RIAK DEBUG] " + error.Error())
		return all_keys, error
	}
	
	defer resp.Body.Close()

	rbody, error := ioutil.ReadAll(resp.Body)

	if debug { 
		fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/users/keys/keys?true => " + string(rbody)) 
	}

	// var ord1 = user { }
	var all_keys_list Keys

	if err := json.Unmarshal(rbody, &all_keys_list); err != nil {		
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return all_keys_list.Keys, err
	}
	fmt.Println("Keys are",all_keys_list.Keys)
	return all_keys_list.Keys, nil
	
	// var output user
	//  msg := json.Unmarshal(body, &output); 
	//  fmt.Println("output", output)

	// if msg != nil {
	// 	fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", msg)
	// 	return ord_nil, msg
	// }	
	// fmt.Println("ord is",output)
	// return msg, nil
}

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "{message: %q}", message)
}
