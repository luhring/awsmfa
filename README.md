# awsmfa

A tool for enabling AWS CLI operations that require MFA

[![CircleCI](https://circleci.com/gh/luhring/awsmfa.svg?style=svg)](https://circleci.com/gh/luhring/awsmfa)
[![Go Report Card](https://goreportcard.com/badge/github.com/luhring/awsmfa)](https://goreportcard.com/report/github.com/luhring/awsmfa)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/luhring/awsmfa/blob/master/LICENSE)

## Usage

Prerequisites for using this tool:

1. You must have already obtained access credentials (access key ID and secret access key) for an IAM user in an AWS account.
1. These credentials should (ideally) be saved in your local `credentials` file —- see ["Configuration and Credential Files"](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) for help setting this up —- but can alternatively be stored in [AWS-specific environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html). Currently, if your credentials are saved to a `credentials` file, they must be stored in the "default" profile within the file.
1. You must have associated a virtual MFA device with your IAM user. If you need help doing this, check out ["Enabling a Virtual Multi-factor Authentication (MFA) Device (Console)"](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html).

**Note:** your experience will be smoother if you store your credentials in a `credentials` file and you _don't_ have AWS-specific environment variables set.

### Syntax

`awsmfa [commands] [mfa-token]`

### Commands

`-h`, `--help`: Show this help text

`-r`, `--restore`: Restore original credentials back to AWS credentials file

`mfa-token` must be the currently displayed numeric MFA token from the device you've configured as a virtual MFA device associated with your IAM user.

### Examples

To obtain temporary session credentials from AWS and save to credentials file:

```bash
$ awsmfa 123456
Created backup of original credentials at /Users/dan/.aws/credentials_backup_by_awsmfa.

Authentication successful!

Saved new session credentials to /Users/dan/.aws/credentials.
You can now perform actions that require MFA.
```

To switch back to using permanent access credentials:

```bash
$ awsmfa --restore
Restored original credentials from backup.
You can no longer perform actions that require MFA.
```

## Road map

- ~~Ability to get a session token via default profile~~
- Ability to specify custom session duration
- Ability to use non-default profiles
- Ability to assume a role
