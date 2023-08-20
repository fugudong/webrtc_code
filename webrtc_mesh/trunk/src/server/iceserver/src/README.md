# 目录说明
ice_server 	: ICE server程序
net_monitor : 实时网络监控

# ICE server编译安装

1. 先安装libevent
wget https://github.com/downloads/libevent/libevent/libevent-2.0.21-stable.tar.gz
tar xf libevent-2.0.21-stable.tar.gz
cd libevent-2.0.21-stable
./configure
make install

2. 然后安装turnserver
cd ice_server
tar xfz turnserver-4.5.1.1.tar.gz
cd turnserver-4.5.1.1
./configure 
make install