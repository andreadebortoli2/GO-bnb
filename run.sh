#!/bin/bash

go build -o GO-bnb cmd/web/*.go && ./GO-bnb -dbname=go_bnb -dbuser=postgres -dbpassword=password -cache=false -production=false