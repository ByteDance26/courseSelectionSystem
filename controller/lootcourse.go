package controller

import (
	modules "courseSelectionSystem/modules"
	_type "courseSelectionSystem/type"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetCourse(c *gin.Context) {
	var GetCourseResponse _type.GetCourseResponse
	//传统获取参数
	var CourseID string
	CourseID = c.Query("CourseID")
	var id int64
	var err2 error
	if id, err2 = strconv.ParseInt(CourseID, 10, 64); err2 != nil {
		//参数错误
		GetCourseResponse = _type.GetCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, GetCourseResponse)
		return
	}

	//还没设置老师不存在

	//正常情况
	course, err := modules.GetCourseById(id)
	if err != nil {
		GetCourseResponse = _type.GetCourseResponse{
			Code: _type.UnknownError,
		}
		c.JSON(http.StatusBadRequest, GetCourseResponse)
		return
	}
	GetCourseResponse = _type.GetCourseResponse{
		Code: _type.OK,
		Data: *course,
	}
	c.JSON(http.StatusOK, GetCourseResponse)
}

func CreateCourse(c *gin.Context) {
	var CreateCourseRequest _type.CreateCourseRequest
	var CreateCourseResponse _type.CreateCourseResponse

	if err := c.ShouldBindJSON(&CreateCourseRequest); err != nil {

		//参数错误
		CreateCourseResponse = _type.CreateCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, CreateCourseResponse)
		//fmt.Println(err)
		return
	}

	course := modules.Course{
		Name: CreateCourseRequest.Name,
		Cap:  CreateCourseRequest.Cap,
	}

	//正常情况
	err := course.Insert()
	if err != nil {
		//创建失败
		CreateCourseResponse = _type.CreateCourseResponse{
			Code: _type.UnknownError,
		}
		c.JSON(http.StatusBadRequest, CreateCourseResponse)
		return
	}
	//创建成功
	CreateCourseResponse = _type.CreateCourseResponse{
		Code: _type.OK,
	}
	CreateCourseResponse.Data.CourseID = strconv.FormatInt(course.CourseId, 10)
	c.JSON(http.StatusOK, CreateCourseResponse)
}

func BindCourse(c *gin.Context) {
	var BindCourseRequest _type.BindCourseRequest
	var BindCourseResponse _type.BindCourseResponse

	if err := c.ShouldBindJSON(&BindCourseRequest); err != nil {
		//参数错误
		BindCourseResponse = _type.BindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, BindCourseResponse)
		return
	}

	var bd modules.BindCourse
	var err2 error
	if bd.CourseId, err2 = strconv.ParseInt(BindCourseRequest.CourseID, 10, 64); err2 != nil {
		//参数错误
		BindCourseResponse = _type.BindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, BindCourseResponse)
		return
	}
	if bd.TeacherId, err2 = strconv.ParseInt(BindCourseRequest.TeacherID, 10, 64); err2 != nil {
		//参数错误
		BindCourseResponse = _type.BindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, BindCourseResponse)
		return
	}
	var errType _type.ErrNo
	errType, err2 = bd.Insert()
	switch errType {
	//正常情况
	case _type.OK:
		BindCourseResponse.Code = _type.OK
		break
	//未知错误
	case _type.UnknownError:
		BindCourseResponse.Code = _type.UnknownError
		break
	//课程已经绑定
	case _type.CourseHasBound:
		BindCourseResponse.Code = _type.CourseHasBound
		break
	//课程不存在
	case _type.CourseNotExisted:
		BindCourseResponse.Code = _type.CourseNotExisted
		break
	default:
		break
	}
	c.JSON(http.StatusBadRequest, BindCourseResponse) //TODO::对OK，http返回码不对
	return
}

