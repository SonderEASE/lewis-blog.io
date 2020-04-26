---
layout: post
title:  "C++工具库分享"
date:   2020-04-03 19:10:00 +0530
tags: tool-list 
categories: [分享 | Share]
---
:star2::star2::star2::star2::star2:



# 开源 C++ 库列表

[这里](https://zh.cppreference.com/mwiki/index.php?title=cpp/links/libs&variant=zh-hans)已经构建了开源 C++ 库的常用列表，使人们在需要特定功能的实现时，不必再浪费时间搜索。赞:star2:

# 我的工具库
## TOC
+ toc of tools
  - <a href="#t1">EasyHttp</a> 

  - <a href="#t2">JSON for Modern C++</a> 

  - <a href="#t3">LRU-Cache</a> 

  - <a href="#t4">CryptoUtils</a>

  - <a href="#t5">Libuv Handle(uv_async_t;uv_signal_t;uv_timer_t)封装</a>

  - <a href="#t6">Thread-SafeQueue</a>

  - <a href="#t7">RWLock</a>

  - <a href="#t8">Http 接口</a>
  
  - <a href="#t9">封装好的 spdlog 日志</a>

+ toc of algorithm
  - <a href="#a1">SHA</a>

  - <a href="#a2">MD5</a>

  - <a href="#a3">...</a>


&nbsp;
### <a name="t1">EasyHttp</a> 

刚入职的时候, 基本上入坑每一个项目都是从数据上报开始做起的:joy:(模块相对独立, 便于弄清楚哪些是关键数据,个人觉得从数据上报入手熟悉项目是非常友好的), EasyHttp对curl进行了简单的封装, 使用起来极其简单~ 

[EasyHttp.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.hpp)&#8195;&#8195;[EasyHttp.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/EasyHttp/EasyHttp.cpp)

&nbsp;
### <a name="t2">JSON for Modern C++</a> 

与EasyHttp是配套使用的, 上报一般都使用的是Json数据格式, 由德国大佬nlohmann为现代C++编写的json库使用起来异常丝滑~ 如果对性能有着极高的要求可以参考这个[benchmark](https://github.com/miloyip/nativejson-benchmark), 作者对各个Json库做了[多个维度的对比](https://www.zhihu.com/question/23654513).

这位老哥写的 [使用攻略](https://blog.csdn.net/fengxinlinux/article/details/71037244) 还是比较详尽的,可以参考~

[Json.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json.hpp)&#8195;&#8195;[JsonFwd.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/nlohmann/json_fwd.hpp)

&nbsp;
### <a name="t3">LRU-Cache</a> 

在设计中继转发的时候有用到, 这个模板类实现的LRU-cache看着异常顺眼, 没有一行碍眼的代码, 不需要动脑经就能看得懂:wink:

[LRUCache.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/LRU-Cache/LRUCache.hpp)

&nbsp;
### <a name="t4">CryptoUtils</a> 

CryptoUtils是加密的好伴侣, 将md5,sha0等转换为16进制的格式. 简单实用~

[CryptoUtils.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/CryptoUtils/CryptoUtils.hpp)&#8195;&#8195;[CryptoUtils.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/CryptoUtils/CryptoUtils.cpp)

&nbsp;
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

Example::~Example() {
    ...
    timer.Close();
}

```
[Async.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Async.hpp)&#8195;&#8195;[Signal.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Signal.hpp)&#8195;&#8195;[Timer.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Event/Timer.hpp)

&nbsp;
### <a name="t6">Thread-SafeQueue</a> 

借助模板实现的线程安全队列, 可以把要处理的事件加入到线程的安全队列中,线程在适当的时候处理线程中的队列, 简单高效的实现异步执行任务的列队~

[SafeQueue.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/SafeQueue/SafeQueue.hpp)

&nbsp;
### <a name="t7">RWlock</a> 
由 [std::condition_variable(信号量)](https://en.cppreference.com/w/cpp/thread/condition_variable) 和 [std::unique_lock\<std::mutex\>](https://en.cppreference.com/w/cpp/thread/unique_lock) 结合实现的读写锁(写优先, 避免**写者饥饿**), 并进行了两层封装, 可以选择自己控制 lock 和 unlock 的时机, 也可以通过外层的封装 WriteGuard 或者 ReadGuard 当作 **scoped_lock** 使用. 

RWLock-pthread.hpp 则是针对基于pthread的多线程程序封装的读写锁~

[RWLock.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/RWLock/RWLock.hpp)
&#8195;&#8195;[RWLock-pthread.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/RWLock/RWLock-pthread.hpp)

&nbsp;
### <a name="t8">Http 接口</a>

在实际业务中, 经常会有增加数据接口的需求, 比如在线读取, 删除, 增添数据等. Monitor 对这一需求进行了完美的封装, 在需要增添新的接口时, 只需要通过RegisterRouter接口结合lambda表达式完成注册即可. 使用起来非常简单, 而且十分独立. 

注: 我这里的实现是基于 libuv 以及自己的Tcp封装库, 使用时记得配套修改, 改动起来比较简单~ 

[MonitorManager.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Monitor/MonitorManager.hpp)&#8195;&#8195;[MonitorManager.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Monitor/MonitorManager.cpp)&#8195;&#8195;[MonitorSession.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Monitor/MonitorSession.hpp)&#8195;&#8195;[MonitorSession.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Monitor/MonitorSession.cpp)&#8195;&#8195;[MonitorUtils.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/Monitor/MonitorUtils.hpp)

&nbsp;
### <a name="t9">封装好的 spdlog 日志</a>

spdlog **非常的快**, **性能** 是它的主要目标, 同时他也支持异步模式(可选)、自定义格式、条件日志、多线程/单线程日志、日志级别在线修改、可每日生成日志文件等等功能. 这里对他做一个简单的封装, 支持通过<a href="#t8">Http 接口</a>对日志级别进行在线修改, 未来如果业务有更多的需求也会更新更多的接口.

[Logger.hpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/spdlog/Logger.hpp)&#8195;&#8195;[Logger.cpp](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/spdlog/Logger.cpp)

&nbsp;
# 个人算法库

### <a name="a1">SHA</a> 

安全散列算法[SHA](https://zh.wikipedia.org/wiki/SHA%E5%AE%B6%E6%97%8F), SHA0用的比较多, 其他系列可以通过[openssl]库(https://www.openssl.org/)中的sha.h获取~

[SHA0.h](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/algorithm/SHA0/SHA0.h)&#8195;&#8195;[SHA0.c](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/tool-list/algorithm/SHA0/SHA0.c)

### <a name="a2">MD5</a> 

信息摘要算法[MD5](https://zh.wikipedia.org/wiki/MD5), 同样使用的是[openssl库](https://www.openssl.org/)中的md5.h中的算法. 终端上敲 **md5sum + filename/string** 可以查看文件或者字符串的md5值.



# The Tail End
> *Always leave the code better than you found it. —— The Boy Scout Rule*
