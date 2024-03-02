# Project Introduction
* **Read this in other languages: [English](README_en.md) | [中文](README.md).**
This project is an e-commerce microservice project built using grpc+gin, and uses consul to register and discover microservices. The e-commerce project is divided into 9 microservices and a front-end client, including: Merchandise microservices, shopping cart microservices, currency microservices, advertising microservices, product recommendation microservices, mail microservices, payment microservices, delivery microservices, settlement microservices, and a front-end.
#### Technology stack
* grpc
* ProtocolBuffer
* gin
* consul
* ...

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/mic.png)


## Use:
* go language environment
* consul Service
* Development tools and other environments <br>
1. Clone the project to a local directory
```
git clone https://github.com/1280019840/Microservice-mall.git
```
2.cd to the Microservice-mall folder
```
cd Microservice-mall
```
3.Open terminal Startup consul
```
consul agent -dev
```
4.Launch each system microservice file:
For example, enter productcatalogservice
```
cd productcatalogservice
```
5.Starting microservices
```
go run main.go
```
Repeat the following steps until each microservice has started successfully
#### Note: The checkoutservice is started last
6.Go to front-end folder
```
cd frotend
```
7.Start the front-end file, handler, and middleware
```
go run middleware.go handler.go rpc.go main.go
```
Therefore, after microservices are started successfully, you can view consul
```
http://localhost:8500/ui/dc1/services
```
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/consul.png)
<br>

## Effect display
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/home1.png)
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/home2.png)

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/details1.png)
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/details2.png)

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/cart_nil.png)

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/pay1.png)
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/pay2.png)
#### There are pages to download source code view<br>
#### Thanks for watching, remember star thank

