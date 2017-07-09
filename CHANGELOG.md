# Changelog

## 0.6.0 (2017-07-09)

* 增加说明在某些情况下 ObjectPutHeaderOptions.ContentLength 必须要指定
* 增加 ObjectUploadPartOptions.ContentLength


## 0.5.0 (2017-06-28)

* 修复 ACL 相关 API 突然失效的问题.
  (因为 COS ACL 相关 API 的 request 和 response xml body 的结构发生了变化)
* 删除调试用的 DebugRequestTransport(把它移动到 examples/ 中)


## 0.4.0 (2017-06-24)

* 去掉 API 中的 authTime 参数，默认不再自动添加 Authorization header
  改为通过自定义 client 的方式来添加认证信息
* 增加 AuthorizationTransport 辅助添加认证信息

## 0.3.0 (2017-06-23)

* 完成所有 API


## 0.2.0 (2017-06-10)

* 调用 bucket 相关 API 时不再需要 bucket 参数, 把参数移到 service 中
* 把参数 signStartTime, signEndTime, keyStartTime, keyEndTime 合并为 authTime


## 0.1.0 (2017-06-10)

* 完成 Service API
* 完成大部分 Bucket API(还剩一个 Put Bucket Lifecycle)
