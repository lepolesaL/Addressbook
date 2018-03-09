package contact

import(
	"regexp"
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

func (c *Contact) Validate() bool {
	regEmail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	regPhoneNumber :=regexp.MustCompile("^[0-9]*$")
	if c.Name == "" || c.Email == "" || c.PhoneNumber == "" || (Address{}) == c.Address {
		return false
	}
	if len(c.Name) < 4 ||  len(c.Name) > 25 || len(c.Email) < 5 || len(c.PhoneNumber) < 10 || len(c.PhoneNumber) > 12 {
		return false
	} 

	if  len(c.Address.Street) < 4 || len(c.Address.Street) > 50 || len(c.Address.City) < 4 || len(c.Address.City) > 25 {
		return false
	}

	if len(c.Address.Country) < 2 || len(c.Address.Country) > 2 {
		return false
	}

	if !regEmail.MatchString(c.Email) {
		return false
	}

	if !regPhoneNumber.MatchString(c.PhoneNumber) {
		return false
	}

	return true
}