.PHONY: benchtrace
benchtrace:
	go test -c & GODEBUG=allocfreetrace=1 ./benchmarks.test -test.run=none -test.bench=^BenchmarkOneLogTrace -test.benchtime=10ms 2>trace.log

.PHONY: benchonelog
benchonelog:
	go test -benchmem -run=^BenchmarkOnelog -bench=^BenchmarkOnelog -benchtime=30ms

.PHONY: benchonelograce
benchonelograce:
	go test -race -benchmem -run=^BenchmarkOnelog -bench=^BenchmarkOnelog -benchtime=30ms

.PHONY: benchonelogcpu
benchonelogcpu:
	go test -benchmem -run=^BenchmarkOnelog -bench=^BenchmarkOnelog -benchtime=30ms -cpuprofile cpu.out
