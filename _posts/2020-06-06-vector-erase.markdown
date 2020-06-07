---
layout: post
title:  "STL的内存管理"
date:   2020-06-06 17:00:00 +0530
tags:   erase  STL 内存管理
categories: [代码 | Coding]
---
vector使用不当导致的coredump :neckbeard: :collision:

&nbsp;
## 错误使用vector导致崩溃的场景

在直播后台的代码里， 有一个专门用来做数据统计的线程， 需要一个数据结构去存储每一路流对应的信息，包括流量统计等。 存储流信息的结构并不需要强调顺序， 所以之前一直用的是 **unordered_map** 以减少 **查询**， **删除** 流信息时的耗时. 一个新的功能是需要将流按它的流量进行排序， 对流量低的一些流进行特殊处理， 并上报给数据统计中心。 那么就需要对流按流量（value）进行排序， 考虑到上报周期是10s, 且存储的流数即使在高峰期也只有1000多条而已， 我决定将存储结构换成**vector**， 这样在排序时候不需要使用 **std::transform** 进行转移再排序。

将 unorderd_map -> vector  崩溃就发生了~ 

## Why? - STL的内存管理

在这个业务场景里， 存储结构的生存周期是 **从始至终** 的， 除非做数据统计的线程退出， 否则这个存储结构将不会结束自己的**生命周期**。 而vector本身相比于map又有一个特点（顺序容器相比于关联容器）。 那就是 **vector 是不会对内存进行自动回收的**。

```c++ 
int main () {
    vector<int> vec = {1,2,3,4,5,6};
    // clear vec的内存并没有被释放，  只是vec的迭代器范围发生了改变。
    vec.clear();
    
    // vec有效迭代器的范围发生了改变
    std::cout<<"遍历vec: "
    for(auto iter = vec.begin(); iter!=vec.end(); ++iter) {
         std::cout<< *(iter)<<" ";
    }
    std::cout<<std::endl;
    std::cout<<"vec的内存："
    // 但是vec的内存并没有被释放
    auto iter = vec.begin();    
    for(int i = 0; i<6; ++i) {
        std::cout<< *(iter+i) << "";
    }
}
```
输出结果:
```
遍历vec: 
vec的内存： 123456
```

那么vector这种内存不会自动回收的机制(erase和clear也是一样的)，就会导致一个问题，用来存储流信息的vec的 **内存会一直增长**，随着时间的推移，内存就会被干爆 :collision: coredump就在内存耗尽的重点等着我们。 那unordered_map为什么就ok？ 因为stl的map以及unordered_map（关联容器）会对内存进行回收管理，总结起来具体的做法就是当容器中的元素占用内存总和小于128字节时， clear和erase不会释放内存，但是当大于128字节时，则系统会直接调用malloc或者free进行内存的分配。具体可以参考[这篇文章](https://cloud.tencent.com/developer/article/1157252).


## 那非要用vector怎么办？

也有办法， 借助 swap ， 具体看代码：

```c++
int main () {
    vector<int> vec = {1,2,3,4,5,6};
    vec.clear();
    //1. 先声明一个和vec一样的临时变量（这个一样就只包含有效迭代器内的元素了）
    //2. 临时变量调用swap（）函数两者进行交换
    //3. 语句结束，临时变量自动析构
    vector<int>(vec).swap(vec);

    auto iter = vec.begin(); 

    cout<<"size="<<vec.size()<<endl;//0
	cout<<"capacity="<<vec.capacity()<<endl;//0
    for(int i = 0; i<6; ++i) {
        // coredump happens here
        std::cout<< *(iter+i) << "";
    }
}
```
输出结果：
```c++
size=0
capacity=0
Segmentation fault (core dumped)
```


### 写在最后

其实有些东西书里还是都提到过的，这就还是基础不牢的问题， 有时间还是要多去了解stl容器内存管理相关的知识。