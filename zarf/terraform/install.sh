#!/bin/bash
yum update -y
amazon-linux-extras install docker
yum install docker
service docker start
usermod -a -G docker ec2-user
echo "${ghcr_token}" | docker login https://ghcr.io -u taraktikos --password-stdin
mkdir -p /home/ec2-user/.docker/
mv /root/.docker/config.json /home/ec2-user/.docker/config.json
chown ec2-user:ec2-user /home/ec2-user/.docker/config.json
