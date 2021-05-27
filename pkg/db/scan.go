package db

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// task表

type Task struct {
	gorm.Model
	Id            int            `gorm:"primary_key" json:"id"`
	Remarks       string      	 `gorm:"column:remarks" json:"remarks"`
	Target   	  string    	 `gorm:"type:longtext" json:"Target"`
	Operator	  string		 `gorm:"type:string" json:"operator"`
}

// result表

type Result struct {
	gorm.Model
	Id            int            `gorm:"primary_key" json:"id"`
	Detail    	  datatypes.JSON `gorm:"column:detail" json:"detail"`
	TaskId		  Task			 `gorm:"foreignkey:Desc"`
}

type ResultSearchField struct {
	Search       string
}

//func GetResultTotal(field *ResultSearchField) (total int64){
//	db := GlobalDB.Model(&Result{})
//	if field.Search != ""{
//		db = db.Where(
//			GlobalDB.Where("remarks like ?", "%"+field.Search+"%").
//				Or("Target like ?", "%"+field.Search+"%"))
//	}
//	db.Count(&total)
//	return
//}
//
//func GetResult(page int, pageSize int, field *TaskSearchField) (tasks []Task) {
//
//	db := GlobalDB.Model(&Task{})
//
//	if field.Search != ""{
//		db = db.Where(
//			GlobalDB.Where("remarks like ?", "%"+field.Search+"%").
//				Or("Target like ?", "%"+field.Search+"%"))
//	}
//	//	分页
//	if page > 0 && pageSize > 0 {
//		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)
//	}
//	return
//}
