package mygomock

import (
	"fmt"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestFoo(t *testing.T) {
	// 创建Mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)
	m.EXPECT().Bar(1).Return(1)
	m.EXPECT().Bar(2).Return(2)

	fmt.Println("--1--", m.Bar(2))
	fmt.Println("--2--", m.Bar(1))
}

func TestFoo2(t *testing.T) {
	// 创建Mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)
	gomock.InOrder( // 必须按照指定顺序调用Bar()
		m.EXPECT().Bar(1).Return(1),
		m.EXPECT().Bar(2).Return(2),
	)

	fmt.Println("--1--", m.Bar(2))
	fmt.Println("--2--", m.Bar(1))
}

func TestFoo3(t *testing.T) {
	// 创建Mock控制器
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)
	m.EXPECT().Bar(1).DoAndReturn(func(x int) int {
		return x + 10
	})
	m.EXPECT().Bar(2).DoAndReturn(func(x int) int {
		return x + 10
	})

	fmt.Println("--1--", m.Bar(1))
	fmt.Println("--2--", m.Bar(2))
}

func TestFoo4(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)
	// Do中传入的函数相当于Hook函数
	m.EXPECT().Bar(1).Do(func(x int) int {
		fmt.Println("x: ", x)
		return -1 // 返回值无意义
	}).Return(10)

	fmt.Println("--1--", m.Bar(1))
}

func TestFoo5(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockFoo(ctrl)
	m.EXPECT().Bar(1).Return(1)
	// gomock.Any() 表示匹配任意参数
	// AnyTimes() 表示匹配任意多次
	m.EXPECT().Bar(gomock.Any()).Return(100).AnyTimes()

	fmt.Println("--1--", m.Bar(1))
	fmt.Println("--2--", m.Bar(2))
	fmt.Println("--3--", m.Bar(3))
	fmt.Println("--4--", m.Bar(4))
}

type CustomMatcher struct {
	Value *Car
}

func NewCustomMatcher(value *Car) *CustomMatcher {
	return &CustomMatcher{Value: value}
}

func (m *CustomMatcher) Matches(x interface{}) bool {
	value2 := x.(*Car)
	fmt.Println("value1:", m.Value, "value2:", value2)
	return m.Value.Age == value2.Age
}

func (m *CustomMatcher) String() string {
	return fmt.Sprintf("CustomMatcher(%v)", m.Value)
}

func TestDealer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := NewMockDealer(ctrl)
	m.EXPECT().Evaluate(NewCustomMatcher(&Car{Color: "red", Age: 1})).Return(100).AnyTimes()

	fmt.Println("--1--", m.Evaluate(&Car{Color: "red", Age: 1}))
	fmt.Println("--2--", m.Evaluate(&Car{Color: "blue", Age: 1}))
}
