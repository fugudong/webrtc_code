房间服务器，客户端加入同一房间实现通话
同一房间通话最大人数可配置。

# 功能

# 编译
go build -o bin src/main.go



# 运行(具体参照《音视频通话-房间服务器详细设计》配置章节)
1. ICE服务器
```
./turnserver -L 172.16.0.3 -a -u lqf:123456 -v -f -r nort.gov
```
2. 房间服务器
```
nohup ./roomserver &
```
