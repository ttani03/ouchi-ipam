# ouchi-ipam

"ouchi-ipam" is a simple IPAM (IP Address Management) service written in Go.

## Description

"ouchi-ipam" allows you to manage subnets and their respective IP addresses. The service provides a REST API for communication and stores the subnet and IP address information in a PostgreSQL database.

## Usage

An easy way to start "ouchi-ipam" is to use Docker Compose:

```
docker compose up --build -d
```

You can access the documentation of the REST API by navigating to `http://localhost:8080/swagger/index.html` from a browser.

First, you need to create a new subnet:

```
$ curl -X POST -H "Content-type: application/json" http://localhost:8080/subnets -d@- <<EOF
{
  "name": "web",
  "cidr": "192.168.0.0/24",
  "gateway": "192.168.0.1",
  "name_server": "8.8.8.8",
  "description": "subnet for web"
}
EOF
{"id":1,"name":"web","cidr":"192.168.0.0/24","gateway":"192.168.0.1","name_server":"8.8.8.8","description":"subnet for web"}
```

After that, reserve an IP address by the subnet ID:

```
$ curl -X POST http://localhost:8080/subnets/1/ip
{"id":1,"subnet_id":1,"address":"192.168.0.1","hostname":""}
```

You can free a reserved IP address by the following way:

```
$ curl -X DELETE http://localhost:8080/subnets/1/ip/192.168.0.1
```