func UnbindCourse(c *gin.Context) {
	var UnbindCourseRequest _type.UnbindCourseRequest //TODO::改类型需要强制 binding:"required" 才能判断参数错误
	var UnbindCourseResponse _type.UnbindCourseResponse

	if err := c.ShouldBindJSON(&UnbindCourseRequest); err != nil {
		//参数错误
		UnbindCourseResponse = _type.UnbindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, UnbindCourseResponse)
		return
	}

	var bd modules.BindCourse //TODO::前面已经判断过参数正误，是否重复
	var err2 error
	if bd.CourseId, err2 = strconv.ParseInt(UnbindCourseRequest.CourseID, 10, 64); err2 != nil {
		//参数错误
		UnbindCourseResponse = _type.UnbindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, UnbindCourseResponse)
		return
	}
	if bd.TeacherId, err2 = strconv.ParseInt(UnbindCourseRequest.TeacherID, 10, 64); err2 != nil {
		//参数错误
		UnbindCourseResponse = _type.UnbindCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, UnbindCourseResponse)
		return
	}
	var errType _type.ErrNo
	errType, err2 = bd.Delete()
	switch errType {
	//正常情况
	case _type.OK:
		UnbindCourseResponse.Code = _type.OK
		break
	//未知错误
	case _type.UnknownError:
		UnbindCourseResponse.Code = _type.UnknownError
		break
	//课程已经绑定
	case _type.CourseNotBind:
		UnbindCourseResponse.Code = _type.CourseNotBind
		break
	//课程不存在
	case _type.CourseNotExisted:
		UnbindCourseResponse.Code = _type.CourseNotExisted
		break
	default:
		break
	}
	c.JSON(http.StatusBadRequest, UnbindCourseResponse)
	return
}

func GetTeacherCourse(c *gin.Context) {
	var GetTeacherCourseResponse _type.GetTeacherCourseResponse
	TeacherID := c.Query("TeacherID")
	var id int64
	var err2 error
	if id, err2 = strconv.ParseInt(TeacherID, 10, 64); err2 != nil {
		//参数错误
		GetTeacherCourseResponse = _type.GetTeacherCourseResponse{
			Code: _type.ParamInvalid,
		}
		c.JSON(http.StatusBadRequest, GetTeacherCourseResponse)
		return
	}

	//还没设置老师不存在

	courses, err := modules.GetCoursesByTeacherId(id)
	if err != nil {
		//未知错误
		GetTeacherCourseResponse = _type.GetTeacherCourseResponse{
			Code: _type.UnknownError,
		}
		c.JSON(http.StatusBadRequest, GetTeacherCourseResponse)
		return
	}
	//正常返回
	GetTeacherCourseResponse = _type.GetTeacherCourseResponse{
		Code: _type.OK,
	}
	GetTeacherCourseResponse.Data.CourseList = courses
	c.JSON(http.StatusOK, GetTeacherCourseResponse)
}

func ScheduleCourse(c *gin.Context) {
	var ScheduleCourseResponse _type.ScheduleCourseResponse
	var m map[string][]string
	if err := c.ShouldBindJSON(&m); err != nil {
		//参数错误
		ScheduleCourseResponse.Code = _type.ParamInvalid
		c.JSON(http.StatusBadRequest, ScheduleCourseResponse)
		return
	}
	tu := Tu(m)
	ScheduleCourseResponse.Data = tu
	ScheduleCourseResponse.Code = _type.OK
	c.JSON(http.StatusBadRequest, ScheduleCourseResponse)
	return
}

func Tu(edge map[string][]string) (ret map[string]string) {
	//表示的是课程是否被选了
	m := make(map[string]string)
	//表示结果,同时老师是否选择次课程
	ret = make(map[string]string)

	for key, _ := range edge {
		if _, ok := ret[key]; !ok {
			visit := make(map[string]bool)
			if dfs(edge, m, visit, ret, key) {
				continue
			}
		}
	}
	return ret
}

func dfs(edge map[string][]string, m map[string]string, visit map[string]bool, n map[string]string, u string) bool {
	var courses []string
	courses = edge[u]
	for _, val := range courses {
		if _, ok := visit[val]; !ok {
			visit[val] = true
			if val1, ok1 := m[val]; !ok1 || dfs(edge, m, visit, n, val1) {
				n[u] = val
				m[val] = u
				return true
			}
		}
	}
	return false
}
