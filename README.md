###Quickstart
docker build -t dns-proxy .
docker run -p 53:53 -p 53:53/udp dns-proxy

###How It works
Simple, single-file, dns proxy written in Go using miekg/dns library.
Listens on port 53 for UDP and TCP. TCP traffic will be encrypted as DNS-Over-TLS
when sent to upstream resolver. You can change the upstream resolver in the config file.
I played around a bit with the Viper library. You can see how it could be further extended
to use environment variables. It also has support for popular config services.

###Resources And Notes
https://routley.io/tech/2017/12/28/hand-writing-dns-messages.html

Looked through java examples and felt like the libraries were lacking

Found a very comprehensive Golang library with a good example and opted for that
https://github.com/fffaraz/microdns-proxy/blob/master/microdns-proxy.go
https://github.com/miekg/dns

Found good examples of using crypto/tls Golang package
https://github.com/denji/golang-tls

Padding in tcp response looked like a bug and I took steps to eliminate until I read this:
https://edns0-padding.org/

Used Viper with TOML files, for config, based on recommendation from this blog
https://medium.com/@felipedutratine/manage-config-in-golang-to-get-variables-from-file-and-env-variables-33d876887152

Researched best way to bundle golang apps into Docker. Decided to use Glide
https://github.com/Masterminds/glide
https://blog.hasura.io/the-ultimate-guide-to-writing-dockerfiles-for-go-web-apps-336efad7012c

####Security Concerns:
I am not a security expert and thus you shouldn’t trust software that I’ve written to be completely secure. Use a large FOSS project like coredns instead. It even has a forwarding plugin for exactly this purpose
Certs are too easily accessible (Literally uploaded them to GitHub)

####In a Microservices Environment
This could be used encrypt dns traffic of applications needing dns resolution
Could assist in service discovery of external/3rd party services

####How to improve:
Grab certs from system instead of bundling them
Add support for DNS-over-HTTPS
Externalize configuration with support for config server
Write tests
Use a better multi-stage docker build to reduce size of final image
Use smaller base image so initial builds don’t take so long
