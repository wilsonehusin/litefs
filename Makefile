CGO_ENABLED?=0
GO?=go

.PHONY: clean
clean:
	rm -f bin/litefs

bin/litefs:
	# Force enable CGO for SQLite driver
	CGO_ENABLED=1 go mod tidy
	CGO_ENABLED=${CGO_ENABLED} ${GO} build -trimpath \
		-o bin/litefs \
		./cmd/litefs

include Makefile.bindl

.PHONY: test/prep
test/prep: DBMATE_PATH=tmp/tests/litefs.db
test/prep: DBMATE_SCHEMA_FILE=tmp/tests/schema.sql
test/prep: bin/dbmate
	@echo "=== TEST DATABASE ==="
	rm -f ${DBMATE_PATH} ${DBMATE_SCHEMA_FILE}
	mkdir -p tmp/tests
	${MAKE} --no-print-directory db/migrate DBMATE_PATH=${DBMATE_PATH} DBMATE_SCHEMA_FILE=${DBMATE_SCHEMA_FILE}
	@echo "--- TEST DATABASE ---"
	@echo 

.PHONY: test/unit
test/unit: test/prep
	${GO} test -v -race ./...

.PHONY: test/bench
test/bench: test/prep
	${GO} test -v -run='.*Benchmark.*' -benchmem -bench=.

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
