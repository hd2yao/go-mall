package library

import (
	"context"
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/hd2yao/go-mall/common/errcode"
	"github.com/hd2yao/go-mall/common/util"
	"github.com/hd2yao/go-mall/common/util/httptool"
	"github.com/hd2yao/go-mall/logic/do"
	"github.com/hd2yao/go-mall/resources"
)

type WxPayLib struct {
	ctx       context.Context
	payConfig WxPayConfig
}

type WxPayConfig struct {
	AppId           string
	MchId           string
	PrivateSerialNo string
	AesKey          string
	NotifyUrl       string
}

func NewWxPayLib(ctx context.Context, payConfig WxPayConfig) *WxPayLib {
	return &WxPayLib{
		ctx:       ctx,
		payConfig: payConfig,
	}
}

const prePayApiUrl = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"

type PrePayParam struct {
	AppId       string `json:"app_id"`
	MchId       string `json:"mch_id"`       // 商户号 ID
	Description string `json:"description"`  // 商品描述
	OutTradeNo  string `json:"out_trade_no"` // 业务的订单号
	NotifyUrl   string `json:"notify_url"`   // 支付成功后，结果回调通知 url
	Amount      struct {
		Total    int    `json:"total"` // 订单总金额，单位为分
		Currency string `json:"currency"`
	} `json:"amount"`
	Payer struct {
		OpenId string `json:"open_id"`
	} `json:"payer"`
}

// WxPayInvokeInfo 前端用 JSAPI 调起支付的参数信息
// 微信支付文档: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
type WxPayInvokeInfo struct {
	AppId     string `json:"appId"`
	TimeStamp string `json:"timeStamp"`
	NonceStr  string `json:"nonceStr"`
	Package   string `json:"package"`
	SignType  string `json:"signType"`
	PaySign   string `json:"paySign"`
}

type WxPayNotifyResponse struct {
	CreateTime string              `json:"create_time"`
	Resource   WxPayNotifyResource `json:"resource"`
}

type WxPayNotifyResource struct {
	Ciphertext     string `json:"ciphertext"`
	AssociatedData string `json:"associated_data"`
	Nonce          string `json:"nonce"`
}

// WxPayNotifyResourceData 微信支付结果通知中解密后的 resource 数据
type WxPayNotifyResourceData struct {
	TransactionID string `json:"transaction_id"`
	Amount        struct {
		PayerTotal    int    `json:"payer_total"`
		Total         int    `json:"total"`
		Currency      string `json:"currency"`
		PayerCurrency string `json:"payer_currency"`
	} `json:"amount"`
	Mchid       string    `json:"mchid"`
	TradeState  string    `json:"trade_state"`
	BankType    string    `json:"bank_type"`
	SuccessTime time.Time `json:"success_time"`
	Payer       struct {
		Openid string `json:"openid"`
	} `json:"payer"`
	OutTradeNo     string `json:"out_trade_no"`
	AppID          string `json:"AppID"`
	TradeStateDesc string `json:"trade_state_desc"`
	TradeType      string `json:"trade_type"`
	Attach         string `json:"attach"`
}

