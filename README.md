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

#### 用户端文件如下：

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/vue-user.png)

#### 后台管理员端文件如下：

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/vue-admin.png)

#### 数据库表如下：

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/MySQL.png)

## 使用：
* Java语言环境
* MySQL数据库
* Vue安装
#### 后端项目：
1.克隆项目到本地
```
gi clone https://github.com/1280019840/Sportswear-mall.git
```
2.进入后端目录：
```
cd admin-main
```
3.更新pox.xml文件的依赖
4.配置数据库，图片存放位置等信息

#### 用户端：
1.进入项目目录
```
cd vue-index
```
2.安装依赖
```
npm install
```
3.启动项目
```
npm run dev
```

#### 后台管理员端：
1.进入项目目录
```
cd vue-admin
```
2.安装依赖
```
npm install
```
3.启动项目
```
npm run serve
```

## 效果展示
![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/home1.png)
![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/home2.png)
![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/home3.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/register.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/login.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/category.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/details.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/cart.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/pay.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/order.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/forum.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/admin_home.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/goods_order.png)

![image](https://github.com/1280019840/Sportswear-mall/raw/main/img/slideshow_admin.png)

#### 还有的页面可下载源码查看<br>
#### 感谢观看，记得star谢谢

