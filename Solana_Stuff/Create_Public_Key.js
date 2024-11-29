const { Keypair } = require("@solana/web3.js");
const bs58 = require("bs58");
function main(value) {
  let walletTypeShit = Keypair.generate();

  console.log("Public Key " + walletTypeShit.publicKey.toBase58());
  console.log("Secret Key : " + bs58.encode(walletTypeShit.secretKey));
}
main();
