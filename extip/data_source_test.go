package extip

import (
	"fmt"
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
				ExpectError: regexp.MustCompile("invalid or unknown key: this_doesnt_exist"),
			},
		},
	})
}

type TestHttpMock struct {
	server *httptest.Server
}

const testDataSourceConfig_basic = `
data "extip" "http_test" {
  resolver = "%s/meta_%d.txt"
}
output "ipaddress" {
  value = "${data.extip.http_test.ipaddress}"
}
`

func TestDataSource_http200(t *testing.T) {
	testHttpMock := setUpMockHttpServer()

	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testDataSourceConfig_basic, testHttpMock.server.URL, 200),
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
	testHttpMock := setUpMockHttpServer()

	defer testHttpMock.server.Close()

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      fmt.Sprintf(testDataSourceConfig_basic, testHttpMock.server.URL, 404),
				ExpectError: regexp.MustCompile("HTTP request error. Response code: 404"),
			},
		},
	})
}

func setUpMockHttpServer() *TestHttpMock {
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

	return &TestHttpMock{
		server: Server,
	}
}
