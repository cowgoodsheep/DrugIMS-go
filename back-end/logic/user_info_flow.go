package logic

import (
	"drugims/middleware"
	"drugims/model"
	"errors"
)

const (
	MaxuserNameSize    = 64
	MaxDescriptionSize = 256
)

// UserInfoFlow 用户信息流
type UserInfoFlow struct {
	UserInfo *model.UserInfo
	UserId   int32  `json:"user_id"`
	Token    string `json:"token"`
}

// 注册用户
func RegisterUser(userInfo *model.UserInfo) (*UserInfoFlow, error) {
	return NewUserInfoFlow(userInfo).registerUserDo()
}

// 用户登录
func LoginUser(userInfo *model.UserInfo) (*UserInfoFlow, error) {
	return NewUserInfoFlow(userInfo).loginUserDo()
}

// 用户修改信息
func UpdateUser(userInfo *model.UserInfo) (*UserInfoFlow, error) {
	return NewUserInfoFlow(userInfo).updateUserDo()
}

func NewUserInfoFlow(userInfo *model.UserInfo) *UserInfoFlow {
	return &UserInfoFlow{UserInfo: userInfo}
}

// 注册
func (u *UserInfoFlow) registerUserDo() (*UserInfoFlow, error) {
	// 对注册信息进行合法性验证
	if err := u.registerUserCheck(); err != nil {
		return nil, err
	}
	// 注册操作
	if err := u.registerUser(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserInfoFlow) registerUserCheck() error {
	if u.UserInfo.UserName == "" {
		return errors.New("用户名为空")
	}
	if len(u.UserInfo.UserName) > MaxuserNameSize {
		return errors.New("用户名长度超出限制")
	}
	// if len(u.telephone) != 11 { todo
	// 	return errors.New("手机号码长度不为11位")
	// }
	for _, v := range u.UserInfo.Telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	//判断手机号是否已被注册
	if model.IsUserExistByTelephone(u.UserInfo.Telephone) {
		return errors.New("该手机号已被注册")
	}
	//判断用户名是否已被注册
	if model.IsUserExistByUserName(u.UserInfo.UserName) {
		return errors.New("该用户名已被注册")
	}

	if u.UserInfo.Password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserInfoFlow) registerUser() error {
	userinfo := &model.UserInfo{
		UserName:  u.UserInfo.UserName,
		Telephone: u.UserInfo.Telephone,
		Password:  u.UserInfo.Password,
		Role:      u.UserInfo.Role,
		Status:    1,
	}
	//上传数据库
	if err := model.CreateUser(userinfo); err != nil {
		return err
	}

	//根据用户手机号生成token
	token, err := middleware.MakeToken(u.UserInfo.Telephone)
	if err != nil {
		return err
	}
	u.Token = token
	u.UserId = userinfo.UserId
	return nil
}

// 登录
func (u *UserInfoFlow) loginUserDo() (*UserInfoFlow, error) {
	//对登录信息进行合法性验证
	if err := u.loginUserCheck(); err != nil {
		return nil, err
	}
	//从数据库中得到登录用户的信息
	if err := u.loginUser(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserInfoFlow) loginUserCheck() error {
	// if len(u.telephone) != 11 { todo
	// 	return errors.New("手机号码长度不为11位")
	// }
	for _, v := range u.UserInfo.Telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	if u.UserInfo.Password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserInfoFlow) loginUser() error {
	// 从数据库中寻找用户
	userInfo := model.QueryUserByTelephone(u.UserInfo.Telephone)
	if userInfo.UserId == 0 { // 手机号找不到就用用户名
		userInfo = model.QueryUserByUserName(u.UserInfo.UserName)
	}
	if userInfo.UserId == 0 { // 都找不到就报错
		return errors.New("找不到该用户")
	}
	if userInfo.Password != u.UserInfo.Password {
		return errors.New("密码错误")
	}
	// 根据用户手机号生成token
	token, err := middleware.MakeToken(u.UserInfo.Telephone)
	if err != nil {
		return err
	}
	u.UserInfo = userInfo
	u.Token = token
	return nil
}

// 更新
func (u *UserInfoFlow) updateUserDo() (*UserInfoFlow, error) {
	//对信息进行合法性验证
	if err := u.updateUserCheck(); err != nil {
		return nil, err
	}
	//从数据库中得到登录用户的信息
	if err := u.updateUser(); err != nil {
		return nil, err
	}
	return u, nil
}

func (u *UserInfoFlow) updateUserCheck() error {
	if len(u.UserInfo.UserName) > MaxuserNameSize {
		return errors.New("用户名长度超出限制")
	}
	// 判断用户名是否已被使用
	if model.IsUserExistByUserName(u.UserInfo.UserName) {
		return errors.New("该用户名已被使用")
	}
	// 对密码进行加密
	if u.UserInfo.Password != "" {
		if password, err := middleware.SHAMiddleWare(u.UserInfo.Password); err != nil {
			return err
		} else {
			u.UserInfo.Password = password
		}
	}
	return nil
}

func (u *UserInfoFlow) updateUser() error {
	tu := model.QueryUserByUserId(u.UserInfo.UserId)
	if u.UserInfo.UserName == "" {
		u.UserInfo.UserName = tu.UserName
	}
	if u.UserInfo.Password == "" {
		u.UserInfo.Password = tu.Password
	}
	if u.UserInfo.Telephone == "" {
		u.UserInfo.Telephone = tu.Telephone
	}
	if u.UserInfo.Address == "" {
		u.UserInfo.Address = tu.Address
	}
	if u.UserInfo.Role == "" {
		u.UserInfo.Role = tu.Role
	}
	model.UpdateUserInfo(u.UserInfo.UserId, u.UserInfo)
	return nil
}
