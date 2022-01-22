package demo

import (
	"errors"
	"net/http"
	"testing"
)

type mockContext struct {
	channel  string
	code     int
	response map[string]string
}

func (m *mockContext) Demo() (Demo, error) {
	return Demo{
		Datas: m.channel,
	}, nil
}

//c *mockContext เพื่อต้องการจะให้ เก็บค่า
func (c *mockContext) JSON(code int, v interface{}) {
	c.code = code
	c.response = v.(map[string]string)
}
func TestDemo(t *testing.T) {
	handlers := &Handler{ // กำหนด Handler
		channel: "Online",
	}

	c := &mockContext{channel: "Offline"} // กำหนด channel ที่ไม่อยากได้
	handlers.Demo(c)

	want := "Offline is not accepted"
	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}

}

func TestDemo2(t *testing.T) {
	handlers := &Handler{
		channel: "Offline",
	}

	c := &mockContext{channel: "Online"}
	handlers.Demo(c)

	want := "Online is not accepted"
	if want != c.response["message"] {
		t.Errorf("%q is expected but got %q\n", want, c.response["message"])
	}
}

type mockContextBadRequest struct {
	code     int
	response map[string]string
}

func (c mockContextBadRequest) Demo() (Demo, error) {
	return Demo{}, errors.New("went wrong")
}

func (c *mockContextBadRequest) JSON(code int, v interface{}) {
	c.code = code
	c.response = v.(map[string]string)
}
func TestDemoBadRequest(t *testing.T) {
	handlers := &Handler{}

	c := &mockContextBadRequest{}
	handlers.Demo(c)

	want := http.StatusBadRequest

	if want != c.code {
		t.Errorf("%d status code is expected but got %d\n", want, c.code)
	}
}

type mockContextBadRequestWithChannel struct {
	channel   string
	jsonCount int
}

func (c *mockContextBadRequestWithChannel) Demo() (Demo, error) {
	return Demo{Datas: c.channel}, errors.New("went to wrong channel")
}

func (c *mockContextBadRequestWithChannel) JSON(code int, v interface{}) {
	c.jsonCount++
}
func TestOnlyCallJSONOneTime(t *testing.T) {
	handlers := &Handler{
		channel: "Offline",
	}

	c := &mockContextBadRequestWithChannel{}
	handlers.Demo(c)

	want := 1

	if want != c.jsonCount {
		t.Errorf("it should call one time %d\n", c.jsonCount)
	}
}

type spyStore struct {
	wasCalled bool
}

func (s *spyStore) Save(Demo) error {
	s.wasCalled = true
	return nil
}
func TestDemoSaved(t *testing.T) {
	spy := &spyStore{}
	handlers := &Handler{
		channel: "Online",
		store:   spy,
	}

	c := &mockContext{channel: "Online"}
	handlers.Demo(c)
	want := true
	if want != spy.wasCalled {
		t.Errorf("is should store data")
	}
}

type failStore struct{}

func (failStore) Save(Demo) error {
	return errors.New("Error")
}
func TestDemoFailSaved(t *testing.T) {
	store := &failStore{}
	handlers := &Handler{
		channel: "Online",
		store:   store,
	}

	c := &mockContext{channel: "Online"}
	handlers.Demo(c)

	want := http.StatusInternalServerError

	if want != c.code {
		t.Errorf("%d is expected but got %d\n", want, c.code)
	}
}

func TestOk(t *testing.T) {
	store := &spyStore{}
	handlers := &Handler{
		channel: "Online",
		store:   store,
	}

	c := &mockContext{channel: "Online"}
	handlers.Demo(c)

	if _, ok := c.response["message"]; !ok {
		t.Errorf("message key is expected")
	}

}
