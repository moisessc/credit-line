# :credit_card: **GOLANG - CREDIT LINE API**

![Build status](https://img.shields.io/badge/build-passing-blue "Build status")
![Test status](https://img.shields.io/badge/test-passing-green "Test status")

## **Description**
This project is a simple API that allows you to simulate the calculation of the credit limit allowed for a Startup or an SME.

---

## :rocket: **Getting Started**
1. Open your terminal and clone this repository in your computer with the command: ```git clone https://github.com/moisessc/credit-line.git```
2. Now go to the cloned repository folder with: ```cd credit-line``` 
3. Create an .env file from the .env.template file and set the environment variables if you need to change them; otherwise you can skip this step, the API will load the default environment values.

### **Run the project with docker :whale:**
1. You can run the project using docker, if you don't have it, you can **[Install Docker Here!](https://www.docker.com/products/docker-desktop)** :point_left: :point_left:
2. Now you can use the following commands in your terminal:

  | **Command** | **Description** |
  | --- | --- |
  |```docker docker build -t credit-line-api . && docker image prune -f```|This command build the project and create a docker image|
  |```docker run -it -p 3000:3000 credit-line-api```|This command execute the API and exposes in port 3000|
  |```docker run -it --env-file ./.env -p 3000:3000 credit-line-api```|This command execute the API and exposes in port 3000 with an .env file configuration|

**:warning: Note:** If you change the application port in environment varibales you need to change the port in the docker run command, for example: ```docker run -it -d --env-file ./.env -p 8080:${MY_NEW_PORT} credit-line-api```

### **Run the project without docker :computer:**
  1. You need go in your computer, this project runs with 1.17 go version, if you don't have it you can download node **[Here!](https://go.dev/dl/)** :point_left: :point_left:
  2. Now you can run the following command: ```go run cmd/credit-line-api/main.go```  

### **Usage :pencil:** 
You can send a request in the following path when the application is running: ```http://localhost:3000/api/v1/credits/calculate/limit```

**Payload:**
```
{
    "foundingType": "Startup",
    "cashBalance": 13435.30,
    "monthlyRevenue": 4235.45,
    "requestedCreditLine": 100,
    "requestedDate": "2021-07-19T16:32:59.860Z"
}
```

### **Testing** 🧪
You can run the tests of the application with the following command: ```go test ./... -v```

---


## :file_folder: **Project structure**
The project has three principal packages, the communication between the layers is through of interfaces to isolate the implementations and promote the dependency injection through constructor, each package represent a general aspect of the application:
| **Package** | **Description** |
| --- | --- |
|**cmd** | Contains the entry points for the application|
|**internal** | Contains the core of the application, here are the bussines rules and domains|
|**pkg** | Contains packages necessary for the application but that do not belong to the core of the application, here are packages such as validators, middlewares, environments loaders etc.|

### :file_folder: **cmd Package**
Inside this package we have a package with the name credit-line-api, the idea is to have a package for each executable of the application, for example, imagine that we have another package with the name: credit-line-cli this package will contain the necessary to build an executable of the application in CLI mode, the structure of this packages are:
- **cmd/credit-line-api**
    - **bootsrap package:** Package that has the necessary to initialize the application in a determinate mode
        - **bootstrap.go** File that concentrating all the initialization of the application (Dependency injection, Server, Router, etc.)
        - **router.go** File that contains the router to declare the endpoints path in this case the echo implementation
        - **server.go** File that contains the logic to initialize a server with any router such as echo or gin
    - **main.go** Main file to initialize the bootstrap

### :file_folder: **internal Package**
Contains the packages core of the application divided in four:
- **internal**
    - **controller package:** This package is the entry point for the application core, communicate the bootstrap layer with the core, in this case, the package contains the echo handlers to communicate the router of the bootstrap layer with the service layer, if we need to change the router from Echo to Gin, we'll need to create the Gin handlers here
    - **service package:** This package contain the business logic (usecases) for the domain, communicates mainly with the infrastructure layers (controllers, repositories, etc.)
    - **calculator package:** This package contain the contract and implementation to calculate the credit line, can be seen as a kind of deposit
    - **model package:** Contains the domain entities

### :file_folder: **pkg Package**
Packages that do not belong to the core of the application and have a specific functionality:
- **pkg**
    - **cache package:** Package to handle a simple cache for the non-functional requirements
    - **env package:** This packages allows to the application read a set environment variables
    - **errors package:** Package to handle all the errors in the application
    - **middleware package:** Contains rate limits and retries middlewares for non-functional requirements
    - **validator package:** Contains the functionality to validate the request