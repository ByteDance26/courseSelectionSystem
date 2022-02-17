package controller

import (
	"courseSelectionSystem/DB"
	"courseSelectionSystem/middle"
	"courseSelectionSystem/modules"
	_type "courseSelectionSystem/type"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"strings"
)

// CreateMember POST创建成员
// TODO::检查当前登录用户是否是管理员 LoginRequired
// TODO::内置管理员账号
func CreateMember(c *gin.Context) {
	var Response _type.CreateMemberResponse
	var Request _type.CreateMemberRequest
	// 获取JSON参数
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	// 将参数绑定到Member里
	NewMember := _type.Member{
		Nickname: Request.Nickname,
		Username: Request.Username,
		Password: Request.Password,
		UserType: Request.UserType,
		Status:   _type.Existed,
	}

	//判断是否登录且是管理员
	s, err := middle.GetSimpleSession(c)
	if err == _type.SessionError {
		Response.Code = _type.UnknownError
		c.JSON(http.StatusOK, Response)
		return
	} else {
		//判断登录
		if s.Value["userId"] == nil {
			Response.Code = _type.LoginRequired
			c.JSON(http.StatusOK, Response)
			return
		}
		//判断是管理员
		var mem modules.Member
		idNum, _ := strconv.Atoi(s.Value["userId"].(string))
		err := mem.GetMemberByUserId(idNum)
		if err == gorm.ErrRecordNotFound || mem.UserType != int(_type.Admin) {
			Response.Code = _type.PermDenied
			c.JSON(http.StatusOK, Response)
			return
		}
	}
	// 判断参数是否合法
	if modules.IsMemberParamValid(NewMember) {
		// 参数合法（Username检查在后面)
		Response.Code = _type.OK
	} else {
		// 参数不合法
		Response.Code = _type.ParamInvalid
		Response.Data.UserID = ""
	}

	// 添加记录
	if Response.Code == _type.OK {
		db := DB.MysqlDB
		result := db.Create(&NewMember)
		// 判断用户名是否重复
		if result.Error == nil {
			// Username 没有重复，成功创建
			Response.Data.UserID = strconv.Itoa(int(NewMember.UserID))
			//updateRedis
			if NewMember.UserType == _type.Student {
				DB.CreateStudent(strconv.Itoa(int(NewMember.UserID)))
			}
		} else if Error := fmt.Sprintf("%v", result.Error); strings.Contains(Error, "Error 1062") {
			// Username 重复
			Response.Code = _type.UserHasExisted
			Response.Data.UserID = ""
		} else {
			// TODO 未知错误 还没有遇到这种情况
			//fmt.Println(Error)
			Response.Code = _type.UnknownError
		}
	}
	c.JSON(http.StatusOK, Response)
}

