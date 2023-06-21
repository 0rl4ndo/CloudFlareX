package TLSX

import (
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	http "github.com/bogdanfinn/fhttp"
	"github.com/bogdanfinn/fhttp/http2"
	tls_client "github.com/bogdanfinn/tls-client"
	tls "github.com/bogdanfinn/utls"
)

type Time struct {
	time.Time
}

type data struct {
	Time Time `json:"time"`
}
type Cookie struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Path        string `json:"path"`
	Domain      string `json:"domain"`
	Expires     time.Time
	JSONExpires Time          `json:"expires"`
	RawExpires  string        `json:"rawExpires"`
	MaxAge      int           `json:"maxAge"`
	Secure      bool          `json:"secure"`
	HTTPOnly    bool          `json:"httpOnly"`
	SameSite    http.SameSite `json:"sameSite"`
	Raw         string
	Unparsed    []string `json:"unparsed"`
}

type TLSX struct {
	TLSClient tls_client.HttpClient
}

func (client *TLSX) Req(URL string, Body string, MethodS string, Cookies []Cookie, Headers map[string]string) (*http.Response, string, error) {

	req, err := http.NewRequest(MethodS, URL, strings.NewReader(Body))
	if err != nil {
		log.Println(err)
		return &http.Response{}, "", nil
	}
	for _, properties := range Cookies {
		req.AddCookie(&http.Cookie{
			Name:       properties.Name,
			Value:      properties.Value,
			Path:       properties.Path,
			Domain:     properties.Domain,
			Expires:    properties.JSONExpires.Time, //TODO: scuffed af
			RawExpires: properties.RawExpires,
			MaxAge:     properties.MaxAge,
			HttpOnly:   properties.HTTPOnly,
			Secure:     properties.Secure,
			Raw:        properties.Raw,
			Unparsed:   properties.Unparsed,
		})
	}

	headerOrder := []string{}
	for key, _ := range Headers {
		lowercasekey := strings.ToLower(key)
		headerOrder = append(headerOrder, lowercasekey)
	}

	for Key, Value := range Headers {
		if Key != "host" {
			req.Header.Set(Key, Value)
		}
	}
	u, _ := url.Parse(URL)

	req.Header.Set("Host", u.Host)
	req.Header.Set("user-agent", Headers["user-agent"])

	req.Header = http.Header{
		http.HeaderOrderKey: headerOrder,
		http.PHeaderOrderKey: []string{
			":method",
			":authority",
			":scheme",
			":path",
		},
	}
	resp, err := client.TLSClient.Do(req)

	if err != nil {
		log.Println(err)
		return &http.Response{}, "", nil
	}

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return &http.Response{}, "", nil
	}
	return resp, string(buff), nil

}

func HTTPClient() TLSX {
	options := []tls_client.HttpClientOption{
		tls_client.WithTimeoutSeconds(60),
		tls_client.WithInsecureSkipVerify(),
		tls_client.WithClientProfile(GetTLS()), // use custom profile here
		tls_client.WithCookieJar(tls_client.NewCookieJar()),
	}

	client, err := tls_client.NewHttpClient(tls_client.NewNoopLogger(), options...)
	if err != nil {
		log.Println(err)
		return TLSX{
			TLSClient: client,
		}
	}
	return TLSX{
		TLSClient: client,
	}
}

func GetTLS() tls_client.ClientProfile {
	settings := map[http2.SettingID]uint32{
		http2.SettingHeaderTableSize:      65536,
		http2.SettingEnablePush:           0,
		http2.SettingMaxConcurrentStreams: 1000,
		http2.SettingInitialWindowSize:    6291456,
		http2.SettingMaxHeaderListSize:    262144,
	}

	settingsOrder := []http2.SettingID{
		http2.SettingHeaderTableSize,
		http2.SettingEnablePush,
		http2.SettingMaxConcurrentStreams,
		http2.SettingInitialWindowSize,
		http2.SettingMaxHeaderListSize,
	}

	pseudoHeaderOrder := []string{
		":method",
		":authority",
		":scheme",
		":path",
	}

	connectionFlow := uint32(15663105)
	signatureScheme := []tls.SignatureScheme{
		tls.ECDSAWithP256AndSHA256,
		tls.PSSWithSHA256,
		tls.PKCS1WithSHA256,
		tls.ECDSAWithP384AndSHA384,
		tls.PSSWithSHA384,
		tls.PKCS1WithSHA384,
		tls.PSSWithSHA512,
		tls.PKCS1WithSHA512,
	}
	specFactory := func() (tls.ClientHelloSpec, error) {
		return tls.ClientHelloSpec{
			CipherSuites: []uint16{
				tls.GREASE_PLACEHOLDER,
				tls.TLS_AES_128_GCM_SHA256,
				tls.TLS_AES_256_GCM_SHA384,
				tls.TLS_CHACHA20_POLY1305_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
				tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			},
			CompressionMethods: []uint8{
				tls.CompressionNone,
			},
			Extensions: []tls.TLSExtension{
				&tls.UtlsGREASEExtension{},
				&tls.SNIExtension{},
				&tls.UtlsExtendedMasterSecretExtension{},
				&tls.RenegotiationInfoExtension{},
				&tls.SupportedCurvesExtension{[]tls.CurveID{
					tls.GREASE_PLACEHOLDER,
					tls.X25519,
					tls.CurveP256,
					tls.CurveP384,
				}},
				&tls.SupportedPointsExtension{SupportedPoints: []byte{
					0x00, // pointFormatUncompressed
				}},
				&tls.SessionTicketExtension{},
				&tls.ALPNExtension{AlpnProtocols: []string{"h2", "http/1.1"}},
				&tls.StatusRequestExtension{},
				&tls.SignatureAlgorithmsExtension{SupportedSignatureAlgorithms: signatureScheme},
				&tls.SCTExtension{},
				&tls.KeyShareExtension{[]tls.KeyShare{
					{Group: tls.CurveID(tls.GREASE_PLACEHOLDER), Data: []byte{0}},
					{Group: tls.X25519},
				}},
				&tls.PSKKeyExchangeModesExtension{[]uint8{
					tls.PskModeDHE,
				}},
				&tls.SupportedVersionsExtension{[]uint16{
					tls.VersionTLS13,
					tls.VersionTLS12,
				}},
				&tls.UtlsCompressCertExtension{
					Algorithms: []tls.CertCompressionAlgo{tls.CertCompressionBrotli},
				},
				&tls.UtlsGREASEExtension{},
				&tls.ApplicationSettingsExtension{
					SupportedProtocols: []string{
						"h2",
					},
				},
				&tls.UtlsPaddingExtension{GetPaddingLen: tls.BoringPaddingStyle},
			},
		}, nil
	}

	CustomTLSClient := tls_client.NewClientProfile(tls.ClientHelloID{
		Client:      "CustomTLSClient",
		Version:     "1",
		Seed:        nil,
		SpecFactory: specFactory,
	}, settings, settingsOrder, pseudoHeaderOrder, connectionFlow, nil, nil)
	return CustomTLSClient
}
