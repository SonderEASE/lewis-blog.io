---
layout: post
title:  "Go-数据监控平台"
date:   2020-05-26 22:03:36 +0530
tags:   Go WebMonitor gRpc
categories: [代码 | Coding]
---
最近刚刚接手了直播调度, 其中有一个webmonitor包十分吸引我, 它可以十分快捷的实现一个实时数据展示界面, 再结合gRpc,与其他模块进行数据的采集与订阅, 一个实用的数据监控平台就建立起来了:stuck_out_tongue_winking_eye:

&nbsp;
# 建立HomePage

```go

web_monitor.Init(
        8888,  // 监听端口
		web_monitor.WithVersion("v1.0"), // 版本号
 		web_monitor.WithServiceName("web-monitor"),  // title
        web_monitor.WithPProf(true),  // 是否开启 pprof
		web_monitor.WithAccount(func(username, pwd string) bool { //设置网页登录账号密码
			if "lewis,112358" == username+","+pwd {
				return true
			}
			return false
        }),
        /*
        web_monitor.WithNoAccount()
        */
)

```

Init 除了可以传入包里已经写好的接口以外完全可以自己再加入新的接口, Init内部还实现了很多常用的信息网页, 比如/sys/load, /sys/cpu, /sys/mem等.

通过这样一个接口,**HomePage**就建立起来了, HomePage的样子,展现什么数据, 也可以在 **handleHomePage(w http.ResponseWriter, r \*http.Request)** 里进行调整. 

![avatar](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/pics/web-monitor-home-page.png)

# 丰富页面

实际上 web_monitor 对 **http.NewServeMux()** 进行了封装, 增加新的页面只需要调用 AddPage() 函数即可, 所有通过AddPage()增加的网页都被添加到了 **"/"** base路由下, 同时传入 **AddPage()** 的 func 都会被保存在一个全局变量 **handles** 中,  当监听到请求时, 建立HomePage时给根路由 **"/"** 注册的 **HandleBase** 内会对路由进行匹配, 并调用传入 AddPage()的 func 完成相关的页面操作.

