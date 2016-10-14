package lib

import "strings"


func Split(path string) []string{
	temp := []byte(path)

	if len(temp) == 0 {
		return []string{}
	}

	if len(temp) == 1 && temp[0] == '/'{
		return []string{}
	}

	var result []string

	last := 0
	current := 0
	length := len(temp)

	for index, value := range temp {
		if value == '/'{
			current = index

			if current > last {
				if current > (last + 1){
					result = append(result, string(temp[last+1: current]))
					last = current
				}

				if current == (last + 1) {
					last = current
				}

			}
		}else{
			if index == (length - 1){
				result = append(result, string(temp[last+1: length]))
			}
		}
	}

	return result
}


func Join(pths []string) string  {
	return "/" + strings.Join(pths, "/")
}


func Clean(path string) string  {
	return Join(Split(path))
}