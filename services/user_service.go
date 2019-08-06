/**
 * 抽奖系统数据处理（包括数据库，也包括缓存等其他形式数据）
 */
package services

import (
	"fmt"
	"log"
	"lottery/conf"
	"sync"

	"github.com/gomodule/redigo/redis"

	"lottery/comm"
	"lottery/dao"
	"lottery/datasource"
	"lottery/models"
)

// 用户信息，可以缓存(本地或者redis)，有更新的时候，可以直接清除缓存或者根据具体情况更新缓存
var cachedUserList = make(map[int]*models.LtUser)
var cachedUserLock = sync.Mutex{}

type UserService interface {
	GetAll(page, size int) []models.LtUser
	CountAll() int
	//Search(country string) []models.LtUser
	Get(id int) *models.LtUser
	//Delete(id int) error
	Update(user *models.LtUser, columns []string) error
	Create(user *models.LtUser) error
}

type userService struct {
	dao *dao.UserDao
}

func NewUserService() UserService {
	return &userService{
		dao: dao.NewUserDao(datasource.InstanceDbMaster()),
	}
}

func (s *userService) GetAll(page, size int) []models.LtUser {
	return s.dao.GetAll(page, size)
}

func (s *userService) CountAll() int {
	return s.dao.CountAll()
}

//func (s *userService) Search(country string) []models.LtUser {
//	return s.dao.Search(country)
//}

func (s *userService) Get(id int) *models.LtUser {
	// 先读Redis缓存
	data := s.getByCache(id)
	if data == nil || data.Id <= 0 {
		// 缓存没有数据,读MySQL数据库
		data = s.dao.Get(id)
		if data == nil || data.Id <= 0 {
			// 构造空值结构只有id有值,写入缓存 ; 后面从缓存读取空值出来,就不读数据库了
			data = &models.LtUser{Id: id}
		}
		// 把数据存到Redis缓存
		s.setByCache(data)
	}
	return data
}

//func (s *userService) Delete(id int) error {
//	return s.dao.Delete(id)
//}

func (s *userService) Update(data *models.LtUser, columns []string) error {
	// 先更新缓存,这里直接是清空该data对应的缓存数据;后面再读取会从数据库更新1个新数据到缓存
	s.updateByCache(data, columns)
	// 再更新数据
	return s.dao.Update(data, columns)
}

func (s *userService) Create(data *models.LtUser) error {
	return s.dao.Create(data)
}

func (s *userService) getByCache(id int) *models.LtUser {
	// key := fmt.Sprintf("info_user_%d", id)
	key := fmt.Sprintf(conf.RdsUserCacheKeyPrefix + "%d", id)
	rds := datasource.InstanceCache()
	// 使用 Redis Hash 结构 HGETALL 读取整个hash结构 HGET 读取特定hash字段
	dataMap, err := redis.StringMap(rds.Do("HGETALL", key))
	if err != nil {
		log.Println("user_service.getByCache HGETALL key = ", key, ", error = ", err)
		return nil
	}
	dataId := comm.GetInt64FromStringMap(dataMap, "Id", 0)
	if dataId <= 0 {
		return nil
	}
	data := &models.LtUser{
		Id:         int(dataId),
		Username:   comm.GetStringFromStringMap(dataMap, "Username", ""),
		Blacktime:  int(comm.GetInt64FromStringMap(dataMap, "Blacktime", 0)),
		Realname:   comm.GetStringFromStringMap(dataMap, "Realname", ""),
		Mobile:     comm.GetStringFromStringMap(dataMap, "Mobile", ""),
		Address:    comm.GetStringFromStringMap(dataMap, "Address", ""),
		SysCreated: int(comm.GetInt64FromStringMap(dataMap, "SysCreated", 0)),
		SysUpdated: int(comm.GetInt64FromStringMap(dataMap, "SysUpdated", 0)),
		SysIp:      comm.GetStringFromStringMap(dataMap, "SysIp", ""),
	}
	return data
}

func (s *userService) setByCache(data *models.LtUser) {
	if data == nil || data.Id <= 0 {
		return
	}
	id := data.Id
	// key := fmt.Sprintf("info_user_%d", id)
	key := fmt.Sprintf(conf.RdsUserCacheKeyPrefix + "%d", id)
	rds := datasource.InstanceCache()

	// params := []interface{}{key}
	// params = append(params, "Id", id)
	params := redis.Args{key}
	params.Add(id)
	if data.Username != "" {
		//params = append(params, "Username", data.Username)
		//params = append(params, "Blacktime", data.Blacktime)
		//params = append(params, "Realname", data.Realname)
		//params = append(params, "Mobile", data.Mobile)
		//params = append(params, "Address", data.Address)
		//params = append(params, "SysCreated", data.SysCreated)
		//params = append(params, "SysUpdated", data.SysUpdated)
		//params = append(params, "SysIp", data.SysIp)

		params.Add(params,"Username",data.Username)
		params.Add(params,"Blacktime",data.Blacktime)
		params.Add(params,"Realname",data.Realname)
		params.Add(params,"Mobile",data.Mobile)
		params.Add(params,"Address",data.Address)
		params.Add(params,"SysCreated",data.SysCreated)
		params.Add(params,"SysUpdated",data.SysUpdated)
		params.Add(params,"SysIp",data.SysIp)
	}
	// _, err := rds.Do("HMSET", params...)
	_, err := rds.Do("HMSET", params)
	if err != nil {
		log.Println("user_service.setByCache HMSET params = ", params, ", error = ", err)
	}
}

// 可以和 setByCache 类似做精确更新,这里直接删除key 简化处理
func (s *userService) updateByCache(data *models.LtUser, columns []string) {
	if data == nil || data.Id <= 0 {
		return
	}
	// key := fmt.Sprintf("info_user_%d", data.Id)
	key := fmt.Sprintf(conf.RdsUserCacheKeyPrefix + "%d", data.Id)
	rds := datasource.InstanceCache()
	rds.Do("DEL", key)
}
