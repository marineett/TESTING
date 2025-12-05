SHELL := /bin/bash

.PHONY: help allure allure-open clean

help:
	@echo "Available targets:"
	@echo "  make allure       - Generate Allure report"
	@echo "  make allure-open  - Open report in browser"
	@echo "  make clean        - Remove generated files"

allure:
	@rm -rf allure-results-full
	@mkdir -p allure-results-full
	@cd backend && go test -v ./service_logic -count=1 -json 2>&1 | go-junit-report -package-name ./service_logic > ../allure-results-full/junit-service_logic.xml || true
	@cd backend && go test -v ./data_base -count=1 -json 2>&1 | go-junit-report -package-name ./data_base > ../allure-results-full/junit-data_base.xml || true
	@find backend/allure -path "*/integration_tests/*.json" -type f -exec cp {} allure-results-full/ \; 2>/dev/null || true
	@find backend/allure -path "*/allure-results-e2e*/*.json" -type f -exec cp {} allure-results-full/ \; 2>/dev/null || true
	@echo '{' > allure-results-full/executor.json
	@echo '  "name": "Test Suite",' >> allure-results-full/executor.json
	@echo '  "type": "local"' >> allure-results-full/executor.json
	@echo '}' >> allure-results-full/executor.json
	@echo "UNIT.TESTS=./service_logic (223), ./data_base (151)" > allure-results-full/environment.properties
	@echo "INTEGRATION.TESTS=DepartmentSuite (16), ClientSuite (14)" >> allure-results-full/environment.properties
	@echo "E2E.TESTS=ChatSuite (13), AuthSuite (10), LessonSuite (6), PersonalDataSuite (4), APISuite (1)" >> allure-results-full/environment.properties
	@echo "Total.Tests=438" >> allure-results-full/environment.properties
	@echo "Date=$(shell date '+%Y-%m-%d %H:%M:%S')" >> allure-results-full/environment.properties
	@if [ -d "allure-report-full/history" ]; then \
		mkdir -p allure-results-full/history; \
		cp -R allure-report-full/history/* allure-results-full/history/; \
	fi
	@allure generate ./allure-results-full --clean -o ./allure-report-full >/dev/null 2>&1
	@echo "Report generated successfully"

allure-open:
	@allure open ./allure-report-full 2>/dev/null || echo "Run 'make allure' first"

clean:
	@rm -rf allure-results-full allure-report-full backend/allure-results backend/allure-report backend/junit*.xml
	@echo "Cleaned"