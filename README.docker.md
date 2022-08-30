# docker分支 #
#### 一个专门用于构建docker镜像的分支 #
此版本是为了编译docker镜像而专门制作的
##### 特点 #
- 因为docker镜像无法更改，所以将配置由配置文件改为了由环境变量决定，具体环境变量有：
	- ServerPort：同conf.json中的Port
	- ServerAddress：同conf.json中的ServerAddr
	- URL：同conf.json中的URL
	- SourceName：同conf.json中的SourceName
		> 具体见README.md
	镜像运行命令:docker run --name chat -p 80:8080 -e URL=192.168.200.128 -e ServerPort=8080 -e ServerAddress=0.0.0.0 -e SourceName='root:123456@tcp(192.168.200.1:3306)/chat?parseTime=true' chat:v1

