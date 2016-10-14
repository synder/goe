package lib

import "testing"

func Test_HtmlEscape(t *testing.T){
	var html = "<html><body><div>my name is 'wiki'</div></body></html>"

	escape := HtmlEscape(html)
	expect := "&lt;html&gt;&lt;body&gt;&lt;div&gt;my name is &#39;wiki&#39;&lt;/div&gt;&lt;/body&gt;&lt;/html&gt;"

	if escape !=  expect{
		t.Errorf("HtmlEscape(html) -> %s but expect is %s", escape, expect)
	}
}