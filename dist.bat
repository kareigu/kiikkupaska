mkdir dist
del /S /Q .\dist

go build -ldflags "-s -w -H=windowsgui"


copy .\kiikkupaskaa.exe .\dist\
Xcopy /E /I .\assets .\dist\assets