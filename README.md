# Proxy List

根据 [v2fly/domain-list-community](https://github.com/v2fly/domain-list-community) 自动生成 Surge 规则和 GFW 列表。

## Surge

在 Surge 配置文件中添加 RULE-SET 即可

```text
RULE-SET,https://raw.githubusercontent.com/lijinglin3/proxylist/master/surge/category-ads-all,REJECT
RULE-SET,https://raw.githubusercontent.com/lijinglin3/proxylist/master/surge/cn,DIRECT
RULE-SET,https://raw.githubusercontent.com/lijinglin3/proxylist/master/surge/geolocation-!cn,PROXY
```

## GFW List

与 `gfwlist` 的 [使用方法](https://github.com/gfwlist/gfwlist/wiki) 一致，下载链接 [https://raw.githubusercontent.com/lijinglin3/proxylist/master/proxylist.txt](https://raw.githubusercontent.com/lijinglin3/proxylist/master/proxylist.txt)
