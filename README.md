# awstools
A few helpful AWS tools.

## Tools
* `assume` - assumes a role on a subaccount behind a bastion.
* `ec2` - searches EC2 instances by pattern.
* `rotate-bastion-key` - rotates the bastion key

## Configuration
The default path to the configuration file is `$HOME/.config/awstools/config.toml`.

Here is an example of a `config.toml`:
```
defaultRegion = "eu-west-1"

[profiles]
bastion = "bastion_account"
bastionMfa = "bastion_mfa"

[accounts]
main = "000000000001"
dev = "000000000002"
prod = "000000000003"
```

* `profiles` section contains profile names that will be saved in `$HOME/.aws/credentials`.
* `accounts` section contains account ids and its names.