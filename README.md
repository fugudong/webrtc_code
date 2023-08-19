# webrtc_mesh

WebRTC mesh模型

# 支持的功能
1.  多方通话，通过room_server的配置文件config.ini进行配置，由于mesh模型上行带宽压力大，一般建议小于6个人。
2.  支持分布式部署coturn服务器
3.  支持不同平台的通话

# 支持的客户端
1.  Web端
2.  Android原生APP
3.  Windows、IOS待开发

# 组件服务
1. loginserver 登录服务器，客户端需要使用appid等信息登录loginserver获取roomserver和iceserver的地址
   目前没有配置
2. roomserver房间服务器，客户端需要加入房间实现通话
3. webserver 网页服务器，客户端操作网页

