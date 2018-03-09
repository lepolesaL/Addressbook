package dao

import(
	"encoding/json"
	"log"
	c "app/contact"
	s "app/shared"
)

type AddressBook struct {
}

//Get a single contact from redis datastore. Returns empty contact if error occurs
func (a AddressBook) GetContact(client s.RedisClient, key string) (c.Contact, bool) {
	val, err := client.Get(key).Result()
	if err != nil {
		log.Printf("Could not get contact from datastore, %v", err)
		return c.Contact{}, false
	}
	log.Printf("value of get is, %v and error is %v", val, err)
	var contact c.Contact
	if err := json.Unmarshal([]byte(val), &contact); err != nil {
		log.Printf("Error occurred during unmarshalling, %v", err)
        return c.Contact{} , false
    }
	return  contact, true
} 

//Puts the value of new contact in datastore. First check if it als
func (a *AddressBook) SetContact(client s.RedisClient, contact c.Contact, cmd int) (bool, string) {
	if &contact != nil && contact.Validate() {
		_, present := a.GetContact(client, contact.Email)
		if present && cmd == s.POST {
			return false, s.CONTACT_ALREADY_EXISTS
		} else if !present && cmd == s.PUT {
			return false, s.CONTACT_DOES_NOT_EXIST
		}
		val, _ := json.Marshal(contact)
		err := client.Set(contact.Email, val , 0).Err()
		if err != nil {
			return false, s.DATA_STORE_WRITE_ERR
		}
		return true, s.SUCCESS
	}
	log.Printf("Contact could be null or invalid")
	return false, s.IMPL_ERR

}

func (a AddressBook) DeleteContact(client s.RedisClient, key string) (bool, string, c.Contact) {
	contact , present := a.GetContact(client, key) 
	if !present {
		return false, s.CONTACT_DOES_NOT_EXIST, contact
	}
	if _,err := client.Del(key).Result(); err != nil {
		log.Printf("Error occurred while deleting contact: %v", err)
        return false, s.DATA_STORE_DELETE_ERR, contact
	}
	return true, s.SUCCESS, contact
}

func (a AddressBook) GetAllContacts(client s.RedisClient) map[string]c.Contact {
	var contacts = make(map[string]c.Contact)
	keys, err:= client.Keys("*").Result()
	if err != nil {
		panic(err)
	}
	for _, key := range keys {
		contacts[key],_ = a.GetContact(client, key)
	}
	return contacts;
}