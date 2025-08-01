@echo off
setlocal ENABLEEXTENSIONS

REM Load .env file manually
for /f "usebackq tokens=1,* delims==" %%i in (".env") do (
    set %%i=%%j
)

REM Menu
if "%1"=="" goto help
if "%1"=="start" goto start
if "%1"=="lint" goto lint
if "%1"=="tests" goto tests
if "%1"=="testsum" goto testsum
if "%1"=="swagger" goto swagger
if "%1"=="migrate-up" goto migrateup
if "%1"=="migrate-down" goto migratedown
if "%1"=="docker" goto docker
if "%1"=="docker-test" goto dockertest
if "%1"=="docker-down" goto dockerdown
if "%1"=="docker-cache" goto dockercache

goto end

:start
echo Starting app...
go run src\main.go
goto end

:lint
echo Running linter...
golangci-lint run
goto end

:tests
echo Running tests...
go test -v ./test/...
goto end

:testsum
echo Running tests with gotestsum...
cd test
gotestsum --format testname
cd ..
goto end

:swagger
echo Generating Swagger docs...
cd src
swag init
cd ..
goto end

:migrateup
echo Running migrations up...
migrate -database "postgres://%DB_USER%:%DB_PASSWORD%@%DB_HOST%:%DB_PORT%/%DB_NAME%?sslmode=disable" -path src/database/migrations up
goto end

:migratedown
echo Running migrations down...
migrate -database "postgres://%DB_USER%:%DB_PASSWORD%@%DB_HOST%:%DB_PORT%/%DB_NAME%?sslmode=disable" -path src/database/migrations down
goto end

:docker
echo Starting Docker containers...
docker-compose up --build
goto end

:dockertest
echo Starting Docker and running tests...
docker-compose up -d
call :tests
goto end

:dockerdown
echo Stopping Docker containers...
docker-compose down --rmi all --volumes --remove-orphans
goto end

:dockercache
echo Cleaning Docker builder cache...
docker builder prune -f
goto end

:help
echo Available commands:
echo    dev.bat start
echo    dev.bat lint
echo    dev.bat tests
echo    dev.bat testsum
echo    dev.bat swagger
echo    dev.bat migrate-up
echo    dev.bat migrate-down
echo    dev.bat docker
echo    dev.bat docker-test
echo    dev.bat docker-down
echo    dev.bat docker-cache
goto end

:end
endlocal
