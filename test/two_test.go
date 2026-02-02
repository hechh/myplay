package entity

import (
	"fmt"
	"myplay/common/pb"
	"reflect"
	"testing"
)

func Two[Actor any, R any, T any, B any, C any, F any](f any) {
	aType := reflect.TypeOf((*Actor)(nil)).Elem()
	fmt.Println("---->", aType.String(), aType.Kind().String())
	rType := reflect.TypeOf((*R)(nil)).Elem()
	fmt.Println("---->", rType.String(), rType.Kind().String(), rType.Elem().String())
	tType := reflect.TypeOf((*T)(nil)).Elem()
	fmt.Println("---->", tType.String(), tType.Kind().String())
	a := new(&T)
	bType := reflect.TypeOf((*B)(nil)).Elem()
	fmt.Println("---->", bType.String(), bType.Kind().String(), bType.Elem().String())
	cType := reflect.TypeOf((*C)(nil)).Elem()
	fmt.Println("---->", cType.String(), "===>", cType.Kind().String())
	fType := reflect.TypeOf((*F)(nil)).Elem()
	fmt.Println("---->", fType.String(), "===>", fType.Kind().String())
	ff, ok := f.(func())
	if ok {
		ff()
	}
}

func TestHandler(t *testing.T) {
	Two[int, *uint64, *pb.AuthReq, []byte, interface{}, func()](func() {
		t.Log("------")
	})
}

/*
type TwoHandler[Actor any, R any, T any] struct {
	framework.ISerialize
	crc32 uint32
	name  string
	f     func(*Actor, framework.IContext, *R, *T) error
}

func NewTwoHandler[Actor any, R any, T any](f func(*Actor, framework.IContext, *R, *T) error) *TwoHandler[Actor, R, T] {
	tType := reflect.TypeOf((*R)(nil)).Elem() // 获取V2的类型反射对象
	rTYPE := reflect.TypeOf((*T)(nil)).Elem() // 获取V2的类型反射对象

	tType.Kind()

	name := framework.ParseActorFunc(reflect.ValueOf(f))
	return &TwoHandler[Actor, R, T]{
		name:  name,
		crc32: framework.GetCrc32(name),
		f:     f,
	}
}
*/
