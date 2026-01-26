package test

import (
	"myplay/server/client/mock"
	"testing"
)

func TestAll(t *testing.T) {
	if err := mock.Init("../../../configure/env/develop/config.yaml", 1); err != nil {
		t.Log(err)
		return
	}
	usr1 := uint64(10)
	if err := mock.Login(usr1, 1); err != nil {
		t.Log(err)
		return
	}
	if err := mock.Login(usr1, 2); err != nil {
		t.Log(err)
		return
	}
}

func TestPlayer1(t *testing.T) {
	if err := mock.Init("../../../configure/env/develop/config.yaml", 1); err != nil {
		t.Log(err)
		return
	}
	usr1 := uint64(10)
	if err := mock.Login(usr1, 1); err != nil {
		t.Log(err)
		return
	}
}

func TestPlayer2(t *testing.T) {
	if err := mock.Init("../../../configure/env/develop/config.yaml", 1); err != nil {
		t.Log(err)
		return
	}
	usr2 := uint64(11)
	if err := mock.Login(usr2, 2); err != nil {
		t.Log(err)
		return
	}
}
