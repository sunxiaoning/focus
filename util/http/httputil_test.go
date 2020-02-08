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
	fmt.Println(PostJsonWithHeader("http://localhost:7001/api/v1/ppay/notify", headers, "{\"payChannel\":\"ALIPAY\",\"payeeAccountId\":3452,\"payer\":\"xiaoning\",\"payAmount\":\"60.00\",\"successTime\":\"2020-01-02 12:00:30\"}", time.Second*10))
}
