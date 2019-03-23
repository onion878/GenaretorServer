package service

import (
	"../structs"
	"../utils"
	"github.com/gin-gonic/gin"
	"time"
)

type Template struct {
}

func (*Template) Create(template structs.Template) string {
	engine := utils.GetCon()
	var templates []structs.Template
	engine.Where("user=? and name=?", template.User, template.Name).Find(&templates)
	if len(templates) > 0 {
		templates[0].Created = time.Now()
		templates[0].Detail = template.Detail
		engine.ID(templates[0].Id).Update(templates[0])
		return templates[0].Id
	}
	template.Created = time.Now()
	total, err := engine.Insert(&template)
	if err != nil || total == 0 {
		return template.Id
	} else {
		return template.Id
	}
}

func (*Template) CreateDetail(detail structs.TemplateDetail) bool {
	engine := utils.GetCon()
	detail.Created = time.Now()
	total, err := engine.Insert(&detail)
	if err != nil || total == 0 {
		return false
	} else {
		return true
	}
}

func (*Template) ListByPage(page structs.Page) interface{} {
	engine := utils.GetCon()
	var template structs.Template
	var templates []structs.Template
	total, _ := engine.Count(&template)
	if len(page.Name) > 0 {
		engine.Where("`name` like ?", "%"+page.Name+"%").Limit(page.Limit, page.Start).OrderBy("`created` desc").Find(&templates)
	} else {
		engine.Limit(page.Limit, page.Start).OrderBy("`created` desc").Find(&templates)
	}
	return gin.H{
		"total": total,
		"rows":  templates,
	}
}

func (*Template) ListDetailByPage(page structs.Page, pid string) interface{} {
	engine := utils.GetCon()
	var template structs.TemplateDetail
	var templates []structs.TemplateDetail
	total, _ := engine.Count(&template)
	if len(page.Name) > 0 {
		engine.Where("pid = ? and `name` like ?", pid, "%"+page.Name+"%").Limit(page.Limit, page.Start).OrderBy("`created` desc").Find(&templates)
	} else {
		engine.Where("pid = ?", pid).Limit(page.Limit, page.Start).OrderBy("`created` desc").Find(&templates)
	}
	return gin.H{
		"total": total,
		"rows":  templates,
	}
}

func (*Template) GetNewestDetail(pid string) structs.TemplateDetail {
	engine := utils.GetCon()
	var templates []structs.TemplateDetail
	engine.Where("pid = ?", pid).Limit(1).OrderBy("`created` desc").Find(&templates)
	return templates[0]
}

func (*Template) CheckNew(pid string, id string) bool {
	engine := utils.GetCon()
	var s []structs.TemplateDetail
	engine.Where("pid = ?", pid).Limit(1).OrderBy("`created` desc").Find(&s)
	println(len(s))
	if s[0].Id == id {
		return false
	} else {
		return true
	}
}
