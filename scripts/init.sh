#!/bin/bash

cat >>/etc/my_init.d/services <<-EOGO
#!/bin/bash
dpkg-reconfigure openssh-server

/etc/init.d/mysql start
/etc/init.d/redis-server start
/etc/init.d/rabbitmq-server start
/etc/init.d/ssh start
EOGO
chmod a+x /etc/my_init.d/services

sed -i "s/#UsePAM/UsePAM/" /etc/ssh/sshd_config

/sbin/my_init