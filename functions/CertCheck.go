package functions

import (
	"flag"
	"fmt"
	"os"
	"time"
)

var (
	host       = flag.String("host", "", "Comma-separated hostnames and IPs to generate a certificate for")
	validFrom  = flag.String("start-date", "", "Creation date formatted as Jan 1 15:04:05 2011")
	validFor   = flag.Duration("duration", 365*24*time.Hour, "Duration that certificate is valid for")
	isCA       = flag.Bool("ca", false, "whether this cert should be its own Certificate Authority")
	rsaBits    = flag.Int("rsa-bits", 2048, "Size of RSA key to generate. Ignored if --ecdsa-curve is set")
	ecdsaCurve = flag.String("ecdsa-curve", "", "ECDSA curve to use to generate a key. Valid values are P224, P256 (recommended), P384, P521")
	ed25519Key = flag.Bool("ed25519", false, "Generate an Ed25519 key")
)

func CertCheck() string {
	certCheckOk := false
	certDir := "/tmp/cert"
	if _, err := os.Stat("/cert/key.pem"); err == nil {
		fmt.Printf("/cert/key.pem exists\n")
		certCheckOk = true
	} else {
		fmt.Printf("/cert/key.pem not exist\n")
		certCheckOk = false
	}
	if certCheckOk == true {
		if _, err := os.Stat("/cert/cert.pem"); err == nil {
			fmt.Printf("/cert/cert.pem exists\n")
			certCheckOk = true
		} else {
			fmt.Printf("/cert/cert.pem not exist\n")
			certCheckOk = false
		}
	}

	if certCheckOk == true {
		certDir = "/cert"
	} else {
		fmt.Printf("Self certs will be used\n")
		certDir = "/tmp/cert"
		SslCertGenerate()
	}
	return certDir

}
