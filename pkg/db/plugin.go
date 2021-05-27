package db

import (
	"gorm.io/datatypes"
)

// plugins 表

type Plugin struct {
	Id            int            `gorm:"primary_key" json:"id"`
	VulId         string         `gorm:"column:vul_id" json:"vul_id"`
	Affects       string         `gorm:"column:affects" json:"affects"`
	JsonPoc       datatypes.JSON `gorm:"column:json_poc" json:"json_poc"`
	Enable        bool           `gorm:"column:enable" json:"enable"`
	Desc          int            `gorm:"column:description" json:"description"`
	Vulnerability *Vulnerability `gorm:"foreignkey:Desc"`
}

type PluginSearchField struct {
	Search       string
	EnableField  int
	AffectsField string
}

func GetPluginsTotal(field *PluginSearchField) (total int64){
	db := GlobalDB.Model(&Plugin{})

	if field.AffectsField != ""{
		db = db.Where("affects = ?", field.AffectsField)
	}
	if field.EnableField != -1{
		db = db.Where("enable = ?", field.EnableField)
	}
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("vul_id like ?", "%"+field.Search+"%").
				Or("json_poc like ?", "%"+field.Search+"%"))
	}
	db.Count(&total)
	return
}

func GetPlugins(page int, pageSize int, field *PluginSearchField) (plugins []Plugin) {

	db := GlobalDB.Preload("Vulnerability")

	if field.AffectsField != ""{
		db = db.Where("affects = ?", field.AffectsField)
	}
	if field.EnableField != -1{
		db = db.Where("enable = ?", field.EnableField)
	}
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("vul_id like ?", "%"+field.Search+"%").
				Or("json_poc like ?", "%"+field.Search+"%"))
	}
	//	分页
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&plugins)
	}
	return
}

func GetPlugin(id int) (plugin Plugin){
	GlobalDB.Model(&Plugin{}).Where("id = ?", id).First(&plugin)
	return
}

func EditPlugin(id int, plugin Plugin) bool {
	GlobalDB.Model(&Plugin{}).Model(&Plugin{}).Where("id = ?", id).Updates(plugin)
	return true
}

func AddPlugin(plugin Plugin) bool {
	GlobalDB.Create(&plugin)
	return true
}

func DeletePlugin(id int) bool {
	GlobalDB.Model(&Plugin{}).Where("id = ?", id).Delete(Plugin{})
	return true
}

func ExistPluginByID(id int) bool {
	var plugin Plugin
	GlobalDB.Model(&Plugin{}).Where("id = ?", id).First(&plugin)
	if plugin.Id >0 {
		return true
	}
	return false
}

func ExistPluginByVulId(vul_id string) bool {
	var plugin Plugin
	GlobalDB.Model(&Plugin{}).Where("vul_id = ?", vul_id).First(&plugin)
	if plugin.Id >0 {
		return true
	}
	return false
}

