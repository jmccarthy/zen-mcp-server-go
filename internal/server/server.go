package server

import (
	"context"
	"encoding/json"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/sourcegraph/jsonrpc2"

	"github.com/BeehiveInnovations/zen-mcp-server-go/internal/tools"
)

// stdioConn implements io.ReadWriteCloser on top of separate reader/writer.
type stdioConn struct {
	r io.Reader
	w io.Writer
}

func (s stdioConn) Read(p []byte) (int, error)  { return s.r.Read(p) }
func (s stdioConn) Write(p []byte) (int, error) { return s.w.Write(p) }
func (s stdioConn) Close() error                { return nil }

// Server implements a minimal JSON-RPC server over stdio.
type Server struct {
	dispatcher *Dispatcher
}

// NewServer returns a new Server with default tools.
func NewServer() *Server {
	disp := NewDispatcher()
	disp.Register(&tools.GetVersionTool{})
	return &Server{dispatcher: disp}
}

// Run starts processing JSON-RPC messages from r and writes responses to w.
func (s *Server) Run(ctx context.Context, r io.Reader, w io.Writer) error {
	conn := jsonrpc2.NewConn(ctx, jsonrpc2.NewBufferedStream(stdioConn{r, w}, jsonrpc2.VSCodeObjectCodec{}), jsonrpc2.HandlerWithError(s.handle))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		conn.Close()
	}()

	<-conn.DisconnectNotify()
	wg.Wait()
	return nil
}

func (s *Server) handle(ctx context.Context, conn *jsonrpc2.Conn, req *jsonrpc2.Request) (any, error) {
	logrus.WithField("method", req.Method).Debug("rpc request")
	var params map[string]any
	if req.Params != nil {
		if err := json.Unmarshal(*req.Params, &params); err != nil {
			return nil, err
		}
	}
	return s.dispatcher.Call(ctx, req.Method, params)
}
