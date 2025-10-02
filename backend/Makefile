
SHELL := /bin/bash
PKGS := ./service_logic ./data_base


.PHONY: help test svc db london classic shuffle offline \
        ensure-junit junit junit-strict \
        prepare-allure allure allure-open allure-html clean \
        p1 p8 par1 par8 cpu1 cpu124 \
		build run build-tests run-tests

build:
	@mkdir -p bin
	@go build -o bin/app ./...
	@echo "Built prod binary: bin/app"

run: build
	@./bin/app

build-tests:
	@mkdir -p bin
	@go test -c ./service_logic -o bin/service_logic.test
	@go test -c ./data_base      -o bin/data_base.test
	@echo "Built test binaries: bin/service_logic.test, bin/data_base.test"

run-tests: build-tests
	@./bin/service_logic.test -test.v
	@./bin/data_base.test -test.v
help:
	@echo "Targets:"
	@echo "  test          - все пакеты (-race)"
	@echo "  svc           - только сервисная логика"
	@echo "  db            - только репозитории/БД"
	@echo "  london        - London-тесты (in-memory фейки)"
	@echo "  classic       - Classic-тесты (DuckDB :memory:)"
	@echo "  shuffle       - запуск в перемешку"
	@echo "  offline       - запуск без сети"
	@echo "  allure        - собрать Allure-отчёт (по пакетам, с трендом)"
	@echo "  allure-open   - открыть отчёт (Allure server)"
	@echo "  allure-html   - открыть локально index.html"
	@echo "  clean         - почистить allure-results / allure-report / junit-*.xml"
	@echo "  p1/p8         - параллелизм по пакетам (-p) для ./..."
	@echo "  par1/par8     - параллелизм внутри пакета (-parallel) для ./service_logic"
	@echo "  cpu1/cpu124   - GOMAXPROCS / -cpu для ./service_logic"

test:
	@go test ./... -count=1 -race

svc:
	@go test ./service_logic -v -count=1

db:
	@go test ./data_base -v -count=1

london:
	@go test ./service_logic -run London -v -count=1

classic:
	@go test ./service_logic -run Classic -v -count=1

shuffle:
	@go test ./... -shuffle=on -count=1

offline:
	@GOPROXY=off go test ./... -count=1

ensure-junit:
	@command -v go-junit-report >/dev/null || { \
		echo "Installing go-junit-report..."; \
		go install github.com/jstemmer/go-junit-report/v2@latest; }

junit: ensure-junit
	@rm -f junit-*.xml
	@set -e; \
	for pkg in $(PKGS); do \
		name=$$(basename $$pkg); \
		echo ">> $$pkg -> junit-$$name.xml"; \
		go test -v $$pkg -count=1 -json \
			| go-junit-report -package-name $$pkg > junit-$$name.xml || true; \
	done
	@echo "OK: junit-*.xml готово"

junit-strict: ensure-junit
	@rm -f junit-*.xml
	@set -e; \
	for pkg in $(PKGS); do \
		name=$$(basename $$pkg); \
		echo ">> STRICT $$pkg -> junit-$$name.xml"; \
		go test -v $$pkg -count=1 -json \
			| go-junit-report -package-name $$pkg -set-exit-code > junit-$$name.xml; \
	done
	@echo "OK: junit-*.xml (STRICT) готово"

prepare-allure:
	@rm -rf allure-results
	@mkdir -p allure-results
	@cp junit-*.xml allure-results/
	@[ -d allure-report/history ] && mkdir -p allure-results/history && cp -R allure-report/history/* allure-results/history/ || true

allure: junit prepare-allure
	@allure generate ./allure-results --clean -o ./allure-report
	@echo "Allure report -> ./allure-report"

allure-open:
	@allure open ./allure-report

allure-html: allure
	@open ./allure-report/index.html 2>/dev/null || xdg-open ./allure-report/index.html || true
	@echo "Opened ./allure-report/index.html"

clean:
	@rm -rf allure-results allure-report junit.xml junit-*.xml

p1:
	@go test ./... -p 1 -run Test -v

p8:
	@go test ./... -p 8 -run Test -v

par1:
	@go test ./service_logic -parallel 1 -run . -v

par8:
	@go test ./service_logic -parallel 8 -run . -v

cpu1:
	@GOMAXPROCS=1 go test ./service_logic -run . -count=1 -v

cpu124:
	@go test ./service_logic -run . -count=1 -cpu=1,2,4 -v
