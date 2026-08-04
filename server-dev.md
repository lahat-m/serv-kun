Create ssh key with ssh-keygen -t ed25519 -b 4096 -C “lahat”
Copies the key to the server
Login into the server
Create new non root user
Modify privileges of the user: usermod -aG sudo lahat
Configure lahat to login into remote server/machine via ssh command and ssh key file
Disable root login
Setup ssh alias
Install java openjdk turmeric via sdkman
Generate and upload jar executable to remote server /opt/ renamed to /opt/kuniel.jar
Setup environment file
Setup systemd service for auto-relaunch
Setup ssl certificate with nginx
Setup github actions
Configure firewall