// GetMember GET获取成员信息
func GetMember(c *gin.Context) {
	var Response _type.GetMemberResponse
	var Request _type.GetMemberRequest
	// 获取JSON参数
	if err := c.ShouldBindQuery(&Request); err != nil {
		fmt.Println(err)
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	// 全局变量DB赋值
	db := DB.MysqlDB
	// 查询UserID对应的记录
	var MemberInfo _type.Member
	result := db.Find(&MemberInfo, Request.UserID)
	// 传入TMember类型参数
	ResponseMemberInfo := _type.TMember{
		UserID:   strconv.Itoa(int(MemberInfo.UserID)),
		Nickname: MemberInfo.Nickname,
		Username: MemberInfo.Username,
		UserType: MemberInfo.UserType,
	}
	// 检查错误
	if result.RowsAffected == 0 {
		// 用户不存在：检查 ErrRecordNotFound 错误
		Response.Code = _type.UserNotExisted
		ResponseMemberInfo.UserID = ""
		ResponseMemberInfo.UserType = -1
	} else if MemberInfo.Status == _type.Deleted {
		// 用户已删除 TODO 还没检查
		Response.Code = _type.UserHasDeleted
	} else {
		// 一切正常
		Response.Code = _type.OK
	}
	// 传入Response
	Response.Data = ResponseMemberInfo
	// 用户未登录 还没做

	// 返回Response
	c.JSON(http.StatusOK, Response)
}

// ListMember GET批量获取成员信息
func ListMember(c *gin.Context) {
	var Response _type.GetMemberListResponse
	var Request _type.GetMemberListRequest
	// 获取JSON参数
	if err := c.ShouldBindQuery(&Request); err != nil {
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	// 获取Request中需要的参数Offset和Limit
	Offset := Request.Offset
	Limit := Request.Limit
	// 全局变量DB赋值
	db := DB.MysqlDB
	// 查找数据库
	var ListedMembers []_type.Member
	result := db.Limit(Limit).Offset(*Offset).Find(&ListedMembers, "status = ?", _type.Existed) // 传入0值binding会报错.需要指针
	// SELECT * FROM members OFFSET 5 LIMIT 10 WHERE status = 1;
	// 处理结果
	Rows := result.RowsAffected // 返回找到的记录数
	if err := result.Error; err != nil {
		Response.Code = _type.UnknownError
		c.JSON(http.StatusOK, Response)
		return
	}
	//fmt.Println(Rows, len(ListedMembers))
	//
	ResponseMemberList := make([]_type.TMember, Rows)
	for i := 0; i < int(Rows); i++ {
		ResponseMemberList[i].UserID = strconv.Itoa(int(ListedMembers[i].UserID))
		ResponseMemberList[i].Nickname = ListedMembers[i].Nickname
		ResponseMemberList[i].Username = ListedMembers[i].Username
		ResponseMemberList[i].UserType = ListedMembers[i].UserType
	}
	Response.Code = _type.OK
	Response.Data.MemberList = ResponseMemberList
	c.JSON(http.StatusOK, Response)
}

// UpdateMember POST更新成员信息
func UpdateMember(c *gin.Context) {
	var Response _type.UpdateMemberResponse
	var Request _type.UpdateMemberRequest
	// 获取JSON参数
	if err := c.ShouldBindJSON(&Request); err != nil {
		fmt.Println(err)
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	// 获取Request中参数
	UserID, err := strconv.Atoi(Request.UserID)
	if err != nil {
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	Nickname := Request.Nickname
	// 判断用户是否存在or已删除
	switch Response.Code, _ = modules.IsMemberOK(UserID); Response.Code {
	case _type.UserNotExisted:
		fallthrough
	case _type.UserHasDeleted:
		c.JSON(http.StatusOK, Response)
		return
	}
	// 更新数据库
	// 条件更新
	db := DB.MysqlDB
	result := db.Model(&_type.Member{}).Where("user_id = ?", UserID).Update("nickname", Nickname)
	if err := result.Error; err != nil {
		fmt.Println(err)
		Response.Code = _type.UnknownError
	} else {
		Response.Code = _type.OK
	}
	c.JSON(http.StatusOK, Response)
}

// DeleteMember POST删除成员信息
func DeleteMember(c *gin.Context) {
	var Response _type.UpdateMemberResponse
	var Request _type.UpdateMemberRequest
	// 获取JSON参数
	if err := c.ShouldBindJSON(&Request); err != nil {
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	// 获取Request中参数
	UserID, err := strconv.Atoi(Request.UserID)
	if err != nil {
		Response.Code = _type.ParamInvalid
		c.JSON(http.StatusOK, Response)
		return
	}
	Status := _type.Deleted
	// 判断用户是否存在or已删除
	switch Response.Code, _ = modules.IsMemberOK(UserID); Response.Code {
	case _type.UserNotExisted:
		fallthrough
	case _type.UserHasDeleted:
		c.JSON(http.StatusOK, Response)
		return
	}
	// 更新数据库，删除用户
	// 条件更新
	db := DB.MysqlDB
	result := db.Model(&_type.Member{}).Where("user_id = ?", UserID).Update("status", Status)
	if err := result.Error; err != nil {
		fmt.Println(err)
		Response.Code = _type.UnknownError
	} else {
		Response.Code = _type.OK
		//update Redis
		DB.DeleteStudent(strconv.Itoa(UserID))
	}
	c.JSON(http.StatusOK, Response)
}
