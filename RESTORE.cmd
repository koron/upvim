@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0
SET TARGET_DIR=%BASEDIR%
GOTO :MAIN

:DELETE
IF EXIST "%1" DEL /F /Q "%1"
EXIT /B 0

:MAIN

CALL :DELETE %TARGET_DIR%_upvim\var\vim74-x86-anchor.txt
CALL :DELETE %TARGET_DIR%_upvim\var\vim74-x86-recipe.txt
CALL :DELETE %TARGET_DIR%_upvim\var\vim74-x64-anchor.txt
CALL :DELETE %TARGET_DIR%_upvim\var\vim74-x64-recipe.txt

upvim.exe %TARGET_DIR%
GOTO :END

:END
ECHO ========================================================================
ECHO This window will be closed automatically in 10 seconds.
PING localhost -n 10 > nul

ECHO PAUSE
