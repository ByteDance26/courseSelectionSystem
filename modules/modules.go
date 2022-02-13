package modules

//需要更改的地方，
//1、为了方便，之后将关系表冗余一下
//2、老师不存在错误的处理，等统一了再去写

import (
	"awesomeProject1/DB"
	_type "awesomeProject1/type"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Course struct {
	CourseId int64 `gorm:"primary_key;AUTO_INCREMENT;"`
	Name     string
	Cap      int
}
type BindCourse struct {
	CourseId  int64 `gorm:"primary_key"`
	TeacherId int64
}
type Teacher struct {
	TeacherId int64 `gorm:"primary_key;column:user_id"`
	Nickname  string
	Username  string
	UserType  _type.UserType
}

//course

func (Course) TableName() string {
	return "course"
}
func (course *Course) Insert() error {
	return DB.MysqlDB.Create(course).Error
}

func GetCourseById(id int64) (*_type.TCourse, error) {
	var course _type.TCourse
	err := DB.MysqlDB.Raw("select c.course_id as course_id,c.name as name,bc.teacher_id as teacher_id from course c left join bind_course bc on c.course_id = bc.course_id where c.course_id= ?", id).
		Scan(&course).Error
	return &course, err
}
func (BindCourse) TableName() string {
	return "bind_course"
}

//BindCourse
//errType，0表示正常，课程不存在，老师不存在和课程已经绑定使用types.go中定义的错误编码

func (bindCourse *BindCourse) Insert() (errType _type.ErrNo, err error) {

	DB.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var (
			course  Course
			bc      BindCourse
			teacher Teacher
		)
		//判断课程是否存在
		if err1 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&course, "course_id = ?", bindCourse.CourseId).Error; err1 != nil {
			if errors.Is(err1, gorm.ErrRecordNotFound) {
				//课程不存在错误
				errType = _type.CourseNotExisted
				err = err1
				return nil
			} else {
				//未知错误
				errType = _type.UnknownError
				err = err1
				return nil
			}
		}
		//判断老师是否存在
		if err1 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&teacher, "user_id = ?", bindCourse.TeacherId).Error; err1 != nil {
			if errors.Is(err1, gorm.ErrRecordNotFound) {
				//老师不存在错误
				errType = _type.UnknownError
				err = err1
				return nil
			} else {
				//其他错误
				errType = _type.UnknownError
				err = err1
				return nil
			}
		}
		//判断课程是否已经绑定
		var err2 error
		err2 = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&bc, "course_id = ?", bindCourse.CourseId).Error

		if errors.Is(err2, gorm.ErrRecordNotFound) {
			//说明可以绑定
			if err3 := tx.Create(bindCourse).Error; err3 != nil {
				//其他错误
				errType = _type.UnknownError
				err = err3
				return nil
			}
			errType = _type.OK
			return nil

		} else if err2 != nil {
			//说明有其他错误
			errType = _type.UnknownError
			err = err2
			return nil
		} else {
			//说明有存在绑定的错误
			errType = _type.CourseHasBound
			err = err2
			return nil
		}

		return nil
	})
	return errType, err
}

func (bindCourse *BindCourse) Delete() (errType _type.ErrNo, err error) {

	DB.MysqlDB.Transaction(func(tx *gorm.DB) error {
		var (
			course  Course
			bc      BindCourse
			teacher Teacher
		)
		//判断课程是否存在
		if err1 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&course, "course_id = ?", bindCourse.CourseId).Error; err1 != nil {
			if errors.Is(err1, gorm.ErrRecordNotFound) {
				//课程不存在错误
				errType = _type.CourseNotExisted
				err = err1
				return nil
			} else {
				//未知错误
				errType = _type.UnknownError
				err = err1
				return nil
			}

		}
		//判断老师是否存在
		if err1 := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&teacher, "user_id = ?", bindCourse.TeacherId).Error; err1 != nil {
			if errors.Is(err1, gorm.ErrRecordNotFound) {
				//老师不存在错误，也就是课程不存在
				errType = _type.CourseNotBind
				err = err1
				return nil
			} else {
				//其他错误
				errType = _type.UnknownError
				err = err1
				return nil
			}
		}
		//判断课程是否已经绑定
		var err2 error
		err2 = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&bc, "course_id = ?", bindCourse.CourseId).Error

		if errors.Is(err2, gorm.ErrRecordNotFound) {
			//说明有未绑定过错误
			errType = _type.CourseNotBind
			err = err2
			return nil
		} else if err2 != nil {
			errType = _type.UnknownError
			err = err2
			return nil
			//说明有其他错误

		} else {
			//说明可以解绑
			if err3 := tx.Delete(bindCourse).Error; err3 != nil {
				errType = _type.UnknownError
				err = err3
				return nil
			}
			errType = _type.OK
			return nil
		}
		return nil
	})
	return errType, err
}

//teacher

func (Teacher) TableName() string {
	return "member"
}
func GetTeacherById(id int64) (*Teacher, error) {

	var teacher Teacher
	err := DB.MysqlDB.First(&teacher, "user_id = ?", id).Error
	return &teacher, err
}

//出现两种错误，
//一找不到
//二找到却出错

func GetCoursesByTeacherId(id int64) (courses []*_type.TCourse, err error) {
	if err2 := DB.MysqlDB.Raw("select c.course_id as course_id,c.name as name,bc.teacher_id as teacher_id from course c join bind_course bc on c.course_id = bc.course_id where bc.teacher_id= ?", id).
		Scan(&courses).Error; err2 != nil {
		err = err2
		return courses, err
	}
	return courses, err
}
