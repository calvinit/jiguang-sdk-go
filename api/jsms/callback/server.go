/*
 *
 * Copyright 2025 cavlabs/jiguang-srv-go authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package callback

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cavlabs/jiguang-sdk-go/api"
	"github.com/cavlabs/jiguang-sdk-go/jiguang"
)

// 回调接口服务核心结构。
type Server struct {
	server    *http.Server
	path      string
	isRunning bool
	mu        sync.RWMutex
	logger    jiguang.Logger
}

// 创建新的 Server 回调接口服务实例。
func NewServer(appKey, masterSecret string, opts ...ConfigOption) (*Server, error) {
	c := config{
		addr:   defaultAddr,
		path:   defaultPath,
		logger: api.DefaultJSmsLogger,
	}

	for _, opt := range opts {
		if err := opt.apply(&c); err != nil {
			return nil, err
		}
	}

	p := loggingDataProcessor{logger: c.logger}      // 需要使用用户可能自定义设置的 logger
	if c.flag&flagReply == 0 {
		c.reply = loggingReplyDataProcessor(p)       // 「用户回复消息」SMS_REPLY
	}
	if c.flag&flagReport == 0 {
		c.report = loggingReportDataProcessor(p)     // 「短信送达状态」SMS_REPORT
	}
	if c.flag&flagTemplate == 0 {
		c.template = loggingTemplateDataProcessor(p) // 「模板审核结果」SMS_TEMPLATE
	}
	if c.flag&flagSign == 0 {
		c.sign = loggingSignDataProcessor(p)         // 「签名审核结果」SMS_SIGN
	}

	if c.handler == nil {
		h := defaultHandler{
			appKey:       appKey,
			masterSecret: masterSecret,
			reply:        c.reply,
			report:       c.report,
			template:     c.template,
			sign:         c.sign,
		}
		c.handler = http.HandlerFunc(h.Callback)
	}

	return &Server{
		server: &http.Server{
			Addr:    c.addr,
			Handler: c.handler,
		},
		path:   c.path,
		logger: c.logger,
	}, nil
}

// 处理回调请求。
func (srv *Server) Handle(w http.ResponseWriter, r *http.Request) error {
	srv.server.Handler.ServeHTTP(w, r)
	return nil
}

// 启动回调接口服务。
func (srv *Server) Run() error {
	if srv.hasStarted() {
		return errors.New("JSMS callback server is already running")
	}

	srv.start()

	var wg sync.WaitGroup
	wg.Add(1)

	srv.logger.Infof(context.TODO(), "正在启动极光短信回调接口服务，监听地址为 %s，回调路径为 %s", srv.server.Addr, srv.path)

	startCh, errorCh := make(chan struct{}), make(chan error, 1)

	go func() {
		defer wg.Done()

		ln, err := net.Listen("tcp", srv.server.Addr)
		if err != nil {
			errorCh <- err
			return
		}
		close(startCh)

		if err = srv.server.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errorCh <- err
		}
	}()

	go srv.autoStop()

	select {
	case <-startCh:
		srv.logger.Infof(context.TODO(), "极光短信回调接口服务启动成功！")
	case err := <-errorCh:
		srv.logger.Errorf(context.TODO(), "极光短信回调接口服务启动失败：%s", err)
		return err
	case <-time.After(time.Second * 5):
		srv.logger.Error(context.TODO(), "极光短信回调接口服务启动超时！")
		return errors.New("JSMS callback server startup timeout")
	}

	wg.Wait()

	return nil
}

// 监听系统信号（如 SIGINT、SIGTERM 等），自动停止回调接口服务。
func (srv *Server) autoStop() {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	sig := <-stopCh // 等待接收到停止信号。

	srv.logger.Infof(context.TODO(), "接收到停止信号：%s！", strings.ToUpper(sig.String()))

	if srv.hasStarted() {
		srv.stop()
	} else {
		srv.logger.Info(context.TODO(), "极光短信回调接口服务已停止！")
		os.Exit(-1)
	}

	srv.logger.Info(context.TODO(), "正在停止极光短信回调接口服务...")
	// 使用 5 秒钟的宽限时间来优雅关闭服务。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.server.Shutdown(ctx); err != nil {
		srv.logger.Warnf(context.TODO(), "极光短信回调接口服务优雅停止失败：%s，正在尝试强制停止...", err)
		if err = srv.server.Close(); err != nil {
			srv.logger.Errorf(context.TODO(), "极光短信回调接口服务强制停止失败：%s，直接退出！", err)
			os.Exit(1)
		}
	}
	srv.logger.Info(context.TODO(), "极光短信回调接口服务已停止！")
	os.Exit(0)
}

func (srv *Server) hasStarted() bool {
	srv.mu.RLock()
	defer srv.mu.RUnlock()
	return srv.isRunning
}

func (srv *Server) start() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	mux := http.NewServeMux()
	mux.Handle(srv.path, srv.server.Handler)
	srv.isRunning = true
}

func (srv *Server) stop() {
	srv.mu.Lock()
	defer srv.mu.Unlock()
	srv.isRunning = false
}
