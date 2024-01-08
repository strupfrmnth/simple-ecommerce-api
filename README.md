# simple-ecommerce-api
## Introduction
This project is my implementation of E-commerce API which uses Golang + Redis + PostgreSQL. It provides a set of RESTful APIs, allowing client applications to perform various operations related to products, orders, and users. This is an ongoing project, and continuous improvements will be made to enhance its functionality and features.

## Getting Started
1. Set up .env file to fit the database configuration
2. Change the configs/config.json file as you need
3. go build ./cmd
4. Run cmd executable file

## Features
* JWT for authencation
* Rate limiting for too many requests
* Limit repeatable requests from same IP during a period
* CRUD operations on users, products and orders
* Password security

## API
A simple description of the API is listed below.

### User
#### Requests
* `POST /user/register`
* `POST /user/login`
* `GET /users`
    * Get all users
* `GET /users/{id}`
* `PUT /users/{id}`
* `DELETE /users/{id}`

### Product
#### Requests
* `GET /products`
    * Get all products
* `GET /products/{id}`
* `POST /products`
    * Add a product
* `PUT /products/{id}`
* `DELETE /products/{id}`

### Order
#### Requests
* `POST /order`
    * Add an order

### Others
#### Request
* `GET /`
    * for testing IP limit
* `GET /rate`
    * for testing too many requests