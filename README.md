# LLOG

> Lightweight log agent.

## 说明
1. 超轻量级日志收集，过滤，上报工具。支持glob选取日志，收集日志上报至指定API或ES，后续支持kafka，redis。
2. 支持filebeat核心功能。
3. 相同运行环境，监控相同日志，比filebeat占用内存少50%以上。

测试机器：mac book pro 系统版本：10.14.5  配置：i9/1TB SSD/32G

|程序|监控文件数|占用内存|线程数|
|:--:|:--:|:--:|:--:|
|llog|4|5.8MB|23|
|llog|20|6.8MB|34|
|llog|50|7.1MB|36|
|filebeat|4|13.9MB|31|
|filebeat|20|16.2MB|37|
|filebeat|50|35.3MB|38|


4. 一键安装，无依赖。

## 安装

### 指定配置

#### 1. 下载对应版本LLA
``` shell script
# 下载linux 64 bit
wget http://qiniucdn.schoeu.com/lla_64bit

# 更改程序名
mv lla_64bit lla
```
或
``` shell script
# 下载linux 32 bit
wget http://qiniucdn.schoeu.com/lla_32bit

# 更改程序名
mv lla_32bit lla
```

#### 2. 创建配置文件，新建lla_conf.yml文件，内容如下

``` yaml
# 日志收集配置块
input:

# 存放各类日志文件的glob匹配路径
- log_path: ["/var/folders/lp/jd6nj9ws5r3br43_y7qw66zw0000gn/T/.nm_logs/*","/path/to/error/log/.log"]

  # 在输入中排除符合正则表达式列表的日志行
  #exclude_lines: ["test"]

  # 包含输入中符合正则表达式列表的日志行
  #include_lines: ["^\\w+"]

  # 忽略掉符合正则表达式列表的文件
  #exclude_files: ["\\d{4}.log"]

  # 默认为false, 从文件开始处重新发送所有内容。设置为true会从文件尾开始监控文件新增内容把新增的每一行文件进行发送
  tail_files: true

  #检测是否有新增日志文件的频率，默认为10秒
  #scan_frequency: 10

  # 最后一次读取文件后，持续时间内没有再写入日志，将关闭文件句柄，默认是 5分钟
  #close_inactive: 300
  
  # 发送自定义字段，默认会放在fields字段下, 当然也可以使用json字符串, 如  '{"a":"b"}'
  #fields: "some field here"  

  # 多行匹配
  #multiline:
    # 多行匹配点
    #pattern: "^error_log"
    # 最多匹配多少行，默认10
    #max_lines: 10

#- log_path: ["/other/log/path"]

  # 在输入中排除符合正则表达式列表的日志行
  #exclude_lines: ["test"]
  # ...

# 输出配置块:
output:

  # 把收集到的日志发送到指定API
  # 请求boby中带有JSON数据，以POST方法发送至指定接口
  #api_server:
    # 是否启用
    #enable: false
    #url: "http://127.0.0.1:9200/nma"

  elasticsearch:
    # 是否启用
    enable: false
    host: ["http://127.0.0.1:9200/"]
    index: "nma"
    # 输出认证.
    #username: "admin"
    #password: "s3cr3t"

# 通用配置块

# 应用名
#name: "llog"
# 是否上报系统级别日志（cpu，内存，磁盘，网络）, 默认为false，不上报
#sys_info: true

# 系统信息上报时间间隔，默认为10秒
#sys_info_during: 10

# 设置最大使用cpu数量, 默认无限制
#max_procs: 8

```

#### 3. 后台启动lla agent
``` shell script
# 默认配置启动
nohup ./lla >> lla_nohup.log 2>&1 &

# 指定配置文件启动
nohup ./lla ./lla_conf.yml >> lla_nohup.log 2>&1 &
```

## 上报数据格式
``` json
{
    "@logId": "cc621467-b53e-4e76-84b5-5679567c986f",
    "@message": "log content here...",
    "@timestamps": 1569751757188,
    "@name": "LLOG",
    "@version": "1.0.0",
    "@type": "normal|error|system"
    "@fields": "{\"key\":\"value\"}"
}

```

## 特性
- [x] 获取系统信息（cpu，内存，磁盘，网络）
- [x] 支持Glob语法批量指定日志
- [x] output支持ElasticSearch
- [x] output支持HTTP API
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
- [x] 自动关闭长期不活动文件句柄
- [x] 可限制cpu最多使用核数
- [x] 支持自定义字段，用于检索
- [x] 保存文件状态
- [ ] 支持socket数据传输
- [ ] 可设置日志上报线程数
- [ ] 支持多套独立配置
