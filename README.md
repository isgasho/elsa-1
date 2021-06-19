### elsa

***

#### 概述

​		*elsa* 是基于 *gRPC* 封装的一套简单、高效、易用的 *golang* 微服务框架。由于 *gRPC* 天生支持跨语言的特性， *elsa* 也天生支持多语言、跨语言的 *RPC* 调用，当前框架仅实现了 *Go*语言版本的 *Client* 支持，*Java* 版 *[Client](http://wwww.github.com/busgo/elsa-java)* 将在近期放出 。具体支持语言种类请参考 *[grpc](https://www.grpc.io/)* 官网。

#### 框架特性

1. 基于 *gRPC* 高性能开源 http2服务框架封装。
2. 采用分布式 *CAP* 理论的 *AP* 实现服务注册、发现服务保障服务高可用。
3. 简单、易用、统一接口通讯采用 *proto* 文件定义接口约束。
4. 统一网关API服务、内部通过 *gRPC* 通讯，外部入口统一通过统一 *API* 网关 进行统一流控、鉴权等。
5. 丰富的多语言 *Client* 支持，按需选择语言进行快速微服务集成。



### gRPC服务注册发现



<img src="/Users/apple/Desktop/wp/code/home/go/elsa/docs/img/arch.png" alt="架构图" style="zoom:60%;" />

#### 快速开始

***

##### 依赖

在使用之前需要先安装相关依赖

* [Go](https://github.com/golang) 请选择最后1.13.5以后[版本](https://golang.org/doc/devel/release.html)。 

* [Protocol buffer](https://developers.google.com/protocol-buffers) 编译二进制文件 *protoc* 采用 [*proto3* ](https://developers.google.com/protocol-buffers/docs/proto3)版本。具体安装请参考 [安装向导](https://www.grpc.io/docs/protoc-installation/) 部分。

* *go-grpc* 自动生成 *grpc go* 代码 插件安装，在终端执行如下命令。将生成的 *protoc-gen-go-grpc* 文件添加环境变量中。

  ```shell
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
  ```

   

  ##### 注册中心服务配置

  下载代码

  ``` shell
  git clone https://www.github.com/busgo/elsa
  ```

  ​      配置注册中心集群地址 (格式 ip1:port1, ip2:port2)

  ```shell
  cd elsa/cmd/registry
  go build .
  nohup ./elsa >> elsa.log 2>&1 &
  ```

  ##### 服务提供者

  

  ##### 服务消费者

  

  ##### TODO

  ***

  1. 统一网关API 
  2. Java Client 
  3. 链路追踪集成
  4. 分布式配置服务集成
  5. 分布式任务调度集成

  

