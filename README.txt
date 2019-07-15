# IMDB Server

RESTful service simple IMDB like database without using any external framework or library.

## Prerequisite

- Docker
- Docker Compose (version that supports "2.1" format of docker-compose file.)

## Usage

There is a docker-compose.yml file that can be used to build the two required components for the api.

- MySQL
- Golang based REST api 

```
$ go build -o imdb ./cmd/api
```
Above command would build a binary within the root folder of the project. Once this is done, One could just do 

```
$ docker-compose up -d
```
After this you should be able to access the GET, PUT POST urls for adding viewing and modifying a movie using REST api on url `http://localhost:8080/movie`.

Example POST request:
```
curl -X "POST" "http://localhost:8080/movie" \
     -H "Content-Type: application/json; charset=utf-8" \
     -d $'{"name": "some movie", "imdb_score":"9.88", "director":"James Cameron", "genre":["Sci fi", "Action"]}'
```

