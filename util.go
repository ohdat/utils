package utils

import (
	"encoding/json"
	"github.com/spf13/viper"
	"reflect"
	"strconv"
	"strings"
)

func StrToInt(str string) (int, error) {
	if str == "" {
		return 0, nil
	}
	id, err := strconv.Atoi(str)
	return id, err
}
func ArrStringToMapInt32(arr []string) map[int32]string {

	var maps = make(map[int32]string)
	for i, v := range arr {
		maps[int32(i)] = v
	}
	return maps
}

func Str2ArrInt(str string) (arr []int, err error) {
	//字符串分割成数组
	arrStr := strings.Split(str, ",")
	for _, v := range arrStr {
		if v == "" {
			continue
		}
		id, errs := StrToInt(v)
		if errs != nil {
			continue
		}
		arr = append(arr, id)
	}
	return
}

// IsDev 是否是开发环境
func IsDev() bool {
	if viper.GetString("environment") == "development" {
		return true
	}
	return false
}

// NotLogin 是否不需要登录
func NotLogin() bool {
	if IsDev() && viper.GetInt("notlogin") == 1 {
		return true
	}
	return false
}

func Struct2struct(in interface{}, out interface{}) {
	s, _ := json.Marshal(in)
	json.Unmarshal(s, out)
}

// StructAssign
//binding type interface 要修改的结构体
// value type interface 有数据的结构体
func StructAssign(binding interface{}, value interface{}) {
	bVal := reflect.ValueOf(binding).Elem() //获取reflect.Type类型
	vVal := reflect.ValueOf(value).Elem()   //获取reflect.Type类型
	vTypeOfT := vVal.Type()
	for i := 0; i < vVal.NumField(); i++ {
		// 在要修改的结构体中查询有数据结构体中相同属性的字段，有则修改其值
		name := vTypeOfT.Field(i).Name
		if ok := bVal.FieldByName(name).IsValid(); ok {
			//验证类型
			if tok := bVal.FieldByName(name).Type().AssignableTo(vTypeOfT.Field(i).Type); tok {
				bVal.FieldByName(name).Set(reflect.ValueOf(vVal.Field(i).Interface()))
			}
		}
	}
}
