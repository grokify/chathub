module github.com/grokify/chathub

go 1.25.5

require (
	github.com/agentplexus/mcpkit v0.3.2
	github.com/grokify/omnistorage v0.2.1
	github.com/grokify/omnistorage-github v0.1.3
	github.com/modelcontextprotocol/go-sdk v1.4.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/google/go-github/v82 v82.0.0 // indirect
	github.com/google/go-querystring v1.2.0 // indirect
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/grokify/gogithub v0.9.1 // indirect
	github.com/grokify/mogo v0.73.2 // indirect
	github.com/inconshreveable/log15 v3.0.0-testing.5+incompatible // indirect
	github.com/inconshreveable/log15/v3 v3.1.0 // indirect
	github.com/jpillora/backoff v1.0.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.3 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.ngrok.com/muxado/v2 v2.0.1 // indirect
	golang.ngrok.com/ngrok v1.13.0 // indirect
	golang.org/x/net v0.50.0 // indirect
	golang.org/x/oauth2 v0.35.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/term v0.40.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// Force log15/v3 to version compatible with ngrok v1.12.0 (has ext.RandId)
replace github.com/inconshreveable/log15/v3 => github.com/inconshreveable/log15/v3 v3.0.0-testing.5
