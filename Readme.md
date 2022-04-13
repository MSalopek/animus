# Animus

Animus is a backend tool (API, proxy) that makes interacting with IPFS easier. It was inspired by services sucha as [Pinata](https://www.pinata.cloud/), [web3.storage](https://web3.storage/) and [Fleek](https://fleek.co/).

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


# Main goal
This is a primarily a productive way for me to learn IPFS. I am hoping to reach feature parity with Pinata or Fleek in the near future. And use this in some other projects.


# How to run
```sh
$ docker-compose -f ./docker-compose.local.yml up
```

This command will initialize your devenv.

**The devenv consists of the following:**
## IPFS Node

* node is not bootstrapped
* node is not private
* to make it private generate add a swarm key via `ipfs-swarm-key-gen > ./devenv/ipfs/swarm.key`; add `LIBP2P_FORCE_PNET=1` env var in docker-compose
* node exposes localhost:5001 so you can use `webui` and has required HTTP headers set

**Using WEBUI**
There are some issues with running webui on non-bootstrapped nodes. You need to set the correct HTTP CORS headers and make the built version of `ipfs/webui` available.
* You can either pull `ipfs/webui` from IPFS and pin it to your node (this is not foolproof, for some reason I'm getting different CID for the build directory)
```sh
$ docker-compose exec ipfs sh
$ wget -qO- https://ipfs.io/ipfs/bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva | tar -xf -
$ ipfs add -r bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva
```

* Or you can pull and run `ipfs/webui` locally and use it that way (no extra config is required, just forward the port `:5001`)
```sh
$ git clone git@github.com:ipfs/ipfs-webui.git
$ cd ipfs-webui
$ npm install
$ npm start
```
The `npm start` command will open `localhost:3000` in your browser - you should see the webui if everything is installed properly.


You can also try building `ipfs-webui` and try to add the resulting `build` dir to your ipfs node.
```sh
$ git clone git@github.com:ipfs/ipfs-webui.git
$ cd ipfs-webui
$ npm run build
$ docker-compose exec ipfs ipfs add ./build
```

More info on: https://github.com/ipfs/ipfs-webui

## PostgreSQL DB
- create the database and required schemas
- default credentials for devenv: `host=localhost user=animus password=animus dbname=animus port=5432 sslmode=disable`

# Manual Testing of API methods
You can find the Postman Collection in `./postman`. Load it into your postman.
- try to login with `Animus API/public/login` to get JWT token using devenv credentials: `{"email": "admin@example.com", "password": "administrator"}`
- take the `token` and add it to the `Animus` env variable called `token`

# Disclamer
This is not production ready. You are using it at your own risk.
