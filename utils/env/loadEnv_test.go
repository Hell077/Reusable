package env

import "testing"

func TestLoadEnv_Error(t *testing.T) {
	err := LoadEnv("config/.env")
	if err == nil {
		t.Error(err)
	}
	t.Log("Success, function is not working if path is wrong")
}

func TestLoadEnv(t *testing.T) {
	err := LoadEnv("../../config/.env")
	if err != nil {
		t.Error(err)
	}
	t.Log("Success load env")
}
