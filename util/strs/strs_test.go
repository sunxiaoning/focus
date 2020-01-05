package strutil

import (
	"fmt"
	"testing"
)

func TestStrs(t *testing.T) {
	if IsValidMoney("0") {
		t.Error("0 none Pass!")
	}
	if IsValidMoney("0.00") {
		t.Error("0.00 none Pass!")
	}
	if IsValidMoney("-1") {
		t.Error("0.00 none Pass!")
	}
	if IsValidMoney("aba") {
		t.Error("aba none Pass!")
	}
	if IsValidMoney("1") {
		t.Log("1 Pass!")
	}
	if IsValidMoney("1.00") {
		t.Log("1.00 pass!")
	}

}

func TestStrLen(t *testing.T) {
	fmt.Println(len(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA+BpAoySuUZ93iOVXYyry
W+a141KdU+ZBlyu4B/ZT5Jl8VBg5f89VeeSYRqT7Xj8JL+bv+HTnaPjfzUc0z6Nr
sJswvlaSkrO8NE/bWajUdBizuPQXg/BKJnFiybFXQEkbiunATEFQNEuGWQuv5dc1
doOnh1MJegorC5ZxGXHktU1wRABrGK/2DQvpw3pYp42O1pTfD21zQgap0fjQ7IIn
Jzco0CU0y6ez3BK4Aq0lK+jD8cVIuyiGSjufclY/wsZogACgtQeAZDJ0O90aACFa
uuK6MleV3oP5YfXjO2kZuBpLwl/t6gW8LmemN6JJFMq6AwcbQxcVSklGXymOtenm
GwIDAQAB
-----END PUBLIC KEY-----`))

	fmt.Println(len(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQD4GkCjJK5Rn3eI
5VdjKvJb5rXjUp1T5kGXK7gH9lPkmXxUGDl/z1V55JhGpPtePwkv5u/4dOdo+N/N
RzTPo2uwmzC+VpKSs7w0T9tZqNR0GLO49BeD8EomcWLJsVdASRuK6cBMQVA0S4ZZ
C6/l1zV2g6eHUwl6CisLlnEZceS1TXBEAGsYr/YNC+nDelinjY7WlN8PbXNCBqnR
+NDsgicnNyjQJTTLp7PcErgCrSUr6MPxxUi7KIZKO59yVj/CxmiAAKC1B4BkMnQ7
3RoAIVq64royV5Xeg/lh9eM7aRm4GkvCX+3qBbwuZ6Y3okkUyroDBxtDFxVKSUZf
KY616eYbAgMBAAECggEAN3bO+mnJ2o92zpDOv2mrcqYaBW7Doyz3fs8UPhtwV8uE
QtyDhjIYnr5e2HQrib8305CiFv4zeYEhryd7A+w5t+qJtBNwgwFRUrSDigC3NhkL
nI5c7275dKymdAaERefSE3T8O0/imT5FE4UMVqqM1ijKe/MxTCoXw2hnclPG2Ey9
2s3Y7weJvPgP7m8B6kfNFI6YadwDJQynEcR23kW98sIACeVibdeEcmcEna2w3rXs
/vH8f3t0N92oF+WcxEN0+iBi6pHIGKYzJQv54mrCHKNQJljKcl5LVTNlu6/sjbbb
GisKN2eHXWxYWrAy4sM9gjzpeR4YMYlA3iP4NvtimQKBgQD+/m3b96y6BNJSYoc4
VCMWdw3UsaFdS8kHH9OsQZpc98efzFkU+L6ysxV3iE7fn/DSOfmJcPDDAeLxkK4a
YEFkXbklXqTmrSOtlXNO+whD0UlrA1cf6IgRQY21QSgiUQCoAS37sTfpVHTE7o6T
LPg4zF7bmzqbEYL8fQGeBUHdRwKBgQD5FNzFvCGtSNyitcz55bBovu/YZwzyPtzr
Ux3Pk+HushD5SDKU33lu27tpmQdDUhVW2urXf+V4H+7ct6yAs7a0rEsz7Nos61mQ
nlsiXJM8R+B7PVA6xr2wmwVsE2vS8UVpy60yCzQQPE997qS6ahC9VJlU4pikGBMK
lqyyCC/KjQKBgQCsPVgfuRCnJhKbK4qC9fItfoWNId8JkeNcOJdWR4npOkVRH0sc
61iEhIr/jscYLoUQu7Besjcuwdt1qHCxyErjbJtfLqrcVh8/ZS/aLZ7LkFazPjJf
j0Y5wbnisPxXEbgLw6A33uERbsbhLvKHX1zboNCCLjxL+mwr+JRbFNoEiQKBgFBj
EkliyT+it0pwACJapc94Z/HgeEYqUaRFI+bdZFpj76R0T5bKdOd5VQfkknqAoFBy
wL4iEc3uCGoFgU/cMEgpHvA4LcW3gyVwZhs143LeA63igOUnRQsdTOevcOoYYf2d
9VykTv46aLFM9q8PEi34gq/pnbe/6U1OiZe/mqT1AoGBAM1196wCvnPmad3zlyzJ
+nexVRPKgMNvfSRO5loGCz98/bZJz7ySW6eIxaRCSbR7xNWcqnPkiD4QLskYtmtT
+ABzNFvZGrv3StZnwB6i+ouOu9eRzTstUANmOyPbb3TH7xPKevPm8HCBr7474iA4
cs0f0DWM1v971EU7xtjb/Gck
-----END PRIVATE KEY-----`))
}
