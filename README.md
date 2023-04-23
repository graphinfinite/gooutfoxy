# Go out foxy!
- test project from https://gist.github.com/foxcpp/0fdf9bad8504fa803e73406bbeffadb3

## content
- grpc Server
- grpc-Gateway to REST
- Swagger UI
- Docker container
- parsing rusprofile.ru

## Run
```sh
go run cmd/rusprofile/main.go
```

## Docker
```sh
docker build --tag "ruspro" .
docker run -p 80:80 -p 81:81 "ruspro"
```


