package controller

import (
	"awesomeProject1/DB"
	"awesomeProject1/modules"
	_type "awesomeProject1/type"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
)

// CreateMember 创建新用户
// TODO 检查当前登录用户是否是管理员 LoginRequired
func CreateMember(c *gin.Context) {
	var Response _type.CreateMemberResponse
	var Request _type.CreateMemberRequest
	// 获取JSON参数
	if err := c.ShouldBind(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Connection error"})
	}
	// 将参数绑定到Member里
	NewMember := _type.Member{
		Nickname: Request.Nickname,
		Username: Request.Username,
		Password: Request.Password,
		UserType: Request.UserType,
		Status:   _type.Existed,
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
		} else if Error := fmt.Sprintf("/v", result.Error); strings.Contains(Error, "Error 1062") {
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

// GetMember 获取当前登录用户信息
func GetMember(c *gin.Context) {
	var Response _type.GetMemberResponse
	var Request _type.GetMemberRequest
	// 获取JSON参数
	if err := c.ShouldBind(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Connection error"})
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

func ListMember(c *gin.Context) {
	var Response _type.GetMemberListResponse
	var Request _type.GetMemberListRequest
	// 获取JSON参数
	if err := c.ShouldBind(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Connection error"})
	}
	// 获取Request中需要的参数Offset和Limit
	Offset := Request.Offset
	Limit := Request.Limit
	// 全局变量DB赋值
	db := DB.MysqlDB
	// 查找数据库
	var ListedMembers []_type.Member
	result := db.Limit(Limit).Offset(Offset).Find(&ListedMembers, "status = ?", _type.Existed)
	// SELECT * FROM members OFFSET 5 LIMIT 10 WHERE status = 1;
	// 处理结果
	Rows := result.RowsAffected // 返回找到的记录数
	if err := result.Error; err != nil {
		Response.Code = _type.UnknownError
		c.JSON(http.StatusOK, Response)
		return
	}
	fmt.Println(Rows, len(ListedMembers))
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
