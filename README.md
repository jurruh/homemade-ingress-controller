# Homemade Ingress controller
An attempt to make my own Kubernetes ingress controller. The purpose of this project was mainly to get started with the Go programming language wich I had no experience in before this project. As an additionial benefit I tought it was cool to know something about how ingress controllers work under the hood.

## Features
- Reading all Ingresses from the Kubernetes api
- Reverse proxy a request to the right service according the `Host:` header

## Todo
- Support for configurable paths, not just `/`
- Support for multiple paths
- Autowatch for Ingress changes
- HTTPS with automatic certificate generation