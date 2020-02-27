package gormkit

import (
	"bytes"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
	"text/template"
)

type Params map[string]interface{}

func ExecScanPlus(db *gorm.DB, sqlTpl string, params *Params) *gorm.DB {
	paramsMap := * params
	// 正则找出，模板分析出中所有的 #{xxxx}
	tplParamsRegexp := regexp.MustCompile(`#{(?P<param>\w*)}`)
	tplParams := tplParamsRegexp.FindAllStringSubmatch(sqlTpl, -1)
	sqlValues := make([]interface{}, 0, len(tplParams))
	for _, tplParam := range tplParams {
		full := tplParam[0]
		short := tplParam[1]
		v, exist := paramsMap[short]
		if !exist {
			panic(full + "：不在params中")
		}
		sqlTpl = strings.Replace(sqlTpl, full, "?", 1)
		sqlValues = append(sqlValues, v)
	}
	// 解析模板
	sql, err := parseSqlTmpl(sqlTpl, params)
	if err != nil {
		panic("模板解析失败")
	}
	return db.Raw(sql, sqlValues...)
}
func parseSqlTmpl(tmpl string, params interface{}) (string, error) {
	t := template.New("sql")
	var err error
	t, err = t.Parse(tmpl)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err = t.Execute(&tpl, params); err != nil {
		return "", err
	}
	sql := tpl.String()
	return sql, nil
}
