package keyvalue

import "strconv"

func parseInt64(str string) (int64, error) {
	if str == "" {
		return 0, nil
	}
	return strconv.ParseInt(str, 10, 64)
}
