# iot_server_device
此项目包含三个服务，组合实现了前端通过阿里云iot交互设备的后端系统。

## iot_server4
本身是一个websocket服务器，转发前端的命令然后用阿里云接口发给设备。再将阿里云返回的设备异步响应推送给客户端。

- 多个客户端同时发送消息给一个设备

因为硬件设备处理并发消息的能力有限（每220毫秒只能处理一个消息），因此需要把客户端发给同一设备的消息进行处理。220毫秒内的重复消息只取第一个，其余的丢弃并通知给其客户端。
维护全局缓存model/device，用读写锁包裹，其中每个设备记录发送消息的时间戳，以此判断延迟。
并且阿里云接口有qps限制（500）也要满足。
这里我用限流消息队列（在cellnet.queue的基础上用令牌桶算法限流）拼装了一个吸管Straw架构，model/device的每个设备缓存都有一个Straw作为向该设备发送消息的唯一通道，以此控制消息延迟。所有的请求设备的Straw会将消息通过限流消息队列发送给阿里云以对接其qps500的限制。
Straw=消息队列+afterfunc回调。

- 响应设备消息给多个订阅客户端

基于cellnet，保存sesid从sessionmanager中查找ses并广播。

- 设备和客户端的交互控制

主要是生命周期，在线状态，使用model/session、client的众多标记位实现。

- 客户端、设备心跳检查+实时轮询

使用cellnet的loop实现。

## iot_server6
使用阿里云iot的amqp客户端订阅指定消费组的设备响应，用pulsar转发给iot_server4。
可以对阿里云iot的异步设备消息进行过滤，比如非实时历史消息。

## iot_server8
基于gin的http服务器，用于登录验证和生成基于jwt的token并返回。

## 对应测试前端
https://github.com/GramYang/iot_server_device_client

