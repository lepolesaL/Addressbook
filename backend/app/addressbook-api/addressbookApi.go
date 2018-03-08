package main

import( 
	"fmt"
	"encoding/json"
    "log"
	"net/http"
   	"github.com/gorilla/mux"
	"github.com/go-redis/redis"
	"github.com/rs/cors"
	"time"
	"regexp"
)
// NB: Needs refactoring
const (
	REDIS_DB = "redis"
	REDIS_PORT = "6379"
)
//Error codes
const (
	SUCCESS = "0"
	CONTACT_ALREADY_EXISTS = "1"
	DATA_STORE_WRITE_ERR = "2"
	CONTACT_DOES_NOT_EXIST = "3"
	DATA_STORE_DELETE_ERR = "4"
	IMPL_ERR = "10"
)

//Api Command codes
const (
	POST = 1
	PUT = 2
)

type Address struct {
	Street string `json:"street"`
	City string	   `json:"city"`
	Country string `json:"country"`
}

type Contact struct {
	Name string `json:"name"`
	Email string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Address Address `json:"address"`
}

type AddressBook struct {
}

type RedisClient interface {
	Get(key string) *redis.StringCmd
	Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(keys ...string) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
}

type Server struct {
	rClient RedisClient
	router *mux.Router
	addressbook AddressBook

}

// Contact validation
func (c *Contact) validate() bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	if c.Name == "" || c.Email == "" || c.PhoneNumber == "" || (Address{}) == c.Address {
		return false
	}
	if len(c.Name) < 5 ||  len(c.Name) > 25 || len(c.Email) < 5 || len(c.PhoneNumber) < 10 || len(c.PhoneNumber) > 12 {
		return false
	} 

	if  len(c.Address.Street) > 4 || len(c.Address.Street) < 50 || len(c.Address.City) > 4 || len(c.Address.City) < 25 {
		return false
	}

	if len(c.Address.Country) < 2 || len(c.Address.Country) > 2 {
		return false
	}

	if !re.MatchString(c.Email) {
		return false
	}

	return true
}

//Get a single contact from redis datastore. Returns empty contact if error occurs
func (a AddressBook) getContact(client RedisClient, key string) (Contact, bool) {
	val, err := client.Get(key).Result()
	if err != nil {
		log.Printf("Could not get contact from datastore, %v", err)
		return Contact{}, false
	}
	log.Printf("value of get is, %v and error is %v", val, err)
	var contact Contact
	if err := json.Unmarshal([]byte(val), &contact); err != nil {
		log.Printf("Error occurred during unmarshalling, %v", err)
        return Contact{} , false
    }
	return  contact, true
} 

//Puts the value of new contact in datastore. First check if it als
func (a *AddressBook) setContact(client RedisClient, contact Contact, cmd int) (bool, string) {
	if &contact != nil && contact.validate() {
		_, present := a.getContact(client, contact.Email)
		if present && cmd == POST {
			return false, CONTACT_ALREADY_EXISTS
		} else if !present && cmd == PUT {
			return false, CONTACT_DOES_NOT_EXIST
		}
		val, _ := json.Marshal(contact)
		err := client.Set(contact.Email, val , 0).Err()
		if err != nil {
			return false, DATA_STORE_WRITE_ERR
		}
		return true, SUCCESS
	}
	log.Printf("Contact is not supposed to be null")
	return false, IMPL_ERR

}

func (a AddressBook) deleteContact(client RedisClient, key string) (bool, string, Contact) {
	contact , present := a.getContact(client, key) 
	if !present {
		return false, CONTACT_DOES_NOT_EXIST, contact
	}
	if _,err := client.Del(key).Result(); err != nil {
		log.Printf("Error occurred while deleting contact: %v", err)
        return false, DATA_STORE_DELETE_ERR, contact
	}
	return true, SUCCESS, contact
}

func (a AddressBook) getAllContacts(client RedisClient) map[string]Contact {
	var contacts = make(map[string]Contact)
	keys, err:= client.Keys("*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		contacts[key],_ = a.getContact(client, key)
	}
	return contacts;
}

// Intialize server
func (s *Server) Init(client RedisClient) bool {
	s.addressbook = AddressBook{}
	s.router = mux.NewRouter()
	s.rClient = client
	return true
}

func (s *Server) RestApiGetContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiGetContact called")
	params := mux.Vars(r)
	email := params["email"]
	contact,_ := s.addressbook.getContact(s.rClient, email)
	json.NewEncoder(w).Encode(contact)
}

func (s *Server) RestApiGetContacts(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiGetContacts called")
	var contacts []Contact
	for _, contact := range s.addressbook.getAllContacts(s.rClient) {
		contacts = append(contacts, contact)	
	}
	json.NewEncoder(w).Encode(contacts)
}

func (s *Server) RestApiCreateContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiCreateContact called")
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
    var contact Contact
    err := decoder.Decode(&contact)
    if err != nil {
		panic(err)
    }
	w.Write(s.setContact(contact, POST))
}

func (s *Server) RestApiDeleteContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiDeleteContact called")
	var message []byte
	params := mux.Vars(r)
	email := params["email"]
	ok, errorCode, contact := s.addressbook.deleteContact(s.rClient, email)
	if ok {
		log.Printf("Successfully deleted contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Success"}`)
	} else {
		log.Printf("Failed to delete contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Failed to delete contact"}`)
	}
	w.Write(message)
}

func (s *Server) setContact(contact Contact, cmd int) []byte {
	var message []byte
	ok, errorCode := s.addressbook.setContact(s.rClient, contact, cmd)
	if  ok {
		log.Printf("Successfully add contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Success"}`)
	} else {
		log.Printf("Failed to add contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Failed to add contact"}`)
	}
	return message
}

func (s *Server) RestApiUpdateContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiDeleteContact called")
	defer r.Body.Close()
	params := mux.Vars(r)
	email :=params["email"]
	decoder := json.NewDecoder(r.Body)
    var contact Contact
    err := decoder.Decode(&contact)
    if err != nil {
		panic(err)
	}
	if _, present := s.addressbook.getContact(s.rClient, email); present {
		w.Write(s.setContact(contact, PUT))
	} else {
		w.Write([]byte(`{"ErrorCode":"`+CONTACT_DOES_NOT_EXIST+`","Message":"Could not find resource"}`))
	}
}

func (s *Server) handleRequest() {
	s.router.HandleFunc("/addressbook/contact", s.RestApiCreateContact).Methods("POST")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiGetContact).Methods("GET")
	s.router.HandleFunc("/addressbook/contacts", s.RestApiGetContacts).Methods("GET")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiUpdateContact).Methods("PUT")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiDeleteContact).Methods("DELETE")
}

func main() {
	fmt.Printf("Server started\n")
	var client = redis.NewClient(&redis.Options{
		Addr:     REDIS_DB+":"+REDIS_PORT,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	log.Printf(pong)
	server := Server{}
	if server.Init(client) {
		fmt.Printf("Server intialization successful\n")
		server.handleRequest()
		handler := cors.AllowAll().Handler(server.router)
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
	} else {
		fmt.Printf("Server intialization failed\n")
	}
}