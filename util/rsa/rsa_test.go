package rsautil

import (
	"strings"
	"testing"
)

func TestSign(t *testing.T) {
	origin := strings.Join([]string{"2020-01-02 12:00:30", "999", "/ppay/notify", "{\"payChannel\":\"WECHATPAY\",\"payeeAccountId\":3451,\"payer\":\"xiaoning\",\"payAmount\":\"80.00\",\"successTime\":\"2020-01-02 12:00:30\"}"}, ",")
	priKey, err := ParseKeyFromString(`-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCsmRZ+euoLSUH4
QSXacwDUsq+G/lFwu/oEHrIn0RR9CbfPIWg0UwOAq4Cg4dDi+SoTByufm39fMm0w
8zeYU8Bo9YbKI+Ml14serl8EK+q7kIAh2BVb31f5n/qb4v0iOWIhiIm85RHhfCKV
5vj3in0YR3wEb4ZlErpnHQrVSmyyCh10chyN50rFBekhojGKa1ugTu++6OauE6UY
8WF1CdPLR2ioYU70iCWQqmgU/WCmEOpAwaMHpxseKzro8xrJySf+e2djiHBgpkAX
ctgZdLBaTxKv3fFdBLbZlC1A/s7qnmQSvcoQQ9yEASVHF5t1ucTi6Kv4B8+d3cHj
/Do1BJFbAgMBAAECggEAGw68ygM24aIQZ3y/FGnm/XEttznzRnsRjMfLIhbS4W8+
k4gnL9y7tdYtCguclZh+EhGTsyfLZf8fwwa99nFavpmPEe1iVLJfquWnG55O5CyZ
CbU4U/jFb+q2fO2uxUi9q4geH7DBhtnRCjL1YMkjJ0U9mexp4zp2YmfZJOrP83EE
ZAqUfME1LoBWR32QZGNIf+kKIGsfuOPQ/aqSxMHsCsGQ1vT3i10lATvM5bR5AYNY
e4zqWTxfEKEO0Rv8ApWFcaw4x6BJORiGOEXsuifaG72el7NIXaKeX8lXqodkQo2t
cBjdv9vXvwcCFBqnQ6FewApWUwnEEqFQISu2+S/JgQKBgQDjihiwrmQuR9u3DUq1
UKWzf4kyOifIMC83aJkiw0d55L/df1jmwh+IZqO0Qc57Aw/DBzIftuCA5fQAWy02
3ebWajFYzuFHqLMeRQ67RygmCqUe2wsHcrQGbetOOu7aNo1jWabqECUiy5zE92XS
CjHAy59l4ObkGM1OkDzxkCxguwKBgQDCL76fROIY9TpsF1XxfM7BHX6s4yjfAJHn
y90+EC2swij3vJkDrCMiqVVWjAkBS5YCUzDHr0tBI+4V3xU/s9T2WjDHkfw0mDYg
dNkZYYJOHwg7KVirAFFJc8osnoeWxEOiipYOel/RuAZC1NWIrYvYAIjag13wiw9J
Di6cSyxX4QKBgAeG8fomSrodNm9/yRDmchTWCzvWIKrGrUkv9KDpxNuLba1uIQrB
MTnZ62BzLNl06HiTLF1QN20MLl40pfJCtPgy2x8M+Pbd5c8CidI4MGPRxlSW+m5s
pPfxeu9Dk9M0Ksk2lgb4McJM6gq6BGxGWg7+rW85WWoCSAhpTRrQoicxAoGBAIS6
U6XTVGNQwtT9Ak5kS4Gt4lbTka0TW6c/LgLs1fteXtguPbxH3WAks+LLJqCPBIKk
UBQ82cg1gdfLOl/nUCnurabLPsLIQz7d/0Ic5w70oRCnCyceuXDmRwtGdFD153Dj
HvvYY0Qab5Ugmq+oR4ylmOUao4v10MXTfsJvk2ihAoGAf508D7FYrcGvL/9s90OR
F43Bj26CXxZLKcY+dnPVpAzZFDyqf0QxKQcv0BlFEEZqIBYpQaRIOfLfDkP24yzS
rMSfhgrJNZDc5V7HaqKezNPhuvfTbwdGlBYa0kQQ+bSnx+ICEBps3zr9+iKLfq9R
lU00uXQGlVtl390h8iKyW90=
-----END PRIVATE KEY-----`)
	if err != nil {
		t.Error("parse priKey err!", err)
	}

	sign, err := Sign(origin, priKey)
	if err != nil {
		t.Error("sign err!", err)
	}
	t.Log(sign)
}
