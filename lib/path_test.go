package lib

import (
	"testing"
	"strings"
)

func Test_Split(t *testing.T){
	var urlPath1 = "/name/synder.me"
	var urlPath2 = "//names///show/synder.me/"

	temp1 := Split(urlPath1)
	temp2 := Split(urlPath2)

	if len(temp1) != 2{
		t.Error("split filed length error")
	}

	if len(temp2) != 3{
		t.Error("split filed length error")
	}


	if strings.Join(temp1, "/") != "name/synder.me" {
		t.Errorf("%s != %s", strings.Join(temp1, "/"), "name/synder.me")
	}

	if strings.Join(temp2, "/") != "names/show/synder.me" {
		t.Errorf("%s != %s", strings.Join(temp2, "/"), "names/show/synder.me")
	}
}

func Test_Join(t *testing.T) {
	var urlPath = "//names///show/synder.me/"

	if Join(Split(urlPath)) != "/names/show/synder.me"{
		t.Errorf("%s != %s", Join(Split(urlPath)), "/names/show/synder.me")
	}
}

func Test_Clean(t *testing.T){
	var urlPath = "//names///show/synder.me/"

	if Clean(urlPath) != "/names/show/synder.me" {
		t.Errorf("%s != %s", Clean(urlPath), "/names/show/synder.me")
	}
}