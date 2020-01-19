package main

import "testing"

import "strings"

func TestJsonDecoder(t *testing.T) {
	json := `{"Issues":null,"Report":{"Linters":[{"Name":"govet","EnabledByDefault":true},{"Name":"bodyclose"},{"Name":"errcheck","EnabledByDefault":true},{"Name":"golint"},{"Name":"staticcheck","EnabledByDefault":true},{"Name":"unused","Enabled":true,"EnabledByDefault":true},{"Name":"gosimple","EnabledByDefault":true},{"Name":"stylecheck"},{"Name":"gosec"},{"Name":"structcheck","EnabledByDefault":true},{"Name":"varcheck","EnabledByDefault":true},{"Name":"interfacer"},{"Name":"unconvert"},{"Name":"ineffassign","EnabledByDefault":true},{"Name":"dupl"},{"Name":"goconst"},{"Name":"deadcode","EnabledByDefault":true},{"Name":"gocyclo"},{"Name":"gocognit"},{"Name":"typecheck","EnabledByDefault":true},{"Name":"gofmt"},{"Name":"goimports"},{"Name":"maligned"},{"Name":"depguard"},{"Name":"misspell"},{"Name":"lll"},{"Name":"unparam"},{"Name":"dogsled"},{"Name":"nakedret"},{"Name":"prealloc"},{"Name":"scopelint"},{"Name":"gocritic"},{"Name":"gochecknoinits"},{"Name":"gochecknoglobals"},{"Name":"godox"},{"Name":"funlen"},{"Name":"whitespace"},{"Name":"wsl"}]}}`
	_, err := decodeJSON(strings.NewReader(json))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
}
