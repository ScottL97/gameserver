package controllers

import (
	"errors"
	"fmt"
	"gameserver/models"
	"sync"

	uuid "github.com/satori/go.uuid"
)

type UserManager struct {
	loginUsers      map[string]string // 在线用户
	loginUsersMutex sync.Mutex
}

var (
	umInstanceOnce sync.Once
	umInstance     *UserManager
)

func init() {
	getUserManager().loginUsers = make(map[string]string)

}

func getUserManager() *UserManager {
	umInstanceOnce.Do(func() {
		umInstance = new(UserManager)
	})
	return umInstance
}

// 登录，生成id，保存在loginUsers中，返回生成的id，客户端的请求都需要包含该id
func loginUser(user models.User) (string, error) {
	// todo: 数据库检查密码
	if user.Password != "password" {
		return "", errors.New("wrong password")
	}
	id := genId()
	getUserManager().loginUsersMutex.Lock()
	defer getUserManager().loginUsersMutex.Unlock()
	getUserManager().loginUsers[id] = user.UserName

	return id, nil
}

func genId() string {
	// 生成id，保存在loginUsers(id: username)
	id := uuid.NewV4().String()
	return id
}

func deleteUser(id string) {
	fmt.Println("delete:", id)
	getUserManager().loginUsersMutex.Lock()
	defer getUserManager().loginUsersMutex.Unlock()
	delete(getUserManager().loginUsers, id)
}

func mapToSlice(m map[string]string) []string {
	s := make([]string, 0, len(m))
	for _, v := range m {
		s = append(s, v)
	}
	return s
}
