module example

replace core-v => ..

go 1.25.0

require core-v v0.0.0-00010101000000-000000000000

require (
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	sigs.k8s.io/kind v0.31.0 // indirect
)
