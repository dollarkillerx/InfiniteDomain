package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
)

var tagUrl = "https://rr2---sn-i3belney.googlevideo.com/videoplayback?c=WEB&clen=58204144&dur=845.902&ei=s3TvYtCOB9OZvcAPxaeYkAM&expire=1659881747&fexp=24001373%2C24007246&fvip=3&gir=yes&id=o-ABkIZjACsTJHLX1TQP8jd95DV0gKp0GIV3N5PfowRWmt&initcwndbps=1086250&ip=43.135.75.195&itag=18&lmt=1657781415028762&lsig=AG3C_xAwRAIgI7DY5EbpmU_tkrMFx0u3hPnroLtDgTFx_MHAFRmqm1UCIH6CItIpbmu5L0_XUTKJDNuXe0EnL82Eqb7icohs-ykA&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&mh=uz&mime=video%2Fmp4&mm=31%2C29&mn=sn-i3belney%2Csn-i3b7knzl&ms=au%2Crdu&mt=1659859727&mv=m&mvi=2&n=DrnOwUYP_JWSAQ&ns=RnqdAFmfk6fIMHrNOTDO_XAH&pl=18&ratebypass=yes&rbqsm=fr&requiressl=yes&sig=AOq0QJ8wRQIgYTqxP6Pievb9JjxtBevto8iCA-HyLxzf7H9JQ7MiPBQCIQDt-mUWmaTOj1JTihK69yCsGlGyx7ArdlA2vfxpN33AmA%3D%3D&source=youtube&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cspc%2Cvprv%2Cmime%2Cns%2Cgir%2Cclen%2Cratebypass%2Cdur%2Clmt&spc=lT-KhvxVXufkE154JEjC7ZKu1tMeL3E&txp=5538434&vprv=1"

func main() {
	parseUrl, err := url.Parse(tagUrl)
	if err != nil {
		log.Fatalln(err)
	}

	http.HandleFunc("/test", func(writer http.ResponseWriter, request *http.Request) {
		req := &http.Request{
			URL:        parseUrl,
			Method:     "GET",
			Header:     make(http.Header),
			Proto:      "HTTP/1.1",
			ProtoMajor: 1,
			ProtoMinor: 1,
			Close:      true,
		}

		req.Header = request.Header

		client := http.Client{}
		do, err := client.Do(req)
		if err != nil {
			log.Fatalln(err)
		}
		defer do.Body.Close()

		for k, v := range do.Header {
			for _, hk := range v {
				writer.Header().Add(k, hk)
			}
		}

		writer.WriteHeader(do.StatusCode)
		io.Copy(writer, do.Body)
	})

	if err := http.ListenAndServe("0.0.0.0:8574", nil); err != nil {
		log.Fatalln(err)
	}
}
