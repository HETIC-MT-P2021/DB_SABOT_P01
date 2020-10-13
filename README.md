# Our school project

The goal of this project is to create a simple golang app with advanced SQL request.
We used Golang for our API, and Mysql as a database.

### Starting With Docker

After cloning the repo, create the .env file according to .env.example, `cd` into the project and run following commands

```bash
docker-compose up --build
```

The app will be accessible at localhost:8080 !

## Stating manualy

The project requires Golang v 1.14.4

Install the dependencies and start the server.

```sh
$ git clone https://github.com/HETIC-MT-P2021/go-wiki-group-4.git
$ go get
$ go get -u github.com/gin-gonic/gin
$ go get -u github.com/cosmtrek/air
$ go mod vendor
$ air
```

# Licence

The code is available under the MIT license.
