CGO_ENABLED?=1
GO_BIN?=go
GO?=CGO_ENABLED=${CGO_ENABLED} ${GO_BIN}

.PHONY: clean
clean:
	rm -f bin/litefs

bin/litefs:
	# Force enable CGO for SQLite driver
	${GO} mod tidy
	${GO} build -trimpath \
		-o bin/litefs \
		./cmd/litefs

include Makefile.bindl

GIT_STAMP?=$(shell git describe --match="" --always --dirty)
TEST_RUN_NAME?=$(shell date +%s)-${GIT_STAMP}

.PHONY: test/prep
test/prep: DBMATE_PATH=tmp/tests/litefs.db
test/prep: DBMATE_SCHEMA_FILE=tmp/tests/schema.sql
test/prep: bin/dbmate
	@echo "=>> TEST DATABASE: INITIALIZING <<="
	rm -f ${DBMATE_PATH} ${DBMATE_SCHEMA_FILE}
	mkdir -p tmp/tests
	${MAKE} --no-print-directory db/migrate DBMATE_PATH=${DBMATE_PATH} DBMATE_SCHEMA_FILE=${DBMATE_SCHEMA_FILE}
	@echo "<<= TEST DATABASE: READY ==>"
	@echo 

.PHONY: test/unit
test/unit: test/prep
	${GO} test -v -race ./...

.PHONY: test/bench
test/bench: test/prep
	${GO} test -v -run='.*Benchmark.*' -benchmem -bench=.

.PHONY: test/bench/benchstat
test/bench/benchstat: test/prep
	@mkdir -p tmp/bench/benchstat
	${GO} test -v -run='.*Benchmark.*' -benchmem -bench=. > tmp/bench/benchstat/${TEST_RUN_NAME}.txt
	@echo
	@ls -d tmp/bench/benchstat/* | tail -n2 | xargs benchstat

.PHONY: test/bench/pprof
test/bench/pprof: test/prep
	@mkdir -p tmp/bench/pprof
	${GO} test -v -run='.*Benchmark.*' \
		-cpuprofile tmp/bench/pprof/${TEST_RUN_NAME}-cpuprofile-litefs.out \
		-memprofile tmp/bench/pprof/${TEST_RUN_NAME}-memprofile-litefs.out \
		-benchmem -bench=ProfLiteFS
	${GO} test -v -run='.*Benchmark.*' \
		-cpuprofile tmp/bench/pprof/${TEST_RUN_NAME}-cpuprofile-osfs.out \
		-memprofile tmp/bench/pprof/${TEST_RUN_NAME}-memprofile-osfs.out \
		-benchmem -bench=ProfOSFS

.PHONY: test/all
test/all: test/unit
test/all: test/bench

DBMATE_PATH?=tmp/litefs.db
DBMATE_MIGRATIONS_DIR?=sql/migrations
DBMATE_SCHEMA_FILE?=sql/schema/schema.sql
DBMATE_FLAGS?=--url $(addprefix sqlite:,${DBMATE_PATH}) --migrations-dir ${DBMATE_MIGRATIONS_DIR} --schema-file ${DBMATE_SCHEMA_FILE}

.PHONY: db/migrate
db/migrate: bin/dbmate
	bin/dbmate ${DBMATE_FLAGS} up

.PHONY: db/reset
db/reset: bin/dbmate
	rm -f ${DBMATE_PATH} ${DBMATE_SCHEMA_FILE}

.PHONY: db/generate-orm
db/generate-orm: bin/sqlc
	bin/sqlc generate
