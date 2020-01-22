---
layout: post
title:  "因string使用不当引起的线上服务崩溃"
date:   2020-01-06 22:03:36 +0530
tags:   std::string char[]
categories: Coding
---
直接使用赋值号将字符数组赋给了string, 而没有使用string.assign(char_array,sizeof(char_array)),以及 string(char_array, sizeof(char_array)) 引起的一场悲剧:persevere:...... 

```c++

/* bad example */
std::string CreateRemotePeerId(){
    char remote_sha[32];
    std::string remote_peerid;

    //do sth to format remote_sha

    remote_peerid = remote_sha; 
    // 就是这里,悲剧发生了, 因为remote_sha中可能有'\0'字符, 直接赋值就会导致string被截断.

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

/* intuitive example */
int main () {
    char src[5] = {'1','\0', '2','3','4'};
    std::string bad_dst = src;
    std::string good_dst1 = std::string(src, sizeof(src));
    std::string good_dst2;
    good_dst2.assign(src, sizof(src));

    std::cout<<bad_dst<<std::endl;
    std::cout<<good_dst1<<std::endl;
    std::cout<<good_dst2<<std::endl;
    
    //out put
    //1
    //1234
    //1234
}

```

<font face="微软雅黑" >上学的时候看C++各种书籍, 强调过字符数组与string之间转换要注意处理方式的数不胜数, 没想到最终还是踩了这个坑. 这次错误的赋值直接导致了错误的删除逻辑, 在积累了一周后线上服务大面积告警, 算是一次事故了, 难受。 大佬们说不要有心理负担，我觉得负担还是要有，但要让这些负担转换为自己的储备。希望自己可以胆大心细加油干~ :muscle:</font>
