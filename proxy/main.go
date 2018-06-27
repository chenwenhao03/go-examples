i
func (m *Service) ReverseProxy(c *gin.Context) {
	svc := c.Param("service")
	target := "127.0.0.1:9090"
	if os.Getenv("PROJECT") != "" {
		target = fmt.Sprintf("%s-%s-%s:8080", os.Getenv("PROJECT"), os.Getenv("APPNAME"), svc)
	}
	log.Infof("reverse proxy with target [%s]", target)
	if IsWebsocket(c.Request) {
		c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/v1/proxy/"+svc)
		WebsocketProxy(target).ServeHTTP(c.Writer, c.Request)
		return
	}
	director := func(req *http.Request) {
		req.Host = target
		req.URL.Scheme = "http"
		req.URL.Host = target
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1/proxy/"+svc)
		if len(req.URL.Path) == 0 {
			req.URL.Path = "/"
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
	return
}

func (m *Service) ReverseZeppelinProxy(c *gin.Context) {
	target := "127.0.0.1:9090"
	if os.Getenv("PROJECT") != "" {
		target = fmt.Sprintf("%s-%s-%s:8080", os.Getenv("PROJECT"), os.Getenv("APPNAME"), "zeppelin")
	}
	reverseProxy(c, target)
	return
}

func (m *Service) ReverseMasterProxy(c *gin.Context) {
	target := "127.0.0.1:8080"
	if os.Getenv("PROJECT") != "" {
		target = fmt.Sprintf("%s-%s-%s:8080", os.Getenv("PROJECT"), os.Getenv("APPNAME"), "master")
	}
	reverseProxy(c, target)
	return
}

func reverseProxy(c *gin.Context, target string) {
	log.Infof("reverse proxy with target [%s]", target)
	if IsWebsocket(c.Request) {
		//c.Request.URL.Path = strings.TrimPrefix(c.Request.URL.Path, "/v1/proxy/"+svc)
		WebsocketProxy(target).ServeHTTP(c.Writer, c.Request)
		return
	}
	director := func(req *http.Request) {
		req.Host = target
		req.URL.Scheme = "http"
		req.URL.Host = target
		//req.URL.Path = strings.TrimPrefix(req.URL.Path, "/v1/proxy/"+svc)
		if len(req.URL.Path) == 0 {
			req.URL.Path = "/"
		}
	}
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(c.Writer, c.Request)
	return
}

