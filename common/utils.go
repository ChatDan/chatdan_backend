package common

import "github.com/jinzhu/copier"

func MustConvert[T any](model any) T {
	var result T
	err := copier.Copy(&result, model)
	if err != nil {
		panic(err)
	}
	return result
}
