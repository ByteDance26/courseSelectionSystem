package main

import (
	DB "courseSelectionSystem/DB"
	"courseSelectionSystem/controller"
	"courseSelectionSystem/modules"
	"courseSelectionSystem/router"
	"fmt"
	"github.com/gin-gonic/gin"
)

func mainTest() {
	DB.MysqlInit()

	r := gin.Default()
	router.RegisterRouter(r)
	r.Run()

}

//DBTest
func testBindCourseInsert() {
	bd := &modules.BindCourse{
		CourseId:  2,
		TeacherId: 1,
	}
	errTy, err := bd.Insert()
	fmt.Println(errTy, err)
}

func testBindCourseDelete() {
	bd := &modules.BindCourse{
		CourseId:  2,
		TeacherId: 1,
	}
	errTy, err := bd.Delete()
	fmt.Println(errTy, err)
}
func testGetCoursesByTeacherId() {
	courses, err := modules.GetCoursesByTeacherId(1)
	fmt.Printf("courses:%v,err:%v", courses[0], err)
}
func testGetTeacherById() {
	teacher, err := modules.GetTeacherById(1)
	if err != nil {
		panic(fmt.Sprintln(err))
	}
	fmt.Println(teacher)
}
func testGetCourseByID() {
	course, err := modules.GetCourseById(1)
	fmt.Printf(" %v courses:%v,err:%v", *course, course, err)
}

//算法Test
func testTu() {
	m := map[string][]string{
		"1": []string{"1", "2"},
		"2": []string{"2", "3"},
		"3": []string{"1", "2"},
		"4": []string{"3"},
	}
	tu := controller.Tu(m)
	fmt.Printf("%v", tu)
}
