package model

import (
	"drugims/dao"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// UserInfo Model
type UserInfo struct {
	UserId       int32           `json:"user_id"`                                                      // 用户id
	UserName     string          `json:"user_name"`                                                    // 用户名
	Password     string          `json:"password"`                                                     // 密码
	Telephone    string          `json:"telephone"`                                                    // 手机号
	Description  string          `json:"description"`                                                  // 描述
	Avatar       string          `json:"avatar"`                                                       // 头像
	Address      string          `json:"address"`                                                      // 地址
	Role         string          `json:"role"`                                                         // 用户角色，管理员;客户;供应商
	Balance      decimal.Decimal `json:"balance" gorm:"type:decimal(10,2);column:balance"`             // 余额
	BlockBalance decimal.Decimal `json:"block_balance" gorm:"type:decimal(10,2);column:block_balance"` // 冻结余额
	Status       int8            `json:"status"`                                                       // 用户状态，1:正常;2:禁用
	Token        string          `json:"token"`
	CreateTime   time.Time       `json:"create_time" gorm:"-"`
	UpdateTime   time.Time       `json:"update_time" gorm:"-"`

	Recharge decimal.Decimal `json:"recharge" gorm:"-"` // 充值金额
	Withdraw decimal.Decimal `json:"withdraw" gorm:"-"` // 提现金额
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
	return dao.DB.Model(&UserInfo{}).Create(u).Error
}

// 根据手机号判断该用户是否存在
func IsUserExistByTelephone(telephone string) bool {
	var u UserInfo
	dao.DB.Model(&UserInfo{}).Where("telephone=?", telephone).Where("status=?", 1).First(&u)
	//如果找不到
	return u.UserId != 0
}

// 根据用户名判断该用户是否存在
func IsUserExistByUserName(userName string) bool {
	var u UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_name=?", userName).Where("status=?", 1).First(&u)
	//如果找不到
	return u.UserId != 0
}

// QueryUserByUserId 用id寻找用户
func QueryUserByUserId(userId int32) *UserInfo {
	var uFind UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_id=?", userId).Where("status=?", 1).First(&uFind)
	return &uFind
}

// QueryBlockUserByUserId 用id寻找被拉黑用户
func QueryBlockUserByUserId(userId int32) *UserInfo {
	var uFind UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_id=?", userId).First(&uFind)
	return &uFind
}

// QueryUserByTelephone 用手机号寻找用户
func QueryUserByTelephone(telephone string) *UserInfo {
	var uFind UserInfo
	dao.DB.Model(&UserInfo{}).Where("telephone=?", telephone).Where("status=?", 1).First(&uFind)
	return &uFind
}

// QueryUserByUserName 用用户名寻找用户
func QueryUserByUserName(userName string) *UserInfo {
	var uFind UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_name=?", userName).Where("status=?", 1).First(&uFind)
	return &uFind
}

// UpdateUserInfo 更新用户数据
func UpdateUserInfo(userId int32, userInfo *UserInfo) {
	dao.DB.Model(&UserInfo{}).Where("user_id = ?", userId).Updates(userInfo)
}

// LikeGetUserListByUserName 获取根据用户名称模模糊查询用户列表
func LikeGetUserListByUserName(userName string) []*UserInfo {
	var uListFind []*UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_name LIKE ?", "%"+userName+"%").Find(&uListFind)
	return uListFind
}

// GetUserByUserId 获取根据用户id获取用户
func GetUserByUserId(userId int32) *UserInfo {
	var uFind UserInfo
	dao.DB.Model(&UserInfo{}).Where("user_id=?", userId).First(&uFind)
	if uFind.UserId == 0 {
		return nil
	}
	return &uFind
}

// DeleteUser 删除用户
func DeleteUser(userId int32) {
	dao.DB.Model(&UserInfo{}).Where("user_id = ?", userId).Delete(&UserInfo{})
}
