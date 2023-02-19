package cmcc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var NeedLoginAgain = errors.New("need login again")
var ValueError = errors.New("userHeTotalIntegral is not string")

type CmccResponse struct {
	ResultCode    int                    `json:"resultCode"`
	ResultMessage string                 `json:"resultMessage"`
	ResultJson    map[string]interface{} `json:"resultJson,omitempty"`
}

type Cmcc struct {
	endpoint string
	ip       string
	client   *http.Client
}

func NewCmcc(ip string) *Cmcc {
	return &Cmcc{
		endpoint: "https://m.jf.10086.cn",
		ip:       ip,
		client: &http.Client{
			Transport:     nil,
			CheckRedirect: nil,
			Jar:           nil,
			Timeout:       0,
		},
	}
}

func (p Cmcc) SendSmsCode(ctx context.Context, mobile string) error {
	data := url.Values{}
	data.Set("msisdn", mobile)
	data.Set("occasionCode", "sms_user_login_request")

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/cmcc-h5-shop/user/buildAuthVerifyCode", p.endpoint), strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("x-forwarded-for", p.ip)
	req.Header.Set("X-Real-IP", p.ip)
	req.Header.Set("Origin", p.endpoint)
	req.Header.Set("Referer", p.endpoint)
	req.Header.Set("Host", p.endpoint)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var cmccRes CmccResponse
	if err := json.Unmarshal(rawData, &cmccRes); err != nil {
		return err
	}

	if cmccRes.ResultCode != 200 {
		return errors.New(cmccRes.ResultMessage)
	}
	return nil
}

func (p Cmcc) LoginMall(ctx context.Context, mobile string, smsCode string) (string, error) {
	data := url.Values{}
	//ciper := cipher.
	data.Set("userMobile", mobile)
	data.Set("loginType", "2")
	data.Set("verifyCode", smsCode)
	data.Set("realLoginChannel", "CMCCJF")
	data.Set("channel", "h5")

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/cmcc-h5-shop/user/phone", p.endpoint), strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("x-forwarded-for", p.ip)
	req.Header.Set("X-Real-IP", p.ip)
	req.Header.Set("Origin", p.endpoint)
	req.Header.Set("Referer", p.endpoint)
	req.Header.Set("Host", p.endpoint)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var cmccRes CmccResponse
	if err := json.Unmarshal(rawData, &cmccRes); err != nil {
		return "", err
	}

	if cmccRes.ResultCode != 200 {
		return "", errors.New(cmccRes.ResultMessage)
	}

	return res.Header.Get("Sessionid"), nil
}

func (p Cmcc) GetPoints(ctx context.Context, mobile, sessionID string) (int64, error) {
	data := url.Values{}
	data.Set("userMobile", mobile)

	req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/cmcc-h5-shop/users/userLoginByUserMobile", p.endpoint), strings.NewReader(data.Encode()))
	if err != nil {
		return 0, err
	}

	req.Header.Set("sessionid", sessionID)
	req.Header.Set("x-forwarded-for", p.ip)
	req.Header.Set("X-Real-IP", p.ip)
	req.Header.Set("Origin", p.endpoint)
	req.Header.Set("Referer", p.endpoint)
	req.Header.Set("Host", p.endpoint)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := p.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()
	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, err
	}

	var cmccRes CmccResponse
	if err := json.Unmarshal(rawData, &cmccRes); err != nil {
		return 0, err
	}

	if cmccRes.ResultCode != 200 {
		if cmccRes.ResultCode == 401 {
			return 0, NeedLoginAgain
		}
		return 0, errors.New(cmccRes.ResultMessage)
	}

	if v, ok := cmccRes.ResultJson["userHeTotalIntegral"]; ok {
		if strV, ok := v.(string); ok {
			return strconv.ParseInt(strV, 10, 64)
		} else {
			return 0, ValueError
		}
	}
	return 0, errors.New("response is wrong")
}
