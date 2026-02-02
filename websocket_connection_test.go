package ocpp_test

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/aliml92/ocpp"
	v16 "github.com/aliml92/ocpp/v16"
	"go.uber.org/zap"
)

var testLog *zap.SugaredLogger

func TestWebsocketConnection(t *testing.T) {
	// init test logger
	logger, _ := zap.NewDevelopment()
	testLog = logger.Sugar()
	defer testLog.Sync()
	ocpp.SetLogger(testLog)

	// init csms server
	go func() {
		csms := ocpp.NewServer()
		csms.AddSubProtocol("ocpp1.6")
		csms.SetCheckOriginHandler(func(r *http.Request) bool { return true })
		csms.SetPreUpgradeHandler(customPreUpgradeHandler)
		csms.SetCallQueueSize(32)

		csms.On("BootNotification", BootNotificationHandler)
		csms.On("Authorize", AuthorizationHandler)

		csms.Start("0.0.0.0:8999", "/ws/", nil)
	}()

	// wait for init csms server
	time.Sleep(1 * time.Second)

	// test for concurrent connections
	t.Run("Multiple Charge Points Connection Concurrent", func(t *testing.T) {
		var wg sync.WaitGroup
		clientCount := 10

		for i := 0; i < clientCount; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				id := fmt.Sprintf("client%d", idx)

				client := ocpp.NewClient()
				client.SetID(id)
				client.AddSubProtocol("ocpp1.6")
				client.SetBasicAuth(id, "dummypass")
				client.SetCallQueueSize(32)

				cp, err := client.Start("ws://localhost:8999", "/ws")
				if err != nil {
					t.Errorf("[%s] error dialing: %v", id, err)
					return
				}
				defer cp.Shutdown()

				if err = sendBootNotification(cp); err != nil {
					t.Errorf("[%s] BootNotification failed: %v", id, err)
					return
				}

				if !cp.IsConnected() {
					t.Errorf("[%s] expected connected status, got disconnected", id)
				}
			}(i)
		}
		wg.Wait()
	})

	// send messages in single connection
	t.Run("Send Messages In Single Connection", func(t *testing.T) {
		id := "msgsender"
		client := ocpp.NewClient()
		client.SetID(id)
		client.AddSubProtocol("ocpp1.6")
		client.SetBasicAuth(id, "dummypass")
		client.SetCallQueueSize(32)

		cp, err := client.Start("ws://localhost:8999", "/ws")
		if err != nil {
			t.Fatalf("error dialing: %v", err)
		}
		defer cp.Shutdown()

		if err = sendBootNotification(cp); err != nil {
			t.Fatalf("BootNotification failed: %v", err)
		}

		if !cp.IsConnected() {
			t.Error("Charge point should be connected")
		}

		if err = sendAuthorize(cp); err != nil {
			t.Errorf("Authorize failed: %v", err)
		}
	})
}

func sendBootNotification(c *ocpp.ChargePoint) error {
	req := &v16.BootNotificationReq{
		ChargePointModel:  "client00",
		ChargePointVendor: "VendorX",
	}
	res, err := c.Call("BootNotification", req)
	if err != nil {
		return err
	}

	testLog.Debugf("BootNotificationConf: %v", res)
	return nil
}

func sendAuthorize(c *ocpp.ChargePoint) error {
	req := &v16.AuthorizeReq{IdTag: "mockrfid"}
	res, err := c.Call("Authorize", req)
	if err != nil {
		return err
	}
	testLog.Debugf("AuthorizeConf: %v", res)
	return nil
}

func customPreUpgradeHandler(w http.ResponseWriter, r *http.Request) bool {
	u, p, ok := r.BasicAuth()
	if !ok {
		testLog.Debug("error parsing basic auth")
		w.WriteHeader(401)
		return false
	}
	path := strings.Split(r.URL.Path, "/")
	id := path[len(path)-1]

	if u != id {
		testLog.Debugf("Username mismatch. connect: %s, auth user: %s", id, u)
		w.WriteHeader(401)
		return false
	}
	testLog.Debugf("%s connected with pass: %s", u, p)
	return true
}

func BootNotificationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.BootNotificationReq)
	testLog.Debugf("Server received BootNotification from %s: %v", cp.Id, req)

	return &v16.BootNotificationConf{
		CurrentTime: time.Now().Format("2006-01-02T15:04:05.000Z"),
		Interval:    60,
		Status:      "Accepted",
	}
}

func AuthorizationHandler(cp *ocpp.ChargePoint, p ocpp.Payload) ocpp.Payload {
	req := p.(*v16.AuthorizeReq)
	testLog.Debugf("Server received Authorize from %s: %v", cp.Id, req)

	return &v16.AuthorizeConf{
		IdTagInfo: v16.IdTagInfo{
			Status: "Accepted",
		},
	}
}
