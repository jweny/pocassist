package db

type Webapp struct {
	Id       	int    		`gorm:"primary_key" json:"id"`
	Name		string		`gorm:"column:name" json:"name"`
	Provider	string		`gorm:"column:provider" json:"provider"`
	Remarks		string		`gorm:"column:remarks" json:"remarks"`
}

func GetWebApps(page int, pageSize int) (apps []Webapp) {
	//	分页
	db := GlobalDB.Model(&Webapp{})
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&apps)
	}
	return
}

func AddWebapp(app Webapp) bool {
	GlobalDB.Create(&app)
	return true
}
