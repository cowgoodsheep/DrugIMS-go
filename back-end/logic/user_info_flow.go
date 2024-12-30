package logic

import (
	"drugims/middleware"
	"drugims/model"
	"errors"
)

const (
	MaxUsernameSize    = 64
	MaxDescriptionSize = 256
)

// UserInfoFlow 用户信息流
type UserInfoFlow struct {
	username  string
	telephone string
	password  string
	UserId    int32  `json:"user_id"`
	Token     string `json:"token"`
}

// 注册用户并得到token和id
func RegisterUser(telephone, username, password string) (*UserInfoFlow, error) {
	return NewUserInfoFlow(telephone, username, password).RegisterUserDo()
}

func NewUserInfoFlow(telephone, username, password string) *UserInfoFlow {
	return &UserInfoFlow{telephone: telephone, username: username, password: password}
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
	if u.username == "" {
		return errors.New("用户名为空")
	}
	if len(u.username) > MaxUsernameSize {
		return errors.New("用户名长度超出限制")
	}
	if len(u.telephone) != 11 {
		return errors.New("手机号码长度不为11位")
	}
	for _, v := range u.telephone {
		if v < '0' || v > '9' {
			return errors.New("手机号码存在非数字字符")
		}
	}
	//判断手机号是否已被注册
	if model.IsUserExistByTelephone(u.telephone) {
		return errors.New("该手机号已被注册")
	}

	if u.password == "" {
		return errors.New("密码为空")
	}
	return nil
}

func (u *UserInfoFlow) registerUser() error {
	userinfo := &model.UserInfo{
		UserName:  u.username,
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
func LoginUser(telephone, password string) (*UserInfoFlow, error) {
	return NewGetUserFlow(telephone, password).loginUserDo()
}

func NewGetUserFlow(telephone, password string) *UserInfoFlow {
	return &UserInfoFlow{telephone: telephone, password: password}
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
	if u.telephone == "" {
		return errors.New("手机号为空")
	}
	if len(u.telephone) != 11 {
		return errors.New("手机号码长度不为11位")
	}
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
	//从数据库中寻找用户
	userInfo, err := model.QueryUserByTelephone(u.telephone)
	if err != nil {
		return err
	}
	if userInfo.Password != u.password {
		return errors.New("密码错误")
	}
	//根据用户手机号生成token
	token, err := middleware.MakeToken(u.telephone)
	if err != nil {
		return err
	}
	u.username = userInfo.UserName
	u.UserId = userInfo.UserId
	u.Token = token
	return nil
}
