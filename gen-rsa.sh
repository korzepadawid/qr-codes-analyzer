#!/bin/bash

echo "Generating RSA key pair..."

openssl genrsa -out ./token/keys/rsa.private 1024 && \
openssl rsa -in ./token/keys/rsa.private -out ./token/keys/rsa.public -pubout -outform PEM

if [ $? -eq 0 ]
then
	echo "Done!"
else
	echo "failed!"
fi