// CreateOrderPay 创建支付信息
// @param order *do.Order 业务的订单信息
// @param userOpenId string 用户的Openid
// @return payInvokeInfo *WxPayInvokeInfo 前端用于调起微信支付的参数
// @return err error 错误信息
func (wpl *WxPayLib) CreateOrderPay(order *do.Order, userOpenId string) (payInvokeInfo *WxPayInvokeInfo, err error) {
	// 创建预支付单
	// 微信支付文档：https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_1.shtml

	// 生成支付描述信息，通常是订单的第一件商品名称
	payDescription := fmt.Sprintf("GOMALL 商场购买 %s 等商品", order.Items[0].CommodityName)

	// 构造微信支付预支付请求参数
	prePayPram := &PrePayParam{
		AppId:       wpl.payConfig.AppId,
		MchId:       wpl.payConfig.MchId,
		Description: payDescription,
		OutTradeNo:  order.OrderNo,
		NotifyUrl:   wpl.payConfig.NotifyUrl,
	}

	// 设置支付金额信息
	prePayPram.Amount.Total = order.PayMoney
	prePayPram.Amount.Currency = "CNY"
	prePayPram.Payer.OpenId = userOpenId
	// 将预支付参数转换为 JSON 格式
	reqBody, _ := json.Marshal(prePayPram)

	// 获取微信支付 API 调用凭证（签名 token）
	token, err := wpl.getToken(http.MethodPost, string(reqBody), prePayApiUrl)
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}

	// 发送 HTTP POST 请求到微信支付 API，创建预支付订单
	_, replyBody, err := httptool.Post(wpl.ctx, prePayApiUrl, reqBody, httptool.WithHeaders(map[string]string{
		"Authorization": "WECHATPAY2-SHA256-RSA2048 " + token,
	}))
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}

	// 解析微信支付接口返回的预支付 ID
	prePayReply := struct {
		PrePayId string `json:"prepay_id"`
	}{}
	if err = json.Unmarshal(replyBody, &prePayReply); err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}

	// 生成前端调起支付需要的参数
	payInvokeInfo, err = wpl.genPayInvokeInfo(prePayReply.PrePayId)
	if err != nil {
		err = errcode.Wrap("WxPayLibCreatePrePayError", err)
		return
	}
	return payInvokeInfo, nil
}

// genToken 生成微信支付的请求签名
// 文档：https://pay.weixin.qq.com/docs/merchant/development/interface-rules/signature-generation.html
func (wpl *WxPayLib) getToken(httpMethod, requestBody, wxApiUrl string) (token string, err error) {
	// wxApiUrl = "https://api.mch.weixin.qq.com/v3/pay/transactions/jsapi"
	// url.Parse(wxApiUrl) 会将 URL 拆解为结构化对象
	// &url.URL{
	//    Scheme: "https",
	//    Host:   "api.mch.weixin.qq.com",
	//    Path:   "/v3/pay/transactions/jsapi",
	//}
	urlPart, err := url.Parse(wxApiUrl)
	if err != nil {
		return token, err
	}
	// 获取 URL 路径部分: /v3/pay/transactions/jsapi
	canonicalUrl := urlPart.RequestURI()
	timestamp := time.Now().Unix()
	nonce := util.RandomString(32)
	// 构造签名原始字符串，符合微信支付 API 规范
	message := fmt.Sprintf("%s\n%s\n%d\n%s\n%s\n", httpMethod, canonicalUrl, timestamp, nonce, requestBody)

	// 商户私有证书放在了 resources 目录下
	// 读取私钥文件
	pemFileReader, err := resources.LoadResourceFile("wxpay.private.pem")
	if err != nil {
		return token, err
	}
	// 读取私钥内容
	privateKey, err := ioutil.ReadAll(pemFileReader)
	if err != nil {
		return token, err
	}

	// 计算 SHA256 哈希值
	sha256MsgBytes := util.SHA256HashBytes(message)
	// 使用 RSA 私钥进行 PKCS1v15 签名
	signBytes, err := util.RsaSignPKCS1v15(sha256MsgBytes, privateKey, crypto.SHA256)
	if err != nil {
		return token, err
	}
	// 对签名结果进行 Base64 编码
	sign := base64.StdEncoding.EncodeToString(signBytes)

	// 生成最终的 Authorization 认证信息（符合微信支付 API 格式）
	token = fmt.Sprintf("mchid=\"%s\",nonce_str=\"%s\",timestamp=\"%d\",serial_no=\"%s\",signature=\"%s\"",
		wpl.payConfig.MchId, nonce, timestamp, wpl.payConfig.PrivateSerialNo, sign)
	return token, nil
}

