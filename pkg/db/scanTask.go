package db

const TaskStatusRunning  = "running"
const TaskStatusDown  = "down"
const TaskStatusError  = "error"

// task表

type Task struct {
	//gorm.Model
	Id            int            `gorm:"primary_key" json:"id"`
	Remarks       string      	 `gorm:"column:remarks" json:"remarks"`
	Target   	  string    	 `gorm:"type:longtext" json:"target"`
	Operator	  string		 `gorm:"type:string" json:"operator"`
	Status        string		 `gorm:"type:string" json:"status"`
	Results		  []Result		 `gorm:"foreignKey:TaskId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"results"`
}

func AddTask(task *Task) bool {
	task.Status = TaskStatusRunning
	GlobalDB.Create(task)
	return true
}

func ErrorTask(id int) bool {
	GlobalDB.Model(&Task{}).Where("id = ?", id).Update("status",TaskStatusError)
	return true
}

func DownTask(id int) bool {
	GlobalDB.Model(&Task{}).Where("id = ?", id).Update("status",TaskStatusDown)
	return true
}

type TaskSearchField struct {
	Search       string
}

func GetTaskTotal(field *TaskSearchField) (total int64){
	db := GlobalDB.Model(&Task{})
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("remarks like ?", "%"+field.Search+"%").
				Or("Target like ?", "%"+field.Search+"%"))
	}
	db.Count(&total)
	return
}

func GetTask(page int, pageSize int, field *TaskSearchField) (tasks []Task) {
	db := GlobalDB.Model(&Task{})
	if field.Search != ""{
		db = db.Where(
			GlobalDB.Where("remarks like ?", "%"+field.Search+"%").
				Or("Target like ?", "%"+field.Search+"%"))
	}
	//	分页
	if page > 0 && pageSize > 0 {
		db = db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&tasks)
	}
	return
}

func DeleteTask(id int) bool {
	GlobalDB.Select("Results").Model(&Task{}).Where("id = ?", id).Delete(&Task{})
	return true
}

func ExistTaskByID(id int) bool {
	var task Task
	GlobalDB.Model(&Task{}).Where("id = ?", id).First(&task)
	if task.Id >0 {
		return true
	}
	return false
}