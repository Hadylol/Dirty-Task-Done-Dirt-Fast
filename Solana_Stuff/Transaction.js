const {
  Keypair,
  Transaction,
  LAMPORTS_PER_SOL,
  SystemProgram,
  sendAndConfirmTransaction,
  Connection,
  clusterApiUrl,
} = require("@solana/web3.js");
const bs58 = require("bs58");

// main(PublickKey, PrivateKey, SenderKey, sol)
async function main(PrivateKey, SenderKey, sol) {
  const PrivateKeyBytes = bs58.default.decode(PrivateKey);
  if (PrivateKeyBytes.length !== 64) {
    throw new Error("Invalid Private Key: Must be 64 bytes long.");
  }
  const keypair = Keypair.fromSecretKey(PrivateKeyBytes);

  const transaction = new Transaction().add(
    SystemProgram.transfer({
      fromPubkey: keypair.publicKey,
      toPubkey: SenderKey,
      lamports: sol * LAMPORTS_PER_SOL,
    })
  );
  const connection = new Connection(clusterApiUrl("devnet"));
  const signature = await sendAndConfirmTransaction(connection, transaction, [
    keypair,
  ]);

  console.log(signature);
}
const args = process.argv.slice(2); // Skip node and script path
if (args.length < 3) {
  console.log(
    "Usage: node script.js <PublicKey> <PrivateKey> <SenderKey> <Sol>"
  );
  process.exit(1);
}
const [PrivateKey, SenderKey, sol] = args;
main(PrivateKey, SenderKey, parseFloat(sol));
