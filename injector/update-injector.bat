@echo off 

:: Build for linux
echo "building injector for linux"
set GOARCH=amd64
set GOOS=linux
go tool dist install -v pkg/runtime
go install -v -a std
go build
echo Build successfully

:: ask to continue
setlocal
:PROMPT
SET /P TOCONTINUE=Continue to update lab? (Y/[N])? 
IF /I "%TOCONTINUE%" NEQ "Y" GOTO END

:: Update lab
echo Updating sender
scp injector idanmadmon@192.168.108.128:~/bin/injector

echo Updating receiver
scp injector idanmadmon@192.168.108.129:~/bin/injector

echo Lab updated successfully

:END
endlocal