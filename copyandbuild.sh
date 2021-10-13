#!/bin/bash

KEY=~/.ssh/obc_bastion_prod
BASTION=bastion.prod.inclusivecareco.org
SSH_USER=ec2-user
GO_PATH=/usr/local/go/bin

scp -i $KEY .env go.mod go.sum main.go $SSH_USER@$BASTION:~/migration/
ssh -i $KEY $SSH_USER@$BASTION "cd migration && $GO_PATH/go build && ./migration"
