package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gcore/app"
)

var UtmGif = []byte{
	0x47, 0x49, 0x46, 0x38, 0x39, 0x61, 0x01, 0x00, 0x01, 0x00, 0xf0, 0x01,
	0x00, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x21, 0xf9, 0x04, 0x01, 0x0a,
	0x00, 0x00, 0x00, 0x2c, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
	0x00, 0x02, 0x02, 0x44, 0x01, 0x00, 0x3b,
}
var UtmGifLen = strconv.Itoa(len(UtmGif))

// 生成访客id
func GenerateVisitId() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano()/1000000, rand.Int63n(1000))
}

func WriteVisitorId(w http.ResponseWriter, r *http.Request, visitorId string) {
	config := app.GetConfig()
	WriteVisitorIdWithDomain(w, r, visitorId, config.Get("cookie/domain"))
}

func WriteVisitorIdWithDomain(w http.ResponseWriter, r *http.Request, visitorId, domain string) {

	config := app.GetConfig()

	if len(visitorId) == 0 {
		return
	}

	expiration := time.Now()
	expiration = expiration.AddDate(10, 0, 0)
	cookie := &http.Cookie{
		Name:    config.Get("cookie/uid"),
		Value:   visitorId,
		Path:    "/",
		Domain:  domain,
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
	// 添加Cookie到请求里面，多次获取的时候不会重复生成
	r.AddCookie(cookie)
	w.Header().Set("P3P", `CP="CURa ADMa DEVa PSAo PSDo OUR BUS UNI PUR INT DEM STA div COM NAV OTC NOI DSP COR"`)
	w.Header().Set("Cache-Control", "no-cache, private, no-store, must-revalidate, max-stale=0, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
}

func GetVisitId(w http.ResponseWriter, r *http.Request) (visitorId string) {
	config := app.GetConfig()
	cookie, err := r.Cookie(config.Get("cookie/uid"))
	if err == nil {
		visitorId = cookie.Value
	}
	return
}

func GetUUID(w http.ResponseWriter, r *http.Request) (visitorId string) {
	config := app.GetConfig()
	cookie, err := r.Cookie(config.Get("cookie/utmUuid"))
	if err == nil {
		visitorId = cookie.Value
	}
	return
}

// 检查访客id是否存在，
// 如果不存在，则生成一个并往浏览器端写访客id.
// @return: 访客id
func CheckAndWriteVisitorId(w http.ResponseWriter, r *http.Request) (visitorId string, isNewVisit bool) {
	visitorId = GetVisitId(w, r)
	if visitorId == "" {
		visitorId = GenerateVisitId()
		WriteVisitorId(w, r, visitorId)

		isNewVisit = true
	} else if w.Header().Get("P3P") != "" {
		// 如果有P3P头，则是这次请求写的Cookie，则为新访客
		isNewVisit = true
	}
	return
}

func CheckAndNewVisitorId(w http.ResponseWriter, r *http.Request) (visitorId string, isNewVisit bool) {
	visitorId = GetVisitId(w, r)
	if visitorId == "" {
		visitorId = GenerateVisitId()
		isNewVisit = true
	} else if w.Header().Get("P3P") != "" {
		// 如果有P3P头，则是这次请求写的Cookie，则为新访客
		isNewVisit = true
	}
	return
}

// return json string result
// WriteJson(obj) or WriteJson(obj, "text/html")
func WriteJson(w http.ResponseWriter, data interface{}, contentType ...string) {
	var ct string
	if len(contentType) == 1 {
		ct = contentType[0]
	} else {
		ct = "application/json"
	}
	w.Header().Set("Content-Type", ct)
	s := ""
	b, err := json.Marshal(data)
	if err != nil {
		s = `{
    "success": false,
    "message": "json.Marshal error"
}`
	} else {
		s = string(b)
	}
	fmt.Fprint(w, s)
}

func WriteJsonp(w http.ResponseWriter, callback string, data interface{}) {
	w.Header().Set("Content-Type", "application/javascript")
	if _, err := w.Write([]byte(callback + "(")); err != nil {
		return
	}
	b, _ := json.Marshal(data)
	if _, err := w.Write(b); err != nil {
		return
	}
	_, _ = w.Write([]byte(");"))
	return
}

func PostJson(message string, registUrl string) (err error) {
	response, err := http.Post(registUrl, "application/json", strings.NewReader(message))
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//read status code
	if response.StatusCode != 200 {
		return err
	}
	return nil
}

func PostProto(message string, registUrl string) (msg []byte, err error) {
	response, err := http.Post(registUrl, "application/octet-stream", strings.NewReader(message))
	if err != nil {
		return
	}
	defer response.Body.Close()

	//read status code
	if response.StatusCode != 200 {
		return
	}
	msg, err = ioutil.ReadAll(response.Body)
	return
}

func GetUrlResp(Url string) (msg []byte, err error) {
	response, err := http.Get(Url)
	if err != nil {
		return
	}
	defer response.Body.Close()

	//read status code
	if response.StatusCode != 200 {
		err = fmt.Errorf("status code error : %d", response.StatusCode)
		return
	}
	msg, err = ioutil.ReadAll(response.Body)
	return
}

func GetUrl(Url string) (err error) {
	response, err := http.Get(Url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//read statu code
	if response.StatusCode != 200 {
		return err
	}
	return nil
}

func PostForm(message url.Values, Url string) (err error) {
	response, err := http.PostForm(Url, message)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//read statu code
	if response.StatusCode != 200 {
		return err
	}
	return nil
}

func GetForm(message url.Values, Url string) (err error) {
	response, err := http.Get(Url + "?" + message.Encode())
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//read statu code
	if response.StatusCode != 200 {
		return err
	}
	return nil
}

// 返回空白图片，强制不缓存
func WriteNoCacheBlankImage(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, private, no-store, must-revalidate, max-stale=0, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "Mon, 26 Jul 1997 05:00:00 GMT")
	//返回一个空gif
	w.Header().Add("Content-Type", "image/gif")
	w.Header().Add("Content-Length", UtmGifLen)
	w.Write(UtmGif)
}

var topDomains string = "ac.cn,ac.jp,ac.uk,ad.jp,adm.br,adv.br,agr.br,ah.cn,am.br,arq.br,art.br,asn.au,ato.br,av.tr,bel.tr,bio.br,biz.tr,bj.cn,bmd.br,cim.br,cng.br,cnt.br,co.at,co.jp,co.uk,com.au,com.br,com.cn,com.eg,com.hk,com.mx,com.ru,com.tr,com.tw,conf.au,cq.cn,csiro.au,dr.tr,ecn.br,edu.au,edu.br,edu.tr,emu.id.au,eng.br,esp.br,etc.br,eti.br,eun.eg,far.br,fj.cn,fm.br,fnd.br,fot.br,fst.br,g12.br,gb.com,gb.net,gd.cn,gen.tr,ggf.br,gob.mx,gov.au,gov.br,gov.cn,gov.hk,gov.tr,gr.jp,gs.cn,gx.cn,gz.cn,ha.cn,hb.cn,he.cn,hi.cn,hk.cn,hl.cn,hn.cn,id.au,idv.tw,imb.br,ind.br,inf.br,info.au,info.tr,jl.cn,jor.br,js.cn,jx.cn,k12.tr,lel.br,ln.cn,ltd.uk,mat.br,me.uk,med.br,mil.br,mil.tr,mo.cn,mus.br,name.tr,ne.jp,net.au,net.br,net.cn,net.eg,net.hk,net.lu,net.mx,net.ru,net.tr,net.tw,net.uk,nm.cn,no.com,nom.br,not.br,ntr.br,nx.cn,odo.br,oop.br,or.at,or.jp,org.au,org.br,org.cn,org.hk,org.lu,org.ru,org.tr,org.tw,org.uk,plc.uk,pol.tr,pp.ru,ppg.br,pro.br,psc.br,psi.br,qh.cn,qsl.br,rec.br,sc.cn,sd.cn,se.com,se.net,sh.cn,slg.br,sn.cn,srv.br,sx.cn,tel.tr,tj.cn,tmp.br,trd.br,tur.br,tv.br,tw.cn,uk.com,uk.net,vet.br,wattle.id.au,web.tr,xj.cn,xz.cn,yn.cn,zj.cn,zlg.br,co.nr,co.nz,com.fr,"

// 获取根域名。
//   realestate.mydomain.com.fr => mydomain.com.fr
//   whois.bizz.cc => bizz.cc
//   www.something.co.cc => something.co.cc
func GetRootDomain(host string) string {
	if strings.HasPrefix(host, "http://") {
		host = strings.Replace(host, "http://", "", 1)
	}

	idx := strings.Index(host, "/")
	if idx > 0 {
		var a []byte
		for i := 0; i < idx; i++ {
			a = append(a, host[i])
		}
		host = string(a)
	}

	partCount := 0
	prePartIndex := 0
	l := len(host)
	for i := l - 1; i > -1; i-- {
		if host[i] != '.' {
			continue
		}
		partCount++
		switch partCount {
		case 1:
			prePartIndex = i
		case 2:
			td := host[i+1 : l]
			if prePartIndex-i > 4 {
				return td
			} else if strings.Index(topDomains, td+",") > -1 {
				continue
			} else {
				return td
			}
		default:
			return host[i+1 : l]
		}
	}
	return host
}

// 获取url的域名
func GetDomain(purl string) string {
	firstIndex := 0
	indexCount := 0
	for i, l := 0, len(purl); i < l; i++ {
		if purl[i] == '/' {
			if indexCount == 1 {
				firstIndex = i
			} else if indexCount == 2 {
				return purl[firstIndex+1 : i]
			}
			indexCount++
		}
	}
	if indexCount == 0 {
		return purl
	}
	if indexCount == 2 {
		return purl[firstIndex+1 : len(purl)]
	}
	return ""
}