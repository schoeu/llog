# LLOG

> Lightweight log agent.

[中文文档](./README_zh.md)

## Instructions
1. Lightwelterweight log collection, filtering, reporting tools.Support glob selected log, collect log report to specify the API or ES, follow-up support kafka, redis.
2. Support filebeat core functionality.
3. The same operating environment, monitoring the same logs, more than 50% less than filebeat takes up memory.

Testing machine：macbookpro. system version：10.14.5（i9/1TB SSD/32G）

|name|monitoring file|memory|thread|
|:--:|:--:|:--:|:--:|
|llog|4|6.1MB|23|
|llog|20|6.8MB|34|
|llog|50|7.1MB|36|
|filebeat|4|13.9MB|31|
|filebeat|20|16.2MB|37|
|filebeat|50|35.3MB|38|


4. A key to install, no dependence.

## Install

### Specified configuration

#### 1. Download the corresponding version
``` shell script
# download linux 64 bit
wget http://qiniucdn.schoeu.com/llog_64bit

```
Or
``` shell script
# download linux 32 bit
wget http://qiniucdn.schoeu.com/llog_32bit

```

#### 2. Create a configuration file, new llog_conf. Yml file, the content is as follows

``` yaml
# log collection configuration block
input:

# to store all kinds of log file glob matching path
- log_path: ["/var/folders/lp/jd6nj9ws5r3br43_y7qw66zw0000gn/T/.nm_logs/nm_apps?/*.log"]
  # in the input to exclude conform to the regular expression list of log line
  #exclude_lines: ["test"]

  # include conform to the regular expression in the input list log line
  #include_lines: ["^\\w+"]

  # ignore conform to the regular expression list file
  #exclude_files: ["\\d{4}.log"]

  # the default is false, beginning to send all the content from a file.Set to true will from the tail to start monitoring file additions send new files on each line
  tail_files: true

  #test whether have increased frequency of log files, the default for 10 seconds
  scan_frequency: 10

  # for the last time, after reading the file last time didn't log, will close the file handle, the default is 5 minutes
  close_inactive: 300

  # to send custom fields, the default will be under the fields fields, it can also use a json string, such as' {" a ":" b "} '
  #fields: "some field here"

  # multi-line matching
  #multiline:
    # multi-line matching points
    #pattern: "^normal_log"
    # up to match how many rows, 10 by default
    #max_lines: 10

- log_path: ["/var/folders/lp/jd6nj9ws5r3br43_y7qw66zw0000gn/T/.nm/*.log"]
  # multi-line matching
  multiline:
    # multi-line matching points
    pattern: "^error_log"
    # up to match how many rows, 10 by default
    max_lines: 5
  scan_frequency: 160
  close_inactive: 30

# output configuration block:
output:

  # the collected log is sent to a designated API
  # request with the JSON data in the boby, sending by POST method to specify the interface
  #api_server:
  # whether to enable
  #enable: false
  #url: "http://127.0.0.1:9200/nma"

  elasticsearch:
    # whether to enable
    enable: true
    host: ["http://127.0.0.1:9200/"]
    index: "nma"
    # output certification.
    #username: "admin"
    #password: "s3cr3t"

# general configuration block

# application name
#name: "llog"
# if system level log (CPU, memory, disk, network), the default is false, is not reported
#sys_info: true

# system information reporting time interval, the default for 10 seconds
#sys_info_during: 10

# set the maximum use of CPU number, unrestricted by default
#max_procs: 8

# file status to keep configuration
#snapshot:
  # file status switch, default is not open
  #enable: false

  # save document status, a snapshot of the current state to a local, a kick-off meeting for next time use snapshot content
  #snapshot_dir: '/path/to/snapshot/file'

  # save the file regularly, defaults to 5 seconds
  #snapshot_during: 5

```

#### 3. Start llog in background
``` shell script
# the default configuration
nohup ./llog_64bit >> llog_nohup.log 2>&1 &

# specified configuration file
nohup ./llog_64bit ./llog_conf.yml >> llog_nohup.log 2>&1 &
```

## report data format
``` json
{
    "@logId": "cc621467-b53e-4e76-84b5-5679567c986f",
    "@message": "log content here...",
    "@timestamps": 1569751757188,
    "@name": "LLOG",
    "@version": "1.0.0",
    "@type": "normal|error|system",
    "@fields": "{\"key\":\"value\"}"
}

```

## Features
- [x] get information system (CPU, memory, disk, network)
- [x] support batch designated log Glob grammar
- [x] The output support ElasticSearch
- [x] The output support HTTP API
- [x] in input to exclude the regular expression list of log line
- [x] a list in line with the regular expression in the input of the log line
- [x] ignore the regular expression list file
- [x] upload at most how many characters in a log event
- [x] replacement for yaml configuration file
- [x] API, ES request Timeout Settings
- [x] log multi-line matching, commonly used in error stack information collection
- [x] log multi-line matching maximum limit line
- [x] can be configured from the log file starting or tail log monitor
- [x] add file test
- [x] automatically shut down long-term inactive file handle
- [x] can limit the CPU use auditing at most
- [x] support custom fields, is used to retrieve
- [x] Save the file status
- [x] Support for multiple sets of independent configuration 
- [ ] can be set up log reports the number of threads
