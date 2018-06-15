
REM - This is the main config script - which takes as parameters the username/password, vhost and endpoing url (without a slash)


curl -s -i -k -u %1:%2 -H "content-type:application/json" -XPUT %4/api/vhosts/%3
curl    -i -k -u %1:%2 -H "content-type:application/json" -d "{""configure"":"".*"",""write"":"".*"",""read"":"".*""}" -XPUT %4/api/permissions/%3/%1


curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""durable"":true}" -XPUT %4/api/queues/%3/programmer.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""durable"":true}" -XPUT %4/api/queues/%3/world.all


curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":false}" -XPUT %4/api/exchanges/%3/programmermic
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":true}" -XPUT %4/api/exchanges/%3/programmer.general
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":true}" -XPUT %4/api/exchanges/%3/programmer.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":false}" -XPUT %4/api/exchanges/%3/worldmic
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":true}" -XPUT %4/api/exchanges/%3/world.general
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""type"":""topic"",""durable"":true,""internal"":true}" -XPUT %4/api/exchanges/%3/world.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""world.general.programmer.#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/programmermic/e/world.general
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/world.general/e/world.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/programmer.all/q/programmer.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""programmer.general.world.#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/worldmic/e/programmer.general
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/programmer.general/e/programmer.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/world.all/q/world.all
curl -i -k -u %1:%2 -H "content-type:application/json" -d"{""routing_key"":""#"",""arguments"":{}}" -XPOST %4/api/bindings/%3/e/programmer.all/e/world.all

                    