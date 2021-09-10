mkdir dist
del /S /Q .\dist

go build -ldflags "-s -w"

copy .\kiikkupaskaa.exe .\dist\
Xcopy /E /I .\assets .\dist\assets