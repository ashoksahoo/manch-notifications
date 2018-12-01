Test the Binary
------------------

creates binary in project directory & installs to `$PATH`
```shell
go build 
go install 
```
make sure working directory is project dir as GCM config file is inside the project.

``` shell
env=production REDIS_HOST="manch.jcrvar.0001.aps1.cache.amazonaws.com MONGO_DB=mongodb://manch_user:<password>@cluster0-shard-00-00-ngxmc.mongodb.net:27017,cluster0-shard-00-01-ngxmc.mongodb.net:27017,cluster0-shard-00-02-ngxmc.mongodb.net:27017/manch?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin" NATS=nats://172.31.54.226:4222,nats://172.31.32.102:4222 notification-service
```    

To Install Service
----------------
Create Service file
``` shell
sudo nano /etc/systemd/system/manch-notification.service
```

Copy to systemd service file

```text
[Unit]
Description=Manch Notification Service

[Service]
Environment="env=production"
Environment="REDIS_HOST=manch.jcrvar.0001.aps1.cache.amazonaws.com"
Environment="MONGO_DB=mongodb://manch_user:<password>@cluster0-shard-00-00-ngxmc.mongodb.net:27017,cluster0-shard-00-01-ngxmc.mongodb.net:27017,cluster0-shard-00-02-ngxmc.mongodb.net:27017/manch?ssl=true&replicaSet=Cluster0-shard-0&authSource=admin"
Environment="NATS=nats://172.31.54.226:4222,nats://172.31.32.102:4222"
WorkingDirectory=/home/ubuntu/manch/manch-notifications
ExecStart=/home/ubuntu/manch/manch-notifications/notification-service
Restart=always

[Install]
WantedBy=multi-user.target
```
Enable and start the service

``` shell
sudo systemctl start manch-notification.service 
sudo systemctl enable manch-notification.service 
sudo systemctl status manch-notification.service 
```
