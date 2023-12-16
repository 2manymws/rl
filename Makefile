export GO111MODULE=on

default: test

ci: depsdev test

test:
	cp go.mod testdata/go_test.mod
	go test ./... -coverprofile=coverage.out -covermode=count

benchmark: depsdev
	go test -bench . -benchmem -benchtime 10000x -run Benchmark | octocov-go-test-bench --tee > custom_metrics_benchmark.json

cachegrind: depsdev
	cd testdata/testbin && go build -o testbin
	setarch `uname -m` -R valgrind --tool=cachegrind --cachegrind-out-file=cachegrind.out --I1=32768,8,64 --D1=32768,8,64 --LL=8388608,16,64 ./testdata/testbin/testbin 10000
	cat cachegrind.out | octocov-cachegrind --tee > custom_metrics_cachegrind.json

lint:
	golangci-lint run ./...
	go vet -vettool=`which gostyle` -gostyle.config=$(PWD)/.gostyle.yml ./...

depsdev:
	go install github.com/Songmu/ghch/cmd/ghch@latest
	go install github.com/Songmu/gocredits/cmd/gocredits@latest
	go install github.com/k1LoW/octocov-go-test-bench/cmd/octocov-go-test-bench@latest
	go install github.com/k1LoW/octocov-cachegrind/cmd/octocov-cachegrind@latest
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
	go mod download
	gocredits -w .
	cat _EXTRA_CREDITS >> CREDITS
	git add CHANGELOG.md CREDITS go.mod go.sum

release:
	git push origin main --tag