// 生成调起支付的参数
// 微信支付文档: https://pay.weixin.qq.com/wiki/doc/apiv3/apis/chapter3_1_4.shtml
func (wpl *WxPayLib) genPayInvokeInfo(prepayId string) (payInvokeInfo *WxPayInvokeInfo, err error) {
	payInvokeInfo = &WxPayInvokeInfo{
		AppId:     wpl.payConfig.AppId,
		TimeStamp: fmt.Sprintf("%v", time.Now().Unix()),
		NonceStr:  util.RandomString(32),
		Package:   "prepay_id=" + prepayId,
		SignType:  "RSA",
	}

	// 签名
	message := fmt.Sprintf("%s\n%s\n%s\n%s\n", payInvokeInfo.AppId, payInvokeInfo.TimeStamp, payInvokeInfo.NonceStr, payInvokeInfo.Package)

	pemFileReader, err := resources.LoadResourceFile("wxpay.private.pem")
	if err != nil {
		return
	}
	privateKey, err := ioutil.ReadAll(pemFileReader)
	if err != nil {
		return
	}

	sha245MsgBytes := util.SHA256HashBytes(message)
	signBytes, err := util.RsaSignPKCS1v15(sha245MsgBytes, privateKey, crypto.SHA256)
	if err != nil {
		return
	}
	payInvokeInfo.PaySign = base64.StdEncoding.EncodeToString(signBytes)

	return payInvokeInfo, nil
}

// ValidateNotifySignature 验证微信支付结果通知的签名
// 微信API文档: https://pay.weixin.qq.com/docs/merchant/development/interface-rules/signature-verification.html
// @param timeStamp 签名生成时间 从 HTTP 头 Wechatpay-Timestamp 获取
// @param nonce HTTP 头 Wechatpay-Nonce 中的应答随机串
// @param signature 微信支付的应答签名, 通过HTTP头Wechatpay-Signature 传递
// @param rawPost 原始请求体
func (wpl *WxPayLib) ValidateNotifySignature(timeStamp, nonce, signature, rawPost string) (verifyRes bool, err error) {
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		err = errcode.Wrap("WxPayLibValidateCallBackSignatureError", err)
		return
	}

	message := fmt.Sprintf("%s\n%s\n%s\n", timeStamp, nonce, rawPost)

	pemFileReader, err := resources.LoadResourceFile("wxp_pub.pem")
	if err != nil {
		err = errcode.Wrap("WxPayLibValidateCallBackSignatureError", err)
		return
	}
	publicKeyStr, err := ioutil.ReadAll(pemFileReader)
	if err != nil {
		err = errcode.Wrap("WxPayLibValidateCallBackSignatureError", err)
		return
	}

	//pem解码
	block, _ := pem.Decode(publicKeyStr)
	//x509解码
	publicKeyInterface, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic(err)
	}
	publicKey := publicKeyInterface.PublicKey.(*rsa.PublicKey)
	//验证数字签名
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, util.SHA256HashBytes(message), signatureBytes) //crypto.SHA1
	verifyRes = nil == err

	return verifyRes, err
}

// DecryptNotifyResourceData 解密微信支付通知中的 resource 数据
func (wpl *WxPayLib) DecryptNotifyResourceData(rawPost string) (notifyResourceData *WxPayNotifyResource, err error) {
	var notifyResponse WxPayNotifyResponse
	if err = json.Unmarshal([]byte(rawPost), &notifyResponse); err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyResourceDataError", err)
	}

	aesKey := []byte(wpl.payConfig.AesKey)
	nonce := []byte(notifyResponse.Resource.Nonce)
	associatedData := []byte(notifyResponse.Resource.AssociatedData)
	ciphertext, err := base64.StdEncoding.DecodeString(notifyResponse.Resource.Ciphertext)
	if err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyResourceDataError", err)
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyResourceDataError", err)
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyResourceDataError", err)
	}
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, associatedData)
	if err != nil {
		return notifyResourceData, errcode.Wrap("WxPayLibDecryptNotifyResourceDataError", err)
	}

	err = json.Unmarshal(plaintext, &notifyResourceData)
	return
}
