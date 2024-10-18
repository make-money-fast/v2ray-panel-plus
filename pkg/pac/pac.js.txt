var domainList = %s ;
var subDomainList = %s;

var regexpList =  %s;

var direct = "DIRECT"
var proxy = "%s"


function FindProxyForURL(url, host) {
    if (isInNet(dnsResolve(host), "127.0.0.1", "255.255.255.255")) {
        return proxy;
    }
    if (domainList.hasOwnProperty(host)) {
        return proxy
    }
    var dst = direct
    subDomainList.forEach(function (e) {
        if (host.endsWith(e)) {
            dst = proxy
            return false
        }
    });
    regexpList.forEach(function (e) {
        if (e.test(host)) {
            dst = proxy
            return false
        }
    });

    return dst;
}