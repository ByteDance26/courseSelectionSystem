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

func CreateMember(c *gin.Context) {
	var Response _type.CreateMemberResponse
	var Request _type.CreateMemberRequest
	// 获取JSON参数
	if err := c.ShouldBind(&Request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error"})
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
		// 参数合法（Username在后面)
		Response.Code = _type.OK
	} else {
		// 参数不合法
		Response.Code = _type.ParamInvalid
		Response.Data.UserID = ""
	}

	// 添加记录
	if Response.Code == _type.OK {
		db := DB.MysqlDB
		tx := db.Create(&NewMember)
		// 判断用户名是否重复

		if tx.Error == nil {
			// Username 没有重复，成功创建
			Response.Data.UserID = strconv.Itoa(int(NewMember.UserID))
		} else if Error := fmt.Sprintf("/v", tx.Error); strings.Contains(Error, "Error 1062") {
			// Username 重复
			Response.Code = _type.UserHasExisted
			Response.Data.UserID = ""
		} else {
			fmt.Println(Error)
		}
	}
	c.JSON(http.StatusOK, Response)
}
