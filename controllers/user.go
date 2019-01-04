package controllers

import (
	"fmt"
	"text/template"
	"webim/models"

	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (c *UserController) PageIndex() {
	c.TplName = "index.html"
}

func (c *UserController) PageSignin() {
	c.TplName = "signin.html"
}

func (c *UserController) PageSignup() {
	c.TplName = "signup.html"
}

func (c *UserController) PageChat() {
	username := c.GetSession("username")
	password := c.GetSession("password")
	fmt.Println("session:", username, " ", password)
	if username == nil || password == nil {
		c.Ctx.WriteString("<script>alert(\"请先登陆\");window.location.href = \"/signin\";</script>")
	} else {
		c.TplName = "main.html"
		c.Data["UserName"] = username
	}
}

func (c *UserController) Signin() {
	username := c.GetString("username")
	password := c.GetString("password")
	username = template.HTMLEscapeString(username)
	password = template.HTMLEscapeString(password)
	user := &models.User{
		Username: username,
		Password: password,
	}
	type result struct {
		Val string
	}
	var res result
	err := user.ReadDB()
	fmt.Println(username + " " + password)
	if err == nil {
		if c.GetSession("username") != username {
			//fmt.Println(c.GetSession("username"))
			c.SetSession("username", username)
			//fmt.Println(c.GetSession("username"))
		}

		if c.GetSession("password") != password {
			c.SetSession("password", password)
			//fmt.Println(c.GetSession("password"))
		}
		res.Val = username
	}
	c.Data["json"] = &res
	c.ServeJSON()
}

func (c *UserController) Signup() {
	username := c.GetString("username")
	password := c.GetString("password")
	username = template.HTMLEscapeString(username)
	password = template.HTMLEscapeString(password)
	user := &models.User{
		Username: username,
		Password: password,
	}
	type result struct {
		Val string
	}
	println("username:", username)
	err := user.Create()
	var res result
	if err == nil {
		res.Val = username
	}
	c.Data["json"] = &res
	c.ServeJSON()
}
