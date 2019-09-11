# FirstMicroService[![Go Report Card](https://goreportcard.com/badge/github.com/DarthSett/Restaurant)](https://goreportcard.com/report/github.com/DarthSett/FirstMicroService) [![CircleCI](https://circleci.com/gh/DarthSett/FirstMicroService/tree/master.svg?style=svg)](https://circleci.com/gh/DarthSett/FirstMicroService/tree/master)

This is an app consisting of two microservices communicating among each other running on seperate docker containers. Kubernetes config files are also available.

The first microservice is a csv-uploader that accepts a csv consisting of amazon links from the user and saves the links to a DB and sendsa signal to the second microservice when its done saving the links.

The second microservice is a scraper that scrapes the name, price and seller from the links in db and stores the scraped info in another table.
