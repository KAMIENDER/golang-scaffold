# golang-scaffold
集成了golang web server快速启动所需要的所有基础组建

# Demo
推荐使用VsCode

## VsCode
* 以conf/template.yaml文件为例子，创建自定义配置
* 配置环境变量，你可以配置在vscode的launch.json里面
```json
{
    "config_path": "your config file path"
}
```
* 点击 Launch Server，你会看到Demo服务已经在你本地的8080接口启动

# Start
* 创建自定义配置，参考template.yaml
* 配置配置文件的环境变量
```
{
    "config_path": "your config file path"
}
```
* 替换脚手架包名为你自己的包名
```
replace github.com/KAMIENDER/golang-scaffold with your pkg name
```
## DataBase
* 增加你的数据库表
* 修改infra/database/mysql/gen_tool/mysql_gen.go, 添加你的数据库对应的生成配置代码, 参考已有的user表的生成代码
* 执行infra/database/mysql/gen_tool/mysql_gen.go, 生成对应的golang结构体以及对应的CURD代码

## User Auth

最基础的用户表: ddl/user.sql
* 全局搜索 `Add user information` , 根据评论提示生成增加你的user信息以及对应操作

# Dependency
**Auth:**

AuthBoss: https://github.com/volatiletech/authboss

**ORM**

gorm: https://github.com/go-gorm/gorm

gen: https://gorm.io/zh_CN/gen/

**Web**

gin: https://github.com/gin-gonic/gin

**Pay**

gopay: https://github.com/go-pay/gopay

**Email**

gmail: gopkg.in/gomail.v2

**Log**

zap: https://github.com/uber-go/zap