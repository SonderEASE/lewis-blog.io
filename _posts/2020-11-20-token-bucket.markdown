---
layout: post
title:  "流量控制 令牌桶限制(TBF,Token Bucket Filter)"
date:   2020-12-08 21:00:00 +0530
tags:   tbf
categories: [分享 | Share]
---
一

&nbsp;
## <a name="t1"> 什么是令牌桶 </a> 

令牌桶过滤器(TBF,Token Bucket Filter)是一个简单的队列规则：只允许以不超过预先设定的速率到来的请求通过，但可能允许短暂的突发请求通过。

TBF的实现在于一个缓冲器（桶），该缓冲器（桶）不断地被一些叫做”令牌”的虚拟数据以特定速率(token rate)填充着。桶最重要的参数就是它的大小，也就是它能够存储令牌的数量。

[wiki](http://en.wikipedia.org/wiki/Token_bucket)
[图文理解推荐](https://studygolang.com/articles/29683?fr=sidebar)

## <a name="t2"> 常用实现 </a>

+ 系统调用

```c++

tc qdisc add dev ens6 root tbf rate 1mbit  latency 20ms burst 5k

/*设置限速 1m bit，令牌桶大小 5k bytes, latency 20ms*/
```

+ "github.com/juju/ratelimit"


## <a name="t3"> ratelimit库优化 </a>

"github.com/juju/ratelimit"库的一些缺点:

+ 初始化后的令牌桶无法进行参数调整(桶容量, 限制速率)
+ 没有延时等待机制

### 动态调整桶参数
&#8195;&#8195;在实际的应用中, 流量控制往往是需要根据资源的负载去进行动态调整的, 当机器负载低, 资源充足的时候, 我们希望服务更多的请求, 当机器负载高, 资源紧张的时候, 我们希望限制更多的请求以防止服务出现异常. 所以一个不能调节限制参数的令牌桶几乎是没有应用空间的.


### 延时等待
&#8195;&#8195;考虑这样一个场景, 我们根据服务的实际运营状况, 决定每秒只接收1个客户, 以防止出现拥塞等问题. 假设某一秒同时有两个请求, 根据我们的算法规定, 我们会接受一个请求, 拒绝一个请求, 但是如果恰好只是这两个请求时间冲突了, 之后几乎没有什么请求来我们的服务器, 我们的服务根本没有压力, 那么我们死板地按照算法拒绝一个请求的做法就显得不是很合理了. 

令牌桶规定了请求的频率, 任何高于这个频率的请求都会被拒绝掉, 为了让算法更好应对刚才我们所说的场景, 我们需要加入一个新的参数, 延迟等待参数,  我们可以在拒绝一个请求时, 等待一个小的时间间隔, 这样当两个请求同时到来时, 如果之后没有请求了, 那们短暂的等待就可以换来一个令牌. 对于一些想尽量少的拒绝请求的服务来说,这样的做法非常的有用.

### ratelimit 优化库分享

[ratelimit.go](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/go-tool-list/ratelimit/ratelimit.go)&#8195;&#8195;[ratelimit_test.go](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/go-tool-list/ratelimit/ratelimit_test.go)
[reader.go](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/go-tool-list/ratelimit/reader.go)&#8195;&#8195;


