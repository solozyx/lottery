package services

import (
	"fmt"
	"strconv"
	"time"

	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

type UserdayService interface {
	GetAll(page, size int) []models.LtUserday
	CountAll() int64
	Search(uid, day int) []models.LtUserday
	Count(uid, day int) int
	Get(id int) *models.LtUserday
	//Delete(id int) error
	Update(user *models.LtUserday, columns []string) error
	Create(user *models.LtUserday) error
	GetUserToday(uid int) *models.LtUserday
}

type userdayService struct {
	dao *dao.UserdayDao
}

func NewUserdayService() UserdayService {
	return &userdayService{
		dao: dao.NewUserdayDao(datasource.InstanceDbMaster()),
	}
}

func (s *userdayService) GetAll(page, size int) []models.LtUserday {
	return s.dao.GetAll(page, size)
}

func (s *userdayService) CountAll() int64 {
	return s.dao.CountAll()
}

func (s *userdayService) Search(uid, day int) []models.LtUserday {
	return s.dao.Search(uid, day)
}

func (s *userdayService) Count(uid, day int) int {
	return s.dao.Count(uid, day)
}

func (s *userdayService) Get(id int) *models.LtUserday {
	return s.dao.Get(id)
}

//func (s *userdayService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

func (s *userdayService) Update(data *models.LtUserday, columns []string) error {
	return s.dao.Update(data, columns)
}

func (s *userdayService) Create(data *models.LtUserday) error {
	return s.dao.Create(data)
}

func (s *userdayService) GetUserToday(uid int) *models.LtUserday {
	y, m, d := time.Now().Date()
	strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
	day, _ := strconv.Atoi(strDay)
	list := s.dao.Search(uid, day)
	if list != nil && len(list) > 0 {
		return &list[0]
	} else {
		return nil
	}
}
