# Address book

This is a single page application that stores contacts.

## Getting Started
  
 Get project by running the following command
 ```
 git clone https://github.com/lepolesaL/Addressbook.git
```

### Prerequisites

  - Docker version 17.12.0
  - Docker compose version 1.18.0
 - Angular CLI: 1.7.2

### Installing

 - Change directory to project directory
 - Build project docker images using the command bellow:
 ```
 docker-compose build
 ```
 - Run project using the following command
 ```
 docker-compose up
 ```
 - After the docker containers are started, go to the following address in the browser
 ```
 http://localhost:4200/
 ```
####Testing
- Test frontend by navigating to frontend/webapp folder and run the following command
```
ng test
```
- Test backend by running the following commands
```
docker-compose build backend
```
```
docker-compose run backend go test -v ./addressbook-api
```

##NOTE
This app is not yet scalable since there is no load balancing in place.
