package lib

import "testing"

func Test_EncodeUrl(t *testing.T){

}

func Test_EncodeUrlComponent(t *testing.T){
	var url = "/path/?query=你哈"

	println(EncodeUrlComponent(url))
}