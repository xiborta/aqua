# aqua
A basic aquarium monitoring toolkit written in go

- Temperature reading using a DS18B20 temperature probe
- Data collection running on a raspberry pi
- Publication to a centralized collector
- Prometheus scraping & display


Setup:

DS18B20 -> Pi -> w1toAqua ---coap--> aqua <--- prometheus


How to

- Connect a DS18B20 to a Raspberry PI
- start collector process: 'go run aqua.go <port>'
- start probe process: 'go run w1toaqua.go <ip:port>'
- build a prometheus container
  - adjust collector scraping ip:port in 'prometheus.yaml'
  - build prometheus container: 'docker build -t prometheus/aqua-monitor .'
  - start prometheus: 'docker run -p 9090:9090 --name prometheus-aqua -d prometheus/aqua-monitor'
  - connect to prometheus host on port 9090 and 



