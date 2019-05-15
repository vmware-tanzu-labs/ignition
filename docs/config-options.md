# Available Options for the `ignition-config` User Provided Service

| Name | Description | Required? | Default Value | Example Value |
| --- | --- | --- | --- | --- |
| `session_secret` | Used to secure the user's session cookie. You should randomly generate the contents of this value and limit access to it. | YES | `null` | `somegeneratedvalue` |
| `system_domain` | The PAS installation's system domain | YES | `null` | `sys.example.net` |
| `uaa_origin` | Where a user in UAA originated from. See [here](./README.md#choose-your-authentication-method) for more information | YES | `null` | `ldap` |
| `api_client_id` | The client created [here](./README.md#create-the-ignition-uaa-client) | YES | `null` | `ignition` |
| `api_client_secret` | The secret for the client created [here](./README.md#create-the-ignition-uaa-client) | YES | `null` | `supersecretvalue` |
| `authorized_domain` | This is the email domain that valid users belong to | YES | `null` | ` @example.net` |
| `auth_variant` | Determines whether you are using the SSO tile or OpenID | NO | `p-identity` | `p-identity` or `openid` |
| `auth_scopes` | Only change this if you have a specific reason to. | NO | `openid,profile,user_attributes` | `openid,profile,user_attributes` |
| `auth_servicename` | The name of the SSO service instance | NO | `ignition-identity` | `myssosi` |
| `auth_url` | The base Authorization URL. If `auth_variant` is `p-identity`, this is provided by the `auth_servicename` service | NO | `null` | `https://login.YOURSYSTEMDOMAIN` |
| `client_id` | The client ID used for logging into ignition. If `auth_variant` is `p-identity`, this is provided by the `auth_servicename` service | NO | `null` | `ignition-login`
| `client_secret` | The client secret used for logging into ignition. If `auth_variant` is `p-identity`, this is provided by the `auth_servicename` service | NO | `null` | `supersecret` |
| `skip_tls_validation` | Please don't. | NO | `false` | `true` |
| `org_prefix` | When set, created orgs will be named `org_prefix-username` | NO | `ignition` | `mycompany` |
| `org_count_update_interval` | How often to poll the Cloud Controller API for stats on created orgs. Formatted as a Go [Duration](https://golang.org/pkg/time/#Duration) | NO | `1m` | `1h` |
| `space_name` | The name of the space to create in each vended org | NO | `playground` | `sandbox` |
| `quota_name` | The name of the quota associated with each vended org | NO | `ignition` | `default` |
| `iso_segment_name` | The name of the Isolation Segment each vended org will be associated with | NO | `shared` | `iso-01` |
