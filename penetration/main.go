package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sync"
)

func main() {
	runtime.GOMAXPROCS(7)
	wg := sync.WaitGroup{}
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("test")
			form := url.Values{}
			form.Set("login_name", "Bayan")
			form.Set("login_password", "123")

			resp, err := http.PostForm("http://127.0.0.1:8080/login-process", form)
			if err != nil {
				fmt.Println(err)
				return
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Println(string(body))

		}()
	}
	wg.Wait()
}
