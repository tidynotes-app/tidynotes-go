package models

type Stores struct {
	FileID      int    `gorm:"primaryKey;autoIncrement;column:fileid" json:"fileid"`
	FileName    string `gorm:"unique;column:filename" json:"filename"`
	SubjectName string `gorm:"column:subjectname" json:"subjectname"`
	GoogleID    string `gorm:"column:googleid" json:"googleid"`
}

func (Stores) TableName() string {
	return "tidynotes.stores"
}
