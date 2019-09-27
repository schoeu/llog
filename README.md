# nma
Log agent for node monitor.

## 说明
本程序为Node监控模块`nm`的代理模块，主要作用为收集各Node实例（nm）的日志，并获取机器系统信息合并上报。

## 安装

### 默认配置
```
wget -qO- http://qiniucdn.schoeu.com/install.sh | sh
```

### 指定配置

#### 1. 下载对应版本nma agent
```
# 下载linux 64 bit
wget http://qiniucdn.schoeu.com/nma_64bit

# 更改程序名
mv nma_64bit nma
```
或
```
# 下载linux 32 bit
wget http://qiniucdn.schoeu.com/nma_32bit

# 更改程序名
mv nma_32bit nma
```

#### 2. 创建配置文件，新建nma_conf.json文件，内容如下

```
{
  "logDir": ["/path/to/normal/log/*.log"],
  "noSysInfo": false,
  "logServer": "http://your_log_server_host",
  "errLogs": ["/path/to/error/log/*.log", "/other/path/to/error/log/*.log"]
}
```

#### 3. 后台启动nma agent
```
# 默认配置启动
nohup ./nma >> nma_nohup.log 2>&1 &

# 指定配置文件启动
nohup ./nma ./config.json >> nma_nohup.log 2>&1 &
```

#### 配置说明

|配置名|示例|说明|默认值|
|--|--|--|--|
|logDir|["/path/to/normal/log/*.log","/path/to/error/log/*.log"]|存放各类日志文件的glob匹配路径|"$tmp/.nm_logs/*"|
|noSysInfo|false|是否上报系统级别日志（cpu，内存，磁盘，网络）|false|
|logServer|http://your_log_server_host|日志上报接口，会以POST方式上报json数据|-|
|exclude|["^\w+"]|在输入中排除符合正则表达式列表的那些行|-|
|include|["^\w+"]|包含输入中符合正则表达式列表的那些行|所有行|
|excludeFiles|["\d{3}.log"]|忽略掉符合正则表达式列表的文件|-|
