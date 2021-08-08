# go-service

export TF_VAR_ghcr_token=$GHCR_TOKEN

terraform apply

chmod 400 .ssh/key.pem 

ssh -i .ssh/key.pem ec2-user@IP
