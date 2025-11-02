# 协议

# RESP协议
## RESP（Redis Serialization Protocol）是Redis使用的二进制协议，用于客户端和服务器之间的通信。
## 它简单、快速、可靠，被广泛用于Redis的客户端库和服务器实现。

## RESP协议详细说明： RESP协议基于文本协议，使用换行符（\n）作为分隔符。每个请求和响应都是一个简单的字符串，每个部分之间用换行符分隔。

### 数据类型
### RESP协议支持以下数据类型：
### 简单字符串（Simple String）
    前缀："+"
    格式："+<字符串>\r\n"
    示例："+OK\r\n"
### 错误（Error）
    前缀："-"
    格式："-<错误消息>\r\n"
    示例："-ERR unknown command 'foo'\r\n"
### 整数（Integer）
    前缀：":"
    格式：":<整数>\r\n"
    示例：":100\r\n"
### 批量字符串（Bulk String）
    前缀："$"
    格式："$<字符串长度>\r\n<字符串数据>\r\n"
    示例："$5\r\nhello\r\n"
### 数组（Array）
    前缀："*"
    格式："*<数组长度>\r\n<数组元素1>\r\n<数组元素2>\r\n..."
    示例："*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n"

### redis 命令格式
### 每个命令由一个或多个部分组成，每个部分都是一个简单的字符串。命令的第一个部分是命令名称，后面是命令的参数。
### 例如，一个简单的GET请求可能如下所示：
### set key1 value1
    *3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n

