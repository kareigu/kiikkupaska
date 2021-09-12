mkdir dist
del /S /Q .\dist

go build -ldflags "-s -w"

"C:\Program Files (x86)\Resource Hacker\ResourceHacker.exe" -open kiikkupaskaa.exe -save kiikkupaskaa.exe -action addskip -res assets/fav.ico -mask ICONGROUP,MAIN

copy .\kiikkupaskaa.exe .\dist\
Xcopy /E /I .\assets .\dist\assets