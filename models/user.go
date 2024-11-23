package models

type Users struct {
	UserID       uint   `gorm:"primaryKey;autoIncrement;column:userid" json:"userid"`
	UserName     string `gorm:"unique;column:username" json:"username"`
	UserEmail    string `gorm:"unique;column:useremail" json:"useremail"`
	UserPassword string `gorm:"column:userpassword" json:"userpassword"`
	IsLoggedIn   bool   `gorm:"column:isloggedin" json:"isloggedin"`
}

func (Users) TableName() string {
	return "tidynotes.users"
}
