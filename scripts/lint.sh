dirs=$(go list -f '{{.Dir}}' ./... | grep -v /vendor/)
for d in $dirs; do goimports -d $d/*.go; done
