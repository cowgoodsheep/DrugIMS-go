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
	userName  string
	telephone string
	password  string

	UserInfo *model.UserInfo
	UserId   int32  `json:"user_id"`
	Token    string `json:"token"`
}

// 注册用户并得到token和id
func RegisterUser(telephone, userName, password string) (*UserInfoFlow, error) {
	return NewUserInfoFlow(telephone, userName, password).RegisterUserDo()
}

func NewUserInfoFlow(telephone, userName, password string) *UserInfoFlow {
	return &UserInfoFlow{telephone: telephone, userName: userName, password: password}
}

func (u *UserInfoFlow) RegisterUserDo() (*UserInfoFlow, error) {
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
	if u.userName == "" {
		return errors.New("用户名为空")
	}
	if len(u.userName) > MaxuserNameSize {
		return errors.New("用户名长度超出限制")
	}
	// if len(u.telephone) != 11 { todo
	// 	return errors.New("手机号码长度不为11位")
	// }
	for _, v := range u.telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	//判断手机号是否已被注册
	if model.IsUserExistByTelephone(u.telephone) {
		return errors.New("该手机号已被注册")
	}
	//判断用户名是否已被注册
	if model.IsUserExistByUserName(u.userName) {
		return errors.New("该用户名已被注册")
	}

	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserInfoFlow) registerUser() error {
	userinfo := &model.UserInfo{
		UserName:  u.userName,
		Telephone: u.telephone,
		Password:  u.password,
		Status:    1,
	}
	//上传数据库
	if err := model.CreateUser(userinfo); err != nil {
		return err
	}

	//根据用户手机号生成token
	token, err := middleware.MakeToken(u.telephone)
	if err != nil {
		return err
	}
	u.Token = token
	u.UserId = userinfo.UserId
	return nil
}

// 用户登录并返回token和id
func LoginUser(telephone, userName, password string) (*UserInfoFlow, error) {
	return NewGetUserFlow(telephone, userName, password).loginUserDo()
}

func NewGetUserFlow(telephone, userName, password string) *UserInfoFlow {
	return &UserInfoFlow{telephone: telephone, userName: userName, password: password}
}

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
	for _, v := range u.telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserInfoFlow) loginUser() error {
	// 从数据库中寻找用户
	userInfo := model.QueryUserByTelephone(u.telephone)
	if userInfo.UserId == 0 {
		userInfo = model.QueryUserByUserName(u.userName)
	}
	if userInfo.Password != u.password {
		return errors.New("密码错误")
	}
	// 根据用户手机号生成token
	token, err := middleware.MakeToken(u.telephone)
	if err != nil {
		return err
	}
	u.UserInfo = userInfo
	u.Token = token
	return nil
}
