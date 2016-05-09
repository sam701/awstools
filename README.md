# awstools
A few helpful AWS tools.

## Tools
* `assume` - assumes a role on a subaccount behind the main account where you have an access key
* `ec2` - searches EC2 instances by pattern
* `rotate-main-account-key` - rotates the main account key
* `cloudformation` - prints stacks parameters and outputs
* `kms` - decrypts base64 encoded text or encrypts and encodes with base64

## Setup
Export reqired environment variables:
```
export GOPATH=$HOME/goprojects
export PATH=$PATH:$GOPATH/bin
```

Install `awstools`:
```
go get -u github.com/sam701/awstools
```

Add to your .bash_profile
```
aws_assume(){
	tmpFile=/tmp/assume.tmp
	awstools assume --export $tmpFile $@ && source $tmpFile
	rm $tmpFile
}
```
Now in order to assume a role on a subaccount, you can run something like this
```
aws_assume AccountName MyRoleOnSubAccount
```


## Configuration
The default path to the configuration file is `$HOME/.config/awstools/config.toml`.

Here is an example of a `config.toml`:
```
defaultRegion = "eu-west-1"
defaultKmsKey = "arn:aws:kms:eu-west-1:000000000001:key/00000000-1111-1111-2222-333333333333"

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
