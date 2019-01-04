package models

import (
	"errors"

	"github.com/astaxie/beego/orm"
)

type User struct {
	Username string `orm:"pk"`
	Password string
}

func (u *User) ReadDB() (err error) {
	o := orm.NewOrm()
	var maps []orm.Params
	num, err := o.QueryTable("User").Values(&maps)
	if num < 1 {
		err = errors.New("此帐号不存在")
		return err
	}
	if err == nil {
		for _, m := range maps {
			if m["Username"] == u.Username && m["Password"] == u.Password {
				return err
			}
		}
		err = errors.New("此帐号不存在")
		return err
	} else {
		return err
	}
}

func (u *User) Create() (err error) {
	o := orm.NewOrm()
	if o.Read(u) == nil {
		err = errors.New("帐号已经存在")
		return err
	}
	_, err = o.Insert(u)
	return err
}
