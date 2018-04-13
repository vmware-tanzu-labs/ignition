## Ignition
[![CircleCI](https://circleci.com/gh/pivotalservices/ignition/tree/master.svg?style=svg)](https://circleci.com/gh/pivotalservices/ignition/tree/master)

A landing page for developers to self-service their way onto your Pivotal Cloud Foundry (PCF) deployment(s).

* Authenticates the user via OpenID Connect (which implicitly uses OAuth 2.0)
* Allows the user to access Apps Manager and view their personal PCF org

### Contribute

This application is a combination of a JavaScript single-page app (built with React) and a Go web app. The JavaScript app is built into a JavaScript bundle that the Go web app serves up. The Go web app also provides an API that the JavaScript app uses to function.

#### Yak Shaving (Developer Setup)

This project uses [`dep`](https://github.com/golang/dep) and [`yarn`](https://yarnpkg.com) for dependency management.

The following setup script shows how to get your MacOS workstation ready for `ignition` development. Don't just blindly execute shell scripts though; [take a thorough look through it](https://raw.githubusercontent.com/pivotalservices/ignition/master/setup.sh) and then run the following:

> `curl -o- https://raw.githubusercontent.com/pivotalservices/ignition/master/setup.sh | bash`

#### Add A Feature / Fix An Issue

We welcome pull requests to add additional functionality or fix issues. Please follow this procedure to get started:

1. Ensure you have `go` `>=1.10.x` and `node` `v8.x.x` installed
1. Ensure your `$GOPATH` is set; this is typically `$HOME/go`
1. Clone this repository: `go get -u github.com/pivotalservices/ignition`
1. Go to the repo root: `cd $GOPATH/src/github.com/pivotalservices/ignition`
1. *Fork this repository*
1. Add your fork as a new remote: `git remote add fork https://github.com/INSERT-YOUR-USERNAME-HERE/ignition.git`
1. Create a local branch: `git checkout -b your initials-your-feature-name` (e.g. `git checkout -b jf-add-logo`)
1. Make your changes, ensure you add tests to cover the changes, and then validate that all changes pass (see `Run all tests` below)
1. Push your feature branch to your fork: `git push fork your initials-your-feature-name` (e.g. `git push fork jf-add-logo`)
1. Make a pull request: `https://github.com/pivotalservices/ignition/compare/master...YOUR-USERNAME-HERE:your-initials-your-feature-name`

### Configure the application
#### Authentication
We recommend you use the [Single Sign-On for PCF](https://network.pivotal.io/products/pivotal_single_sign-on_service) service to authenticate and authorize users for the application. It is possible that an OpenID Connect compliant provider can be used directly, but it is not recommended.

To authenticate users with the Single Sign-On for PCF service:
1. Configure the Single Sign-On for PCF tile in your PCF foundation http://docs.pivotal.io/p-identity/

   Single Sign-On service is a multi tenancy openid connect server. It can integrate with other openid connect, ldap and saml servers as identity providers
1. Create a Single Sign-On service instance named `ignition-identity` in your space.
1. Create a user provided service `cf cups ignition-config -p /path/to/ignition-config.json` (see below for an `ignition-config.json` template)
1. Build ignition `./build_local.sh`
1. Push Ignition to Cloud Foundry and bind both services

  ```
  cf push ignition -p build -c './ignition-linux' -b binary_buildpack --no-start
  cf bind-service ignition ignition-config
  cf bind-service ignition ignition-identity
  cf start ignition
  ```

#### Templates For `ignition-config.json`

When you have a bound Single Sign-On service instance:

```json
{
  "session_secret": "encrypt the secure cookie used to store a user's session",
  "system_domain": "system domain of PAS",
  "uaa_origin": "okta|saml|ldap",
  "api_username": "cloud controller username",
  "api_password": "cloud controller password",
  "authorized_domain": "@example.net make sure open id token user profile email domain"
}
```

When you want to use Google (or an equivalent OpenID Connect provider) to authenticate users (i.e. you do _not_ have a bound Single Sign-On service instance):

Generate a Google [OAuth2 Client ID and Secret](https://console.developers.google.com/apis/credentials), and use them below:

```json
{
  "session_secret": "",
  "system_domain": "run.example.net",
  "uaa_origin": "okta",
  "api_username": "ignition",
  "api_password": "password",
  "authorized_domain": "@example.net",
  "auth_variant": "openid",
  "auth_scopes": "openid,profile,email",
  "auth_url": "https://accounts.google.com",
  "client_id": "your-client-id-here",
  "client_secret": "your-client-secret-here"  
}
```

### Run the application locally

You will need to ensure your environment contains the relevant variables for the app to run. Here is an example `$GOPATH/src/github.com/pivotalservices/ignition/credentials/export.sh`:

```sh
#!/bin/sh

### Server ###
export IGNITION_SCHEME="http" # IGNITION_SCHEME allows you to use http for local development; it is always set to HTTPS on PCF
export IGNITION_DOMAIN="localhost" # IGNITION_DOMAIN allows you to set the domain that will be used to access the app
export IGNITION_PORT="3000" # IGNITION_PORT is the port used to access ignition; this is always set to 443 on PCF
export IGNITION_SERVE_PORT="3000" # IGNITION_SERVE_PORT is the port that ignition listens on; this is usually different to IGNITION_PORT except during development
# export IGNITION_WEB_ROOT="" # IGNITION_WEB_ROOT can be used to store JS / CSS / image resources at a non-default path
export IGNITION_SESSION_SECRET="insert-a-random-session-secret-here" # IGNITION_SESSION_SECRET is used to encrypt the contents of the secure cookie used to store a user's session information
export IGNITION_COMPANY_NAME="Company Name" # IGNITION_COMPANY_NAME is used to white label the UX for ignition

### Your CF Deployment ###
export IGNITION_SYSTEM_DOMAIN="run.example.net" # IGNITION_SYSTEM_DOMAIN is what you get when you take the "api." away from the Cloud Controller API URL
export IGNITION_UAA_ORIGIN="okta" # IGNITION_UAA_ORIGIN is the origin for a user that logs in to Cloud Foundry with your single sign on solution of choice
export IGNITION_API_CLIENT_ID="cf" # IGNITION_API_CLIENT_ID is almost always cf
export IGNITION_API_CLIENT_SECRET="" # IGNITION_API_CLIENT_SECRET is almost always blank
export IGNITION_API_USERNAME="ignition" # IGNITION_API_USERNAME is the username for the user that can create Cloud Foundry organizations
export IGNITION_API_PASSWORD="password" # IGNITION_API_PASSWORD is the password for the user that can create Cloud Foundry organizations

### Developer Experimentation ###
export IGNITION_ORG_PREFIX="ignition" # IGNITION_ORG_PREFIX is used to generate a developer's org name (e.g. ignition-testuser)
export IGNITION_QUOTA_NAME="ignition" # IGNITION_QUOTA_NAME is used to generate a developer's org with the appropriate quota
export IGNITION_SPACE_NAME="playground" # IGNITION_SPACE_NAME is used to create the initial space in a developer's org

### Authorization ###
export IGNITION_AUTHORIZED_DOMAIN="@example.net" # IGNITION_AUTHORIZED_DOMAIN is used to validate that users are allowed to access the application

### Authentication ###
### Single Sign-On ###
export IGNITION_AUTH_VARIANT="openid" # IGNITION_AUTH_VARIANT is openid when you're working locally because you don't have a bound sso service instance
export IGNITION_CLIENT_ID="your-service-instance-client-id"
export IGNITION_CLIENT_SECRET="your-service-instance-client-secret"
export IGNITION_AUTH_URL="https://ignition.login.run.example.net"

### Google ###
# export IGNITION_AUTH_VARIANT="openid"
# export IGNITION_CLIENT_ID="your-client-id-here"
# export IGNITION_CLIENT_SECRET="your-client-secret-here"
# export IGNITION_AUTH_URL="https://accounts.google.com"
# export IGNITION_AUTH_SCOPES="openid,profile,email" # IGNITION_AUTH_SCOPES is not the same for Google as it is for the a Single Sign-On instance, and this allows you to override it with a comma separated list of values
```

1. Make sure you're in the repository root directory: `cd $GOPATH/src/github.com/pivotalservices/ignition && . ./credentials/export.sh`
1. Ensure the web bundle is built: `pushd web && yarn install && yarn build && popd`
1. Start the go web app: `go run cmd/ignition/main.go`
1. Navigate to http://localhost:3000

### Run all tests

1. Make sure you're in the repository root directory: `cd $GOPATH/src/github.com/pivotalservices/ignition`
1. Run go tests: `go test ./...`
1. Run web tests: `pushd web && yarn ci && popd`

### Support

`ignition` is a community supported Pivotal Cloud Foundry add-on. [Opening an issue](https://github.com/pivotalservices/ignition/issues/new) for questions, feature requests and/or bugs is the best path to getting "support". We strive to be active in keeping this tool working and meeting your needs in a timely fashion.
