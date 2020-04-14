---
layout: post
title:  "delete this"
date:   2020-04-13 13:39:00 +0530
tags: delete
categories: [代码 | Coding]
---
一道"古怪"的面试题, delete this 有没有必要的使用场景 :open_mouth:


# delete this

看项目代码的时候,发现有用到delete this的地方, 觉得有一点困惑, 顺手查了一下, 就看到了这道["古怪"的面试题](https://www.v2ex.com/t/559047), 完整问题: 有什么要使用 delete this 的 **必要应用场景**, 即: 不使用 delete this 会使代码变得 confused, 异常冗余, 或难以实现功能.

&nbsp;
## 如何安全的使用 delete this

先试着弄清楚, delete this的使用条件.  首先 delete this 在**非析构函数中是合法的**, 然后再看如何安全的使用:
+ 确保对象是在 **堆** 上的 (eg : new)
+ 确保 delete this 后 **不会再用该对象调用其他(非静态)成员函数**.
+ 确保 delete this 后 **不再访问对象的任何部分**.
+ 确保 delete this 后 **this 指针不会再被访问**.

参考 [C++ Standard Memory Management F&Q](https://isocpp.org/wiki/faq/freestore-mgmt#delete-this)

&nbsp;
## 使用 delete this 的场景

首先来说, 我认为不存在必须使用 delete this 的场景, 只能说有的场景使用 delete this 会更好一些.

+ UI界面关闭

当按下某个按钮关闭窗口时, 在触发的回调函数中调用delete this(窗口关闭后不会再对该窗口进行访问等操作), 不再需要使用额外的操作去管理窗口的生命周期.

+ 异步事件队列

```c++
/* 伪代码 */
// create event thread 
new op = Operation();
op->post(dst_thread_id);

// handl event thread
Operation* op = nullptr;

/* 使用 delete this */  
while(op_queue->Get(&op)) {
    op->DoIt(); 
}

void Operation::DoIt() {
    // do sth ...
    if(condition)
        delete this;
    else 
        this->post(dst_thread_id);
}

/* 不使用 delete this */  
while(op_queue->Get(&op)) {
    op->DoIt();
    delete op; 
}

void Operation::DoIt() {
    // do sth ...
    if(condition)
        return;
    else {
        new op = Operation();
        op->post(this->dst_thread_id);
    } 
}

```

在队列中的事件执行动作之后, 如果满足一些条件(比如不需要再次回调), 就可以使用 delete this 来对事件进行清理, 如果将事件的生命周期交给外界处理, 比如从队列中取出并执行完DoIt后析构事件, 可能会导致有一些 **需要多次回调** 的事件 **反复创建和析构**, 显得不够灵活. 

+  ...

&nbsp;
# The Tail End
以后遇到其他的 delete this 会在这里进行补充~