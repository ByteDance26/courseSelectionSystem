# courseSelectionSystem
选课大作业
# 功能
成员管理的功能实现
# 内容
## 2022.02.15更新
- types.go的成员管理部分的Request结构体都加入了`binding:"required"`的tag
- modules/members.go中改了IsMemberOK函数，增加了_type.Member类型的返回值

用法1：
```go
// 判断用户是否存在or已删除
switch Code, _ = modules.IsMemberOK(UserID); Code {
case _type.UserNotExisted:
// 用户不存在
Response.Code = _type.UserNotExisted
case _type.UserHasDeleted:
// 用户已删除
Response.Code = _type.UserHasDeleted
default:
// 用户存在
}
```
用法2：
```go
Code, MemberInfo := modules.IsMemberOK(UserID)
if Code == _type.OK {
UserType := MemberInfo.UserType
Password := MemberInfo.Password
}
```
## 2022.02.14更新
数据库初始化加入了AutoMigrate部分:
- DB/mysql

另外修改了moudles文件夹的名字为modules

新增三个代码文件：
- controller/members.go
- modules/members.go
- type/types_members.go

tmp和.idea都是Goland自动生成的可执行文件可以删掉

# 边界case
- POST请求的某字段为空/缺少某字段：
    - binding不成功，返回Code = ParamInvalid