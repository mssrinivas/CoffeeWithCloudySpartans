package main

import(
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"time"
	"strings"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"regexp"
)

var debug = true

// var riak1 = "http://10.0.1.218:8098"
// var riak2 = "http://10.0.1.53:8098"
// var riak3 = "http://10.0.1.90:8098"
// var riak4 = "http://10.0.3.50:8098"
// var riak5 = "http://10.0.3.237:8098"
var elb_ca = "http://internal-riak-cart-elb-489823230.us-west-1.elb.amazonaws.com:80"
var elb_or = "http://internal-Riak-project-oregon-961584637.us-west-2.elb.amazonaws.com:80"

var tr = &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}

func NewClient(server string) *Client {
	return &Client{
		Endpoint:  	server,
		Client: 	&http.Client{Transport: tr},
	}
}

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

func (c *Client) CreateOrder(bucket string, key string, value string) (CartItem,error){
	
	nil_cartitem := CartItem {}
	response, err := c.Post(c.Endpoint + "/buckets/"+bucket+"/keys/"+key+"?returnbody=true", "application/json", strings.NewReader(value) )
		
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return nil_cartitem, err
	}	

	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)

	if debug { fmt.Println("[RIAK DEBUG] POST: " + c.Endpoint + "/buckets/"+bucket+"/keys/"+key+"?returnbody=true => "  + string(respBody)) }
	
	var cartitem CartItem
	err = json.Unmarshal(respBody, &cartitem); 
	if err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return nil_cartitem, err
	}
	fmt.Println("Order", cartitem)
	return cartitem, nil
}


func (c *Client) FetchOrder(bucket string, key string) (int, CartItem){

	nil_cartitem := CartItem {}
	count := 0
	response, err := c.Get(c.Endpoint + "/buckets/"+bucket+"/keys/"+key )
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return count, nil_cartitem
	}

	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/"+bucket+"/keys/"+key +" => " + string(respBody)) }
	var cartitem = CartItem { }
	if err := json.Unmarshal(respBody, &cartitem); err != nil {
		fmt.Println("[RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return count, nil_cartitem
	}
	return cartitem.Count, cartitem

}

func (c *Client) UpdateOrder(bucket string, key string, value string) (CartItem,error){

	nil_cartitem := CartItem {}
	request, _  := http.NewRequest("PUT", c.Endpoint + "/buckets/"+bucket+"/keys/"+key+"?returnbody=true", strings.NewReader(value) )
	request.Header.Add("Content-Type", "application/json")

	response, err := c.Do(request)	
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return nil_cartitem, err
	}	
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if debug { fmt.Println("[RIAK DEBUG] PUT: " + c.Endpoint + "/buckets/"+bucket+"/keys/"+key+"?returnbody=true => " + string(respBody)) }
	var cartitem CartItem
	if err := json.Unmarshal(respBody, &cartitem); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return nil_cartitem, err
	}
	return cartitem, nil


}

func (c *Client) DeleteOrder(bucket string, key string) error{

	request, err := http.NewRequest("DELETE", c.Endpoint+"/buckets/"+bucket+"/keys/"+key, nil)
	request.Header.Add("Content-Type", "application/json")

	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return err
	}

	_, err = c.Do(request)
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return err
	}
	return nil
	
}

func (c *Client) FetchKeys(bucket string) ([]string, error){
	
	var nil_keys []string
	response, err := c.Get(c.Endpoint+ "/buckets/"+bucket+"/keys?keys=true")
	if err != nil {
		fmt.Println("[RIAK DEBUG] " + err.Error())
		return nil_keys, err
	}
	defer response.Body.Close()
	respBody, err := ioutil.ReadAll(response.Body)
	if debug { fmt.Println("[RIAK DEBUG] GET: " + c.Endpoint + "/buckets/"+bucket+"/keys?keys=true => " + string(respBody)) }
	var key_list Keys
	if err := json.Unmarshal(respBody, &key_list); err != nil {
		fmt.Println("RIAK DEBUG] JSON unmarshaling failed: %s", err)
		return nil_keys, err
	}
	fmt.Println("Inside FetchKeys %v", key_list.Keys)
	return key_list.Keys, nil
}

func getELB(userid string)(c *Client){
	var validID = regexp.MustCompile(`^[a-n]|^[A-N]`)
	if validID.MatchString(userid){
		c := NewClient(elb_ca)
		return c 
	} else {
		c := NewClient(elb_or)
		return c
	}

}


func init(){

	// c1 := NewClient(riak1)
	// msg, err := c1.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Riak1: ", msg)		
	// }

	// c2 := NewClient(riak2)
	// msg, err = c2.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Riak2: ", msg)		
	// }

	// c3 := NewClient(riak3)
	// msg, err = c3.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Riak3: ", msg)		
	// }

	// c4 := NewClient(riak4)
	// msg, err = c4.Ping()
	// if err != nil {
	// 	log.Fatal(err)
	// } else {
	// 	log.Println("Riak Ping Riak4: ", msg)		
	// }

	 c5 := NewClient(elb_or)
	 msg, err := c5.Ping()
	 if err != nil {
	 	log.Fatal(err)
	 } else {
	 	log.Println("Riak Ping Riak5: ", msg)		
	 }

	c6 := NewClient(elb_ca)
	msg, err = c6.Ping()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("elb Ping Riak: ", msg)		
	}
}


