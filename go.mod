module github.com/soluble1/mweb

go 1.19

require (
	github.com/google/uuid v1.3.0
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/redis/go-redis/v9 v9.0.2
	github.com/soluble1/mcache v0.0.0-20230330080450-480dbd29c7f1
)

require (
	github.com/bits-and-blooms/bitset v1.5.0 // indirect
	github.com/bits-and-blooms/bloom/v3 v3.3.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gotomicro/ekit v0.0.6 // indirect
	golang.org/x/sync v0.1.0 // indirect
)

replace github.com/soluble1/mcache v0.0.0-20230330080450-480dbd29c7f1 => ../cache
