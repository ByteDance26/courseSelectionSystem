package DB

import (
	_type "courseSelectionSystem/type"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
)

var RedisDB *redis.Pool

//选课,原子操作,
//返回值为1表示的是抢课成功
//返回值为2表示的是已经选了课
//返回值为3表示的是抢课失败课程已经满了
//返回值为0表示的其他错误

const (
	SCRIPT_INCR = `
	local studentID=KEYS[1]
	local courseID=KEYS[2]
	local stuKey="student_"..studentID.."_course"
	local couKey="course_"..courseID.."_take"
	local s=redis.call("sismember",stuKey,courseID)
	if tonumber(s) == 1 then
	return 2
	end
	local num=redis.call("get",couKey)
	if tonumber(num)<=0 then
	return 3
	end
	redis.call("decr",couKey)
	redis.call("sadd",stuKey,courseID)
	return 1
`
)

var newScript *redis.Script

func RedisInit() {
	RedisDB = &redis.Pool{ //实例化一个连接池
		MaxIdle: 32, //最初的连接数量
		// MaxActive:1000000,    //最大连接数量
		MaxActive:   0,   //连接池最大连接数量,不确定可以用0（0表示自动定义），按需分配
		IdleTimeout: 300, //连接关闭时间 300秒 （300秒不使用自动关闭）
		Dial: func() (redis.Conn, error) { //要连接的redis数据库
			return redis.Dial("tcp", "180.184.70.231:6379")
		},
	}
	//脚本预加载
	newScript = redis.NewScript(2, SCRIPT_INCR)
}

//redis中存的数据
//一个为course，用于进行 类型：hash
//一个为courseID:take，课程剩余容量 类型：string
//一个为studentID:course 类型:hash key:courseID value:teacherID,name  结果使用json利于相互的转换
//创建课程

func CreateCourse(course _type.TCourse, cap int) {
	c := RedisDB.Get()
	defer c.Close()
	str, err := json.Marshal(course)
	if err != nil {
		fmt.Println(err)
		return
	}
	//课程的详细信息
	_, err = c.Do("hset", "course", course.CourseID, str)
	if err != nil {
		fmt.Println(err)
		return
	}
	//用于抢课的课程信息
	_, err = c.Do("set", "course_"+course.CourseID+"_take", cap)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

//用于解绑的时候进行课程的删除

func DeleteCourse(courseID string) {
	c := RedisDB.Get()
	defer c.Close()
	var err error
	//用于保存课程信息
	_, err = c.Do("srem", "course", courseID)
	if err != nil {
		fmt.Println(err)
		return
	}
	//用于抢课的课程信息
	_, err = c.Do("del", "course_"+courseID+"_take", cap)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}

//选课,原子操作,
//返回值为1表示的是抢课成功
//返回值为2表示的是已经选了课
//返回值为3表示的是抢课失败课程已经满了
//返回值为0表示的其他错误

//选择课程

func AddCourse(studentID string, courseID string) int {
	c := RedisDB.Get()
	defer c.Close()

	ret, err := redis.Int(newScript.Do(c, studentID, courseID))
	if err != nil {
		return 0
	}
	return ret
}

//添加学生
//一个student 类型：set

func CreateStudent(studentID string) {
	c := RedisDB.Get()
	defer c.Close()
	_, err := c.Do("sadd", "student", studentID)
	if err != nil {
		fmt.Println(err)
	}
}

//删除学生,改善删除学生的函数

func DeleteStudent(studentID string) {
	c := RedisDB.Get()
	defer c.Close()
	//"student_"..studentID.."_course"
	_, err := c.Do("srem", "student", studentID)
	_, err = c.Do("del", "student_"+studentID+"_course")
	if err != nil {
		fmt.Println(err)
	}
}

//判断学生是否存在

func BoolStudent(studentID string) bool {
	c := RedisDB.Get()
	defer c.Close()
	do, err := c.Do("sismember", "student", studentID)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if i, ok := do.(int64); !ok || i == 0 {
		return false
	}
	return true
}

//判断课程是否存在

func BoolCourse(courseId string) bool {
	c := RedisDB.Get()
	defer c.Close()
	do, err := c.Do("hexists", "course", courseId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	if i, ok := do.(int64); !ok || i == 0 {
		return false
	}
	return true
}

//获取学生课表

func GetCourses(studentID string) (ret []_type.TCourse) {
	c := RedisDB.Get()
	var do interface{}
	var err error
	ret = make([]_type.TCourse, 0, 0)
	do, err = c.Do("smembers", "student_"+studentID+"_course")
	if err != nil {
		fmt.Println(err)
	}
	//先判断是否选了课程

	//获取所以的课程
	if courses, ok := do.([]interface{}); ok {
		for _, val := range courses {
			reply, err1 := c.Do("hget", "course", val)
			if err != nil {
				fmt.Println(err1)
			}
			if str, ok1 := reply.([]uint8); ok1 {
				var tCourse _type.TCourse
				json.Unmarshal(str, &tCourse)
				ret = append(ret, tCourse)
			}
		}
	}
	return ret
}
