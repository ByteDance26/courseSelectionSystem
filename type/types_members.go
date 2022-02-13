package _type

type StatusType int

const (
	Existed StatusType = 1
	Deleted StatusType = -1
)

// Member 对应数据库表格
type Member struct {
	// binding:"required" 绑定的时候发现是空值就会返回错误
	UserID   int64      `gorm:"primaryKey;autoIncrement:true"`
	Nickname string     `form:"Nickname" json:"Nickname" binding:"required"`
	Username string     `form:"Username" json:"Username" binding:"required" gorm:"index:,unique"`
	Password string     `form:"Password" json:"Password" binding:"required"`
	UserType UserType   `form:"UserType" json:"UserType" binding:"required" gorm:"type:tinyint"`
	Status   StatusType `gorm:"type:tinyint;default:1"`
}
