module github.com/ethereum/go-ethereum

go 1.13

require (
	cuelang.org/go v0.1.1 // indirect
	github.com/Azure/azure-storage-blob-go v0.8.0
	github.com/VictoriaMetrics/fastcache v1.5.4
	github.com/alecthomas/jsonschema v0.0.2
	github.com/aristanetworks/goarista v0.0.0-20170210015632-ea17b1a17847
	github.com/aws/aws-sdk-go v1.25.48
	github.com/btcsuite/btcd v0.0.0-20171128150713-2e60448ffcc6
	github.com/cespare/cp v0.1.0
	github.com/cloudflare/cloudflare-go v0.10.7
	github.com/davecgh/go-spew v1.1.1
	github.com/deckarep/golang-set v0.0.0-20180603214616-504e848d77ea
	github.com/docker/docker v1.4.2-0.20180625184442-8e610b2b55bf
	github.com/dop251/goja v0.0.0-20200219165308-d1232e640a87
	github.com/edsrzf/mmap-go v0.0.0-20160512033002-935e0e8a636c
	github.com/elastic/gosigar v0.8.1-0.20180330100440-37f05ff46ffa
	github.com/emicklei/proto v1.9.0 // indirect
	github.com/etclabscore/go-openrpc-reflect v0.0.26
	github.com/fatih/color v1.6.0
	github.com/fjl/memsize v0.0.0-20180418122429-ca190fb6ffbc
	github.com/gballet/go-libpcsclite v0.0.0-20190607065134-2772fd86a8ff
	github.com/go-stack/stack v1.8.0
	github.com/go-test/deep v1.0.4
	github.com/golang/protobuf v1.3.2-0.20190517061210-b285ee9cfc6c
	github.com/golang/snappy v0.0.1
	github.com/google/uuid v1.1.1 // indirect
	github.com/gorilla/websocket v1.4.1
	github.com/graph-gophers/graphql-go v0.0.0-20191115155744-f33e81362277
	github.com/gregdhill/go-openrpc v0.0.1
	github.com/hashicorp/golang-lru v0.0.0-20160813221303-0a025b7e63ad
	github.com/huin/goupnp v0.0.0-20161224104101-679507af18f3
	github.com/iancoleman/strcase v0.0.0-20191112232945-16388991a334
	github.com/influxdata/influxdb v1.2.3-0.20180221223340-01288bdb0883
	github.com/jackpal/go-nat-pmp v1.0.2-0.20160603034137-1fa385a6f458
	github.com/julienschmidt/httprouter v1.2.0
	github.com/karalabe/usb v0.0.0-20191104083709-911d15fe12a9
	github.com/mattn/go-colorable v0.1.0
	github.com/mattn/go-isatty v0.0.5-0.20180830101745-3fb116b82035
	github.com/mitchellh/go-homedir v1.1.0
	github.com/naoina/toml v0.1.2-0.20170918210437-9fafd6967416
	github.com/olekukonko/tablewriter v0.0.2-0.20190409134802-7e037d187b0c
	github.com/onsi/ginkgo v1.8.0 // indirect
	github.com/onsi/gomega v1.5.0 // indirect
	github.com/open-rpc/meta-schema v0.0.42
	github.com/pborman/uuid v1.2.0
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/peterh/liner v1.1.1-0.20190123174540-a2c9a5303de7
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/tsdb v0.7.1
	github.com/rjeczalik/notify v0.9.1
	github.com/rs/cors v0.0.0-20160617231935-a62a804a8a00
	github.com/spf13/cobra v1.0.0
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2
	github.com/status-im/keycard-go v0.0.0-20190316090335-8537d3370df4
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570
	github.com/stretchr/testify v1.4.0
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	github.com/tidwall/gjson v1.6.0
	github.com/tidwall/pretty v1.0.0
	github.com/tyler-smith/go-bip39 v1.0.1-0.20181017060643-dbb3b84ba2ef
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd
	golang.org/x/text v0.3.2
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	golang.org/x/tools v0.0.0-20200414211825-33e937220d8f // indirect
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce
	gopkg.in/olebedev/go-duktape.v3 v3.0.0-20190213234257-ec84240a7772
	gopkg.in/urfave/cli.v1 v1.20.0
	gopkg.in/yaml.v3 v3.0.0-20200313102051-9f266ea9e77c // indirect

)

// see https://github.com/golang/lint/issues/436#issuecomment-482066447
replace github.com/golang/lint v0.0.0-20190409202823-959b441ac422 => github.com/golang/lint v0.0.0-20190409202823-5614ed5bae6fb75893070bdc0996a68765fdd275

// Use a local development version, managed as a submodule.
replace github.com/gregdhill/go-openrpc => github.com/etclabscore/go-openrpc v0.0.1

replace github.com/alecthomas/jsonschema => github.com/etclabscore/go-jsonschema-reflect v0.0.2

replace github.com/open-rpc/meta-schema => github.com/meowsbits/meta-schema v0.0.42

// replace github.com/etclabscore/go-openrpc-reflect => /home/ia/go/src/github.com/etclabscore/go-openrpc-reflect
