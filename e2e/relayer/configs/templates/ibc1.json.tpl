{
  "chain": {
    "@type": "/relayer.chains.ethereum.config.ChainConfig",
    "chain_id": "ibc1",
    "eth_chain_id": 3018,
    "rpc_addr": "http://localhost:8645",
    "signer": {
      "@type": "/relayer.chains.ethereum.signers.hd.SignerConfig",
      "mnemonic": "math razor capable expose worth grape metal sunset metal sudden usage scheme",
      "path": "m/44'/60'/0'/0/0"
    },
    "ibc_address": "$IBC_ADDRESS",
    "initial_send_checkpoint": 1,
    "initial_recv_checkpoint": 1,
    "enable_debug_trace": false,
    "average_block_time_msec": 1000,
    "max_retry_for_inclusion": 5,
    "allow_lc_functions": {
      "lc_address": "$IBFT2_CLIENT_ADDRESS",
      "allow_all": true
    }
  },
  "prover": {
    "@type": "/relayer.provers.ibft2.config.ProverConfig",
    "trust_level_numerator": 1,
    "trust_level_denominator": 3,
    "trusting_period": 1209600
  }
}
