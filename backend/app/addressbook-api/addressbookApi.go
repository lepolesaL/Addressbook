package main

import( 
	"fmt"
    "log"
	"net/http"
	"github.com/go-redis/redis"
	"github.com/rs/cors"
	s "app/addressbookserver"
)
// NB: Needs refactoring
const (
	REDIS_DB = "redis"
	REDIS_PORT = "6379"
)

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
	server := s.Server{}
	if server.Init(client) {
		fmt.Printf("Server intialization successful\n")
		server.HandleRequest()
		handler := cors.AllowAll().Handler(server.GetRouter())
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", handler))
	} else {
		fmt.Printf("Server intialization failed\n")
	}
}