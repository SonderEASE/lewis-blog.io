---
layout: post
title:  "C++工具库分享"
date:   2020-03-20 19:10:36 +0530
tags: tool-list 
categories: [分享 | Share]
---
:star2::star2::star2::star2::star2:

# Tool List
+ **EasyHttp**

刚入职的时候, 基本上入坑每一个项目都是从数据上报开始做起的:joy:(模块相对独立, 便于弄清楚哪些是关键数据,个人觉得从数据上报入手熟悉项目是非常友好的~), EasyHttp对curl进行了简单的封装, 使用起来极其简单~ 

[EasyHttp.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.hpp)&#8195;&#8195;[EasyHttp.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.cpp)

+ **JSON for Modern C++**

与EasyHttp是一套的, 上报一般都使用的是Json数据格式, 由德国大佬nlohmann为现代C++编写的json库使用起来异常丝滑~ 如果对性能有着极高的要求可以参考这个[benchmark](https://github.com/miloyip/nativejson-benchmark),作者对各个json库做了[多个维度的对比](https://www.zhihu.com/question/23654513).

[json.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json.hpp)&#8195;&#8195;[json_fwd.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json_fwd.hpp)

# <font face="微软雅黑" >the-tail-end</font>
> *Always leave the code better than you found it. —— The Boy Scout Rule*
