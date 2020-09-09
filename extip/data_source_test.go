package extip

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const testDataSourceParameterValue = `
data "extip" "parameter_tests" {
  %s = "%s"
}
`

var parametertests = []struct {
	parameter  string
	value      string
	errorRegex string
}{
	{"this_doesnt_exist", "foo", "An argument named \"this_doesnt_exist\" is not expected here."},
	{"resolver", "not-a-valid-url", "config is invalid: expected \"resolver\" to have a host, got not-a-valid-url"},
	{"resolver", "https://notrealsite.fakeurl", "lookup notrealsite.fakeurl.+no such host"},
}

func TestParameterErrors(t *testing.T) {
	for _, tt := range parametertests {
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config:      fmt.Sprintf(testDataSourceParameterValue, tt.parameter, tt.value),
					ExpectError: regexp.MustCompile(tt.errorRegex),
				},
			},
		})
	}
}

const testDataSourceConfigBasic = `
data "extip" "http_test" {
  resolver = "%s/meta_%s.txt"
}
output "ipaddress" {
  value = data.extip.http_test.ipaddress
}
`

var mockedtestserrors = []struct {
	path       string
	errorRegex string
}{
	{"404", "HTTP request error. Response code: 404"},
	{"timeout", "context deadline exceeded"},
	{"hijack", "transport connection broken|unexpected EOF"},
	{"body_error", "unexpected EOF"},
}

func TestMockedResponsesErrors(t *testing.T) {

	for _, tt := range mockedtestserrors {
		var TestHTTPMock *httptest.Server
		if tt.path == "body_error" {
			// For some reason I cant get this to work as a specific path, so creating a different server for it
			TestHTTPMock = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Length", "1")
			}))
			defer TestHTTPMock.Close()
		} else {
			TestHTTPMock = setUpMockHTTPServer()
			defer TestHTTPMock.Close()
		}
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config:      fmt.Sprintf(testDataSourceConfigBasic, TestHTTPMock.URL, tt.path),
					ExpectError: regexp.MustCompile(tt.errorRegex),
				},
			},
		})
	}
}

var mockedtestssuccess = []struct {
	path string
}{
	{"200"},
}

func TestMockedResponsesSuccess(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.Close()

	for _, tt := range mockedtestssuccess {
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: fmt.Sprintf(testDataSourceConfigBasic, TestHTTPMock.URL, tt.path),
					Check: func(s *terraform.State) error {
						_, ok := s.RootModule().Resources["data.extip.http_test"]
						if !ok {
							return fmt.Errorf("missing data resource")
						}

						outputs := s.RootModule().Outputs

						if outputs["ipaddress"].Value != "127.0.0.1" {
							return fmt.Errorf(
								`'ipaddress' output is %s; want '127.0.0.1'`,
								outputs["ipaddress"].Value,
							)
						}

						return nil
					},
				},
			},
		})
	}
}

const testDataSourceConfigTimeout = `
data "extip" "timeout_tests" {
  resolver 			 = "%s/meta_%s.txt"
  client_timeout = %s
}
output "ipaddress" {
  value = data.extip.timeout_tests.ipaddress
}
`

var timeouttests = []struct {
	path          string
	clienttimeout string
}{
	{"timeout", "3000"},
	{"timeout", "0"},
}

func TestTimeouts(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.Close()

	for _, tt := range timeouttests {
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config: fmt.Sprintf(testDataSourceConfigTimeout, TestHTTPMock.URL, tt.path, tt.clienttimeout),
					Check: func(s *terraform.State) error {
						_, ok := s.RootModule().Resources["data.extip.timeout_tests"]
						if !ok {
							return fmt.Errorf("missing data resource")
						}

						outputs := s.RootModule().Outputs

						if outputs["ipaddress"].Value != "127.0.0.1" {
							return fmt.Errorf(
								`'ipaddress' output is %s; want '127.0.0.1'`,
								outputs["ipaddress"].Value,
							)
						}

						return nil
					},
				},
			},
		})
	}
}

var timeouttesterror = []struct {
	path          string
	clienttimeout string
	errorRegex    string
}{
	{"timeout", "100", "Timeout exceeded while awaiting headers"},
}

func TestTimeoutErrors(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.Close()

	for _, tt := range timeouttesterror {
		resource.UnitTest(t, resource.TestCase{
			Providers: testProviders,
			Steps: []resource.TestStep{
				resource.TestStep{
					Config:      fmt.Sprintf(testDataSourceConfigTimeout, TestHTTPMock.URL, tt.path, tt.clienttimeout),
					ExpectError: regexp.MustCompile(tt.errorRegex),
				},
			},
		})
	}
}

func setUpMockHTTPServer() *httptest.Server {
	Server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "text/plain")
			if r.URL.Path == "/meta_200.txt" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("127.0.0.1"))
			} else if r.URL.Path == "/meta_404.txt" {
				w.WriteHeader(http.StatusNotFound)
			} else if r.URL.Path == "/meta_hijack.txt" {
				w.WriteHeader(100)
				w.Write([]byte("Hello3"))
				hj, _ := w.(http.Hijacker)
				conn, _, _ := hj.Hijack()
				conn.Close()
			} else if r.URL.Path == "/meta_timeout.txt" {
				time.Sleep(2000 * time.Millisecond)
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("127.0.0.1"))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	return Server
}

const testDataSourceConfigReal = `
data "extip" "default_test" {
}
output "ipaddress" {
  value = data.extip.default_test.ipaddress
}
`

func IsIpv4Net(host string) bool {
	return net.ParseIP(host) != nil
}

func TestDataSource_DefaultResolver(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testDataSourceConfigReal,
				Check: func(s *terraform.State) error {
					_, ok := s.RootModule().Resources["data.extip.default_test"]
					if !ok {
						return fmt.Errorf("missing data resource")
					}

					outputs := s.RootModule().Outputs

					if !IsIpv4Net(fmt.Sprintf("%v", outputs["ipaddress"].Value)) {
						return fmt.Errorf(
							`'ipaddress' output was not a valid IP address: %s`,
							outputs["ipaddress"].Value,
						)
					}

					return nil
				},
			},
		},
	})
}
