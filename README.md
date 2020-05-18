## Zendea

![Screenshot](http://static.zendea.com/zendea.jpg)

`zendea`是一个使用Go语言开发的开源社区系统，采用前后端分离技术，Go语言提供api进行数据支撑，用户界面使用Nuxt.js进行渲染，后台界面基于element-ui。

## 功能特性

* 快速、简单
* 界面美观、渐进响应式布局
* 基于OAuth实现第三方帐号登录，目前支持Github/Gitee等
* 图片上传
* 自定义头像/文本
* 用户积分体系
* 普通用户/超级管理员角色划分
* 通知
* Markdown语法支持
* 标签
* 公告/小贴士
* RSS订阅
* 前后端完全分离
* 支持MySQL和Sqlite

## 模块

### backend

> 基于`Go`语言开发，提供Restful风格接口。

*技术栈*
- gin (https://github.com/gin-gonic/gin) Go web 框架
- JWT (https://github.com/appleboy/gin-jwt) JWT Middleware for Gin framework
- gorm (http://gorm.io/) Go 语言 orm 框架

### frontend

> 前端页面渲染服务，基于`nuxt.js`实现。

*技术栈*
- Nuxt.js (https://nuxtjs.org) 基于 Vue 的服务端渲染框架
- Element-UI (https://element.eleme.cn) 饿了么开源的基于 vue.js 的前端库
- Vditor (https://github.com/Vanessa219/vditor) Markdown 编辑器

## Demo
[Zendea](http://zendea.com/).

## 鸣谢
- Homeland (https://github.com/ruby-china/homeland)
- bbs-go (https://github.com/mlogclub/bbs-go)
- zeus-admin (https://github.com/bullteam/zeus-admin)

## License
Zendea is open-sourced software licensed under the [MIT license](http://opensource.org/licenses/MIT)
