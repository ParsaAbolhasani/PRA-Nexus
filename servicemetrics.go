package service

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    // Blockchain metrics
    LatestBlock = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "pra_latest_block",
        Help: "Latest block height",
    })
    
    MempoolSize = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "pra_mempool_size",
        Help: "Number of pending transactions in mempool",
    })
    
    P2PPeers = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "pra_p2p_peers",
        Help: "Number of connected P2P peers",
    })
    
    TransactionsProcessed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_transactions_processed_total",
        Help: "Total number of processed transactions",
    })
    
    BlockTime = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "pra_block_time_seconds",
        Help: "Time between blocks in seconds",
    })
    
    GasPriceGwei = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "pra_gas_price_gwei",
        Help: "Current gas price in Gwei",
    })
    
    // Security metrics
    InvalidSignatures = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_invalid_signature_total",
        Help: "Total number of invalid signature attempts",
    })
    
    DoubleSpendAttempts = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_double_spend_attempts_total",
        Help: "Total number of double-spend attempts detected",
    })
    
    RateLimitHits = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_rate_limit_hits_total",
        Help: "Total number of rate limit hits",
    })
    
    ContractReverts = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_contract_reverts_total",
        Help: "Total number of smart contract reverts",
    })
    
    ContractCalls = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_contract_calls_total",
        Help: "Total number of smart contract calls",
    })
    
    ChainReorgs = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_chain_reorgs_total",
        Help: "Total number of chain reorganizations",
    })
    
    OrphanBlocks = promauto.NewCounter(prometheus.CounterOpts{
        Name: "pra_orphan_blocks_total",
        Help: "Total number of orphan blocks received",
    })
)
