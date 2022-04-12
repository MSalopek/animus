#! /bin/sh
# downloads WebUI assets from IPFS
wget -qO- https://ipfs.io/api/v0/get/QmdbrYy1BtSg2Kg6K5VQbLH5vPcNdSLYe8bb2YvXsrxx5V | tar -xf -
ipfs add -r bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva
