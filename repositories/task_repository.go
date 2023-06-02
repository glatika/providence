package repositories

import (
	"time"

	"github.com/glatika/providence/model"

	"gorm.io/gorm"
)

type taskRepo struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) model.TaskRepository {
	return taskRepo{
		db: db,
	}
}

func (r taskRepo) Create(n *model.Task) (int32, error) {
	err := r.db.Table("tasks").Create(&n).Error
	return n.Id, err
}

func (r taskRepo) FindTaskById(id int) (*model.Task, error) {
	task := model.Task{}
	err := r.db.First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r taskRepo) FindUndeliveredTaskByStockId(stockid int32) (*model.Task, error) {
	task := model.Task{}
	err := r.db.Table("tasks").First(&task, "stock_id = ? and delivered = ?", stockid, false).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r taskRepo) FindUnreportedTaskById(taskid int32) (*model.Task, error) {
	task := model.Task{}
	err := r.db.Table("tasks").First(&task, "id = ? and reported = ?", taskid, false).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r taskRepo) GetAllTasks(page, limit int) ([]model.Task, error) {
	tasks := []model.Task{}
	err := r.db.Offset((page - 1) * limit).Limit(limit).Find(&tasks).Error
	return tasks, err
}

func (r taskRepo) UpdateDeliveredByID(id int) error {
	task, err := r.FindTaskById(id)
	if err != nil {
		return err
	}

	if task == nil {
		return model.ErrRecordNotFound
	}

	now := time.Now()
	task.Delivered = true
	task.DeliveredAt = now

	return r.db.Save(task).Error
}

func (r taskRepo) UpdateReportedByID(id int, success bool, report string) error {
	task, err := r.FindTaskById(id)
	if err != nil {
		return err
	}

	if task == nil {
		return model.ErrRecordNotFound
	}
	now := time.Now()
	task.Reported = true
	task.Success = success
	task.ReportedAt = now
	task.Report = report

	return r.db.Save(task).Error
}
