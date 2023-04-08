[中文文档](https://github.com/KAMIENDER/golang-scaffold/blob/main/README_CN.md)
# golang-scaffold
This is a scaffold that integrates all the basic components needed for fast startup of a Golang web server.

# Demo
Recommend using VsCode.

## VsCode
* Create custom configuration based on conf/template.yaml.
* Configure environment variables, which can be configured in vscode's launch.json:
```json
{
    "config_path": "your config file path"
}
```
* Click "Launch Server" and you will see the Demo service start on your local 8080 port.

# Start
* Create custom configuration, refer to template.yaml.
* Configure environment variables for the configuration file:
```
{
    "config_path": "your config file path"
}
```
* Replace the scaffold package name with your own package name:
```
replace github.com/KAMIENDER/golang-scaffold with your pkg name
```
## DataBase
* Add your database table.
* Modify infra/database/mysql/gen_tool/mysql_gen.go and add the code for generating configuration corresponding to your database. Refer to the user table generation code that already exists.
* Execute infra/database/mysql/gen_tool/mysql_gen.go. It will connect to your database, generate corresponding Golang structures and CURD code based on your database table information.

## User Auth

The most basic user table: ddl/user.sql
* Global search "Add user information" and generate the code for adding your user information and corresponding operations based on the annotation prompts.
* After modifying user.sql, do not forget to execute the corresponding table creation SQL in your database, and then execute infra/database/mysql/gen_tool/mysql_gen.go to update your user structure.
* By default, /auth/* is used as the path for authentication.
* For interfaces that require user login verification, call the auth.WrapHandler method to encapsulate them before configuring gin.
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