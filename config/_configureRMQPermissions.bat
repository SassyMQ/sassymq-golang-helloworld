
ECHO ON                
IF "%4"==""  GOTO SYNTAX_ERROR


REM SETUP Virtual Machine
curl -i -k -u %1:%2 -H "content-type:application/json" -d "{""configure"":"".*"",""write"":"".*"",""read"":"".*""}" -XPUT %4/api/permissions/%3/smqPublic
curl -i -k -u %1:%2 -H "content-type:application/json" -d "{""configure"":"".*"",""write"":"".*"",""read"":"".*""}" -XPUT %4/api/permissions/%3/smqUser
curl -i -k -u %1:%2 -H "content-type:application/json" -d "{""configure"":"".*"",""write"":"".*"",""read"":"".*""}" -XPUT %4/api/permissions/%3/smqAdmin


:SYNTAX_ERROR
ECHO Syntax _ConfigureRMQPermissions {Username} {Password} {VirtualHostName} {broker}
ECHO i.e. _ConfigureRMQPermissions admin admin VHost1 NewUser1 Password1 http://localhost:15672
GOTO END


:END
ECHO Done.
                 