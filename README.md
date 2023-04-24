
# Golang fiber todo

This is api for todo app written in golang & fiber & gorm. I create it for educational perpuses, to learn golang a litle bit.



## Installation


```sh
    go mod tidy
    go mod vendor
```
    
## API Reference

[![Run in Postman](https://run.pstmn.io/button.svg)](https://app.getpostman.com/run-collection/11446645-2697eb37-2718-4937-bf81-8999517ed170?action=collection%2Ffork&collection-url=entityId%3D11446645-2697eb37-2718-4937-bf81-8999517ed170%26entityType%3Dcollection%26workspaceId%3De4868777-80e9-40b2-8f55-c370d83d4c78)

## Run Locally

Clone the project

```bash
  git clone git@github.com:gregor-tokarev/golang-fiber-todo.git
```

Go to the project directory

```bash
  cd golang-fiber-todo
```

Install dependencies

```bash
  go mod tidy
  go mod vendor
```

Init enviroment variables

```sh
cp .local.env .env
```

Provide google oauth client `id` and `secret`

```env
GOOGLE_OAUTH_CLIENT_ID=
GOOGLE_OAUTH_CLIENT_SECRET=
```

Start the server

```bash
  go run main.go
```


## Authors

- [@gregor-tokarev](https://www.github.com/gregor-tokarev)

