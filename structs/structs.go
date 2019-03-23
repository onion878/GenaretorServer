package structs

import "time"

type Index struct {
	Id      string    `xorm:"not null pk VARCHAR(40)"`
	Name    string    `xorm:"not null VARCHAR(100)"`
	Chapter string    `xorm:"not null VARCHAR(100)"`
	Total   int       `xorm:"not null int"`
	Update  time.Time `xorm:"TIMESTAMP"`
	Index   int       `xorm:"null int"`
	Url     string    `xorm:"null VARCHAR(100)"`
	Created time.Time `xorm:"TIMESTAMP created"`
	Flag    bool
}

type Chapter struct {
	Id      string    `xorm:"not null pk VARCHAR(40)"`
	Pid     string    `xorm:"not null VARCHAR(40)"`
	Name    string    `xorm:"not null VARCHAR(150)"`
	Path    string    `xorm:"text"`
	Num     int       `xorm:"not null int"`
	Created time.Time `xorm:"TIMESTAMP created"`
}

type DetailData struct {
	Id  string `json:"id"`
	Pid string `json:"pid"`
}

type UrlData struct {
	Type    string `json:"type"`
	File    string `json:"file"`
	Label   string `json:"label"`
	Default string `json:"default"`
}

type Cookies struct {
	Id    string `xorm:"not null pk int"`
	Value string `xorm:"text"`
}

type User struct {
	UserName string `xorm:"not null pk VARCHAR(40)"`
	Password string `not null VARCHAR(40)`
}

type Page struct {
	Page  int    `json:"page"`
	Start int    `json:"start"`
	Limit int    `json:"limit"`
	Name  string `json:"name"`
}

type Template struct {
	Id      string    `xorm:"not null pk VARCHAR(40)"`
	Name    string    `not null VARCHAR(200)`
	User    string    `not null VARCHAR(40)`
	Info    string    `text`
	Detail  string    `text`
	Created time.Time `xorm:"TIMESTAMP" time_format:"2006-01-02 15:04:05"`
}

type TemplateDetail struct {
	Id      string    `xorm:"not null pk VARCHAR(40)"`
	Pid     string    `not null VARCHAR(40)`
	Name    string    `not null VARCHAR(200)`
	User    string    `not null VARCHAR(40)`
	Info    string    `text`
	Created time.Time `xorm:"TIMESTAMP" time_format:"2006-01-02 15:04:05"`
}
