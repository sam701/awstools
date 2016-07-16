# awstools
[![Build Status](https://travis-ci.org/sam701/awstools.svg?branch=master)](https://travis-ci.org/sam701/awstools) A few helpful AWS tools.

## Tools
* `assume` - assumes a role on a subaccount behind the main account where you have an access key
* `ec2` - searches EC2 instances by pattern
* `rotate-main-account-key` - rotates the main account key
* `cloudformation` - prints stacks parameters and outputs, deletes stacks
* `kms` - decrypts base64 encoded text or encrypts and encodes with base64
* `kinesis` - lists streams, grabs kinesis stream for patterns

## Install

[We provide binaries for all releases through GitHub](https://github.com/sam701/awstools/releases). The latest release is [0.8.1](https://github.com/sam701/awstools/releases/latest).

To install `awstools` choose the binary for your architecture (either OSX or Linux), run a download and use `chmod` to make it executable.

### OSX

```sh
$ curl -o awstools -SsL https://github.com/sam701/awstools/releases/download/0.8.1/awstools_darwin_amd64
$ chmod +x awstools
```

### Linux

```sh
$ curl -o awstools -SsL https://github.com/sam701/awstools/releases/download/0.8.1/awstools_linux_amd64
$ chmod +x awstools
```

## Build
Export reqired environment variables:
```sh
export GOPATH=$HOME/goprojects
export PATH=$PATH:$GOPATH/bin
```

Install `awstools`:
```sh
go get -u github.com/sam701/awstools
```

Add to your .bash_profile
```sh
aws_assume(){
	tmpFile=/tmp/assume.tmp
	awstools assume --export $tmpFile $@ && source $tmpFile
	rm $tmpFile
}
```
Now in order to assume a role on a subaccount, you can run something like this
```sh
aws_assume AccountName MyRoleOnSubAccount
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

# License

This project is licensed under the MIT license. You can find a copy of the license at the top level of the repository.
