package cf

import (
	"CloudFlareX/internal/TLSX"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type CfSession struct {
	Client TLSX.TLSX
	Header map[string]string
}
type CFCVParams struct {
	R string
	M string
}

type InvisibleJs struct {
	Password string
	S        string
}

type Payload struct {
	Wp string `json:"wp"`
	S  string `json:"s"`
}

var (
	re = regexp.MustCompile(`[0-9]*\.[0-9]+:[0-9]+:`)
)

func GetCfSession() CfSession {
	Cf := CfSession{
		Client: TLSX.HTTPClient(),
	}
	return Cf
}

func (cf *CfSession) GetInvisible() string {
	resp, err := http.Get("https://discord.com/cdn-cgi/challenge-platform/scripts/invisible.js")
	if err != nil {
		fmt.Println("Error:", err)
		return "https://discord.com/cdn-cgi/challenge-platform/h/b/scripts/jsd/19ad4730/invisible.js"
	}
	defer resp.Body.Close()

	url := resp.Request.URL.String()
	print("Invisible.js: " + url)
	return url
}

func (cf *CfSession) CFCVParams() CFCVParams {
	// __CF$cv$params
	// var js = "window['__CF$cv$params']={r:'7daa852a9e5d5bf1',m:'re409kQ.Fanao1evE6qArOWN4h7KyqaKmZ16Es3.ERM-1687332353-0-Ae6DnBeNJj08NnGOzrqnWwcayeUXLCYPLxTSpw/6imOR'};_cpo=document.createElement('script');_cpo.nonce='NTUsMjAyLDE5NSwxNzksMjI3LDEzLDIyNSwyMjk=',_cpo.src='/cdn-cgi/challenge-platform/scripts/invisible.js',document.getElementsByTagName('head')[0].appendChild(_cpo);";

	_, data, Err := cf.Client.Req("https://discord.com", "", "GET", []TLSX.Cookie{}, map[string]string{})
	if Err != nil {
		println(Err)
		return CFCVParams{}
	}
	R := strings.Split(strings.Split(data, "={r:'")[1], "',m:'")[0]
	println("[R] " + R)
	M := strings.Split(strings.Split(data, "',m:'")[1], "'};")[0]
	println("[M] " + M)
	return CFCVParams{
		R: R,
		M: M,
	}
}

func (cf *CfSession) InvisibleJs(URL string) InvisibleJs {
	// invisible.js
	// bigint;1850185ocGoXx;open;fromCharCode;error on cf_chl_props;function;ontimeout;addEventListener;_cf_chl_opt;symbol;boolean;isArray;/invisible/jsd;[native code];isNaN;charCodeAt;clientInformation;join;getOwnPropertyNames;POST;%2b;6JkUxrm;ActiveXObject;number;contentDocument;body;tabIndex;string;splice;Error object: ;charAt;removeChild;cFPWv;push;146xgrrtc;includes;Set;readyState;XMLHttpRequest;length;getPrototypeOf;/0.7270492140176802:1687334781:eiyN8e47udKcrMPasGV0_nakebs8zRgPqD8M2IEK7Gw/;random;438552XCcJcF;43796XqJzwx;stringify;application/json;0.7270492140176802:1687334781:eiyN8e47udKcrMPasGV0_nakebs8zRgPqD8M2IEK7Gw;1989970iaJpUo;10118187QYlxte;2wyOeX75lQT8AR$P3cuIHgGWLob9hnU-jtkBKx6CfYdi+zraZ1sSpNM4Jmv0EFqVD;timeout;call;contentWindow;Microsoft.XMLHTTP;Message: ;document;send;7914cdQrHm; - ;159ZReQQy;/cv/result/;undefined;Array;pow;indexOf;onreadystatechange;bgUc;from;bweHEMUExl;__CF$cv$params;Content-type;replace;display: none;/cdn-cgi/challenge-platform/h/;keys;loading;navigator;appendChild;style;iframe;object;createElement;Function;application/x-www-form-urlencoded;33236PEhZDL;11lZaBtG;Object;d.cookie;setRequestHeader;msg;hasOwnProperty;prototype;Content-Type;concat

	_, data, Err := cf.Client.Req(URL, "", "GET", []TLSX.Cookie{}, map[string]string{})
	if Err != nil {
		println(Err)
		return InvisibleJs{}
	}
	Base := re.FindString(data)
	S := Base + strings.Split(data, Base)[1][:43]
	println("[S] " + S)

	/*var index int
	for _, K := range {
		println(index)
		println(K)
		index = index +1
	} */
	js := strings.Split(strings.Split(data, "'.split(';')")[0], "='")
	A := strings.Split(js[len(js)-1], ";")
	var Pass string

	for _, K := range A {
		if len(K) == 65 {
			Pass = K
		}
	}
	println("[K] " + Pass)
	return InvisibleJs{
		Password: Pass,
		S:        S,
	}
}

func GetWP(Pass string) string {
	passs := Pass
	url := "http://127.0.0.1:3000/wp"

	data := map[string]string{
		"pass": passs,
	}

	payload, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return ""
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("Error posting request:", err)
		return ""
	}

	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(buff)
}

func (cf *CfSession) GetResult(Invisible InvisibleJs, cfparams CFCVParams) string {

	Payload, _ := json.Marshal(Payload{
		S:  Invisible.S,
		Wp: GetWP(Invisible.Password),
	})
	Response, _, Err := cf.Client.Req(fmt.Sprintf("https://discord.com/cdn-cgi/challenge-platform/h/b/cv/result/%s", cfparams.R), string(Payload), "POST", []TLSX.Cookie{}, map[string]string{})
	if Err != nil {
		println(Err)
		return ""
	}
	if Cookie, Ok := Response.Header["Set-Cookie"]; Ok {
		if strings.Contains(Cookie[0], "__cf_bm") {
			Cb_Bm := strings.Split(strings.Split(Cookie[0], "__cf_bm=")[1], ";")[0]
			println("🔑 __cf_bm: " + Cb_Bm)
			return Cb_Bm
		}
	}
	return ""

}
