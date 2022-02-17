package controller

import (
	"courseSelectionSystem/DB"
	_type "courseSelectionSystem/type"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//redis中存的数据
//一个为courseId:detail，用于进行 类型：string
//一个为courseID:take，课程剩余容量 类型：string
//一个为studentID:course 类型:map key:courseID value:teacherID,name  结果使用json利于相互的转换

func BookCourse(c *gin.Context) {
	var BookCourseRequest _type.BookCourseRequest
	var BookCourseResponse _type.BookCourseResponse
	if err := c.ShouldBindJSON(&BookCourseRequest); err != nil {
		//参数错误
		BookCourseResponse = _type.BookCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusOK, BookCourseResponse)
		return
	}

	//初始化
	BookCourseResponse = _type.BookCourseResponse{
		Code: _type.OK,
	}

	//判断课程是否存在
	ok := DB.BoolCourse(BookCourseRequest.CourseID)
	if !ok {
		BookCourseResponse.Code = _type.CourseNotExisted
		c.JSON(http.StatusOK, BookCourseResponse)
		return
	}
	//学生不存在
	ok = DB.BoolStudent(BookCourseRequest.StudentID)
	if !ok {
		BookCourseResponse.Code = _type.StudentNotExisted
		c.JSON(http.StatusOK, BookCourseResponse)
		return

	}
	//判断学生是否存在

	//选课,原子操作,
	//返回值为1表示的是抢课成功
	//返回值为2表示的是已经选了课
	//返回值为3表示的是抢课失败课程已经满了
	//返回值为0表示的其他错误
	ret := DB.AddCourse(BookCourseRequest.StudentID, BookCourseRequest.CourseID)
	switch ret {
	case 1:
		BookCourseResponse.Code = _type.OK
		fmt.Println("学生:", BookCourseRequest.StudentID, "抢到", BookCourseRequest.CourseID, "课")
		break
	case 2:
		BookCourseResponse.Code = _type.StudentHasCourse
		fmt.Println("学生:", BookCourseRequest.StudentID, "未抢到", BookCourseRequest.CourseID, "课,已经选了次课")
	case 3:
		BookCourseResponse.Code = _type.CourseNotAvailable
		fmt.Println("学生:", BookCourseRequest.StudentID, "未抢到", BookCourseRequest.CourseID, "课,课程已经满了")
	case 0:
		BookCourseResponse.Code = _type.UnknownError
		fmt.Println("学生:", BookCourseRequest.StudentID, "未抢到", BookCourseRequest.CourseID, "课,未知错误")
	}
	c.JSON(http.StatusOK, BookCourseResponse)
}

func GetStudentCourse(c *gin.Context) {

	var GetStudentCourseResponse _type.GetStudentCourseResponse
	studentID := c.Query("StudentID")
	if _, err := strconv.ParseInt(studentID, 10, 64); err != nil {
		GetStudentCourseResponse = _type.GetStudentCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusOK, GetStudentCourseResponse)
		return

	}
	xx := DB.BoolStudent(studentID)
	if !xx {
		//学生不存在
		GetStudentCourseResponse = _type.GetStudentCourseResponse{
			Code: _type.StudentNotExisted,
		}
		c.JSON(http.StatusOK, GetStudentCourseResponse)
		return
	}
	courses := DB.GetCourses(studentID)
	if len(courses) == 0 {
		//学生没有课程
		GetStudentCourseResponse = _type.GetStudentCourseResponse{
			Code: _type.StudentHasNoCourse,
		}
		c.JSON(http.StatusOK, GetStudentCourseResponse)
		return
	}
	GetStudentCourseResponse.Data.CourseList = courses
	c.JSON(http.StatusOK, GetStudentCourseResponse)
	return

}
