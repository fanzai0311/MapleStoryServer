MapleStoryServer
==========
A MapleStory Server in Go For CMS V079
Go 编写的079版本冒险岛服务端

Target
------

实现一个079版本的冒险岛服务器.
能够实现原Java的部分功能.

Files
-----

- `main.go`    —— 程序入口，负责初始化CPU分配、打开控制台、准备消息分发routines，侦听端口，并且接收客户端的连接。

Packages
--------

- `connection` —— 用于建立与客户端的连接，进行读写数据。
- `message`    —— 用于负责服务器向客户端方向的信息的递送，主要功能是广播和点对点通信；
- `types`      —— 包含整个服务器所用到的各种数据结构，比如用户信息结构体等；
- `users`      —— 包含用户管理相关函数，比如用户进入游戏与离开。
- `console`    —— 给服务端添加一个控制台，便于输入指令调试。
