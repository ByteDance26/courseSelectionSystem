package router

import (
	"courseSelectionSystem/controller"
	"courseSelectionSystem/middle"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	g := r.Group("/api/v1")

	// 登录
	g.POST("/auth/login", middle.HandleSimpleSession("camp-session"), controller.LoginHandle)
	g.POST("/auth/logout", middle.HandleSimpleSession("camp-session"), controller.LogoutHandle)
	g.GET("/auth/whoami", middle.HandleSimpleSession("camp-session"), controller.WhoamiHandle)

	// 成员管理
	g.POST("/member/create", middle.HandleSimpleSession("camp-session"), controller.CreateMember)
	g.GET("/member", middle.HandleSimpleSession("camp-session"), controller.GetMember)
	g.GET("/member/list", middle.HandleSimpleSession("camp-session"), controller.ListMember)
	g.POST("/member/update", middle.HandleSimpleSession("camp-session"), controller.UpdateMember)
	g.POST("/member/delete", middle.HandleSimpleSession("camp-session"), controller.DeleteMember)

	// 排课
	g.POST("/course/create", middle.HandleSimpleSession("camp-session"), controller.CreateCourse)
	g.GET("/course/get", middle.HandleSimpleSession("camp-session"), controller.GetCourse)

	g.POST("/teacher/bind_course", middle.HandleSimpleSession("camp-session"), controller.BindCourse)
	g.POST("/teacher/unbind_course", middle.HandleSimpleSession("camp-session"), controller.UnbindCourse)
	g.GET("/teacher/get_course", middle.HandleSimpleSession("camp-session"), controller.GetTeacherCourse)
	g.POST("/course/schedule", middle.HandleSimpleSession("camp-session"), controller.ScheduleCourse)

	// 抢课
	g.POST("/student/book_course", controller.BookCourse)
	g.GET("/student/course", controller.GetStudentCourse)

}
