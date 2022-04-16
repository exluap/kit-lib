//+build integration

package grpc

import (
	"context"
	kitContext "github.com/exluap/kit/context"
	"github.com/exluap/kit/er"
	"github.com/exluap/kit/log"
	kitTest "github.com/exluap/kit/test"
	"github.com/magiconair/properties/assert"
	"google.golang.org/grpc/codes"
	"testing"
	"time"
)

var Logger = log.Init(&log.Config{Level: log.TraceLevel})

func lf() log.CLoggerFunc {
	return func() log.CLogger {
		return log.L(Logger).Srv("test")
	}
}

type srvImpl struct {
	UnimplementedTestServiceServer
}

func (s *srvImpl) WithError(ctx context.Context, rq *WithErrorRequest) (*WithErrorResponse, error) {
	e := er.WithBuilder("TST-123", "%s happens", "shit").GrpcSt(uint32(codes.AlreadyExists)).C(ctx).F(er.FF{"id": "123"}).Err()
	return nil, e
}

func (s *srvImpl) WithPanic(ctx context.Context, rq *WithPanicRequest) (*WithPanicResponse, error) {
	panic("JUST A PANIC, BRO.....")
}

func Test(t *testing.T) {

	srv, _ := NewServer("test", lf(), &ServerConfig{Port: "55556"})
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(); err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(time.Millisecond * 200)

	cl, err := NewClient(&ClientConfig{Host: "localhost", Port: "55556"})
	if err != nil {
		t.Fatal(err)
	}

	ctx := kitContext.NewRequestCtx().Test().ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	_, err = svc.WithError(ctx, &WithErrorRequest{})
	if err != nil {
		if appErr, ok := er.Is(err); ok {
			ctx := appErr.Fields()["ctx"].(map[string]interface{})
			assert.Equal(t, ctx["_ctx.cl"], "test")
			assert.Equal(t, appErr.Fields()["id"], "123")
			lf()().E(err).Err()
		} else {
			t.Fatal("not app error")
		}
	}
}

func TestPanicRecover(t *testing.T) {
	port := "55557"
	srv, _ := NewServer("test", lf(), &ServerConfig{Port: port})
	RegisterTestServiceServer(srv.Srv, &srvImpl{})

	go func() {
		if err := srv.Listen(); err != nil {
			t.Fatal(err)
		}
	}()

	cl, err := NewClient(&ClientConfig{Host: "localhost", Port: port})
	if err != nil {
		t.Fatal(err)
	}

	ctx := kitContext.NewRequestCtx().Test().ToContext(context.Background())
	svc := NewTestServiceClient(cl.Conn)
	_, err = svc.WithPanic(ctx, &WithPanicRequest{})
	kitTest.AssertAppErr(t, err, ErrCodeGrpcPanic)
}
