package main

import "fmt"

func displayHelpText() {
	const helpText = `Syntax: awsmfa [commands] [mfa-token]

Commands:

-h, --help          Show this help text
-r, --restore       Restore original credentials back to AWS credentials file

'mfa-token' must be the currently displayed numeric MFA token from the device you've configured as a virtual MFA device associated with your IAM user. In addition, active IAM access credentials must already have been stored in your local 'credentials' file or in the AWS-specific environment variables. For help with enabling a virtual MFA device, see https://docs.aws.amazon.com/IAM/latest/UserGuide/id_credentials_mfa_enable_virtual.html.

Examples:

To obtain temporary session credentials from AWS and save to credentials file:

$ awsmfa 123456
[You can now perform AWS actions that require MFA...]

To switch back to using permanent access credentials:

$ awsmfa --restore
[You can no longer perform AWS actions that require MFA...]

For more information: https://github.com/luhring/awsmfa

`
	fmt.Print(helpText)
}
