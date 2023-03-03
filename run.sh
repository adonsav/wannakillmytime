#!/bin/bash

go build -o fgoapp cmd/fgoapp/*.go
./fgoapp -dbname=wannakillmytime -dbuser=adonsav -cache=false -production=false