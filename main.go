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
	list, err := loadV2rayDomainList()
	if err != nil {
		log.Fatalln(err)
	}
	if err := genSurgeRules(list); err != nil {
		log.Fatalln(err)
	}
	if err := genProxyList(list); err != nil {
		log.Fatalln(err)
	}
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

func genSurgeRules(list *router.GeoSiteList) error {
	for _, entry := range list.Entry {
		if strings.HasPrefix(entry.CountryCode, "CATEGORY-") ||
			strings.HasPrefix(entry.CountryCode, "GEOLOCATION-") ||
			strings.HasPrefix(entry.CountryCode, "TLD-") ||
			entry.CountryCode == "CN" {
			data := make([]byte, 0)
			for _, domain := range entry.Domain {
				switch domain.Type {
				case router.Domain_Plain:
					data = append(data, []byte("DOMAIN-KEYWORD,"+domain.Value+"\n")...)
				case router.Domain_Regex:
					data = append(data, []byte("URL-REGEX,"+domain.Value+"\n")...)
				case router.Domain_Domain:
					data = append(data, []byte("DOMAIN-SUFFIX,"+domain.Value+"\n")...)
				case router.Domain_Full:
					data = append(data, []byte("DOMAIN,"+domain.Value+"\n")...)
				}
			}
			err := ioutil.WriteFile("surge/"+strings.ToLower(entry.CountryCode), data, 0644)
			if err != nil {
				return fmt.Errorf("failed to generate surge rule %s", strings.ToLower(entry.CountryCode))
			}
		}
	}
	return nil
}

func genProxyList(list *router.GeoSiteList) error {
	src := []byte("[AutoProxy 0.2.9]\n")
	for _, entry := range list.Entry {
		prefix := ""
		if entry.CountryCode == "CN" {
			prefix = "@@||"
		} else if entry.CountryCode == "GEOLOCATION-!CN" {
			prefix = "||"
		} else {
			continue
		}
		for _, domain := range entry.Domain {
			src = append(src, []byte(prefix+domain.Value+"\n")...)
		}
	}
	data := base64.StdEncoding.EncodeToString(src)
	err := ioutil.WriteFile("proxylist.txt", []byte(data), 0644)
	if err != nil {
		return fmt.Errorf("failed to generate proxy list")
	}
	return nil
}
