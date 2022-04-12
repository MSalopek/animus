#! /bin/sh
# downloads WebUI assets from IPFS
wget -qO- bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva | tar -xf -
ipfs add -r bafybeihcyruaeza7uyjd6ugicbcrqumejf6uf353e5etdkhotqffwtguva
