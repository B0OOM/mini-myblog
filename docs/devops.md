# 将自己的代码发布到线上，支持MySQL，更专业的服务发布上线

当项目写完之后，大家都很期待能发布到线上，那么应该如何发布呢？这里将介绍如何将自己的代码发布到线上，支持MySQL，更专业的服务发布上线。


## 支持MySQL数据库
得益于gorm，这个博客可以支持三种数据库：sqlite3、mysql、postgres。这里我们将介绍如何支持MySQL数据库。

代码已经支持MySQL数据库，只需要启动的时候指定数据库类型即可。

```bash
go run . -dbdriver=mysql -dbfile="root:123456@tcp(127.0.0.1)/miniblog?charset=utf8mb4&parseTime=True&loc=Local"
```

也可以看一下models.go中的代码：
```go
func ConnectDatabase(dbfile, dbdriver string) (*gorm.DB, error) {
	var dialector gorm.Dialector
	if dbdriver == "" || dbdriver == "sqlite" {
		dialector = sqlite.Open(dbfile)
	} else if dbdriver == "mysql" {
		dialector = mysql.Open(dbfile)
	}

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}
}
```
根据这个代码可以看到，如果`dbdriver`是`mysql`，那么就会使用`mysql.Open`来连接数据库

## 如何部署到自己的服务器

一般的云服务器都是用Linux系统，所以我们可以使用Linux系统来部署我们的代码。
1. 编译代码：
   ```shell
    cd cmd
    go build -o miniblog
    ```
2. 上传到服务器：
    ```shell
    cd cmd
    ssh ubuntu@your-server-ip mkdir -p ~/miniblog
    scp miniblog ubuntu@your-server-ip:~/miniblog/
    scp -r static ubuntu@your-server-ip:~/miniblog/
    scp -r templates ubuntu@your-server-ip:~/miniblog/
    ```
3. 运行代码：
    ```shell
    ssh ubuntu@your-server-ip
    cd ~/miniblog
    nohup ./miniblog > miniblog.log 2>&1 &
    ```

这个的方案就是简单运行在8080端口上，如果你有nginx， 那么将nginx配置文件配置一下，就可以通过域名访问了。

如果要实现更专业的服务发布上线，可以使用docker，k8s等技服，这个需要更多的学习和实践。

`nohup` 只是一个简单的服务工具，如果要实现比较稳定，那么就建议用`systemd`来管理服务。
可以写一个简单的`miniblog.service`文件，然后放到`/etc/systemd/system`目录下，然后启动服务。

```shell
[Unit]
Description=Miniblog Service
Description = frp server
After = network.target syslog.target
Wants = network.target

[Service]
Type=simple
WorkingDirectory=/home/ubuntu/miniblog
ExecStart=/home/ubuntu/miniblog/miniblog -dbdriver mysql -dbfile "root:123456@tcp(...)"

[Install]
WantedBy = multi-user.target
```

然后启动服务：
```shell
sudo systemctl start miniblog
```

这样就可以通过`systemd`来管理服务, 也可以通过`systemctl`来查看服务的状态，重启服务等。

课程内容学习遇到问题，可以联系我们老师进行沟通（微信号：jinti2000），我们会及时更新课程内容。

<div style="display: inline-block;text-align: center;">
   <div style="display: inline-block;">
     <h3>入职啦实战项目</h3>
     <img src="../cmd/static/projectQrcode.jpg" width="200" margin-right="100" alt="入职啦实战项目二维码" >
   </div>
   <div style="display: inline-block; margin-left: 30px;">
     <h3>入职啦微信公众号</h3>
     <img src="../cmd/static/weixinQrcode.png" width="200" alt="入职啦公众号二维码" />
   </div>
 </div>

[课程介绍](./README.md)