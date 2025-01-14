module my-modus-app

go 1.23.3

require (
	github.com/google/uuid v1.6.0
	github.com/hypermodeinc/modus/sdk/go v0.16.0
)

replace my-modus-app/graphgen => ./src/graphgen

require (
	github.com/tidwall/gjson v1.18.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tidwall/sjson v1.2.5 // indirect
	golang.org/x/crypto v0.32.0
)
