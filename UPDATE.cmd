@ECHO OFF
SETLOCAL
SET BASEDIR=%~dp0
SET TARGET_DIR=%BASEDIR%
GOTO :MAIN

:MAIN

upvim.exe %TARGET_DIR%
GOTO :END

:END
ECHO ========================================================================
ECHO This window will be closed automatically in 10 seconds.
PING localhost -n 10 > nul

ECHO PAUSE
