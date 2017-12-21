# aqua
A basic aquarium monitoring toolkit written in go

- Temperature reading using a DS18B20 temperature probe
- Data collection running on a raspberry pi
- Publication to a centralized collector
- Prometheus scraping & display


Setup:

DS18B20 ---gpio---> Pi ----> w1toAqua ---coap POST--> aqua <--http GET-- prometheus

How to

- Connect a DS18B20 to a Raspberry PI
- start collector process on any machine: `go run aqua.go <coap-port> <scraping-port>`
- start probe process on the PI: `go run w1toaqua.go <ip:port>`
- build a prometheus container
  - adjust collector scraping ip:port in `prometheus.yaml`
  - build prometheus container: `docker build -t prometheus/aqua-monitor .`
  - start prometheus: `docker run -p 9090:9090 --name prometheus-aqua -d prometheus/aqua-monitor`
  - connect to prometheus host on port 9090 to view data



