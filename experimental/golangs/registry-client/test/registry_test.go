package test

import (
	"crypto/tls"
	"fmt"
	"github.com/aiziyuer/registry/client/common"
	"github.com/aiziyuer/registry/client/registry"
	"github.com/aiziyuer/registry/client/util"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"testing"
)

// 带认证的, 从环境变量中获取账号信息
var clientWithAuth *registry.Registry

// 不带认证的, 只能访问公共信息
var clientWithOutAuth *registry.Registry

func init() {

	// 测试环境以.env文件为准, 生产环境改成godotenv.Load()
	err := godotenv.Overload()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	clientWithAuth = registry.NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}, &registry.Endpoint{
		Schema: util.GetEnvAnyWithDefault("https", "REGISTRY_SCHEMA"),
		Host:   util.GetEnvAnyWithDefault("registry-1.docker.io", "REGISTRY_HOST"),
	}, &common.Auth{
		UserName: util.GetEnvAnyWithDefault("aiziyuer", "REGISTRY_USERNAME"),
		PassWord: util.GetEnvAnyWithDefault("Changeme_123", "REGISTRY_PASSWORD"),
	})

	clientWithOutAuth = registry.NewClient(&http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}, &registry.Endpoint{
		Schema: util.GetEnvAnyWithDefault("https", "REGISTRY_SCHEMA"),
		Host:   util.GetEnvAnyWithDefault("registry-1.docker.io", "REGISTRY_HOST"),
	}, nil)
}

func TestClient(t *testing.T) {
	err := clientWithAuth.Ping()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func TestTagsWithAuth(t *testing.T) {
	output, err := clientWithAuth.Tags("aiziyuer/centos")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(util.PrettyFormat(output))
}

func TestTagsWithoutAuth(t *testing.T) {
	output, err := clientWithOutAuth.Tags("library/centos")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(util.PrettyFormat(output))
}
