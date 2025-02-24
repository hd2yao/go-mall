package library

import (
	"context"
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/stretchr/testify/assert"

	"github.com/hd2yao/go-mall/common/util/httptool"
	"github.com/hd2yao/go-mall/library"
)

func TestMain(m *testing.M) {
	client := &http.Client{Transport: &http.Transport{}}
	gock.InterceptClient(client)
	// 把框架的 httptool 使用的 http client 换成 gock 拦截的 client
	httptool.SetUTHttpClient(client)
	os.Exit(m.Run())
}

func TestWhoisLib_GetHostDetail(t *testing.T) {
	defer gock.Off()
	gock.New("https://ipwho.is").
		MatchHeader("User-Agent", "curl/7.77.0").Get("").
		Reply(200).
		BodyString("{\"ip\":\"127.126.113.220\",\"success\":true}")
	ipDetail, err := library.NewWhoisLib(context.TODO()).GetHostIpDetail()
	assert.Nil(t, err)
	assert.Equal(t, "127.126.113.220", ipDetail.Ip)
}
