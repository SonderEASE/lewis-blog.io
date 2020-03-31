---
layout: post
title:  "了解回调函数"
date:   2020-01-20 20:40:36 +0530
tags: callback
categories: [代码 | Coding]
---
让我来试着讲清楚,什么是回调函数:stuck_out_tongue_closed_eyes:.


# 为什么写这篇文章

&#8195;&#8195;在学习回调函数的过程中, 我经历了三个阶段: **不知道** -> **知道** -> **使用**, 现在我希望自己还可以更进一步, 把回调函数向别人解释清楚. 

# 什么是回调函数

回调函数是
+ **可执行代码** 作为**参数**传入**其他可执行代码**
+ 并由 **其他可执行代码** 执行这段 **可执行代码**
的过程

这段理解是基于wikipedia对于callback的解释, 涉及到同步与异步的概念稍后会给大家做解释.
> *In computer programming, a callback, also known as a "call-after" function, is any executable code that is passed as an argument to other code that is expected to call back (execute) the argument at a given time. This execution may be immediate as in a synchronous callback, or it might happen at a later time as in an asynchronous callback. Programming languages support callbacks in different ways, often implementing them with subroutines, lambda expressions, blocks, or function pointers.*

来看伪代码理解下这个过程~
```c++
void A() {
    std::cout<<"output"<<std::endl;
}

void B(function fn) {
    fn();
}

int main () {
    B(A); //时刻T
}
```
+ 函数A作为参数 fn 传入函数B
    + B不知道A, 只知道参数fn
    + fn对于B来说和普通的参数没有区别, 只是一个局部变量.
+ B通过fn的语法,调用A.
    + 在某一时刻T, B调用了函数fn
    + 如果B传入的参数是A,那么得到的结果就是 "output"


&#8195;&#8195;这里对于函数作为参数,有两点可能需要解释 : 函数是程序设计语言的[一等公民(firt-class function)](https://en.wikipedia.org/wiki/First-class_function)(A), 将函数作为参数的函数,叫做[高阶函数(higher-order function)](https://en.wikipedia.org/wiki/Higher-order_function)(B).

# 回调场景

&#8195;&#8195;一家旅馆提供**叫醒服务**:alarm_clock:(对应上面的B)，但是要求旅客自己决定叫醒的方法。可以是**打客房电话**，也可以是**派服务员去敲门**，睡得死怕耽误事的，还可以要求往自己**头上浇盆水**。这里，"叫醒"这个行为是旅馆提供的，但是叫醒的方式是由旅客决定并告诉旅馆的，在上面的例子中可以将A理解为叫醒方式的一种(比如A就是打客房电话),也就是传入B的参数, 如果传入的参数是A,那么旅店就会在某个时刻T, 叫醒旅客(B), B回调A, 给旅客打电话.

# 场景实现

```c++
int main() {
    auto hotel = std::make_shared<Hotel>("wanda");
    auto Mia = std::make_shared<Passenger>("Mia", *hotel, 601);

    // Mia选择的叫醒方式是大喊三声Mia
    Mia->m_hotel.OrderWakeUpServer(Mia->room_id, []()->void {
        std::cout<<"Mia!!!\nMia!!!\nMia!!!"<<std::endl;
    });
    // 当然Mia也可以选择酒店提供的叫醒方式, 比如敲门
    Mia->m_hotel.OrderWakeUpServer(Mia->room_id, Hotel::Knock);

    //这个wanda酒店非常的low, 他只提供8点整的叫醒服务.
    //8点了~ 叫醒住在601的mia,用她告诉酒店的方式
    hotel->WakeUp(601);

    return 0;
}

//OrderWakeUpServer的实现 Hotel.hpp
using WakeUpMode = std::function<void()>;
void OrderWakeUpServer(int room_id, const WakeUpMode& mode) {
    wake_lists.insert(std::make_pair(room_id, mode));
};
```
[Hotel类的代码链接](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/What-is-callback/Hotel.hpp)

[Passenger类的代码连接](https://raw.githubusercontent.com/SonderEASE/lewis-blog.io/master/BlogCode/What-is-callback/Passenger.hpp)

在这个例子里我们也可以看到[上面](#什么是回调函数)说的5个要素, Mia自己定义的lambda表达式就是A, Hotel类的OrderWakeUpServer就相当于的B, B的参数WakeUpMode就相当于fn, 8点就是时刻T, 大喊的三声Mia就是回调结果"output".

# 同步回调与异步回调

&#8195;&#8195;无论是什么程序设计语言, 回调函数都分为 **同步** 和 **异步** 两种, 这里以c++为例, 分析下两者的区别.

## 同步回调

下面的代码展示了如何给一个字符串数组按照长度排序. 代码调用了[sort](https://en.cppreference.com/w/cpp/algorithm/sort)进行排序, 并通过指定isShorter为参数,实现长度的比较.
```c++
bool isShorter(const string &s1, const string &s2) {
    return s1.size() < s2.size();
}

...

vector<string> words = {"doudou", "xiaoxiaoxiong", "mona", "lol"};
sort(words.begin(), words.end(), isShorter);

```
实际上对于algorithm的定制操作就是很典型的同步回调, 代码对应了[什么是回调函数](#什么是回调函数)中的5个要素:
+ sort 对应 B 
+ isShorter 对应 A
+ sort的第三个参数 对应 fn
+ 在比较发生的时候 对应 时刻T.
+ isShorter比较返回的结果 对应 回调结果"output" 

由于调用isShort的时刻T,均是在调用sort结束之前,所以这样的回调被称为**同步回调**.

## 异步回调

下面的代码展示了如何在Linux下,组织用户使用ctrl+c退出程序, 并打印提示, 代码调用了signal进行了回调函数的注册(和同步回调不同,这里仅是注册).

```c++
void block_interrupt (int code) { printf("\rSorry We Can't Let You Shut Down This by Press ^C\n"); }

...

signal (SIGINT, block_interrupt);
```

同样的代码对应了[什么是回调函数](#什么是回调函数)中的5个要素:
+ signal 对应 B 
+ block_interrupt 对应 A
+ signal的第二个参数 对应 fn
+ 用户按下ctrl+c的时候 对应 时刻T
+ block_interrupt的打印提示 对应 回调结果"Sorry We Can't Let You Shut Down This by Press ^C" 

由于调用block_interrupt的时刻T是用户按下组合键的时刻, 均是在调用singal之后, 所以这样的回调方式被称为**异步回调**

[回调场景](#回调场景)的例子也是采用的异步回调方式.

**##同步与异步的不同**
+ 同步方式通过 **参数传递** 的方法传递回调函数, 调用者直接使用回调函数,从而完成回调.
+ 异步方式通过 **注册** 的方式, 告知未来的调用者, 并将回调函数存储下来, 调用者在未来的某个时刻T, 取出并调用回调函数, 从而完成回调. 比如我们将Mia定义的叫醒方式和房间号相关联存放在wake_lists中, 由酒店决定在适当的时机调用存放在wake_lists中的方法来叫醒Mia.


# The Tail End
本文只是试图讲清楚回调函数的基本概念与使用方式, 如果有什么问题, 欢迎与我进行交流~ :speech_balloon:  后续我还会对回调相关内容做多的分享, 敬请期待.:ghost: