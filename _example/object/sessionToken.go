package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sts "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sts/v20180813"

	"github.com/mozillazg/go-cos"
	"github.com/mozillazg/go-cos/debug"
)

type TmpAuth struct {
	SecretID     string
	SecretKey    string
	SessionToken string
}

// https://cloud.tencent.com/document/product/598/33416
// https://console.cloud.tencent.com/api/explorer?Product=sts&Version=2018-08-13&Action=GetFederationToken&SignVersion=
// https://cloud.tencent.com/document/product/436/31923
func getTmpAuth() TmpAuth {
	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	parts := strings.Split(u.Host, ".")
	// bucketName := parts[0]
	// bucketParts := strings.Split(bucketName, "-")
	// appID := bucketParts[len(bucketParts)-1]
	region := parts[2]
	regex := regexp.MustCompile("-\\d")
	region = regex.ReplaceAllString(region, "")

	credential := common.NewCredential(
		os.Getenv("COS_SECRETID"),
		os.Getenv("COS_SECRETKEY"),
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sts.tencentcloudapi.com"
	client, _ := sts.NewClient(credential, region, cpf)

	request := sts.NewGetFederationTokenRequest()

	// 没搞明白怎么为单个文件或目录设置 Policy，按照文档的示例尝试总是不对，所以这里 resource 的值设置为 * 以便可以顺利验证程序功能。
	params := "{\"Name\":\"test\",\"Policy\":\"{   \\\"version\\\": \\\"2.0\\\",   \\\"statement\\\": [     {       \\\"action\\\": [         \\\"name/cos:GetObject\\\"       ],       \\\"effect\\\": \\\"allow\\\",       \\\"resource\\\": [         \\\"*\\\"       ]     }   ] }\"}"
	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.GetFederationToken(request)
	if err != nil {
		panic(err)
	}

	cres := response.Response.Credentials
	return TmpAuth{
		SecretID:     *cres.TmpSecretId,
		SecretKey:    *cres.TmpSecretKey,
		SessionToken: *cres.Token,
	}
}

func main() {
	tmpAuth := getTmpAuth()
	fmt.Printf("%#v\n\n", tmpAuth)

	u, _ := url.Parse(os.Getenv("COS_BUCKET_URL"))
	b := &cos.BaseURL{BucketURL: u}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			// 使用临时密钥
			SecretID:     tmpAuth.SecretID,
			SecretKey:    tmpAuth.SecretKey,
			SessionToken: tmpAuth.SessionToken,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  true,
				RequestBody:    true,
				ResponseHeader: true,
				ResponseBody:   true,
			},
		},
	})

	name := "test/hello.txt"
	resp, err := c.Object.Get(context.Background(), name, nil)
	if err != nil {
		panic(err)
	}
	bs, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Printf("%s\n", string(bs))
}
