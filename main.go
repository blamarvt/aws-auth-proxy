package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/smartystreets/go-aws-auth"
)

var proxyHost = flag.String("proxyHost", "", "hostname we're proxying to")
var proxySecure = flag.Bool("proxySecure", true, "if true we're proxying to https")
var listenHost = flag.String("listenHost", "localhost", "host/ip to listen on")
var listenPort = flag.Int("listenPort", 8080, "port to listen on")
var awsAccessKey = flag.String("accessKey", "", "AWS access key for signing request")
var awsSecretKey = flag.String("secretKey", "", "AWS secret key for signing request")

func require(errList *[]error, value *string, message string) {
	if value == nil || *value == "" {
		*errList = append(*errList, errors.New(message))
	}
}

func main() {
	flag.Parse()

	errList := []error{}

	require(&errList, proxyHost, "Please provide -proxyHost!")
	require(&errList, awsAccessKey, "Please provide -accessKey!")
	require(&errList, awsSecretKey, "Please provide -secretKey!")

	if len(errList) > 0 {
		for _, err := range errList {
			fmt.Println(err)
		}
		fmt.Println("\nFix the above error(s)!")
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		director := func(req *http.Request) {
			if *proxySecure {
				req.URL.Scheme = "https"
			}
			req.URL.Host = *proxyHost
			req.Host = *proxyHost
			awsauth.Sign4(req, awsauth.Credentials{
				AccessKeyID:     *awsAccessKey,
				SecretAccessKey: *awsSecretKey,
			})
		}
		proxy := &httputil.ReverseProxy{Director: director}
		proxy.ServeHTTP(w, r)
	})

	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf("%s:%d", listenHost, listenPort),
			nil,
		),
	)
}
