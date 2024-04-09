import {
    createMint,
    createAssociatedTokenAccount,
    mintTo,
    TOKEN_PROGRAM_ID,
  } from "@solana/spl-token";
  describe("Test transfers", () => {
    it("transferLamports", async () => {
      // Generate keypair for the new account
      const newAccountKp = new web3.Keypair();
      // Send transaction
  
      interface TransferData {
        amount: BN; // Assuming 'amount' is a number
        targetAddr: string; // Assuming 'message' is a string
      }
  
      const transferData: TransferData = {
        amount: new BN(1000), // Assuming 'amount' is a BigNumber
        targetAddr: "Hello, world!", // Assuming 'message' is a string
      };
  
   
      const txHash = await pg.program.methods
        .transferLamports(transferData)
        .accounts({
          from: pg.wallet.publicKey,
          to: newAccountKp.publicKey,
        })
        .signers([pg.wallet.keypair])
        .rpc();
      console.log(`https://explorer.solana.com/tx/${txHash}?cluster=devnet`);
      await pg.connection.confirmTransaction(txHash, "finalized");
      const newAccountBalance = await pg.program.provider.connection.getBalance(
        newAccountKp.publicKey
      );
      assert.strictEqual(
        newAccountBalance,
        transferData.amount.toNumber(),
        "The new account should have the transferred lamports"
      );
    });
  
    it("transferSplTokens", async () => {
      // Generate keypairs for the new accounts
      const fromKp = pg.wallet.keypair;
      const toKp = new web3.Keypair();
  
      // Create a new mint and initialize it
      const mintKp = new web3.Keypair();
      const mint = await createMint(
        pg.program.provider.connection,
        pg.wallet.keypair,
        fromKp.publicKey,
        null,
        0
      );
  
      // Create associated token accounts for the new accounts
      const fromAta = await createAssociatedTokenAccount(
        pg.program.provider.connection,
        pg.wallet.keypair,
        mint,
        fromKp.publicKey
      );
      console.log("fromAta" , fromAta.toBase58())
      const toAta = await createAssociatedTokenAccount(
        pg.program.provider.connection,
        pg.wallet.keypair,
        mint,
        toKp.publicKey
      );
      console.log("toAta" , fromAta.toBase58())
      // Mint tokens to the 'from' associated token account
      const mintAmount = 1000;
      await mintTo(
        pg.program.provider.connection,
        pg.wallet.keypair,
        mint,
        fromAta,
        pg.wallet.keypair.publicKey,
        mintAmount
      );
  
      // Send transaction
  
  
      interface TransferData {
        amount: BN; // Assuming 'amount' is a number
        targetAddr: string; // Assuming 'message' is a string
      }
  
      const transferData: TransferData = {
        amount: new BN(999), // Assuming 'amount' is a BigNumber
        targetAddr: "Hello, world!", // Assuming 'message' is a string
      };
  
      const txHash = await pg.program.methods
        .transferSplTokens(transferData)
        .accounts({
          from: fromKp.publicKey,
          fromAta: fromAta,
          toAta: toAta,
          tokenProgram: TOKEN_PROGRAM_ID,
        })
        .signers([pg.wallet.keypair, fromKp])
        .rpc();
      console.log(`https://explorer.solana.com/tx/${txHash}?cluster=devnet`);
      await pg.connection.confirmTransaction(txHash, "finalized");
      const toTokenAccount = await pg.connection.getTokenAccountBalance(toAta);
      assert.strictEqual(
        toTokenAccount.value.uiAmount,
        transferData.amount.toNumber(),
        "The 'to' token account should have the transferred tokens"
      );
    });
  });
  