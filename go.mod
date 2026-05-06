module github.com/envault/envault

go 1.22

require filippo.io/age v1.2.0

require (
	golang.org/x/crypto v0.26.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
)

retract (
	// Retract versions prior to module path correction
	v0.0.0
)
