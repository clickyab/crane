mysql-setup: needroot
	echo 'UPDATE user SET plugin="";' | mysql mysql | true
	echo 'UPDATE user SET password=PASSWORD("$(DBPASS)") WHERE user="$(DB_USER)";' | mysql mysql | true
	echo 'FLUSH PRIVILEGES;' | mysql mysql | true
	echo 'DROP DATABASE IF EXISTS $(DB_NAME); CREATE DATABASE $(DB_NAME);' | mysql -u $(DB_USER) -p$(DBPASS)

rabbitmq-setup: needroot
	[ "1" -eq "$(shell rabbitmq-plugins enable rabbitmq_management | grep 'Plugin configuration unchanged' | wc -l)" ] || (rabbitmqctl stop_app && rabbitmqctl start_app)
	rabbitmqctl add_user $(RUSER) $(RPASS) || rabbitmqctl change_password $(RUSER) $(RPASS)
	rabbitmqctl set_user_tags $(RUSER) administrator
	rabbitmqctl set_permissions -p / $(RUSER) ".*" ".*" ".*"
	wget -O /usr/bin/rabbitmqadmin http://127.0.0.1:15672/cli/rabbitmqadmin
	chmod a+x /usr/bin/rabbitmqadmin
	rabbitmqadmin declare queue name=dlx-queue
	rabbitmqadmin declare exchange name=dlx-exchange type=topic
	rabbitmqctl set_policy DLX ".*" '{"dead-letter-exchange":"dlx-exchange"}' --apply-to queues
	rabbitmqadmin declare binding source="dlx-exchange" destination_type="queue" destination="dlx-queue" routing_key="#"
