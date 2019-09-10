# nma
Log agent for node monitor.

## 说明
本程序为Node监控模块`nm`的代理模块，主要作用为收集各Node实例（nm）的日志，并获取机器系统信息合并上报。

## 安装

### 默认配置
```
wget -qO- https://github.com/schoeu/nma/blob/master/scripts/install.sh?raw=true | sh
```

### 指定配置

#### 1. 下载对应版本nma agent
```
# 下载linux 64 bit
wget https://github.com/schoeu/nma/blob/master/bin/nma_64bit?raw=true

# 更改程序名
mv nma_64bit nma
```
或
```
# 下载linux 32 bit
wget https://github.com/schoeu/nma/blob/master/bin/nma_32bit?raw=true

# 更改程序名
mv nma_32bit nma
```

#### 2. 创建配置文件，新建nma_conf.json文件，内容如下

```
{
  "logDir": "/Users/schoeu/Downloads/git/nm/test/nm_logs",
  "noSysInfo": false,
  "logServer": "http://your_log_server_host"
}
```

#### 3. 后台启动nma agent
```
nohup ./nma >> nma_nohup.log 2>&1 &
```

#### 配置说明

|配置名|示例|说明|默认值|
|--|--|--|--|
|logDir|"/Users/schoeu/Downloads/git/nm/test/nm_logs"|存放Node实例日志文件夹|"$home/.nm_logs/"|
|noSysInfo|false|是否上报系统级别日志（cpu，内存，磁盘，网络）|false|
|logServer|http://your_log_server_host|日志上报接口，会以POST方式上报json数据|-|

