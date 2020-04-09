---
layout: post
title:  "深入回调函数"
date:   2020-04-08 20:40:36 +0530
tags: callback
categories: [代码 | Coding]
---
让我们稍微深入的讨论一下回调函数的相关知识:stuck_out_tongue_closed_eyes:.


背景阅读

+ 如果你还不了解什么是回调(callback), 欢迎读这篇文章:point_right:[了解回调函数](https://sonderease.github.io/lewis-blog.io/%E4%BB%A3%E7%A0%81%20%7C%20coding/2020/01/20/what-is-callback.html)

+ 如果对对象的所有权和生命周期的管理还有些陌生, 欢迎阅读:point_right:[资源管理小记](https://sonderease.github.io/lewis-blog.io/%E7%AC%94%E8%AE%B0%20%7C%20notes/2020/03/20/RAII.html)

从概念上看, 回调是一个调用函数的过程, 这个过程中有两个角色, 计算和数据, 其中计算就是函数, 而数据有两类:
+ **绑定** 的数据, 即回调的 **上下文** (*context*)
+ **未绑定** 的数据, 即执行计算(回调函数)时 **传入的参数**

捕获了上下文(绑定的数据)的回调函数就成为了 [闭包](https://zh.wikipedia.org/zh-cn/%E9%97%AD%E5%8C%85_(%E8%AE%A1%E7%AE%97%E6%9C%BA%E7%A7%91%E5%AD%A6)#%E9%97%AD%E5%8C%85%E7%9A%84%E7%94%A8%E9%80%94), 即闭包 = 函数 + 上下文.

在面向对象语言中, 对象是一等公民, 而函数不是, 所以在实现上:
+ **闭包**(*closure*) 一般通过对象实现(例如 std::function)
+ **上下文** 一般作为 **闭包对象** 的数据成员存在.

从 **对象所有权的角度** 来看, **上下文**又进一步分为:
+ **不变上下文**: 数值/字符串/结构体等基本数据类型, 永远不会失效. 使用起来不需要担心生命周期的问题.
+ **弱引用上下文(可变上下文)**: 这类上下文不为闭包所拥有, 所以执行回调的时候上下文可能失效,如果使用时没有检查,就会导致崩溃.
+ **强引用上下文(可变上下文)**: 这类上下文为闭包所拥有, 在执行回调时可以保证上下文始终有效,但是如果使用后忘记释放,就会导致内存泄露.

可能你已经熟悉了[std::bind](https://en.cppreference.com/w/cpp/utility/functional/bind)/[lambda](https://en.cppreference.com/w/cpp/language/lambda) + [std::function](https://en.cppreference.com/w/cpp/utility/functional/function), 但是你在设计C++回调时, 是否考虑过**以下几个问题?**

1. 回调是同步的还是异步的?
    + 回调时(弱引用)上下文会不会失效?
    + 如何处理失效的(弱引用)上下文?
2. 回调只能执行一次还是多次?
    + 为什么区分一次还是多次?
    + 何时销毁(强引用)上下文?

&nbsp;
&nbsp;
# <a name="t1">回调是同步的还是异步的</a> 
**同步回调** (sync callback) 在 **构造闭包** 的 **调用栈** (call stack) 里 **局部执行**。例如，累加一组得分（使用 lambda 表达式捕获上下文 total）


```c++
int total = 0;
std::for_each(std::begin(scores), std::end(scores), 
              [&total](auto score){ total+=score; });
// 这里的total作为上下文是始终有效的.
```
+ **绑定的数据**: total, 局部变量上下文(弱引用, 所有权在闭包外, 但是生命周期比闭包长)
+ **未绑定的数据**: score, 每次迭代传递的值

![avatar](https://github.com/SonderEASE/lewis-blog.io/blob/master/pics/callback-sync.png?raw=true)

**异步回调** (async callback) 在构造后存储起来，在 **未来某个时刻（不同的调用栈里）非局部执行**。例如，用户界面为了不阻塞 **UI 线程** 响应用户输入，在 **后台线程** 异步加载背景图片，加载完成后再从 **UI 线程** 显示到界面上：

```c++
// call back code
void View::LoadImageCallBack(const Image& image) {
     // 这里this不一定有效了
     if (background_image_view) {
         background_image_view->SetImage(image);
     }
}
// client code 
FetchImageAsync(filename, std::bind(&View::LoadImageCallback, this));

// 用lambda可以等效为以下代码:
FetchImageAsync(filename, std::bind([this](const Image& image) {
    if (background_image_view) {
         background_image_view->SetImage(image);
     }
}));

```
+ **绑定的数据** bind绑定了View对象的this指针(弱引用)
+ **未绑定的数据** View::LoadImageCallBack 的参数 const Image& image

![avatar](https://github.com/SonderEASE/lewis-blog.io/blob/master/pics/callback-async.png?raw=true)

> 注: 
View::FetchImageAsync是基于Chromium的多线程任务模型(参考:[Keeping the Browser Responsive \| Threading and Tasks in Chrome](https://github.com/chromium/chromium/blob/master/docs/threading_and_tasks.md#keeping-the-browser-responsive))

&nbsp;
## <a name="t1.1">回调时(弱引用)上下文会不会失效</a>

前面已经说了, 闭包并不拥有 **弱引用上下文**, 所以上下文可能失效.
+ 对于 **同步回调**, 一般上下文的生命周期长于闭包的生命周期, 所以一般不会失效
+ 对于 **异步回调**, 闭包 并不知道 上下文 的生命周期是否已经结束. 

例如在上面的例子中 **异步加载图片** 的场景, 在等待加载时, 用户可能已经退出了界面. 所以, 在执行View::LoadImageCallback 时:
+ 如果界面还在显示, View对象依然有效, 则执行 SetImage 显示背景图片
+ 如果界面已经退出, View对象已经失效, 那面 background_image_view_ 就成了 **野指针**, 调用 SetImage 就会导致 **崩溃**.

在Chromium的设计中(Base::bind 替换 std::bind))上述的代码都是无法通过编译的(Chromium做了对应的 静态断言 ), 因为传给Base::bind的参数是不安全的, 传递普通对象的 **裸指针(比如this)**, 容易导致悬垂引用, 而且传递了上下文的lambda表达式, 无法检测lambda表达式捕获的 **弱引用的有效性**.

C++核心指南(C++ Core Guidelines)中对此也有讨论:
+ [F.52: Prefer capturing by reference in lambdas that will be used locally, including passed to algorithms](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#Rf-reference-capture) 
+ [F.53: Avoid capturing by reference in lambdas that will be used non-locally, including returned, stored on the heap, or passed to another thread](https://isocpp.github.io/CppCoreGuidelines/CppCoreGuidelines#Rf-value-capture)

&nbsp;
## <a name="t1.2">如何处理失效的(弱引用)上下文</a>

处理的方法就是在弱引用失效的时候, 及时的 **取消回调**, 例如 异步加载图片 的代码, 可以给std::bind传递 View 对象的 **弱引用指针**, 也就是std::weak_ptr<View> :

```c++
FetchImageAsync(filename, base::Bind(&View::LoadImageCallback, AsWeakPtr()));
// 传递weakptr还不是裸指针this
```

在执行 View::LoadImageCallback 时:
+ 如果界面未退出, View对象有效, 执行SettingImage显示背景图片.
+ 如果界面已经退出, 那么弱引用指针也失效了(weak_p.expired())), 那么就退出回调函数.

> 注: 
> + AsWeakPtr是Chromium的实现, base::WeakPtr属于侵入式的智能指针, 非 线程安全.
> + base::Bind 针对 base:WeakPtr扩展了base::IsWeakReceiverM<>检查, 调用前判断弱引用有效性, 可参考: [Binding A Class Method With Weak Pointers \| Callback<> and Bind()](https://github.com/chromium/chromium/blob/master/docs/callback.md#binding-a-class-method-with-weak-pointers)
> + 也可以基于std::weak_ptr表示弱引用所有权, 有一些需要注意的地方可以学习这篇文章: [弱回调 \|《当析构函数遇到多线程 —— C++ 中线程安全的对象回调》陈硕](https://github.com/downloads/chenshuo/documents/dtor_meets_mt.pdf)

...

# The Tail End
本文只是试图讲清楚回调函数的基本概念与使用方式, 如果有什么问题, 欢迎与我进行交流~ :speech_balloon:  后续我还会对回调相关内容做多的分享, 敬请期待.:ghost: