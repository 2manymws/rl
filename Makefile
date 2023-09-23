export GO111MODULE=on

default: test

ci: depsdev test

test:
	cp go.mod testdata/go_test.mod
	go mod tidy -modfile=testdata/go_test.mod
	go test ./... -modfile=testdata/go_test.mod -coverprofile=coverage.out -covermode=count

benchmark: depsdev
	go mod tidy -modfile=testdata/go_test.mod
	go test -modfile=testdata/go_test.mod -bench . -benchmem -benchtime 10000x -run Benchmark | octocov-go-test-bench --tee > custom_metrics_benchmark.json

cachegrind: depsdev
	go mod tidy -modfile=testdata/go_test.mod
	valgrind --tool=cachegrind go test -modfile=testdata/go_test.mod -bench . -benchmem -benchtime 10000x -run Benchmark

lint:
	go mod tidy
	golangci-lint run ./...
	go vet -vettool=`which gostyle` -gostyle.config=$(PWD)/.gostyle.yml ./...
	git restore go.*

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/k1LoW/octocov-go-test-bench/cmd/octocov-go-test-bench@latest
	go install github.com/k1LoW/gostyle@latest

prerelease:
	git pull origin main --tag
	go mod download
	ghch -w -N ${VER}
	gocredits -w .
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum
	git commit -m'Bump up version number'
	git tag ${VER}

prerelease_for_tagpr: depsdev
	gocredits -w .
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag
