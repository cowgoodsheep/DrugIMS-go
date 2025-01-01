package model

import (
	"drugims/dao"
	"errors"
)

// UserInfo Model
type UserInfo struct {
	UserId      int32  `json:"user_id"`     // 用户id
	UserName    string `json:"user_name"`   // 用户名
	Password    string `json:"password"`    // 密码
	Telephone   string `json:"telephone"`   // 手机号
	Description string `json:"description"` // 描述
	Avatar      string `json:"avatar"`      // 头像
	Address     string `json:"address"`     // 地址
	Role        string `json:"role"`        // 用户角色，1:管理员;2:客户;3:供应商
	Status      int8   `json:"status"`      // 用户状态，1:正常;2:禁用
}

// 指定User结构体迁移表user
func (u *UserInfo) TableName() string {
	return "user_info"
}

// CreateUser 创建用户
func CreateUser(u *UserInfo) error {
	if u == nil {
		return errors.New("空指针错误")
	}
	return dao.DB.Debug().Create(u).Error
}

// 根据手机号判断该用户是否存在
func IsUserExistByTelephone(telephone string) bool {
	var u UserInfo
	dao.DB.Where("telephone=?", telephone).Where("status=?", 1).First(&u)
	//如果找不到
	return u.UserId != 0
}

// 根据用户名判断该用户是否存在
func IsUserExistByUserName(userName string) bool {
	var u UserInfo
	dao.DB.Where("user_name=?", userName).Where("status=?", 1).First(&u)
	//如果找不到
	return u.UserId != 0
}

// QueryUserByUserId 用id寻找用户
func QueryUserByUserId(userId int32) *UserInfo {
	var uFind UserInfo
	dao.DB.Where("user_id=?", userId).Where("status=?", 1).First(&uFind)
	return &uFind
}

// QueryUserByTelephone 用手机号寻找用户
func QueryUserByTelephone(telephone string) *UserInfo {
	var uFind UserInfo
	dao.DB.Where("telephone=?", telephone).Where("status=?", 1).First(&uFind)
	return &uFind
}

// QueryUserByUserName 用用户名寻找用户
func QueryUserByUserName(userName string) *UserInfo {
	var uFind UserInfo
	dao.DB.Where("user_name=?", userName).Where("status=?", 1).First(&uFind)
	return &uFind
}

// UpdateUserInfo 更新用户数据
func UpdateUserInfo(userId int32, userInfo *UserInfo) {
	dao.DB.Model(&UserInfo{}).Where("user_id = ?", userId).Updates(userInfo)
}
