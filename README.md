# What the FUN architecture

**F**unctionality

**U**ntility

**N**on-Complexity

## Introduction

This is a template for creating a Go web service API project. Intended to be used as a starting point for creating a new Go web service API project and be a guideline for the project structure.

## Objective

- **Standardized Go project structure:** Offer a pre-configured project structure, essential tools, and best practices to reduce setup time and accelerate development.
- **Facilitate quality:** Enforce coding standards, Provide a testing framework and examples to encourage thorough testing and maintain code quality such as _golangci-lint_, _govulncheck_, and _pre-commit_ etc.
- **Enhance security:** Implement security measures from the start, such as vulnerability scanning with govulncheck and secure coding practices.
- **Support best practices:** Adhere to industry-standard design patterns and architectural principles for Go web services.
- **Be extensible:** Allow for easy customization and integration with various tool and frameworks.

**NOTE** you can **DELETE** anything that you **DON'T NEED** in this template.

## Table of Contents

- [Pre-requisite: before getting started or git clone this repository](#pre-requisite-before-getting-started-or-git-clone-this-repository)
- [How to copy this template](#how-to-copy-this-template)
- [Project structural components](#project-structural-components)
- [Guideline How To Code](#guideline-how-to-code)
- [Conventions](#conventions)
- [When to use log.Fatal/log.Panic](#when-to-use-log.Fatal/log.Panic)
- [Gracefully shutting down](#gracefully-shutting-down)
- [Libraies suggestion](#libraies-suggestion)
- [Bad Practices](#bad-practices)
- [Contribution-Guide](#contribution-guide)
- [References](#references)
- [Cross Functional Requirements](#cross-functional-requirements)
- [Technical Requirement Guidelines](#technical-requirement-guidelines)

## Pre-requisite: before getting started or git clone this repository
<details>
<summary>see contents</summary>

- Go 1.23 or later [https://go.dev/dl/](https://go.dev/dl/)
- set up GOPRIVATE environment variable to allow go get private repository
  - `go env -w GOPRIVATE=xxx.com`
- config ssh for private repo [GitLab SSH Keys](https://api.github.com/user/keys). otherwise, you will not be able to install use `gonew` command to create a new project.
  - or use personal access token set in `~/.netrc` file

if you are using ssh then you can add the following configuration to your git config file.

- install `gonew` to create a new project
  - `go install golang.org/x/tools/cmd/gonew@latest`
- install **direnv** to manage environment variables [direnv](https://direnv.net/)
  - `brew install direnv`
    - `echo "eval \"$(direnv hook bash)\"" >> ~/.zshrc` if you are using zsh. otherwise, use `~/.bashrc`
    - `source ~/.zshrc`
    - config direnv look at `.env` file by config `~/.config/direnv/direnv.toml`

```toml
[global]
load_dotenv = true
```
</details>

## How to copy this template
<details>
<summary>see contents</summary>

you also can just copy this repository and rename it to your project name and rename or the import module name to your project name.
However, you can use `gonew` command to archive the same result.

1. make sure you set `go env -w GOPRIVATE=xxx.com` (to allow go get private repository)
1. copy this repository and rename it to your project name or use `gonew github.com/pallat/wtf github.com/pallat/api`
1. cd to your project directory `cd xxxx` the
1. make your project under version control by running `git init`
1. and run `make setup` to install necessary tools
1. test the installation by running `make test` to run all the tests
1. load environment variables by running `direnv allow`
1. start project by running `make run` to run the project
1. add your git remote repository by running `git remote add origin github.com/your/xxxx` (DO NOT FORGET TO CREATE YOUR REPOSITORY FIRST IN GitLab)
1. commit your project to git repository by running `git add . && git commit -m "initial commit" && git push -u origin main`
1. start coding your project by adding new package in `app/` directory and look at the [Guideline How To Code](#guideline-how-to-code) below

</details>

## Project structural components

```sh
./
├── .env
├── .env.template
├── .local/
├── openapi/
├── .scripts/
├── Dockerfile
├── Makefile
├── README.md
├── VERSION
├── app/
│   ├── app.go
│   ├── ticket/
│   │   ├── handler.go
│   │   ├── http_service.go
│   │   ├── interpermit.go
│   │   └── storage_service.go
│   ├── response.go
│   ├── response_writer_middleware.go
│   └── response_writer_middleware_test.go
├── config/
├── database/
├── docker-compose.yml
├── gitlabci.yml
├── go.mod
├── go.sum
├── httpclient/
├── kafka/
├── logger/
├── main.go
├── note.txt
├── redis/
└── serror/
```

- `app/`: Business logic and application layer (handler, service, storage)
- `config/`: Configuration files and read from environment variables
- `database/`: Database connection and schema migration, and database operation general insert, update, delete, and select not specific to any table
- `httpclient/`: HTTP client for external API
- `kafka/`: Kafka producer
- `logger/`: Logging configuration
- `redis/`: Redis client
- `serror/`: Custom error package for handling error for the application and business error code
- `main.go`: Entry point of the application
- `Makefile`: Contains all the commands to run the application, test, and build the application
- `Dockerfile`: Dockerfile for building the application
- `docker-compose.yml`: Docker compose file for running the application and database for local development
- `VERSION`: Version file for the application version. This file will be updated by the CI/CD pipeline
- `openapi/`: API documentation and swagger file
- `.scripts/`: Script for CI/CD and other automation (e.g. CI/CD pipeline, development command)
- `.local/`: Local development configuration (e.g. git config pre-commit, pre-push, etc.)
- `.env`: Environment variables file for local development (copy from .env.template and fill in the value)
- `.env`.template: Environment variables template file for application configuration

## Guideline How To Code

### package `app/` is the primary focus for you

- all business logic and application layer should be in this directory. It should be separated by module/package name for example:
  - `app/register`
  - `app/booking`
  - `app/loan`
  - `app/credit`
  - `app/payment`
  - `app/transfer`
  - `app/exchange`
  - `app/purchase`

  - each package represent a module/domain in the system follow ideal of DDD (Domain Driven Design)
- each package SHOULD have one file name same as package name e.g. package `app/register` should have `register.go`
  - this file should contain the struct type and method to interact with the package e.g `type User struct {}` (or `type R struct {}` [why?](#convention) ) , `func New() *User {}` or `func (r *User) Verify(code string) bool () {...}`
- each package SHOULD have [file name](#file-naming-convention) to **represent responsibility** of the file e.g.
  - `handler.go` should contain the handler struct and method to handle the request and response
  - `get_handler.go` or `get_user_handler.go` should contain the handler method `(h *Handler) Get(w http.ResponseWriter, r *http.Request) {...}` to handle get a user data.
  - `gets_handler.go` or `get_all_user_handler.go` should contain the handler method `(h *Handler) GetAllUsers(w http.ResponseWriter, r *http.Request) {...}` to handle get all user data.
  - `create_handler.go` or `create_user_handler.go` should contain the handler method `(h *Handler) Create(w http.ResponseWriter, r *http.Request) {...}` to handle create a user data.
  - `dopa_service.go` should contain function to interact with **external service**
  - `verify_service.go` should contain function to verify the user data base on the **business logic**
  - `storage.go` should contain the storage service struct and method to interact with the database
- you **CAN DELETE** folder `app/interpermit`, `app/license` and `app/order` then delete all related endpoint in `main.go` then start create your own feature in `app/` directory.

### package `config/` is the custom error package

- all configuration should read then store in the struct in this package

### package `database/` is the database connector

- all database connection and schema migration should be in this package
- all database operation general insert, update, delete, and select not specific to any table should be in this package
- separate file for each database type e.g.:
  - `database/postgres.go`
  - `database/mysql.go`
  - `database/redis.go`
  - `database/kafka.go`
- **DO NOT** put database operation specific to any table/entity in this package e.g. `database/user.go`, please DON'T DO THIS. the `user_storage.go` should be in `app/user/user_storage.go` instead. so `user_storage.go` will import `database/postgres.go` or `database/mysql.go` etc. to interact with the database.

### package `httpclient/` is the HTTP client for external API

- this is a helper package to interact with external API using HTTP client.
- contain a function to call external API and return the ressult
- it should be use in `app/xxxx` package to interact with external API e.g. `app/license/http_service.go` using `httpclient.Post(url, body, header)`
- **DO NOT** put business logic in this package. the business logic should be in `app/xxxx` package instead.

### package `kafka/` is the Kafka producer

- this is a helper package to produce message to Kafka
- contain a function to produce message to Kafka
- it should be use in `app/xxxx` package to produce message to Kafka e.g. `app/license/storage_kafka.go` using `kafka.Produce(topic, message)`
- **DO NOT** put business logic in this package. the business logic should be in `app/xxxx` package instead.

### package `redis/` is the Redis client

- this is a helper package to interact with Redis
- contain a function to interact with Redis
- it should be use in `app/xxxx` package to interact with Redis e.g. `app/license/storage_redis.go` using `redis.Get(key)`, `redis.Set(key, value)`
- **DO NOT** put business logic in this package. the business logic should be in `app/xxxx` package instead.

### package `serror/` is the custom error package

- this is a helper package to handle error.
- contain a function to handle error and business error code and standardize error response.
- censor the error message before return to the client for security and privacy reason.

### package `logger/` is the logging configuration

- this is a helper package to configure the logger
- config the logger level, format, and output for standardize logging across the application
- contain a mask for sensitive data in the log such as password, token, PII data
- contain a function to log the request and response with PDPA compliance in mind
- it meant to be use in `app/xxxx` package to log the request and response

### package `openapi/` is the API documentation and swagger file

- this is a helper package to generate API documentation and swagger file
- contain a function to generate API documentation and swagger file
- what is difference between OpenAPI and Swagger?
  - The easiest way to understand the difference is:
    - OpenAPI = Specification
    - Swagger = Tools for implementing the specification

### package `.scripts/` is the script for CI/CD and other automation

- this is a helper package to run the script for CI/CD and other automation
- contain a script to run the CI/CD pipeline and for all development command
- if you have any command that using frequently you can put it here
- then add your command to the `Makefile` to run the script

### application entry point is `main.go`

- just a simple entry point to start the application. it should be simple and clean as possible leave it in root directory of the project. (same level as go.mod)
- it should be use to initialize the application and start the server
- it should be use to initialize the database connection, redis connection, kafka connection, logger, and other configuration
- it should be use to inject the dependency to the `app/xxxx` package
- it should be use to start the server and gracefully shutdown the server
- how to handle multiple main.go?
  - single `main.go` scenario: you start just web service as restful API.
    - you can put `main.go` in the root directory of the project (same level as go.mod)
  - multiple `main.go` scenario: you start multiple services e.g. web service, batch service, cron service, etc.
    - you can put `main.go` in the `cmd/` directory e.g. `cmd/batch/main.go`, `cmd/cron/main.go`, `cmd/web/main.go`
    - so each `main.go` will start different service e.g. `cmd/batch/main.go` will start batch service, `cmd/cron/main.go` will start cron service, `cmd/web/main.go` will start web service
  - we don't have folder `cmd/` in this template. you can create it if you have multiple services.

### `Makefile` contains all the commands to run the application, test, and build the application

- Makefile is a simple way to organize commands and tasks for your project. so that when new developers join the project, they can easily run the command to start the project, run the test, build the project, etc. by just looking at the Makefile.
- it should contain all the commands to run the application, test, and build the application
- it should contain all the commands to run the CI/CD pipeline
- it should contain all the commands to run the development command
- **Ground rules** if you have **any command** that you use frequently you can **put it in the Makefile**

### folder `.scripts/` is the script for CI/CD and other automation

- this is a helper package to run the script for CI/CD and other automation
- contain a script to run the CI/CD pipeline and for all development command
- if you have any command that using frequently you can put it here
- then add your command to the `Makefile` to run the script

### folder `.local/` is the local development configuration

- this is a helper folder to store the local development configuration
- contain a local development configuration file e.g. git config `pre-commit`, `pre-push`, etc.
- any local development configuration that your team agreed to use can be put here

### Available Commands in Makefile

- `make help`: show all available commands in the Makefile
- `make setup`: setup environment git hook, go mod download and install tools (golangci-lint, govulncheck and swagger)
- `make all`: default (mod inst lint test)
- `make mod`: tidy up the go.mod and go.sum file
- `make lint`: run golangci-lint to check the code quality and code smell, some basic security check and prevent to use fmt.Print and fmt.Println in the code
- `make test`: run all the unit test
- `make coverage`: run all the test with coverage
- `make vuln`: run govulncheck to check the vulnerability
- `make clean`: remove files created during build pipeline and cache
- `make precommit`: validate the branch before commit
- `make docker`: build the docker image
- `make run`: run the docker container
- `make diff`: check if there are any uncommitted changes
- `make bump-version`: bump the version of the project
- `make commit-msg`: run pre-commit to check the commit message
- `make swagger`: generate swagger file
- `make openapi`: serve the swagger file

## Conventions

- Any folder that **NOT** a Go's package please name it with __prefix__ `.` e.g. `.scripts/`, `.local/` etc.
- Any new business logic should be in `app/` directory and should be separated by module/package name.
  - e.g. if you have a new business logic for `register` module, you should create a new package `app/register` and put all the business logic in this package.
- **SHOULD NOT** call `controller` in the package name. It should be `handler` instead.
- **SHOULD NOT** create package name by functionally/layers e.g. `app/controller`, `app/service`, `app/repository`, `app/model`
  - Because in Go, package name should be by domain or business features not by functionally/layers.
  - Whenever you facing a **CYCLE IMPORT (circular dependency)**, it's a sign that you are doing it **WRONG way**.(sometime it cause by package name by functionally/layers)
- **Avoid meaningless package names** create package name with `common`, `utils`, `helper`, `enums` etc.
  - Because it's too generic and not specific to any domain or business features.

### Naming Convention

- variable name should be in **MixedCaps** (PascalCase) or **mixedCaps** (camelCase) e.g. `var MaxUsers = 100`, `var maxUsers = 100`
- local variable name should be in **mixedCaps** (camelCase) e.g. `var maxUsers = 100` or `maxUsers := 100`
- In Go is to use **MixedCaps** (PascalCase) or **mixedCaps** (camelCase) rather than underscores to write multiword names.
- a **constant name** should be in **MixedCaps** (public constant) or **mixedCaps** (private constant)
  - GOOD: `const MaxUsers = 100`, `const maxUsers = 100`
  - BAD: `const MAX_USERS = 100`, `const max_users = 100` (because it's not idiomatic in Go)
- function name should be in **MixedCaps** (PascalCase) or **mixedCaps** (camelCase) e.g. `GetUser`, `getOrder`, `CreateOrder`, `createOrder`
- struct name should be in **MixedCaps** (PascalCase) or **mixedCaps** (camelCase) e.g. `type Order struct {}`, `type order struct {}`

### File Naming Convention

- file name should be in **snake_case** e.g. `user_handler.go` or `user_storage.go`
  - GOOD: `user_handler.go`, `user_storage.go`, `user_http_service.go`, `user_http_service_test.go`, `user.go`
  - BAD: `UserHandler.go`, `userHandler.go`, `UserHandler_test.go`, `User.go`
- file name should represent the **responsibility** of the file what it does.
  - example logger:
  - `logger/logger.go`
	  - `logger/level.go`
	  - `logger/mask.go`
	  - `logger/request_response.go`
	  - etc.

		when you want to __make change about log level__ which file you should look at? `logger/level.go` right? so it's easy to find the file.
  - example auth:
	  - `auth/cors.go`
	  - `auth/basic_auth.go`
	  - `auth/jwt.go`
	  - `auth/middleware.go`
	  - etc.

		when you want to __make change about basic authenication__ which file you should look at? `auth/basic_auth.go` right? so it's easy to find the file.
	- example register:
		- `register/storage.go`
		- `register/dopa_service.go`
		- `register/verify_account.go`
		- etc.

		when you want to __make change about Department of Provincial Administration (DOPA) service__ which file you should look at? `register/dopa_service.go` right? so it's easy to find the file.
- file name must create at least one file that same name as package e.g. `package order` should contain atleast `order.go`
- file name `order/xxx_handler.go` it should start with the action of the handler or just `order/handler.go` if it's hold all the handlers.
- we can split file to contain all functionality that it reposible.
	- you can use prefix with action e.g.
		- `get_handler.go`
		- `get_all_handler.go`
		- `create_handler.go`
		- `update_handler.go`
		- `delete_handler.go`
		- etc.

		when you read it with their package name it's easy to understand what it does. e.g.
		- `order/get_handler.go`
		- `order/get_all_handler.go`
		- `order/create_handler.go`
		- `order/update_handler.go`
		- `order/delete_handler.go`

- the **test file** also follow the same pattern and end with `_test.go`
- the **integration test file** that test handler-> service/usecase -> storage that end with `_it_test.go`


the responsibility of each layer should be separated as file name don't need to create a package for each layer.
we can archive layer separation by file name instead of package name for example:

- order.go
  - starting point of the package and model.
  - method `s.OrderStore() OrderStore {}` to translate the app (skill) data structure to database data structure

- order_handler.go or _order_get_handler.go_ or _get_handler.go_
  - struct handler that have [**Storage Interface**](#interface-naming-convention) as dependency
  - Initialize order Storage instance at the `main.go` then inject it into __NewOrderHandler(orderStorage)__ to create a new OrderHandler instance.

- order_storage.go
  - A **Storage struct** implements a **Storage interface**.
  - The Storage struct holds the database connection directly (not through an interface).
  - The Storage struct is initialized in the main.go file.
  - It provides methods for interacting with the database, _order_ entity specific such as:

    - FindOrderByKey(key string) (Order, error)
    - FindAllOrders(orderBy string, page, limit uint) ([]Order, uint, error)
    - CreateOrder(key string, ord OrderStore) error
    - UpdateOrder(key string, ord OrderStore) error

### Package Naming Convention

Go code is organized into packages. Within a package, code can refer to any identifier (name) defined within, while clients of the package may only reference the package’s exported types, functions, constants, and variables. Such references always include the package name as a prefix: foo.Bar refers to the exported name Bar in the imported package named foo.

Good package names make code better. A package’s name provides context for its contents, making it easier for clients to understand what the package is for and how to use it. The name also helps package maintainers determine what does and does not belong in the package as it evolves. Well-named packages make it easier to find the code you need.

**Good package names are short and clear.** They are **lower case**, with **NO under_scores, NO mixedCaps**. They are often **simple nouns**, such as:
- package **time** (provides functionality for measuring and displaying time)
- package **list** (implements a doubly linked list)
- package **http** (provides HTTP client and server implementations)

these are **NOT GOOD for Go package names** The style of names typical of another language might not be idiomatic in a Go program. Here are examples of names that might be good style in other languages but do not fit well in Go
- package **dopaServiceClient** (mixedCaps)
- package **http-handler** (hyphen)
- package **priority_queue** (under_score)

**Abbreviate judiciously** Package names may be abbreviated when the abbreviation is familiar to the programmer. Widely-used packages often have compressed names:

- package **strconv** (string conversion)
- package **syscall** (system call)
- package **fmt** (formatted I/O)

On the other hand, if abbreviating a package name makes it **ambiguous or unclear**, **don’t do it**.

__Don’t steal good names from the user__. Avoid giving a package a name that is commonly used in client code. For example, the buffered I/O package is called *bufio*, not `buf`, since `buf` is a good _variable name_ for a buffer.

**Avoid repetition.** Since client code uses the _package name as a prefix_ when referring to the package contents, the names for those contents need not repeat the package name.
for example The HTTP server type in **package http**.

- GOOD name: **Server** is a good name `type Server struct {...}` when combined with the package name, it becomes **http.Server**, which is clear and unambiguous when used in client code.
- BAD name:, **HTTPServer** is a bad name `type HTTPServer struct {...}` when combined with the package name, it becomes **http.HTTPServer**, which is redundant `http`.`HTTP` verbose when used in client code.

**Simplify function names**. When a function in package pkg returns a value of type pkg.Pkg (or *pkg.Pkg), the function name can often omit the type name without confusion:

```go
	start := time.Now()                                  // start is a time.Time
	t, err := time.Parse(time.Kitchen, "6:06PM")         // t is a time.Time
	ctx = context.WithTimeout(ctx, 10*time.Millisecond)  // ctx is a context.Context
	ip, ok := userip.FromContext(ctx)                    // ip is a net.IP
```

A function named New in package pkg returns a *value* of type pkg.Pkg. This is a standard entry point for client code using that type:

```go
	q := list.New()  // q is a *list.List
```

When a function returns a value of type pkg.T, where T is not Pkg, the function name may include T to make client code easier to understand. A common situation is a package with multiple New-like functions:

```go
	d, err := time.ParseDuration("10s")  // d is a time.Duration
	elapsed := time.Since(start)         // elapsed is a time.Duration
	ticker := time.NewTicker(d)          // ticker is a *time.Ticker
	timer := time.NewTimer(d)            // timer is a *time.Timer
```

### Interface Naming Convention

- could be named with the **-er suffix** e.g. `Reader`, `Writer`, `Formatter`, `Parser`, `Comparer`, `Marshaler`, `Unmarshaler`, `Closer`, `Seeker`, `Flusher`, `Scanner` etc. when it have a **single method** and **sound make sense**.
- keep the interface minimal as possible. **Single method** is the best. **Multiple method** is **OK**. the perfect example of interface in Go is `io.Reader` and `io.Writer` interface. it's **simple** and **easy to understand**. when you compose it also **simple** such as `io.ReadWriter` or `io.ReadCloser` etc.
- could be named with the **-able suffix** e.g. `Comparable`, `Printable`, `Readable`, `Writable`, `Seekable`, `Flushable`, `Scannable` etc. when it **sound make sense**.
- could **NOT** be name with prefix `I` e.g. `IReader`, `IWriter`, `IFormatter`, `IParser`, `IComparer`, `IMarshaler`, `IUnmarshaler`, `ICloser`, `ISeeker`, `IFlusher`, `IScanner` etc. because it's not idiomatic in Go.
- if you have many methods in the interface you should **consider** to **split** it to **smaller interface** or try to name it with **more meaningful name**. e.g. `io.ReadWriter` instead of `io.ReaderWriter` or `io.ReadCloser` instead of `io.ReaderCloser` etc. if you can't find the meaningful name you should **consider** to question your design and **rethink** about it. maybe follow **SOLID** principle. however, defining the interface and naming it is **hard**. I leave you with Rob Pike quote __"Don't design with interfaces, discover them"__.

- for example Storage interface it okay to name it with `Storage` e.g. `orderStorage`, `userStorage`, `LicenseStorage`, `InterpermitStorage` etc.

```go
type userStorage interface {
	SaveUser(context.Context, u User) error
}
```

struct that implement the interface should be named with `stroage` and be private.

```go
type storage struct {
	db *sql.DB
}

func (s *storage) SaveUser(ctx context.Context, u User) error {
	// implementation
}
```

then you export way to create the instance of the storage. by using `NewStorage` function.

```go
func NewStorage(db *sql.DB) storage {
	return storage{db: db}
}
```

basically you don't need to add word `interface` or `Interface` into the interface name. it's already clear that it's an interface at compile time.


## When to use log.Fatal/log.Panic

- **log.Fatal** or **log.Panic** should be used **ONLY** in the **main** function before the application starts for all hard dependencies/configuration that MUST be available for the application to run.

## Gracefully shutting down
this repository already have a code to gracefully shutting down the server. you can look at `main.go` file to see how it's done.

## Libraies suggestion

### Recommended Libraries

- [Env : Configuration](github.com/caarlos0/env/v11) for environment variables because it's simple and no required external dependencies use only Go standard library.

### Unrecommended Libraries but still can use with caution

- [Viper : Configuration](github.com/spf13/viper) for environment variables because it's complex and required external dependencies use Go standard library instead.
- [GORM : ORM](github.com/go-gorm/gorm) for database operation because it's complex and required external dependencies use Go standard library instead. however, you can still use it in simple operation and take your own risk.

## Bad Practices

- bad naming erros as `Exception` e.g. `BaseException.go`, `NotFoundException.go` there is **NO** Exception in **Go**
- bad naming package with `controller` or `service`, `repository`
- bad naming folder start with src with `src/`, `src/common`, `src/utils`
- bad naming too generic name `common`, `utils`, `helper`, `enums`

## Contribution-Guide

- **Fork** this repository
- **Clone** your forked repository
- **Create** a new branch for your feature or bug fix `git checkout -b feature/xxxx` or `git checkout -b bugfix/xxxx`
- **Update** update your implementation or changes any design.
- **Test** your changes `make test` make sure your changes pass all the tests and have test on it if it's possible.
- **Commit** your changes `git add . && git commit -m "your message"`
- **Push** your changes to your forked repository `git push origin feature/xxxx`
- **Create** a new pull request to this repository with **explanation of your changes**
- **Wait** for the review and approval from the reviewer
- **Merge** your pull request

## Reference

## Cross Functional Requirements

- SLA สำหรับตั้งค่า timeout ต่างๆ เช่น Read/Write database
  - มีค่า default ในกรณีที่ไม่ได้ตั้งค่า config
- Log
  - level อะไรบ้าง และแต่ละ level ต้องการข้อมูลอะไร
- Gracefully shutting down
- tracing ID คือค่าอะไร มี span หรือไม่
  - install `direnv` to manage environment variables [direnv](https://direnv.net/)

- when getting started just `make setup` to install necessary tools
  - Install gonew `go install golang.org/x/tools/cmd/gonew@latest`
  - Install golangci-lint `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`
  - Install govulncheck `go install github.com/ossf/scorecard/v2/cmd/govulncheck@latest`
  - Install swager `go install github.com/go-swagger/go-swagger/cmd/swagger@latest`

## Technical Requirement Guidelines

- Standardized Response structure (JSON)
- Auto logging on request and response with PDPA compliance in mind
- Standardized HTTP status code selection(int/out)
- Testable code based
- Code qualities checking on local/CI (golangci-lint, go vet, go test)
- Microservice simple tracing with correlation ID (header X-Ref-Id)
- 12factor, configuration only in environment variables
