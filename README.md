# gin-example
* gin框架开发示例

# 介绍
* 封装底层Dao操作类，支持主从数据库操作
* 不包含第三方orm，使用[gendry](https://github.com/didi/gendry)辅助操作数据库
* 使用redis封装session管理
* error统一使用exception管理
* 支持swagger文档
* 使用[endless](https://github.com/fvbock/endless)优雅重启