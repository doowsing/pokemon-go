# pokemon-go

> 本项目为纪念以前玩的一个2012年停服的游戏所仿，游戏分为h5前端与服务器后端两个部分，本项目为服务器后端部分，采用Go语言开发
> 目前实现的功能包括原游戏的绝大部分功能，包括战斗系统、道具系统、宠物系统、任务系统、NPC系统、聊天系统
> 未完成部分-NPC情人岛、多人战斗神圣战场、多人战斗家族战场

## 技术组成
- 游戏前端界面基于vue完成
- 后端服务采用Go语言开发，使用数据库包括mysql、redis
- 主程序部分使用Gin框架开发，开放http接口给前端
- 聊天与组队系统使用websocket库完成
- 服务间使用rpc-json、http进行远程过程调用与交互

## 后端项目结构
- 后端分为主控制服务、主程序服务、聊天服务、定时任务服务、组队服务
- 主控制服务主要是监控各个服务的状态
- 主程序服务完成大部分的游戏功能，并开放http接口给前端
- 聊天服务提供玩家聊天与公告信息的发布
- 定时任务负责服务器功能性程序调用
- 组队服务负责玩家的组队操作

## 游戏体验
目前可在[服务器](http://139.199.181.53:100/)体验已提交的程序
