/**
 * 线程是否安全的测试
 * 参考 https://www.cnblogs.com/Detector/p/9769840.html
 *      http://tool.oschina.net/commons/
 *		https://www.cnblogs.com/mafeng/p/7068837.html
 * slice切片在并发请求时不会出现异常，但也是线程不安全的，数据更新会有问题
 * go test -v
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"testing"

	"github.com/kataras/iris/httptest"
)

func TestMVC_http(t *testing.T) {
	home := "http://localhost:8080/"
	importApi := "http://localhost:8080/import"
	luckyApi := "http://localhost:8080/lucky"

	urlStr := home
	resp, err := http.Get(urlStr)
	if err != nil {
		t.Errorf("server / err = %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("server / resp is didn't 200 OK:%s", resp.Status)
	}
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	body := string(bodyBytes)
	if body != "当前总共参与抽奖的用户数: 0\n" {
		t.Errorf("server / request err")
	}

	// 并发测试
	// wg 保证子协程全部运行完成
	var wg sync.WaitGroup
	// 启动100个协程并发来执行用户导入操作
	// 如果是线程安全，预期导入成功100个用户
	for i := 0; i < 100; i++ {
		// 每个子协程 wg + 1
		wg.Add(1)
		go func(i int) {
			// 该子协程执行完成,要执行 wg - 1 操作
			defer wg.Done()
			urlStr = importApi
			u := url.Values{}
			u.Set("users", fmt.Sprintf("test_u%d", i))
			resp, err = http.PostForm(urlStr, u)
			if resp.StatusCode != http.StatusOK {
				t.Errorf("server /import resp is didn't 200 OK:%s", resp.Status)
			}
		}(i) // 在go func(){}(i) 需要传参i,在go func(){中直接引用for的i会出现并发安全问题 }
	}
	wg.Wait()

	urlStr = home
	resp, err = http.Get(urlStr)
	bodyBytes, _ = ioutil.ReadAll(resp.Body)
	body = string(bodyBytes)
	if body != "当前总共参与抽奖的用户数: 100\n" {
		t.Errorf("当前总共参与抽奖的用户数应该是: 100\n")
	}

	urlStr = luckyApi
	resp, _ = http.Get(urlStr)
	if resp.StatusCode != http.StatusOK {
		t.Errorf("server /lucky resp is didn't 200 OK:%s", resp.Status)
	}

	urlStr = home
	resp, err = http.Get(urlStr)
	bodyBytes, _ = ioutil.ReadAll(resp.Body)
	body = string(bodyBytes)
	if body != "当前总共参与抽奖的用户数: 99\n" {
		t.Errorf("当前总共参与抽奖的用户数应该是: 99\n")
	}
}

func TestMVC_iris(t *testing.T) {
	e := httptest.New(t, newApp())
	// 并发测试
	// wg 保证子协程全部运行完成
	var wg sync.WaitGroup

	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前总共参与抽奖的用户数: 0\n")

	// 启动100个协程并发来执行用户导入操作
	// 如果是线程安全的时候，预期倒入成功100个用户
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e.POST("/import").WithFormField("users", fmt.Sprintf("test_u%d", i)).Expect().Status(httptest.StatusOK)
		}(i)
	}

	wg.Wait()

	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前总共参与抽奖的用户数: 100\n")
	e.GET("/lucky").Expect().Status(httptest.StatusOK)
	e.GET("/").Expect().Status(httptest.StatusOK).
		Body().Equal("当前总共参与抽奖的用户数: 99\n")
}
