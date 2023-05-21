package utils

import (
	"gorm.io/gorm"
	"testing"
)

// 测试 StructAssign
func TestStructAssign(t *testing.T) {
	type C struct {
		A string `json:"a"`
	}
	type C2 struct {
		A string `json:"a"`
	}
	type A struct {
		gorm.Model
		Name     string `json:"name"`
		Age      int    `json:"age"`
		Birthday string `json:"birthday"`
		C        []C    `json:"c"`
	}

	type B struct {
		ID   uint     `json:"id"`
		Name string   `json:"name"`
		Age  int      `json:"age"`
		Tags []string `json:"tags"`
		C    []C2     `json:"c"`
	}

	a := A{
		Model: gorm.Model{ID: 5555},
		Name:  "John", Age: 30, Birthday: "1990-01-01", C: []C{{A: "a"}}}
	b := new(B)
	err := StructAssign(a, b)
	if err != nil {
		t.Errorf("StructAssign返回错误: %v", err)
	}
	err = StructAssign(a.C, &b.C)
	if err != nil {
		t.Errorf("StructAssign返回错误: %v", err)
	}

	t.Log("b:", b)
	c := "John"
	d := 30
	err = StructAssign(c, &d)
	if err == nil {
		t.Error("StructAssign在a和b不是结构体类型的时候没有返回错误")
	}
}
