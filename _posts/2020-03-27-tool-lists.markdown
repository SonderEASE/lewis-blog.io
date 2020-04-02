---
layout: post
title:  "C++工具库分享"
date:   2020-03-20 19:10:36 +0530
tags: tool-list 
categories: [分享 | Share]
---
:star2::star2::star2::star2::star2:



# 开源 C++ 库列表

[此页面](https://zh.cppreference.com/mwiki/index.php?title=cpp/links/libs&variant=zh-hans)的目的是构建开源 C++ 库的比较列表，使得人们在需要特定功能的实现时，不必浪费时间在网上（ DuckDuckGo 、谷歌、必应等）搜索。赞:star2:


# TOC
+ toc of tools
  - <a href="#t1">EasyHttp</a> 

  - <a href="#t2">JSON for Modern C++</a> 

  - <a href="#t3">LRU-Cache</a> 

  - <a href="#t4">CryptoUtils</a>

  - <a href="#t5">libuv handle封装</a>

  - <a href="#t6">...</a>

+ toc of algorithm
  - <a href="#a1">SHA</a>
  - <a href="#a2">MD5</a>
  - <a href="#a3">...</a>


# 个人工具库

[JSON For Modern C++](###JSON for Modern C++)

### <a name="t1">EasyHttp</a> 

刚入职的时候, 基本上入坑每一个项目都是从数据上报开始做起的:joy:(模块相对独立, 便于弄清楚哪些是关键数据,个人觉得从数据上报入手熟悉项目是非常友好的), EasyHttp对curl进行了简单的封装, 使用起来极其简单~ 

[EasyHttp.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.hpp)&#8195;&#8195;[EasyHttp.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.cpp)

### <a name="t2">JSON for Modern C++</a> 

与EasyHttp是配套使用的, 上报一般都使用的是Json数据格式, 由德国大佬nlohmann为现代C++编写的json库使用起来异常丝滑~ 如果对性能有着极高的要求可以参考这个[benchmark](https://github.com/miloyip/nativejson-benchmark), 作者对各个Json库做了[多个维度的对比](https://www.zhihu.com/question/23654513).

[Json.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json.hpp)&#8195;&#8195;[JsonFwd.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json_fwd.hpp)

### <a name="t3">LRU-Cache</a> 

在设计中继转发的时候有用到, 这个模板类实现的LRU-cache看着异常顺眼, 没有一行碍眼的代码, 不需要动脑经就能看得懂:wink:

[LRUCache.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/LRU-Cache/LRUCache.hpp)

### <a name="t4">CryptoUtils</a> 

CryptoUtils是加密的好伴侣, 将md5,sha0等转换为16进制的格式. 简单实用~

[CryptoUtils.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/CryptoUtils/CryptoUtils.hpp)&#8195;&#8195;[CryptoUtils.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/CryptoUtils/CryptoUtils.cpp)

### <a name="t5">libuv handle封装</a> 

借助模板类对经常用到的uv_async_t, uv_signal_t以及uv_timer_t进行了封装, 结合using使用起来十分干脆利落. 来个例子秀一下~ 有时间我会把libevent和libev也封装一下的.

```c++
#include "Event/Timer.hpp"

class Example {
    using TimerType = uv::Timer<Example>

public:
    ...
    void DoSth();
private:
    void timer1_cb() {}

private:
    uv_loop_t* loop_{};
    TimerType timer1;
}

void Example::DoSth() {
    timer.Init(loop_, this, &Example::timer1_cb);
    timer.Start(10, /*10*/);
    timer.Stop();
}

```

[Async.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Async.hpp)&#8195;&#8195;[Signal.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Signal.hpp)&#8195;&#8195;[Timer.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Timer.hpp)



# 个人算法库

### <a name="a1">SHA</a> 

安全散列算法[SHA](https://zh.wikipedia.org/wiki/SHA%E5%AE%B6%E6%97%8F), SHA0用的比较多, 其他系列可以通过[openssl]库(https://www.openssl.org/)中的sha.h获取~

[SHA0.h](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/algorithm/SHA0/SHA0.h)&#8195;&#8195;[SHA0.c](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/algorithm/SHA0/SHA0.c)

### <a name="a2">MD5</a> 

信息摘要算法[MD5](https://zh.wikipedia.org/wiki/MD5), 同样使用的是[openssl库](https://www.openssl.org/)中的md5.h中的算法. 终端上敲 **md5sum + filename/string** 可以查看文件或者字符串的md5值.



# The Tail End
> *Always leave the code better than you found it. —— The Boy Scout Rule*
