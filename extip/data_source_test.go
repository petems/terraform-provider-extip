package extip

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

const testDataSourceUnknownKeyError = `
data "extip" "fail_compilation_unknown_key" {
	this_doesnt_exist = "foo"
}
`

func TestDataSource_compileUnknownKeyError(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceUnknownKeyError,
				ExpectError: regexp.MustCompile("An argument named \"this_doesnt_exist\" is not expected here."),
			},
		},
	})
}

type TestHTTPMock struct {
	server *httptest.Server
}

const testDataSourceConfigBasic = `
data "extip" "http_test" {
  resolver = "%s/meta_%d.txt"
}
output "ipaddress" {
  value = data.extip.http_test.ipaddress
}
`

func TestDataSource_http200(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testDataSourceConfigBasic, TestHTTPMock.server.URL, 200),
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
func TestDataSource_http404(t *testing.T) {
	TestHTTPMock := setUpMockHTTPServer()

	defer TestHTTPMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      fmt.Sprintf(testDataSourceConfigBasic, TestHTTPMock.server.URL, 404),
				ExpectError: regexp.MustCompile("HTTP request error. Response code: 404"),
			},
		},
	})
}

func setUpMockHTTPServer() *TestHTTPMock {
	Server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Content-Type", "text/plain")
			if r.URL.Path == "/meta_200.txt" {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("127.0.0.1"))
			} else if r.URL.Path == "/meta_404.txt" {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
		}),
	)

	return &TestHTTPMock{
		server: Server,
	}
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

const testDataSourceNonExistant = `
data "extip" "not_real" {
	resolver = "https://notrealsite.fakeurl"
}
output "ipaddress" {
  value = data.extip.not_real.ipaddress
}
`

func TestDataSource_NonExistant(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      testDataSourceNonExistant,
				ExpectError: regexp.MustCompile("Error requesting external IP: Get \"https://notrealsite.fakeurl\": dial tcp: lookup notrealsite.fakeurl.+no such host"),
			},
		},
	})
}
