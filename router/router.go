package router

import (
	"courseSelectionSystem/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 成员管理
	g.Use()
	g.POST("/member/create")
	g.GET("/member")
	g.GET("/member/list")
	g.POST("/member/update")
	g.POST("/member/delete")

	// 登录

	g.POST("/auth/login", controller.LoginHandle)
	g.POST("/auth/logout", controller.LogoutHandle)
	g.GET("/auth/whoami", controller.WhoamiHandle)

	// 排课
	g.POST("/course/create", controller.CreateCourse)
	g.GET("/course/get", controller.GetCourse)

	g.POST("/teacher/bind_course", controller.BindCourse)
	g.POST("/teacher/unbind_course", controller.UnbindCourse)
	g.GET("/teacher/get_course", controller.GetTeacherCourse)
	g.POST("/course/schedule", controller.ScheduleCourse)

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

}
