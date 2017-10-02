package goinjection

import (
	"testing"
)

type testObject struct {
	DummyValue    string `injectValue:test`
	DummyIntValue int    `injectValue:testInt`
	DummyIntBool  bool   `injectValue:testBool`
}

type testSecondObject struct {
	DummyValue testObject `inject:`
}

func TestBasicInjection(t *testing.T) {

	app := NewApplication()

	app.AddDefaultValue("test", "123")
	app.AddDefaultValue("testInt", 1)
	app.AddDefaultValue("testBool", true)

	test := &testObject{}

	app.AddService(test)

	err := app.Wire()

	if err != nil {
		t.Fail()
	}

	if test.DummyValue != "123" {
		t.Fail()
	}

	if test.DummyIntValue != 1 {
		t.Fail()
	}

	if !test.DummyIntBool {
		t.Fail()
	}

}

func TestBasicInjectionWrongType(t *testing.T) {

	app := NewApplication()

	app.AddDefaultValue("test", "123")
	app.AddDefaultValue("testInt", "Tree")
	app.AddDefaultValue("testBool", true)

	test := &testObject{}

	app.AddService(test)

	err := app.Wire()

	if err == nil { // wrong type should cause error!
		t.Fail()
	}

}

func TestBasicObjectInjection(t *testing.T) {

	app := NewApplication()

	app.AddDefaultValue("test", "123")

	test2 := &testSecondObject{}

	test := &testObject{}

	app.AddService(test)
	app.AddService(test2)

	err := app.Wire()

	if err == nil {
		t.Fail()
	}

}

func TestNonePointerObjectInjection(t *testing.T) {

	app := NewApplication()

	app.AddDefaultValue("test", "123")

	test2 := testSecondObject{}

	test := testObject{}

	app.AddService(test) // passing none pointers will fail
	app.AddService(test2)

	err := app.Wire()

	if err != nil {
		t.Fail()
	}

}
