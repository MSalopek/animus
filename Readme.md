# Animus

Animus is an tool that makes interactions with IPFS easier. It was inspired by services sucha as Pinata, web3.storage and Fleek.

## Purpose
The main purpose of the tool is to allow users easier control over their IPFS data. Not everyone wanting to use IPFS is a software engineer and most users do not need to know the details of IPFS to be able to use it effectievly.

## Main features
Animus is essentially a API that accepts file uploads, caches them locally, uploads them to IPFS and allows users to seamlessly access their files for read and write operations.

Although it brings about some centralization and removes some anonimity (users must be registered) it handles abstracts some difficulties when working with IPFS.

## How it works
TLDR; Animus is a type of a proxy.

1. User uploads a file the "normal way" using HTTP
- users must register

2. Animus uploads the file to IPFS
- uploads can be done to public or to private IPFS

3. Users can access their files via the Animus API

TODO:
- gateways
- private IPFS Nodes
- advanced file distribution schemes (pin on multiple nodes)
- develop CDN leveraging IPFS and "normal" cloud storage


## Main goal
This is a primarily a productive way for me to learn IPFS. I am hoping to reach feature parity with Pinata or Fleek in the near future. And use this in some other projects.


# How to run
```sh
$ docker-compose up
```

This command will create the database and required schemas and run the proxy.
Update the variables in the `docker-compose.yml` to make is suitable to your needs.
