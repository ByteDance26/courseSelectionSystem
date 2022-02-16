package controller

import (
	"courseSelectionSystem/middle"
	"courseSelectionSystem/modules"
	_type "courseSelectionSystem/type"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func LoginHandle(c *gin.Context) {
	var r _type.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusOK, _type.LoginResponse{
			Code: _type.ParamInvalid,
		})
	} else {
		var mem modules.Member
		err := mem.GetMemberByUsername(r.Username)
		if err == gorm.ErrRecordNotFound { //check user exsist
			c.JSON(http.StatusOK, _type.LoginResponse{
				Code: _type.WrongPassword,
			})
		} else {
			if r.Password != mem.Password { // check password
				c.JSON(http.StatusOK, _type.LoginResponse{
					Code: _type.WrongPassword,
				})
			} else { //ok login
				userIdStr := strconv.Itoa(mem.UserId)
				middle.SetUserId(c, userIdStr)
				s, err := middle.GetSimpleSession(c)
				_ = err
				c.SetCookie("camp-session", s.SessionID, 999, "/", "localhost", false, true)
				c.JSON(http.StatusOK, _type.LoginResponse{
					Code: _type.OK,
					Data: struct{ UserID string }{UserID: userIdStr},
				})
			}
		}
	}
}

func LogoutHandle(c *gin.Context) {
	//TODO 参数未检查
	id, _ := middle.GetUserId(c)
	err := middle.DelUserId(c)
	if err == _type.LoginRequired {
		c.JSON(http.StatusOK, _type.LogoutResponse{
			Code: _type.LoginRequired,
		})
	} else {
		c.SetCookie("camp-session", "", -1, "/", "localhost", false, true)
		c.JSON(http.StatusOK, _type.LoginResponse{
			Code: _type.OK,
			Data: struct{ UserID string }{UserID: id},
		})
	}
}

func WhoamiHandle(c *gin.Context) {
	//TODO 参数未检查
	id, err := middle.GetUserId(c)
	if err == _type.LoginRequired {
		c.JSON(http.StatusOK, _type.WhoAmIResponse{
			Code: _type.LoginRequired,
		})
	} else {
		var mem modules.Member
		atoi, err2 := strconv.Atoi(id)
		_ = err2
		err := mem.GetMemberByUserId(atoi)
		if err == gorm.ErrRecordNotFound { //check user exsist
			c.JSON(http.StatusOK, _type.LoginResponse{
				Code: _type.UserHasDeleted,
			})
		} else {
			c.JSON(http.StatusOK, _type.WhoAmIResponse{
				Code: _type.OK,
				Data: struct {
					UserID   string
					Nickname string
					Username string
					UserType _type.UserType
				}{UserID: strconv.Itoa(mem.UserId), Nickname: mem.Nickname, Username: mem.Username, UserType: _type.UserType(mem.UserType)},
			})
		}

	}
}
