package flowdata

import (
	"testing"
)

type IpInfoResp struct {
	IP       string
	Hostname string
	City     string
	Region   string
	Country  string
	TimeZone string
}

func TestRequest(t *testing.T) {
	req := NewRequest()
	req.Query = map[string]any{"test": "value"}

	resp, err := req.GET("http://ipinfo.io")
	if err != nil {
		t.Fatal(err)
	}

	var ipinfo IpInfoResp
	if err := resp.JSON(&ipinfo); err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ipinfo)
}
