
func IsWebsocket(req *http.Request) bool {
    connHeader := ""
    connHeaders := req.Header["Connection"]
    if len(connHeaders) > 0 {
        connHeader = connHeaders[0]
    }

    isws := false
    if strings.ToLower(connHeader) == "upgrade" {
        upgradeHdrs := req.Header["Upgrade"]
        if len(upgradeHdrs) > 0 {
            isws = (strings.ToLower(upgradeHdrs[0]) == "websocket")
        }
    }

    return isws
}

// proxy for websocket
func WebsocketProxy(target string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        down, err := net.Dial("tcp", target)
        if err != nil {
            log.Errorf("Error dialing websocket backend %s: %v", target, err)
            return
        }
        hj, ok := w.(http.Hijacker)
        if !ok {
            log.Error("http hijacker assert error")
            return
        }
        up, _, err := hj.Hijack()
        if err != nil {
            log.Errorf("Hijack error: %v", err)
            return
        }
        defer up.Close()
        defer down.Close()

        err = r.Write(down)
        if err != nil {
            log.Errorf("Error copying request to target: %v", err)
            return
        }

        errChan := make(chan error, 2)
        copy := func(dst io.Writer, src io.Reader) {
            _, err := io.Copy(dst, src)
            errChan <- err
        }
        go copy(down, up)
        go copy(up, down)
        <-errChan
    })

}

