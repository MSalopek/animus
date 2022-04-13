#! /bin/sh
# downloads WebUI assets from IPFS and adds them to locally available node
wget -qO- https://ipfs.io/ipfs/bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva | tar -xf -
ipfs add -r bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva
