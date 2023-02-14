# http-cert-provisioner

A simple HTTP/HTTPS request handler to provision certs based on source IP address.

## Example:

Server side:
```
$ ./http-cert-provisioner
Path set to certs/
2023/02/14 07:52:03 Listening with HTTP on :1443 at /
2023/02/14 07:52:05 Trying to open file certs/localhost.pem
2023/02/14 07:52:05 retrieving file "certs/localhost.pem"...
2023/02/14 07:52:05 successfully transferred "certs/localhost.pem" for "[::1]:53952"
```

Client side request
```
$ curl -O http://localhost:1443/server.pem
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100   665    0   665    0     0  85049      0 --:--:-- --:--:-- --:--:-- 95000
```

## Usage

```
# http-cert-provisioner -h
HTTP-Cert-Provisioner (github.com/pschou/http-cert-provisioner, version: 0.1.20230214.0813)

This utility listens on a port and handles HTTP/HTTPS GET requests for
certificates; the listenPath determines which URL prefix is required, and the
extension from the filename will determine which CERTIFICATE file to send (PEM,
JKS, P12).  Note the CERTIFICATE file must have the file naming convention of
/certs/my.fqdn.com.pem to be able to be served up with a /server.pem request
coming from a server with the rDNS entry for the given ip pointing back to the
correct FQDN.  For security, this verifies forward DNS records to prevent rogue
rDNS entries pointing to unauthorized FQDNs.


Usage: ./http-cert-provisioner [options]
  -CA string
    	A PEM encoded CA's certificate file. (default "someCertCAFile")
  -cert string
    	A PEM encoded certificate file. (default "someCertFile")
  -key string
    	A PEM encoded private key file. (default "someKeyFile")
  -listen string
    	Where to listen to incoming connections (example 1.2.3.4:1443) (default ":1443")
  -listenPath string
    	Where to expect requests to be made ("/" -> "/server.pem") (default "/")
  -path string
    	Directory which to serve certs from (default "certs/")
  -tls
    	Enable TLS for secure transport
```
