package types

type PayChannel struct {
	Name      string
	Value     string
	AccountId int
}

var ALIPAY = &PayChannel{"支付宝支付", "ALIPAY", 3452}
var WECHATPAY = &PayChannel{"微信支付", "WECHAT", 3451}

func PayChannels() map[string]*PayChannel {
	return map[string]*PayChannel{
		ALIPAY.Value:    ALIPAY,
		WECHATPAY.Value: WECHATPAY,
	}
}
