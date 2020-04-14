ssh-web-1:
	ssh -i ~/Downloads/WebServer1.pem ec2-user@$$WEB_SERVER_1_IP
ssh-web-2:
	ssh -i ~/Downloads/WebServer1.pem ec2-user@$$WEB_SERVER_2_IP