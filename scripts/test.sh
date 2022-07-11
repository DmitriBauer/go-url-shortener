go test -count=1 $(go list ./... | grep -v /vendor/ | grep -v /cmd/)
