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
