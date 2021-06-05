package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/v2fly/v2ray-core/v4/app/router"
	"google.golang.org/protobuf/proto"
)

func main() {
	if err := genSurgeRules(); err != nil {
		log.Fatalln(err)
	}
	if err := genProxyList(); err != nil {
		log.Fatalln(err)
	}
}

func genSurgeRules() error {
	v2rayDomainList, err := loadV2rayDomainList()
	if err != nil {
		return err
	}
	for _, entry := range v2rayDomainList.Entry {
		if strings.HasPrefix(entry.CountryCode, "CATEGORY") ||
			strings.HasPrefix(entry.CountryCode, "GEO") ||
			entry.CountryCode == "CN" {
			data := parseV2rayDomainListToSurgeList(entry.Domain)
			err = ioutil.WriteFile("./surge/"+strings.ToLower(entry.CountryCode), data, 0644)
			if err != nil {
				return fmt.Errorf("failed to generate surge rule %s", strings.ToLower(entry.CountryCode))
			}
		}
	}
	return nil
}

func loadV2rayDomainList() (*router.GeoSiteList, error) {
	resp, err := http.Get("https://github.com/v2fly/domain-list-community/releases/latest/download/dlc.dat")
	if err != nil {
		return nil, fmt.Errorf("failed to download dlc.dat, err: %+v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read dlc.dat, err: %+v", err)
	}
	list := new(router.GeoSiteList)
	err = proto.Unmarshal(body, list)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dlc.dat, err: %+v", err)
	}
	return list, nil
}

func parseV2rayDomainListToSurgeList(domains []*router.Domain) []byte {
	ret := make([]byte, 0)
	for _, domain := range domains {
		switch domain.Type {
		case router.Domain_Plain:
			ret = append(ret, []byte("DOMAIN-KEYWORD,"+domain.Value+"\n")...)
		case router.Domain_Regex:
			ret = append(ret, []byte("URL-REGEX,"+domain.Value+"\n")...)
		case router.Domain_Domain:
			ret = append(ret, []byte("DOMAIN-SUFFIX,"+domain.Value+"\n")...)
		case router.Domain_Full:
			ret = append(ret, []byte("DOMAIN,"+domain.Value+"\n")...)
		default:
		}
	}
	return ret
}

func genProxyList() error {
	gfwList, err := loadGfwList()
	if err != nil {
		return err
	}
	proxyList, err := ioutil.ReadFile("proxy.list")
	if err != nil {
		return err
	}
	src := append(gfwList, proxyList...)
	data := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(data, src)
	return ioutil.WriteFile("proxylist.txt", data, 0644)
}

func loadGfwList() ([]byte, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/gfwlist/gfwlist/master/gfwlist.txt")
	if err != nil {
		return nil, fmt.Errorf("failed to download gfwlist.txt, err: %+v", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read gfwlist.txt, err: %+v", err)
	}
	return base64.StdEncoding.DecodeString(string(body))
}
