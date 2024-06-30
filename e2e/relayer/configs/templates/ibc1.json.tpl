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
    "gas_estimate_rate": {
      "numerator": 1,
      "denominator": 1
    },
    "max_gas_limit": 10000000,
    "tx_type": "legacy",
    "abi_paths": ["./abis"],
    "allow_lc_functions": {
      "lc_address": "$QBFT_CLIENT_ADDRESS",
      "allow_all": true
    }
  },
  "prover": {
    "@type": "/relayer.provers.qbft.config.ProverConfig",
    "consensus_type": "$CONSENSUS_TYPE",
    "trusting_period": "336h",
    "max_clock_drift": "30s"
  }
}
