const {
  Keypair,
  Connection,
  clusterApiUrl,
  PublicKey,
  LAMPORTS_PER_SOL,
} = require("@solana/web3.js");
const bs58 = require("bs58");
async function main(value) {
  if (value !== null) {
    const connection = new Connection(clusterApiUrl("devnet"));
    try {
      const pubkey = new PublicKey(value);
      const balance = await connection.getBalance(pubkey);

      console.log(`Balance of ${value}: ${balance / LAMPORTS_PER_SOL} sol`);
    } catch (e) {
      console.error(`Failed to get balance of account ${value}:`, e.message);
    }
    return;
  }
  let walletTypeShit = Keypair.generate();

  console.log("Public Key " + walletTypeShit.publicKey.toBase58());
  console.log("Secret Key : " + bs58.default.encode(walletTypeShit.secretKey));
}
const args = process.argv.slice(2);
let PublicKEy = args[0] || null;

main(PublicKEy);
