# awsmfa

A tool for enabling AWS CLI operations that require MFA

[![CircleCI](https://circleci.com/gh/luhring/awsmfa.svg?style=svg)](https://circleci.com/gh/luhring/awsmfa)
[![Go Report Card](https://goreportcard.com/badge/github.com/luhring/awsmfa)](https://goreportcard.com/report/github.com/luhring/awsmfa)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/luhring/awsmfa/blob/master/LICENSE)

## Background

MFA has become an extremely popular and successful security mechanism to defend against situations where passwords or secret keys are unexpectedly exposed to an attacker.

AWS allows IAM policies to specify that the listed permissions are available to a user (or the user's group/role) only when the user has first authenticated with an MFA device. AWS provides more information on this setup in ["Configuring MFA-Protected API Access"](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_configure-api-require.html).

To access the AWS API using local command line tools, while needing to perform actions that require that you authenticate with MFA, you must first obtain temporary credentials via the [GetSessionToken](https://docs.aws.amazon.com/STS/latest/APIReference/API_GetSessionToken.html) or [AssumeRole](https://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html) API methods.

awsmfa makes this process easier for users by providing a simple syntax for providing your MFA device's token code and saving your temporary credentials to disk for use in future commands. awsmfa also makes it easy to discard the temporary credentials and restore your original credentials back to their original location.

## Usage

Prerequisites for using this tool:

1. You must have already obtained access credentials (access key ID and secret access key) for an IAM user in an AWS account.
1. These credentials should (ideally) be saved in your local `credentials` file —- see ["Configuration and Credential Files"](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html) for help setting this up —- but can alternatively be stored in [AWS-specific environment variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html). Currently, if your credentials are saved to a `credentials` file, they must be stored in the "default" profile within the file.
1. You must have associated a virtual MFA device with your IAM user. If you need help doing this, check out ["Enabling a Virtual Multi-factor Authentication (MFA) Device (Console)"](https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html).

**Note:** your experience will be smoother if you store your credentials in a `credentials` file and you _don't_ have AWS-specific environment variables set.

### Syntax

`awsmfa [commands] [mfa-token]`

(`mfa-token` must be the currently displayed numeric MFA token from the device you've configured as a virtual MFA device associated with your IAM user.)

### Commands

`-h`, `--help`: Show this help text. _(Don't specify an `mfa-token` with this command.)_

`-r`, `--restore`: Restore original credentials back to AWS credentials file. _(Don't specify an `mfa-token` with this command.)_

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
