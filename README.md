# awstools
[![Build Status](https://travis-ci.org/sam701/awstools.svg?branch=master)](https://travis-ci.org/sam701/awstools) A few helpful AWS tools.

```
NAME:
   awstools - AWS tools

USAGE:
   awstools [global options] command [command options] [arguments...]

VERSION:
   0.11.3

COMMANDS:
     assume                      assume role on a specified account
     accounts                    print known accounts
     ec2                         print EC2 instances and ELBs
     cloudformation, cf          print CloudFormation stacks information
     rotate-main-account-key, r  create a new access key for main account and delete the current one
     dynamodb, ddb               dynamodb commands
     kms                         encrypt/decrypt text
     kinesis                     print records from kinesis streams
     cloudwatch, cw              search in cloudwatch logs
     help, h                     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config value, -c value  path to config.toml file (default: ~/.config/awstools/config.toml)
   --no-color                turn off color output
   --help, -h                show help
   --version, -v             print the version
```

## Install

[We provide binaries for all releases through GitHub](https://github.com/sam701/awstools/releases). The latest release is [0.11.3](https://github.com/sam701/awstools/releases/latest).

To install `awstools` choose the binary for your architecture (either OSX or Linux), run a download and use `chmod` to make it executable.

### OSX

On Mac you can use Homebrew to install the binary:

```sh
$ brew tap sam701/awstools
$ brew install awstools
```

### Linux

```sh
$ curl -o awstools -SsL https://github.com/sam701/awstools/releases/download/0.11.3/awstools_linux_amd64
$ chmod +x awstools
```

## Build
Export reqired environment variables:
```sh
export GOPATH=$HOME/goprojects
export PATH=$PATH:$GOPATH/bin
```

Install [glide](https://glide.sh).

Install `awstools`:
```sh
go get -d -u github.com/sam701/awstools
cd $GOPATH/src/github.com/sam701/awstools
glide install
go install
```

## Configuration
The default path to the configuration file is `$HOME/.config/awstools/config.toml`.

Here is an example of a `config.toml`:
```toml
defaultRegion = "eu-west-1"
defaultKmsKey = "arn:aws:kms:eu-west-1:000000000001:key/00000000-1111-1111-2222-333333333333"

# Rotate the main account access key every week
keyRotationIntervalMinutes = 10080

[profiles]
mainAccount = "main_account"
mainAccountMfaSession = "main_account_mfa_session"

[accounts]
main = "000000000001"
dev = "000000000002"
prod = "000000000003"
```

* `profiles` section contains profile names that will be saved in `$HOME/.aws/credentials`.
* `accounts` section contains account ids and its names.

Add to your `.bash_profile`
```sh
aws_assume(){
	tmpFile=/tmp/assume.tmp
	awstools assume --export $tmpFile $@ && source $tmpFile
	rm $tmpFile
}
```
or to your `~/.config/fish/config.fish`
```sh
function aws_assume
	set tmp /tmp/aws_assume.tmp
	awstools assume --export $tmp $argv; and source $tmp
	rm $tmp
end
```
Now in order to assume a role on a subaccount, you can run something like this
```sh
aws_assume AccountName MyRoleOnSubAccount
```

### Required IAM permissions

#### AssumeRole

For assuming a role in another account `awstools` needs the following permissions:

- `iam:GetUser`
- `iam:ListAccessKeys`

*Note: `awstools` is using the MFA authenticated sessions for operations on your AWS access key.*

#### Access Key Rotation

For rotating access keys on the relevant account `awstools` needs the following permissions:

- `iam:GetUser`
- `iam:CreateAccessKey`
- `iam:DeleteAccessKey`
- `iam:ListAccessKeys`
- `iam:UpdateAccessKey`

*Note: `awstools` is using the MFA authenticated sessions for operations on your AWS access key.*

# License

This project is licensed under the MIT license. You can find a copy of the license at the top level of the repository.
