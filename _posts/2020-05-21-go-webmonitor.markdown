---
layout: post
title:  "Go-数据监控平台"
date:   2020-05-26 22:03:36 +0530
tags:   Go WebMonitor
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

![avatar](https://github.com/SonderEASE/lewis-blog.io/blob/master/pics/web-monitor-home-page.png?raw=true)

# 丰富页面

实际上 web_monitor 对 **http.NewServeMux()** 进行了封装, 增加新的页面只需要调用 AddPage() 函数即可, 所有通过AddPage()增加的网页都被添加到了 **"/"** base路由下, 同时传入 **AddPage()** 的 func 都会被保存在一个全局变量 **handles** 中,  当监听到请求时, 建立HomePage时给根路由 **"/"** 注册的 **HandleBase** 内会对路由进行匹配, 并调用传入 AddPage()的 func 完成相关的页面操作.


```go

var (
	infoList = map[string]PersonalInfo {"张三": {Name: "张三", PhoneNumber:"18706802993", Age: 16, Hobbies: "网游", Occupation: "公务员"},
		"李四": {Name:"李四", PhoneNumber:"17792384601", Age: 28, Hobbies: "篮球", Occupation: "码农"}}
)


web_monitor.AddPage("/monitor/demo", "Demo", func(page *web_monitor.Page, params *web_monitor.Params){
		// params.GetString和LinkAWithParams对应,
		if _, ok := params.GetString("通讯录"); ok {
			// 增加json界面
			page.SetJson(infoList)
			return
		}
		if name, ok := params.GetString("Name"); ok {
			if v, ok := infoList[name]; ok {
				pane := page.AddPanel()
				pane.AddInput("Name", "Name", name)
				pane.AddSubmitButton("查询")

				Table := page.AddTable("详细信息")
				Table.AddHeader("Hobbies", "Occupation")
				Table.AddRow(v.Hobbies, v.Occupation)

			} else {
				page.AddError(fmt.Errorf("error name"))
			}
			return
		}
		// 增加空白区域
		panel := page.AddPanel()

		// 空白区域中增加html中的H4元素, 元素内容是一个超链接.
		// 例：LinkAWithParams("/user", "用户", "name", "张三", "age", 13) => <a target="_blank" href="/user?name=张三&age=13">用户</a>
		panel.AddH4(web_monitor.LinkAWithParams("/monitor/demo", "通讯录", "通讯录", 1))
		contactTable := page.AddTable("电话本")
		contactTable.AddHeader("name", "age", "tel")
		var name string
		for k, v := range infoList {
			name = web_monitor.LinkAWithParams("/monitor/demo", k , "Name", k)
			contactTable.AddRow(name, v.Age, v.PhoneNumber)
		}
		// 加编号
		contactTable.AddNo()
		// 按第一列升序
		contactTable.Sort(0, true)
	})
```


![avatar](https://github.com/SonderEASE/lewis-blog.io/blob/master/pics/web-monitor-demo.png?raw=true)


除了增加表, 空白区域以外, 还有很多丰富的元素, 比如曲线, 自定义曲线, 脚本, MapTable等, 应有尽有.

# 展示更有价值的信息

分布式服务, 某一个模块自身的数据很难反应线上问题, 这时候就需要与其他的模块进行数据的交互（ **订阅** 和 **发布**）, 将 **多个模块的实时数据** 进行 **整合**, 才能展现出更有价值的数据. 为了达到这个目的我们还需要掌握一下 [**gRpc**]() ~. 
