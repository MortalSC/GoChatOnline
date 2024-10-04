# 项目开发文档说明



## 开发环境说明

- chat online v1 版本：windows vscode
- chat online v2 版本：windows wsl + vscode / Linux（Ubuntu）



## 版本说明

- chat online v1 版本是：参考于刘冰单Aecid的b站网课所写。

- chat online v2 版本是：在保持v1版本功能逻辑的基础上，对代码的更迭，已到达设计的合理性，尽可能解耦合，引入接口等写法调整代码逻辑。



## 脚本说明

>  脚本支持一键编译出客户端和服务端。

- chat online v1 版本提供了三个脚本文件：
  - build.ps1：PowerShell
  - buildforwin.bat：windows 操作系统下的编译脚本
  - buildforlinux.bat：Linux操作系统下的编译脚本