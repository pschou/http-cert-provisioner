# http-cert-provisioner

A simple HTTPS request handler to provision certs based on source IP address.

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
