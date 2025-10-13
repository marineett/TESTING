SHELL := /bin/bash

PKGS := ./service_logic ./data_base

.PHONY: help test svc db london classic shuffle offline \
        ensure-junit junit junit-strict \
        prepare-allure allure allure-open allure-html \
        p1 p8 par1 par8 cpu1 cpu124

help:
	@echo "Targets:"
	@echo "  test         - все пакеты (-race)"
	@echo "  svc          - только сервисная логика"
	@echo "  db           - только репозитории/БД"
	@echo "  london       - London-тесты (in-memory)"
	@echo "  classic      - Classic-тесты (DuckDB :memory:)"
	@echo "  shuffle      - запуск в перемешку"
	@echo "  offline      - запуск без сети"
	@echo "  allure       - junit (мягко) -> отчет (с трендом)"
	@echo "  allure-open  - открыть отчет через http-сервер Allure"
	@echo "  allure-html  - открыть отчет как html (без сервера Allure)"
	@echo "  p1/p8        - параллелизм по пакетам (-p)"
	@echo "  par1/par8    - параллелизм внутри пакета (-parallel)"
	@echo "  cpu1/cpu124  - прогон с GOMAXPROCS/(-cpu)"


test:        
	@go test ./... -count=1 -race

svc:          
	@go test ./service_logic -v -count=1

db:           
	@go test ./data_base -v -count=1

london:      
	@go test ./service_logic -run London -v

classic:      
	@go test ./service_logic -run Classic -v

shuffle:      
	@go test ./... -shuffle=on -count=1

offline:      
	@GOPROXY=off go test ./... -count=1


ensure-junit:
	@command -v go-junit-report >/dev/null || { \
		echo "Installing go-junit-report..."; \
		go install github.com/jstemmer/go-junit-report/v2@latest; }

# для локальной генерации отчета
junit: ensure-junit
	@echo "Running tests -> junit.xml (non-strict)"
	@go test -v ./... -count=1 -json | go-junit-report > junit.xml || true

junit-strict: ensure-junit
	@echo "Running tests -> junit.xml (STRICT)"
	@go test -v ./... -count=1 -json | go-junit-report -set-exit-code > junit.xml

prepare-allure:
	@rm -rf allure-results
	@mkdir -p allure-results
	@cp junit.xml allure-results/
	@[ -d allure-report/history ] && mkdir -p allure-results/history && cp -R allure-report/history/* allure-results/history/ || true

allure: junit prepare-allure
	@allure generate ./allure-results --clean -o ./allure-report

allure-open:
	@allure open ./allure-report

# Открыть без блокировки терминала 
allure-html: allure
	@open ./allure-report/index.html
	@echo "Opened ./allure-report/index.html"


# Параллелизм по пакетам (каждый пакет — свой тестовый бинарь/процесс)
p1:  ; @go test ./service_logic -p 1 -run Test -v
p8:  ; @go test ./service_logic -p 8 -run Test -v

# Параллелизм внутри пакета (работает для тестов с t.Parallel())
par1: ; @go test ./service_logic -parallel 1 -run . -v
par8: ; @go test ./service_logic -parallel 8 -run . -v

# Системная параллельность (сколько OS-тредов выделено рантайму Go)
cpu1:   ; @GOMAXPROCS=1 go test ./service_logic -run . -count=1 -v
cpu124: ; @go test ./service_logic -run . -count=1 -cpu=1,2,4 -v
