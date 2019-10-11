# LLOG

> Lightweight log agent.

## 说明
超轻量级日志收集，过滤，上报工具。支持glob选取日志，收集日志上报至指定API或ES，后续支持kafka，redis。
运行RSS只有8MB左右。
一键安装，无依赖。

## 安装

### 默认配置
```
wget -qO- http://qiniucdn.schoeu.com/install.sh | sh
```

### 指定配置

#### 1. 下载对应版本LLA
```
# 下载linux 64 bit
wget http://qiniucdn.schoeu.com/lla_64bit

# 更改程序名
mv lla_64bit lla
```
或
```
# 下载linux 32 bit
wget http://qiniucdn.schoeu.com/lla_32bit

# 更改程序名
mv lla_32bit lla
```

#### 2. 创建配置文件，新建lla_conf.yml文件，内容如下

```
# 输入配置:
# 是否上报系统级别日志（cpu，内存，磁盘，网络）, 默认为false，不上报
#sys_info: true
# 存放各类日志文件的glob匹配路径
#log_path: ["/var/folders/lp/jd6nj9ws5r3br43_y7qw66zw0000gn/T/.nm_logs/*","/path/to/error/log/.log"]
# 在输入中排除符合正则表达式列表的日志行
#exclude_lines: ["test"]
# 包含输入中符合正则表达式列表的日志行
#include_lines: ["^\\w+"]
# 忽略掉符合正则表达式列表的文件
#exclude_files: ["\\d{4}.log"]
# 默认为false, 从文件开始处重新发送所有内容。设置为true会从文件尾开始监控文件新增内容把新增的每一行文件进行发送
#tail_files: false
#检测是否有新增日志文件的频率，默认为10秒
scan_frequency: 10
# 多行匹配
#multiline:
  # 多行匹配点
  #pattern: "^error_log"
  # 最多匹配多少行，默认500
  #max_lines: 500

# 输出配置:
# 把收集到的日志发送到指定API
# 请求boby中带有JSON数据，以POST方法发送至指定接口
api_server: "http://127.0.0.1:9200/nma/logs"
#elasticsearch:
  #host: ["http://127.0.0.1:9200/nma"]
  # 输出认证.
  #username: "admin"
  #password: "s3cr3t"
  # elasticsearch请求超时事件。默认90秒.
  #timeout: 90

# redis配置
#redis:
#  enabled: true
#  hosts: ["192.168.10.188"]
#  port: 6379
#  datatype: list
#  key: "llog"
#  db: 0

```

#### 3. 后台启动lla agent
```
# 默认配置启动
nohup ./lla >> lla_nohup.log 2>&1 &

# 指定配置文件启动
nohup ./lla ./config.json >> lla_nohup.log 2>&1 &
```

## 上报数据格式
```json
{
"@logId": "cc621467-b53e-4e76-84b5-5679567c986f",
"@message": "log content here...",
"@sysInfo": "{\"dataTime\":\"2019-09-29T18:09:17\",\"logicalCores\":16...}",
"@timestamps": 1569751757188,
"@type": "LLOG",
"@version": "1.0.0"
}

```

## TODO
- [x] 获取系统信息（cpu，内存，磁盘，网络）
- [x] 支持Glob语法批量指定日志
- [x] output支持ElasticSearch
- [x] 在输入中排除符合正则表达式列表的日志行
- [x] 包含输入中符合正则表达式列表的日志行
- [x] 忽略掉符合正则表达式列表的文件
- [x] 一次日志事件中最多上传多少个字符
- [x] 更换配置文件为yaml
- [x] API, ES请求Timeout设置
- [x] 多行日志匹配，一般用于错误堆栈信息收集
- [x] 多行日志匹配限制行上限
- [x] 可配置从日志文件起始或尾部进行日志监听
- [x] 新增文件检测
- [ ] 文件状态保存
- [ ] output支持redis
- [ ] output支持kafka