package handler

import (
	"errors"
	"net/http"
	"regexp"
)

type Router struct {
	Client   *http.Client
	Patterns map[string]Handler
}

// 通用执行函数
func (r *Router) Do(req *http.Request) (*http.Response, error) {

	// 截取request获得真正的api进行处理函数的查找并执行
	for pattern, handler := range r.Patterns {

		p, _ := regexp.Compile(pattern)
		if p.MatchString(req.URL.Path) {

			for _, handler := range handler.Requests {
				err := handler(req)
				if err != nil {
					return nil, err
				}
			}

			resp, err := r.Client.Do(req)
			if err != nil {
				return nil, err
			}

			for _, handler := range handler.Responses {
				err := handler(resp)
				if err != nil {
					return nil, err
				}
			}

			return resp, nil
		}

	}

	return nil, errors.New("not match any pattern, current pattern: " + req.URL.Path)
}
