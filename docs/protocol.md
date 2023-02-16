
# MaxMind GeoIPUpdate Protocol Design

This article analyze <https://github.com/maxmind/geoipupdate> protocol stack design

## HTTP API

### Request

```http
GET /geoip/databases/<edition-id>/update?db_md5=<local-hash>
Host: <host>
User-Agent: geoipupdate/<version>
Authentication: Basic <account-id:license-key>
```

Variables:

- `<host>`, The field is from [configuration file][conf-host]
- `<edition-id>`, The field is escaped path, see [configuration file][conf-edition]
- `<local-hash>`, The field is requester local latest hash ([MD5][hash-md5])
- `<version>`, The field is updater version
- `<account-id:license-key>`, The field is from [configuration file][conf-license]

### Response: [200][status-200] OK

Response Body:

**MUST** use gzipped compression

Response Headers:

- `Content-Encoding: gzip`: The updater **NOT CHECK** this field, but suggest set the field.
- `X-Database-Md5: <new-md5>`: This is **MUST** field, The is new hash ([MD5][hash-md5]), use hexadecimal encoded
- `Last-Modified: <last-modified>`: This is **MUST** field, The is last-modified time, format [RFC1123][rfc1123]

Response Body: ip database file binary \
(_NOT VERIFY PAYLOAD FORMAT_ is [MaxMind Database file format][maxmind-spec])

### Response: [304][status-304] Not Modified

Not Modified

### Response: [30x][status-30x] Redirection

Use new location download update

see Golang builtin HTTP client supports auto follow redirects: <https://pkg.go.dev/net/http#Client>

### Response: Not OK

Download first 256 bytes on error messsage

[conf-edition]: https://github.com/maxmind/geoipupdate/blob/v6.0.0/conf/GeoIP.conf.default#L11-L13
[conf-host]: https://github.com/maxmind/geoipupdate/blob/v6.0.0/conf/GeoIP.conf.default#L20-L21
[conf-license]: https://github.com/maxmind/geoipupdate/blob/v6.0.0/conf/GeoIP.conf.default#L5-L9
[hash-md5]: https://en.wikipedia.org/wiki/MD5
[maxmind-spec]: https://github.com/maxmind/MaxMind-DB
[rfc1123]: https://datatracker.ietf.org/doc/html/rfc1123
[status-200]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/200
[status-304]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/304
[status-30x]: https://developer.mozilla.org/en-US/docs/Web/HTTP/Redirections
