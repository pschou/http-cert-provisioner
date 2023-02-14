package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var (
	about = `HTTP-Cert-Provisioner

This utility listens on a port and handles HTTP/HTTPS GET requests for
certificates; the listenPath determines which URL prefix is required, and the
extension from the filename will determine which CERTIFICATE file to send (PEM,
JKS, P12).  Note the CERTIFICATE file must have the file naming convention of
/certs/my.fqdn.com.pem to be able to be served up with a /server.pem request
coming from a server with the rDNS entry for the given ip pointing back to the
correct FQDN.  For security, this verifies forward DNS records to prevent rogue
rDNS entries pointing to unauthorized FQDNs.
`
	basePath   = flag.String("path", "certs/", "Directory which to serve certs from")
	listen     = flag.String("listen", ":1443", "Where to listen to incoming connections (example 1.2.3.4:1443)")
	listenPath = flag.String("listenPath", "/", "Where to expect requests to be made (\"/\" -> \"/server.pem\")")
	enableTLS  = flag.Bool("tls", false, "Enable TLS for secure transport")
	version    = ""
)

func main() {
	flag.Usage = func() {
		lines := strings.SplitN(about, "\n", 2)
		fmt.Fprintf(os.Stderr, "%s (github.com/pschou/http-cert-provisioner, version: %s)\n%s\n\nUsage: %s [options]\n",
			lines[0], version, lines[1], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()
	if *enableTLS {
		loadTLS()
	}
	fmt.Println("Path set to", *basePath)

	http.HandleFunc("/", requestHandler)
	if *enableTLS {
		log.Println("Listening with HTTPS on", *listen, "at", *listenPath)
		server := &http.Server{Addr: *listen, TLSConfig: tlsConfig}
		log.Fatal(server.ListenAndServeTLS(*certFile, *keyFile))
	} else {
		log.Println("Listening with HTTP on", *listen, "at", *listenPath)
		log.Fatal(http.ListenAndServe(*listen, nil))
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	var filename string
	defer func() {
		if err != nil {
			log.Println("Error:", err)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		r.Body.Close()
	}()

	switch r.Method {
	case "GET":
	default:
		err = fmt.Errorf("Method not handled %q", r.Method)
		return
	}

	// Flatten any of the /../ junk
	filename = filepath.Clean(r.URL.Path)

	// Verify that the right path is being hit on the request
	if !strings.HasPrefix(filename, *listenPath) || strings.HasPrefix(filename, "..") {
		err = fmt.Errorf("Path not allowed %q", filename)
		return
	}

	// Build the CN from the remote IP address
	var host string
	host, _, err = net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		err = fmt.Errorf("Could not determine IP from %q", r.RemoteAddr)
		return
	}

	// Look up FQDN from IP
	var hosts []string
	hosts, err = net.LookupAddr(host)
	if err != nil {
		err = fmt.Errorf("Could not determine host from %q", r.RemoteAddr)
		return
	}

	// Do a forward DNS lookup to verify the IP is a sane destination, to avoid rDNS spoofing!
	var any bool
hosts_lookup:
	for _, hostPrefix := range hosts {
		addrs, _ := net.LookupIP(hostPrefix)
		for _, addr := range addrs {
			if addr.String() == host {
				any = true
				break hosts_lookup
			}
		}
	}
	if !any {
		err = fmt.Errorf("Host missing forward DNS entry %q, disallowing request", r.RemoteAddr, "->", hosts)
		return
	}

	// Get file extension from request
	ext := filepath.Ext(r.URL.Path)
	if ext == "" {
		err = fmt.Errorf("Could not determine ext from request %q", r.URL.Path)
		return
	}

	// Open the file for reading
	var fh *os.File
	for _, hostPrefix := range hosts {
		// Build the exact path to where the cert should be
		filename = path.Join(*basePath, strings.TrimSuffix(hostPrefix, ".")+ext)

		log.Println("Trying to open file", filename)

		// Try opening the file
		if fh, err = os.Open(filename); err == nil {
			break
		}
	}
	if err != nil {
		return
	}
	defer fh.Close()

	log.Printf("retrieving file %q...\n", filename)

	// Copy the stream to disk
	if _, err = io.Copy(w, fh); err != nil {
		return
	}

	log.Printf("successfully transferred %q for %q\n", filename, r.RemoteAddr)

	return
}
