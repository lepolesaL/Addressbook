package main

import (
	"testing"
	"time"
	"github.com/go-redis/redis"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"bytes"
	"log"
	"io/ioutil"
	s "app/addressbookserver"
	d "app/dao"
	c "app/contact"
)

// Error responses from rest api calls
type FakeApiError struct {
	ErrorCode string `json:"ErrorCode"`
	Message string `json:"Message"`
}

// Mock redis.Client impements RedisClient interface
type RedisFakeClient struct {

}

// Implements radis error interface
type FakeRedisError struct {

}

var server s.Server 

func (e FakeRedisError) Error() string {
	return "Error"
}

func (c *RedisFakeClient) Set(key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	return redis.NewStatusResult("Success", nil);
}

func (c *RedisFakeClient) Get(key string) *redis.StringCmd {
	if key == "jone@example.com" {
		return redis.NewStringResult("{\"name\":\"Jone\",\"email\":\"jone@example.com\",\"phonenumber\":\"0791234567\",\"address\":{\"street\":\"12 Fake Street\",\"city\":\"Fake City\",\"country\":\"FA\"}}", nil)
	} else {
		return redis.NewStringResult("Error", FakeRedisError{})
	}
}

func (c *RedisFakeClient) Del(keys ...string) *redis.IntCmd {
	return redis.NewIntResult(1, nil)
}

func (c *RedisFakeClient) Keys(pattern string) *redis.StringSliceCmd {
	return redis.NewStringSliceCmd("jone@example.com")
}

// ************* Test interaction with redis datastore using fakeRedis client *************
func TestSetContact(t *testing.T) {
	var client *RedisFakeClient
	addressbook := d.AddressBook{}
	contact := c.Contact{"Jonas","jonas@example.com", "0791234567", c.Address{"12 Fake Street", "Fake City", "FA"}}
	ok, errorCode := addressbook.SetContact(client, contact, 1);
	if !ok {
		t.Fail()
	}
	t.Logf("Can set contact: %s, ok: %t", errorCode, ok)
}

func TestGetContact(t *testing.T) {
	var client *RedisFakeClient
	contact := c.Contact{"Jone","jone@example.com", "0791234567", c.Address{"12 Fake Street", "Fake City", "FA"}}
	addressbook := d.AddressBook{}
	val, ok := addressbook.GetContact(client, "jone@example.com")
	if ok {
		if val.Name != contact.Name {
			t.Errorf("Expected contact names to be equal")
		} else {
			t.Logf("Can get contact")
		}
	} else {
		t.Fail()
	}
}

func TestDeleteContact(t *testing.T) {
	var client *RedisFakeClient
	addressbook := d.AddressBook{}
	ok,_,_ := addressbook.DeleteContact(client, "jone@example.com")
	if ok {
		t.Logf("Can delete contact")
	} else {
		t.Fail()
	}
}

// ************** Test Server rest api *************
func TestMain(m *testing.M) {
	var client *RedisFakeClient
    server = s.Server{}
	server.Init(client)
	server.HandleRequest()
    code := m.Run()
    os.Exit(code)
}
// Execute http request and records the response from server
func executeRequest(method string, path string, contact c.Contact) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	if (c.Contact{})==contact {
		req, _ := http.NewRequest(method, path, nil)
		server.GetRouter().ServeHTTP(response, req)
	}  else {
		payload,_:= json.Marshal(contact)
		req, _ := http.NewRequest(method, path, bytes.NewBuffer(payload))
		server.GetRouter().ServeHTTP(response, req)
	}
	return response
}

func TestRestApiGetContact(t *testing.T) {
	response := executeRequest("GET", "/addressbook/contact/jone@example.com", c.Contact{})
	body, _ := ioutil.ReadAll(response.Body)
	var contact c.Contact
	if err := json.Unmarshal([]byte(string(body)), &contact); err != nil {
		t.Error("Invalid response")
	}
	if contact.Name != "Jone" {
		t.Error("Invalid response")
	}

}


func TestRestApiCreateContact(t *testing.T) {
	contact := c.Contact{"Jona","jona@example.com", "0791234567", c.Address{"12 Fake Street", "Fake City", "FA"}}
	response := executeRequest("POST", "/addressbook/contact", contact)
	body, _ := ioutil.ReadAll(response.Body)
	var fakeApiError FakeApiError
	if err := json.Unmarshal([]byte(string(body)), &fakeApiError); err != nil {
		t.Error("Invalid response")
	}
	if fakeApiError.ErrorCode != "0" {
		t.Error("Invalid response")
	}
}

func TestRestApiUpdateContact(t *testing.T) {
	contact := c.Contact{"Jone","jone@example.com", "0791234567", c.Address{"12 Fake Street", "Fake City", "FA"}}
	response := executeRequest("PUT", "/addressbook/contact/jone@example.com", contact)
	body, _ := ioutil.ReadAll(response.Body)
	var fakeApiError FakeApiError
	if err := json.Unmarshal([]byte(string(body)), &fakeApiError); err != nil {
		t.Error("Invalid response")
	}
	if fakeApiError.ErrorCode != "0" {
		t.Error("Invalid response")
	}
}

func TestApiGetContacts(t *testing.T) {
	response := executeRequest("GET", "/addressbook/contact", c.Contact{})
	body, _ := ioutil.ReadAll(response.Body)
	log.Println(string(body))
	if contacts := response.Body.String(); contacts == "[]" {
		t.Fail()
	}
}

func TestApiDeleteContact(t *testing.T) {
	response := executeRequest("DELETE", "/addressbook/contact/jone@example.com", c.Contact{})
	body, _ := ioutil.ReadAll(response.Body)
	var fakeApiError FakeApiError
	if err := json.Unmarshal([]byte(string(body)), &fakeApiError); err != nil {
		t.Error("Invalid response")
	}
	if fakeApiError.ErrorCode != "0" {
		t.Error("Invalid response")
	}
}