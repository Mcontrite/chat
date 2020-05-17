package service

import (
	"chat/model"
	"chat/util"
	"errors"
	"fmt"
	"math/rand"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type UserService struct{}

func (s *UserService) Register(mobile, plainpwd, nickname, avatar, sex string) (user model.User, err error) {
	//检测手机号码是否存在,
	tmp := model.User{}
	_, err = DbEngin.Where("mobile=? ", mobile).Get(&tmp)
	if err != nil {
		return tmp, err
	}
	if tmp.Id > 0 {
		return tmp, errors.New("该手机号已经注册")
	}
	//否则拼接插入数据
	tmp.Mobile = mobile
	tmp.Avatar = avatar
	tmp.Nickname = nickname
	tmp.Sex = sex
	tmp.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	tmp.Passwd = util.MakePasswd(plainpwd, tmp.Salt)
	tmp.Createat = time.Now()
	tmp.Token = fmt.Sprintf("%08d", rand.Int31())
	_, err = DbEngin.InsertOne(&tmp)
	//前端恶意插入特殊字符?
	//数据库连接操作失败?
	return tmp, err
}

func (s *UserService) Login(mobile, plainpwd string) (user model.User, err error) {
	tmp := model.User{}
	_, err = DbEngin.Where("mobile = ?", mobile).Get(&tmp)
	if tmp.Id == 0 {
		return tmp, errors.New("该用户不存在")
	}
	if !util.ValidatePasswd(plainpwd, tmp.Salt, tmp.Passwd) {
		return tmp, errors.New("密码不正确")
	}
	//刷新token
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	tmp.Token = token
	_, err = DbEngin.Where(" id = ?", tmp.Id).Cols("token").Update(&tmp)
	return tmp, err
}

//查找某个用户
func (s *UserService) Find(userId int64) (user model.User) {
	tmp := model.User{}
	DbEngin.Where("id = ?", userId).Get(&tmp)
	return tmp
}
