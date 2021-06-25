package mymock

import (
	"github.com/stretchr/testify/mock"
	"testing"
)

type FakePeople struct {
	mock.Mock
}

func (m *FakePeople) Say(str string) string {
	args := m.Called(str)
	return args.String(0)
}

func TestSomething(t *testing.T) {

	// create an instance of our test object
	testObj := new(FakePeople)

	// setup expectations
	testObj.On("Say", "Jack").Return("hello Jack")
	testObj.On("Say", "Lucy").Return("hello Lucy")

	// call the code we are testing
	expected := "People Say:hello Jack"
	got := BizLogic(testObj, "Jack")
	if expected == got {
		t.Logf("got = %v; expected = %v", got, expected)
	} else {
		t.Errorf("got = %v; expected = %v", got, expected)
	}
}
