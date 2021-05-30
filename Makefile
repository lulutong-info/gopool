test:
	@go test -v -count=1 ./

test_10:
	@go test -v -count=10 ./

test_bench:
	@go test -v -run=none -bench=. -benchmem ./