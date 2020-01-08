---
layout: post
title:  "因string使用不当引起的线上服务崩溃"
date:   2019-03-23 21:03:36 +0530
categories: std::string
---
直接使用赋值号将字符数组赋给了string, 而没有使用string.assign(char_array,sizeof(char_array)),以及 string(char_array, sizeof(char_array)) 引起的一场悲剧...... 

```c++

/* bad example */
std::string CreateRemotePeerId(){
    char remote_sha[32];
    std::string remote_peerid;

    //do sth to format remote_sha

    remote_peerid = remote_sha; 
    // 就是这里,悲剧发生了, 因为remote_sha中可能有/0字符, 直接赋值就会导致string被截断.

    return remote_peerid;  

}

/* good example */
std::string CreateRemotePeerId(){
    char remote_sha[32];
    std::string remote_peerid;

    //do sth to format remote_sha

    remote_peerid.assign(remote_sha, sizeof(remote_sha)); 
    //or   remote_peerid = std::string(remote_sha, sizeof(remote_sha));
    return remote_peerid;

}

```

上学的时候看C++各种书籍, 说过字符数组与string之间转换要注意处理方式的不在少数, 没想到最终还是踩了这个坑. 这次错误的赋值直接导致了错误的删除逻辑, 在积累了一周后线上服务大面积告警, 算是一次事故了, 难受呀, 以后这块不会再出错了.
