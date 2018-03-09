package addressbookserver

import( 
	"encoding/json"
    "log"
	"net/http"
   	"github.com/gorilla/mux"
	"app/shared"
	d "app/dao"
	c "app/contact"
)

type Server struct {
	rClient shared.RedisClient
	router *mux.Router
	addressbook d.AddressBook
}

func (s *Server) Init(client shared.RedisClient) bool {
	s.addressbook = d.AddressBook{}
	s.router = mux.NewRouter()
	s.rClient = client
	return true
}

func (s *Server) RestApiGetContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiGetContact called")
	params := mux.Vars(r)
	email := params["email"]
	contact,_ := s.addressbook.GetContact(s.rClient, email)
	json.NewEncoder(w).Encode(contact)
}

func (s *Server) RestApiGetContacts(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiGetContacts called")
	var contacts []c.Contact
	for _, contact := range s.addressbook.GetAllContacts(s.rClient) {
		contacts = append(contacts, contact)	
	}
	json.NewEncoder(w).Encode(contacts)
}

func (s *Server) RestApiCreateContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiCreateContact called")
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
    var contact c.Contact
    err := decoder.Decode(&contact)
    if err != nil {
		panic(err)
    }
	w.Write(s.setContact(contact, shared.POST))
}

func (s *Server) RestApiDeleteContact(w http.ResponseWriter, r *http.Request) {
	log.Printf("RestApiDeleteContact called")
	var message []byte
	params := mux.Vars(r)
	email := params["email"]
	ok, errorCode, contact := s.addressbook.DeleteContact(s.rClient, email)
	if ok {
		log.Printf("Successfully deleted contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Success"}`)
	} else {
		log.Printf("Failed to delete contact: %v", contact)
		message = []byte(`{"ErrorCode":"`+errorCode+`","Message":"Failed to delete contact"}`)
	}
	w.Write(message)
}

func (s *Server) setContact(contact c.Contact, cmd int) []byte {
	var message []byte
	ok, errorCode := s.addressbook.SetContact(s.rClient, contact, cmd)
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
    var contact c.Contact
    err := decoder.Decode(&contact)
    if err != nil {
		panic(err)
	}
	if _, present := s.addressbook.GetContact(s.rClient, email); present {
		w.Write(s.setContact(contact, shared.PUT))
	} else {
		w.Write([]byte(`{"ErrorCode":"`+shared.CONTACT_DOES_NOT_EXIST+`","Message":"Could not find resource"}`))
	}
}

func (s *Server) HandleRequest() {
	s.router.HandleFunc("/addressbook/contact", s.RestApiCreateContact).Methods("POST")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiGetContact).Methods("GET")
	s.router.HandleFunc("/addressbook/contacts", s.RestApiGetContacts).Methods("GET")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiUpdateContact).Methods("PUT")
	s.router.HandleFunc("/addressbook/contact/{email}", s.RestApiDeleteContact).Methods("DELETE")
}

func (s *Server) GetRouter() *mux.Router {
	return s.router
}