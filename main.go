package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/miekg/dns"
	"github.com/spf13/viper"
)

var certs = &tls.Config{}

//grab certs, grab config, and start dns servers
func main() {
	loadConfig()
	certs = loadCerts()

	go serve("udp", nil, dnsUDPHandler)
	go serve("tcp", certs, dnsTCPHandler)

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case s := <-sig:
			log.Fatalf("Signal (%d) received, stopping\n", s)
		}
	}
}

func serve(net string, tls *tls.Config, handlerFunc func(w dns.ResponseWriter, m *dns.Msg)) {
	srv := &dns.Server{Addr: ":53", Net: net, TLSConfig: tls, Handler: dns.HandlerFunc(handlerFunc)}
	log.Printf("Starting %s dns server", net)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to set %v listener %s\n", net, err.Error())
	}
}

func dnsHandler(w dns.ResponseWriter, m *dns.Msg, net string, target string) {
	log.Printf("Request for %s\n", m.Question[0].Name)
	m.Question[0].Name = strings.ToUpper(m.Question[0].Name)
	c := new(dns.Client)
	c.Net = net
	if strings.Contains(net, "tls") {
		c.TLSConfig = certs
	}
	r, _, _ := c.Exchange(m, target)
	r.Question[0].Name = strings.ToLower(r.Question[0].Name)
	for i := 0; i < len(r.Answer); i++ {
		r.Answer[i].Header().Name = strings.ToLower(r.Answer[i].Header().Name)
	}
	w.WriteMsg(r)
}

func loadCerts() *tls.Config {
	cert, err := tls.LoadX509KeyPair(viper.GetString("tls.crt"), viper.GetString("tls.key"))
	if err != nil {
		log.Fatalf("unable to build certificate: %s\n", err)
	}
	certs = &tls.Config{
		Certificates: []tls.Certificate{cert},
	}
	return certs
}

func loadConfig() {
	if os.Getenv("ENVIRONMENT") != "PRODUCTION" {
		viper.SetConfigFile("config.toml")
		viper.AddConfigPath(".")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s", err))
		}
	} else {
		viper.AutomaticEnv() //use environment variables instead
	}
}

func dnsUDPHandler(w dns.ResponseWriter, m *dns.Msg) {
	dnsHandler(w, m, "udp", viper.GetString("udp.addr")+":"+strconv.Itoa(viper.GetInt("udp.port")))
}

func dnsTCPHandler(w dns.ResponseWriter, m *dns.Msg) {
	dnsHandler(w, m, "tcp-tls", viper.GetString("tcp.addr")+":"+strconv.Itoa(viper.GetInt("tcp.port")))
}
