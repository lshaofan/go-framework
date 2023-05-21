package utils

import (
	"errors"
	"reflect"
)

// StructAssign 将两个不同的结构体中的字段进行映射赋值，目前只支持基本类型的字段赋值
// 这个方法可以将结构体A中与结构体B相同的字段映射赋值给结构体B。它将遍历结构体A的所有字段，查找和结构体B中名称和类型完全相同的字段，并将结构体A中的值赋给结构体B中的字段。
// 请注意，此方法只能用于基本类型的字段赋值。如果结构体中包含嵌套结构体或者切片、映射等非基本类型的字段，将会报错。
// 参数 a: 为source结构体，b: 为target结构体
// StructAssign函数用于将一个结构体的值赋值到另一个结构体中
func StructAssign(sourceStruct interface{}, targetStruct interface{}) error {
	// 获取传入的结构体的 reflect.Value 类型
	sourceValue := reflect.ValueOf(sourceStruct)
	targetValue := reflect.ValueOf(targetStruct)
	// 判断sourceStruct是否是结构体类型，如果不是，返回错误  targetStruct是指针类型的结构体
	if sourceValue.Kind() != reflect.Struct || targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		return errors.New("sourceStruct必须是结构体类型，targetStruct必须是结构体指针类型")
	}
	// 获取第二个结构体的 reflect.Value 类型
	targetValue = targetValue.Elem()

	// 获取第一个结构体的类型
	sourceType := reflect.TypeOf(sourceStruct)

	// 遍历第一个结构体的所有字段
	for i := 0; i < sourceType.NumField(); i++ {
		// 获取第一个结构体的字段
		sourceField := sourceType.Field(i)

		// 在第二个结构体中查找与第一个结构体该字段同名的字段 第二个结构体的类型是指针类型
		targetField, ok := targetValue.Type().FieldByName(sourceField.Name)

		// 如果没有找到，跳过该字段
		if !ok {
			continue
		}
		// 判断改字段是否是结构体
		if targetField.Type.Kind() == reflect.Struct {
			//	 第一个结构体的与第二个结构体是否相等
			if sourceField.Type != targetField.Type {
				continue
			}

		}
		// 判断改字段是否是切片
		if targetField.Type.Kind() == reflect.Slice {
			continue
		}

		var targetTag string
		// 获取该字段的 json 标签
		sourceTag := sourceField.Tag.Get("json")

		// 获取该字段的 json 标签
		targetTag = targetField.Tag.Get("json")
		// 如果没有 json 标签，使用字段名作为标签
		if targetTag == "" || sourceTag == "" {
			sourceTag = sourceField.Name
			targetTag = sourceTag
		}
		// 如果两个结构体的对应字段类型相同，则将第一个结构体该字段的值赋值给第二个结构体的该字段
		if sourceTag == targetTag && sourceField.Type.Kind() == targetField.Type.Kind() {
			targetValue.FieldByName(sourceField.Name).Set(sourceValue.Field(i))
		}
	}

	// 如果成功，返回 nil
	return nil
}
