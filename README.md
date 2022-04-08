# 功能

## admin
* 插入模板
* 查找模板
* 查找记录
* 登录后可操作（v2待开发）

## email
* 从kafka读取并发送email
* 插入发送记录

# 数据库

## 邮件模板

| 字段 | 类型 | 说明 |
| --- | :---: | --- |
| id | string | mongo主键 |
| name | string | 模板名称 |
| subject | string | 邮件主题 |
| content | string | 内容 |

## 邮件发送记录

| 字段 | 类型 | 说明 | 
| --- | :---: | --- | 
| id | string | 主键 | 
| send_time | int64 | 发送时间（时间戳） |
| receiver | string | 接受者邮箱地址 |
| is_successful | int | 是否成功，0失败，1成功 |
| template_id | string | 模板id |
| name | string | 邮件模板名称 |
| content | string | 邮件内容,包括邮件主题 |
