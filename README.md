# 项目简介
#### 其他语言版本: [English](README_en.md) | [中文](README.md).
#### 系统开发采用了前后端分离：
该项目是使用grpc+gin来构建的一个电商微服务项目，使用consul来注册和发现微服务，该电商项目被划分了9个微服务和一个前端客户端，包括有：商品微服务、购物车微服务、货币微服务、广告微服务、商品推荐微服务、邮件微服务、付款微服务、配送微服务结算微服务、以及一个前端。
#### 技术栈
* grpc
* ProtocolBuffer
* gin
* consul等等
  
## 文件介绍
* 商品微服务（productcatalogservice）
* 购物车微服务（cartservice）
* 货币微服务（currencyservice）
* 广告微服务（adservice）
* 商品推荐微服务（recommendationservice）
* 邮件微服务（emailservice）
* 付款微服务（paymentservice）
* 配送微服务（shippingservice）
* 结算微服务（checkoutservice）
* 前端（frotend）

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/mic.png)


## 使用：
* go语言环境
* consul服务
* 开发工具等环境<br>
1.克隆项目到本地
```
gi clone https://github.com/1280019840/Microservice-mall.git
```
2.cd到Microservice-mall文件夹
```
cd Microservice-mall
```
3.打开终端启动consul
```
consul agent -dev
```
4.启动每一个系统微服务文件：
比如：进入商品微服务productcatalogservice
```
cd productcatalogservice
```
5.启动微服务
```
go run main.go
```
重复以下步骤，直到每个微服务都启动成功
#### 注意：结算微服务（checkoutservice）最后启动
6.进入前端文件夹
```
cd frotend
```
7.启动前端文件、handler、中间件
```
go run middleware.go handler.go rpc.go main.go
```
所以微服务都启动成功后可以查看consul
```
http://localhost:8500/ui/dc1/services
```
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/consul.png)

## 效果展示
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/home1.png)
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/home2.png)

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/details1.png)
![image](https://github.com/1280019840/Microservice-mall/raw/main/img/details2.png)

![image](https://github.com/1280019840/Microservice-mall/raw/main/img/cart_nil.png)
#### 还有的页面可下载源码查看<br>
#### 感谢观看，记得star谢谢

