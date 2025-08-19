package acloudapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	HeaderAccept      = "Accept"
	HeaderContentType = "Content-Type"
	HeaderUserAgent   = "User-Agent"

	ContentTypeApplicationJson = "application/json"

	DefaultPublicAPIUrl = "https://api.avisi.cloud"
	DefaultUserAgent    = "client-go"
	// DefaultRequestTimeout is the default request timeout for API calls
	DefaultRequestTimeout = 15 * time.Second
)

type RestyClient struct {
	opts ClientOpts

	resty *resty.Client
}

func NewRestyClient(authenticator Authenticator, opts ClientOpts) *RestyClient {
	opts = SetMissingOpts(opts)
	return &RestyClient{
		opts:  opts,
		resty: GetRestyClient(authenticator, opts),
	}
}

type ClientOpts struct {
	// Debug is ... TODO
	Debug bool

	// Trace is ... TODO
	Trace bool

	// DebugShowAuthorizationHeader is  ... TODO
	DebugShowAuthorizationHeader bool

	// APIUrl is the URL of the public api
	APIUrl string

	// UserAgent
	UserAgent string

	// CustomResty sets a custom Resty client
	CustomResty *resty.Client

	// CustomDialer sets a custom net.Dialer on the http.Transport
	//
	// note: only used if CustomResty and CustomTransport are not provided
	CustomDialer *net.Dialer

	// CustomTransport set a custom http.Transport on the Resty client
	//
	// note: only used if CustomResty is not provided
	CustomTransport *http.Transport
}

func SetMissingOpts(opts ClientOpts) ClientOpts {
	if opts.UserAgent != "" {
		opts.UserAgent = fmt.Sprintf("%s (%s)", opts.UserAgent, DefaultUserAgent)
	} else {
		opts.UserAgent = DefaultUserAgent
	}
	if opts.APIUrl == "" {
		opts.APIUrl = DefaultPublicAPIUrl
	}
	return opts
}

func (c *RestyClient) Resty() *resty.Client {
	return c.resty
}

func (c *RestyClient) R() *resty.Request {
	return c.resty.R()
}

func GetRestyClient(authenticator Authenticator, opts ClientOpts) *resty.Client {
	client := opts.CustomResty
	if client == nil {
		client = NewDefaultRestyClient(authenticator, opts, client)
	}
	return client
}

func NewDefaultRestyClient(authenticator Authenticator, opts ClientOpts, client *resty.Client) *resty.Client {
	client = resty.New().
		SetBaseURL(opts.APIUrl).
		SetTimeout(DefaultRequestTimeout).
		SetHeader(HeaderAccept, ContentTypeApplicationJson).
		SetHeader(HeaderContentType, ContentTypeApplicationJson).
		SetHeader(HeaderUserAgent, opts.UserAgent).
		OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
			if authenticator != nil {
				return authenticator.Authenticate(c, r)
			}
			return nil
		})

	dialer := opts.CustomDialer
	if dialer == nil {
		dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}
	}

	transport := opts.CustomTransport
	if transport == nil {
		transport = &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           dialer.DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			MaxIdleConnsPerHost:   runtime.GOMAXPROCS(0) + 1,
			MaxConnsPerHost:       10,
		}
	}
	client.SetTransport(transport)

	if opts.Debug {
		client.SetDebug(true)

		if !opts.DebugShowAuthorizationHeader {
			client.OnRequestLog(func(r *resty.RequestLog) error {
				if r.Header.Get("Authorization") != "" {
					r.Header.Set("Authorization", "<omitted>")
				}
				return nil
			})
		}
	}
	if opts.Trace {
		client.EnableTrace()
	}
	return client
}

func (c *RestyClient) CheckResponse(response *resty.Response, err error) error {
	if response == nil {
		if err != nil {
			return err
		}
		return fmt.Errorf("no response received")
	}

	if c.opts.Trace {
		fmt.Printf("~~~ TRACE %s %s%s ~~~\n", response.Request.Method, c.opts.APIUrl, response.Request.URL)
		ti := response.Request.TraceInfo()
		fmt.Println("  DNSLookup     :", ti.DNSLookup)
		fmt.Println("  ConnTime      :", ti.ConnTime)
		fmt.Println("  TCPConnTime   :", ti.TCPConnTime)
		fmt.Println("  TLSHandshake  :", ti.TLSHandshake)
		fmt.Println("  ServerTime    :", ti.ServerTime)
		fmt.Println("  ResponseTime  :", ti.ResponseTime)
		fmt.Println("  TotalTime     :", ti.TotalTime)
		fmt.Println("  IsConnReused  :", ti.IsConnReused)
		fmt.Println("  IsConnWasIdle :", ti.IsConnWasIdle)
		fmt.Println("  ConnIdleTime  :", ti.ConnIdleTime)
		fmt.Println("  RequestAttempt:", ti.RequestAttempt)
		fmt.Println("  RemoteAddr    :", ti.RemoteAddr.String())
		fmt.Println("==============================================================================")
	}
	if err != nil {
		return err
	}

	body := response.Body()
	var errorResult Error
	if body != nil {
		err = json.Unmarshal(body, &errorResult)
		if err != nil {
			errorResult.Message = string(body)
		}
	}
	if !response.IsSuccess() {
		return fmt.Errorf("%d: %s", response.StatusCode(), errorResult.Message)
	}
	return nil
}

func (c *RestyClient) GetPaged(ctx context.Context, url string) (PagedResult, error) {
	finalPagedResult := PagedResult{}

	morePages := true
	currentPage := 0

	pg := RestyPageGetter{
		client: c,
	}

	for morePages {
		pagedResult, err := pg.Get(ctx, url, currentPage)
		if err != nil {
			return pagedResult, err
		}

		prevContent := finalPagedResult.Content
		finalPagedResult = pagedResult
		finalPagedResult.Content = append(prevContent, pagedResult.Content...)

		if pagedResult.Last {
			morePages = false
			break
		}

		// Check for possible unending loop when "PagedResult.Last" is not set
		if currentPage > MAX_PAGING_LOOPS {
			return PagedResult{}, ErrMaximumPagingLoopsExceeded
		}

		currentPage++
	}

	return finalPagedResult, nil
}
