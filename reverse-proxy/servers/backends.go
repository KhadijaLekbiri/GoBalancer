package servers

import(
	"net/url"
	"reverse-proxy/services/models"
	"log")


//  This Must function get ride of the err and only return the url we want

func Must(rawURL string) (*url.URL){
	url, err := url.Parse(rawURL);
	if err != nil {
		log.Fatal("hhh",rawURL,err)
	} 
	return url
}

var backends = []*models.Backend{{
				URL: Must("http://localhost:8082"),
				Alive: false,
				CurrentConns: 0,
			},
			{
				URL: Must("http://localhost:8083"),
				Alive: true,
				CurrentConns: 0,
			},
			{
				URL: Must("http://localhost:8084"),
				Alive: false,
				CurrentConns: 0,
			},
			{
				URL: Must("http://localhost:8085"),
				Alive: true,
				CurrentConns: 0,
			},
		}

func Backends() ([]*models.Backend){
	return backends
}