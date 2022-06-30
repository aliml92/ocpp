module github.com/aliml92/ocpp

go 1.18

require (
	github.com/google/uuid v1.3.0
	github.com/gorilla/websocket v1.5.0
	gopkg.in/go-playground/validator.v9 v9.31.0
)

require (
	github.com/go-playground/locales v0.14.0 // indirect
	github.com/go-playground/universal-translator v0.18.0 // indirect
	github.com/leodido/go-urn v1.2.1 // indirect
	gopkg.in/go-playground/assert.v1 v1.2.1 // indirect
)

retract (
    v0.2.0 // Published accidentally.
    v0.1.0 // Contains retractions only.
	v1.0.0 // Contains retractions only.
)