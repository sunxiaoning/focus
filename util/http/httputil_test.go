package httputil

import (
	"fmt"
	userconsts "focus/types/consts/user"
	"testing"
	"time"
)

func TestPostJsonWithHeader(t *testing.T) {
	headers := make(map[string]string)
	headers[userconsts.AccessToken] = "JRCASt7GYl0d5g5OAKFgiA=="
	fmt.Println(PostJsonWithHeader("http://localhost:7001/api/v1/ppay/notify", headers, "{\"test\": 1}", time.Second*10))
}
