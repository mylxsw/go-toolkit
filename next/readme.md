
# Next

该模块实现了与php-fpm进程通信的功能，将http请求转发给php-fpm进程处理。

    // CreateHandler 创建服务处理器
    func CreateHandler(params *Config, network, address string) http.Handler {
        fpmAddr := address
        if network == "unix" {
            fpmAddr = fmt.Sprintf("%s:%s", network, address)
        }

        rootDir := filepath.Dir(params.EndpointFile)

        conf := next.Config{
            EndpointFile:    params.EndpointFile,
            ServerIP:        params.ServiceIP,
            ServerPort:      params.ServicePort,
            SoftwareName:    "php-server",
            SoftwareVersion: "0.0.1",
            Rules:           []next.Rule{next.NewPHPRule(rootDir, []string{fpmAddr})},
            RequestLogHandler: func(rc *next.RequestContext) {
                var message bytes.Buffer
                if err := params.AccessLogTemplate.Execute(&message, rc); err != nil {
                    log.Module("server").Errorf("invalid log format: %s", err.Error())
                } else {
                    if params.Debug {
                        log.Module("server.request").
                            WithContext(rc.ToMap()).Debugf(message.String())
                    } else {
                        log.Module("server.request").Debugf(message.String())
                    }
                }
            },
        }

        return next.CreateHttpHandler(&conf)
    }


    http.HandleFunc("/", CreateHandler(config, network, address).ServeHTTP)
    srv := &http.Server{Handler: http.DefaultServeMux}
    go func() {
        if err := srv.Serve(listener); err != nil {
            log.Debugf("The http server has stopped: %v", err)
        }
    }()


> 该模块大部分代码是从Caddy中提取出来的，参考 https://github.com/mholt/caddy/tree/master/caddyhttp/fastcgi