func initRoutes(mx *mux.Router, formatter *render.Render) {

	mx.HandleFunc("/ping", PingHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart/{userid}", GetCartUserHandler(formatter)).Methods("GET")
	mx.HandleFunc("/cart", CreateOrderHandler(formatter)).Methods("POST")
	mx.HandleFunc("/cart", UpdateCartHandler(formatter)).Methods("PUT")
	mx.HandleFunc("/cart/{userid}", DeleteCartHandler(formatter)).Methods("DELETE")
	mx.HandleFunc("/checkout/{userid}", CheckoutCartHandler(formatter)).Methods("GET")
}

func PingHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		formatter.JSON(w, http.StatusOK, struct{ Test string }{"Cart API alive!"})
	}
}


func GetCartUserHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var order_list []CartItem
		input_params := mux.Vars(req)
		user_id := input_params["userid"]
		if user_id == ""{
			formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"UserID missing. Please use /cart/<userid>"})
		} else {
			// c := NewClient(riak1)
			c := getELB(user_id)
			// c := NewClient(elb)
			keys, err := c.FetchKeys(user_id)
			fmt.Println(keys)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"Cart is empty"})
			} else {

				for i:=0; i<len(keys);i++{
					fmt.Println(len(keys))
					_,cartitem := c.FetchOrder(user_id, keys[i])
					order_list = append(order_list,cartitem)
					fmt.Print("%v",order_list)
					}
				formatter.JSON(w, http.StatusOK, order_list)
			}
		}
	}
}


func CreateOrderHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		var cart_order Cart
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&cart_order)
		
		if err != nil {
			log.Fatal(err)
			formatter.JSON(w, http.StatusBadRequest, err)
		}		

		cart_item := cart_order.CartItem

		reqBody, err := json.Marshal(cart_item)

		if err != nil {
			log.Fatal(err)
			formatter.JSON(w, http.StatusBadRequest, err)
		}

		// c := NewClient(riak1)
		// c := NewClient(elb)
		c := getELB(cart_order.UserID)

		item_count, _ := c.FetchOrder(cart_order.UserID, cart_item.ProductID)
		if item_count == 0 {
			response, err := c.CreateOrder(cart_order.UserID, cart_item.ProductID, string(reqBody))
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, response)
			}
		} else {
			cart_item.Count = item_count + 1
			reqBody, err := json.Marshal(cart_item)

			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			}
			response, err := c.UpdateOrder(cart_order.UserID, cart_item.ProductID, string(reqBody))

			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, response)
			}
		}
	}
}

func UpdateCartHandler(formatter *render.Render) http.HandlerFunc {
 	return func(w http.ResponseWriter, req *http.Request) {
	var cart_order Cart
		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&cart_order)
		
		if err != nil {
			log.Fatal(err)
			formatter.JSON(w, http.StatusBadRequest, err)
		}		

		cart_item := cart_order.CartItem

		//reqBody, err := json.Marshal(cart_item)

		//if err != nil {
		//	log.Fatal(err)
		//	formatter.JSON(w, http.StatusBadRequest, err)
		//}

		c := getELB(cart_order.UserID)

		item_count, _ := c.FetchOrder(cart_order.UserID, cart_item.ProductID)
		if item_count == 0 {
				formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"Decrement not allowed"})
		} else {
			cart_item.Count = item_count - 1
			reqBody, err := json.Marshal(cart_item)

			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			}
			response, err := c.UpdateOrder(cart_order.UserID, cart_item.ProductID, string(reqBody))
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, err)
			} else {
				formatter.JSON(w, http.StatusOK, response)
			}
		}
 	}
}

func DeleteCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		input_params := mux.Vars(req)
		user_id := input_params["userid"]
		if user_id == ""{
			formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"UserID missing. Please use /cart/<userid>"})
		} else {
			// c := NewClient(riak1)
			// c := NewClient(elb)
			c := getELB(user_id)
			keys, err := c.FetchKeys(user_id)
			fmt.Println(keys)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"Cart is empty"})
			} else {

				for i:=0; i<len(keys);i++{
					fmt.Println(len(keys))
					err:=c.DeleteOrder(user_id, keys[i])
					if err != nil{
						log.Fatal(err)
						formatter.JSON(w, http.StatusInternalServerError, struct{ Test string }{"Unable to delete cart"})
					}
				}
				formatter.JSON(w,http.StatusOK,struct{ Test string }{"Cart is empty"})
			}
		}

	}
}


func CheckoutCartHandler(formatter *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {

		order_count := 0
		var order_summary OrderSummary
		input_params := mux.Vars(req)
		user_id := input_params["userid"]
		if user_id == ""{
			formatter.JSON(w, http.StatusBadRequest, struct{ Test string }{"UserID missing. Please use /checkout/<userid>"})
		} else {
			// c := NewClient(riak1)
			// c := NewClient(elb)
			c := getELB(user_id)
			keys, err := c.FetchKeys(user_id)
			fmt.Println(keys)
			if err != nil {
				log.Fatal(err)
				formatter.JSON(w, http.StatusOK, struct{ Test string }{"Cart is empty"})
			} else {

				for i:=0; i<len(keys);i++{
					fmt.Println(len(keys))
					item_count,_ := c.FetchOrder(user_id, keys[i])
					order_count = order_count + item_count
					fmt.Print("%v",order_count)
					}
				order_summary.UserID = user_id
				order_summary.Order_count = order_count
				formatter.JSON(w, http.StatusOK, order_summary)
			}
		}

	}
}
