package test_helpers

import (
	"crypto/tls"
	"errors"
	"github.com/rh-messaging/shipshape/pkg/framework/log"
	"io/ioutil"
	"net/http"
	"net/url"
)

type HttpWrapper struct {
	Password string
	User     string
	Method   string
	Header   *http.Header
}

func (hw *HttpWrapper) AddHeader(key, value string) *HttpWrapper {
	if hw.Header == nil {
		hw.Header = &http.Header{}
	}
	hw.Header.Add(key, value)
	return hw
}

func (hw *HttpWrapper) WithMethod(method string) *HttpWrapper {
	hw.Method = method
	return hw
}

func (hw *HttpWrapper) WithPassword(password string) *HttpWrapper {
	hw.Password = password
	return hw
}

func (hw *HttpWrapper) WithUser(user string) *HttpWrapper {
	hw.User = user
	return hw
}

func NewWrapper() *HttpWrapper {
	hw := &HttpWrapper{
		Password: hardcodedCredentials,
		User:     hardcodedCredentials,
		Method:   "GET",
	}
	return hw
}

func (hw *HttpWrapper) PerformHttpRequest(address string) (string, error) {

	//address := test.FormUrl(Protocol, DeployName, "0", SubdomainName, ctx1.Namespace, Domain, AddressBit, Port) //nope.
	// there should be only single address in return in this case.
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	request, err := http.NewRequest(hw.Method, address, nil)
	if err != nil {
		return "", err
	}
	actualPath, _ := url.QueryUnescape(request.URL.Path)
	request.URL = &url.URL{
		Scheme: request.URL.Scheme,
		Host:   request.URL.Host,
		Opaque: actualPath,
	}
	request.Header = *hw.Header
	request.SetBasicAuth(hw.Password, hw.User)
	request.URL.User = url.UserPassword(hw.User, hw.Password)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		log.Logf("body: %s", string(bodyBytes))
		return "", errors.New(resp.Status)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	bodyString := string(bodyBytes)
	// Checking for single item should be enough here.
	return bodyString, nil
}
