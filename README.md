# chat

一个用于聊天的web应用

- 已实现功能
	- 用户
		- 注册用户
		- 以用户名和密码登录
		- 登出
		- 以cookies实现自动登录，并且七天后清除记录
	- 聊天
		- 登录后才能浏览主页、聊天群等页面
		- 可通过在输入文本框内输入文字并点击发送按钮
		- 可通过点击刷新按钮刷新并查看新消息
		- 在搜索框可通过名称和ID搜索群聊
	- 代码
		- 关闭时询问
		- 关闭时断开连接
		- 定时（每晚24时）生成日志文件

- 说明
	- 本项目依赖GO,MySQL
	- GO111MODULE=on
	- 本项目使用gin-gonic/gin为框架，go-sql-driver/mysql作为MySQL驱动
	- 本项目依赖 github.com/gin-gonic/gin , github.com/go-sql-driver/mysql , github.com/robfig/cron/v3 等多个第三方项目，请自行go mod
	  tidy
	- 本项目MySQL数据结构是项目根目录下的chat.sql
	- 数据库名为chat

- 使用
	- 将本项目拉取后，设置GO111MODULE=on
	- 安装github.com/gin-gonic/gin , github.com/go-sql-driver/mysql , github.com/robfig/cron/v3等包
	- 运行sql文件：mysql -u用户名 -p -D chat < chat.sql，构建数据库
	- 将conf.json更改为自己需要的配置
		- Port：服务的端口
		- ServerAddr：服务地址（127.0.0.1为本地，0.0.0.0为公开）
		- URL：用户浏览的地址（即浏览器输入的地址，如果错误会导致cookies无法使用）
		- SourceName：MySQL连接的地址
	- Windows下运行runserver.bat开启服务
	- Linux运行sudo sh runserver.sh开启服务
