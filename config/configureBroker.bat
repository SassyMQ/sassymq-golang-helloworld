
if "%2"=="" GOTO SYNTAXERROR
CALL _configureRMQVirtualHost.bat %rmq_admin_username% %rmq_admin_password% %1 %2
CALL _configureRMQPermissions.bat %rmq_admin_username% %rmq_admin_password% %1 %2
GOTO END


:SYNTAXERROR
ECHO Syntax:  > ConfigureBroker {VHost} {Broker}
ECHO Example: > ConfigureBroker SAMPLE_VHOST http://localhost:15672
GOTO END

:END
                    