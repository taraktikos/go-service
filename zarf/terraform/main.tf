terraform {
  required_version = ">= 1.0.4"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 3.53"
    }
  }

  backend "s3" {
    bucket = "taras-aws-terraform-state"
    key    = "go-service.tfstate"
    region = "eu-central-1"
  }
}

provider "aws" {
  profile = "default"
  region  = "eu-central-1"
}

# 1. Create a VPC
resource "aws_vpc" "go_service_vpc" {
  cidr_block = "10.0.0.0/16"
  tags = {
    Name = "Go Service VPC"
  }
}

# 2. Create Internet Gateway
resource "aws_internet_gateway" "go_service_gateway" {
  vpc_id = aws_vpc.go_service_vpc.id
  tags = {
    Name = "Go Service Internet Gateway"
  }
}

# 3. Create Custom Route Table
resource "aws_route_table" "go_service_route_table" {
  vpc_id = aws_vpc.go_service_vpc.id
  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.go_service_gateway.id
  }
  route {
    ipv6_cidr_block = "::/0"
    gateway_id      = aws_internet_gateway.go_service_gateway.id
  }
  tags = {
    Name = "Go Service Route Table"
  }
}

# 4. Create a Subnet
resource "aws_subnet" "go_service_subnet-1" {
  vpc_id            = aws_vpc.go_service_vpc.id
  cidr_block        = "10.0.1.0/24"
  availability_zone = "eu-central-1a"
  tags = {
    Name = "Go Service Subnet 1a"
  }
}

# 5. Associate subnet with Route Table
resource "aws_route_table_association" "a" {
  subnet_id      = aws_subnet.go_service_subnet-1.id
  route_table_id = aws_route_table.go_service_route_table.id
}

# 6. Create Security Group to allow port 22, 80,443
resource "aws_security_group" "go_service_allow_web" {
  name        = "go_service_allow_web_traffic"
  description = "Allow web inbound traffic for go service"
  vpc_id      = aws_vpc.go_service_vpc.id

  ingress {
    description      = "HTTPS"
    from_port        = 443
    to_port          = 443
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    description      = "HTTP"
    from_port        = 80
    to_port          = 80
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  ingress {
    description      = "SSH"
    from_port        = 22
    to_port          = 22
    protocol         = "tcp"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  egress {
    from_port        = 0
    to_port          = 0
    protocol         = "-1"
    cidr_blocks      = ["0.0.0.0/0"]
    ipv6_cidr_blocks = ["::/0"]
  }

  tags = {
    Name = "Go Service Allow Web"
  }
}

# 7. Create a network interface with an ip in the subnet that was created in step 4
resource "aws_network_interface" "go_service_web_server_nic" {
  subnet_id       = aws_subnet.go_service_subnet-1.id
  private_ips     = ["10.0.1.50"]
  security_groups = [aws_security_group.go_service_allow_web.id]
}

# 8. Assign an Elastic IP to the network interface created in step 7
resource "aws_eip" "go_service_eip" {
  vpc                       = true
  network_interface         = aws_network_interface.go_service_web_server_nic.id
  associate_with_private_ip = "10.0.1.50"
  depends_on                = [aws_internet_gateway.go_service_gateway]
}

# 9. Create app server and install docker form install.sh
resource "aws_instance" "go_service_app_server" {
  ami               = "ami-0453cb7b5f2b7fca2"
  instance_type     = "t2.micro"
  availability_zone = "eu-central-1a"
  key_name          = "taras-aws-thinkpad-x220"
  network_interface {
    device_index         = 0
    network_interface_id = aws_network_interface.go_service_web_server_nic.id
  }
  user_data = templatefile("install.sh", {ghcr_token = var.ghcr_token})
  tags = {
    Name = "Go Service App Server"
  }
}

variable "ghcr_token" {
  description = "The token for the docker registry ghcr.io"
  type        = string
}

output "server_public_ip" {
  value = aws_eip.go_service_eip.public_ip
}
