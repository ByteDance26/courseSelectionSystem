package router

import (
	"courseSelectionSystem/controller"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 登录
	g.POST("/auth/login", controller.LoginHandle)
	g.POST("/auth/logout", controller.LogoutHandle)
	g.GET("/auth/whoami", controller.WhoamiHandle)

	// 成员管理
	g.POST("/member/create", controller.CreateMember)
	g.GET("/member", controller.GetMember)
	g.GET("/member/list", controller.ListMember)
	g.POST("/member/update", controller.UpdateMember)
	g.POST("/member/delete", controller.DeleteMember)

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
