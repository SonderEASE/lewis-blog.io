---
layout: post
title:  "谁都想不到的open files（文件句柄数）"
date:   2020-01-13 13:40:36 +0530
tags: ev::timer pstack open-files
categories: 事记|Logging&#8195;Event
---
看着要过年了， 线上的有几台机器突然出现卡死的现象:scream:，比如原来1分钟的uv::timer，莫名奇妙地变成6分钟，各种诡异事件，最后竟是因为....

最初遇到的问题就是 **ev::timer失效**， 围绕ev一通操作， 查看代码逻辑，增添日志，有同事之前也遇到相同的问题，解决的方法是把定时器周期调长，但只能缓解问题，并没有完全解决，而且这个方法在这里并不适用，
因为这个定时器涉及到关键的上报信息，所以周期没办法调长。

最终通过 **pstack** 发现不止这个定时器，在问题复现时，很多基本的sock句柄都被卡死了， 这才定位到是因为这批新机器在交接的时候没有设置文件句柄数，系统采用的默认文件句柄数是**1024**， ev，uv这些库都是 ***事件驱动** 的，对于我们的程序来说1024个句柄是远远不够的。
这也解释了为什么之前有的同事增加了定时器的周期就有所缓解， 相当于减少了单位事件内的事件数量。

```c++

ulimit -a  //列出所有机器的参数选项

core file size          (blocks, -c) 0
data seg size           (kbytes, -d) unlimited
scheduling priority             (-e) 0
file size               (blocks, -f) unlimited
pending signals                 (-i) 514316
max locked memory       (kbytes, -l) 64
max memory size         (kbytes, -m) unlimited
open files                      (-n) 1024         # 单个进程可用的最大文件句柄数（系统默认1024）
pipe size            (512 bytes, -p) 8
POSIX message queues     (bytes, -q) 819200
real-time priority              (-r) 0
stack size              (kbytes, -s) 10240
cpu time               (seconds, -t) unlimited
max user processes              (-u) 10240
virtual memory          (kbytes, -v) unlimited
file locks                      (-x) unlimited

```

:bulb: **问题发生在个别机器上， 回退稳定版本问题仍然出现。 控制变量，这种时候就不应该只局限在代码上找问题了。**
