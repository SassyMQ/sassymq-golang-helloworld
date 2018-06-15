
ECHO ON                
IF "%4"==""  GOTO SYNTAX_ERROR

:SYNTAX_ERROR
ECHO Syntax _ConfigureRMQPermissions {Username} {Password} {VirtualHostName} {broker}
ECHO i.e. _ConfigureRMQPermissions admin admin VHost1 NewUser1 Password1 http://localhost:15672
GOTO END


:END
ECHO Done.
                 