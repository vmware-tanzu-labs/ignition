# Ignition SSO without SSO Tile
While you can use ignition without the [Single Sign-On for PCF](https://network.pivotal.io/products/pivotal_single_sign-on_service) service, we recommend you use the Single Sign-On for PCF service to authenticate and authorize users for the application. It is possible that an OpenID Connect compliant provider can be used directly, but it is not recommended.

When you want to use Google (or an equivalent OpenID Connect provider) to authenticate users (i.e. you do _not_ have a bound Single Sign-On service instance), follow these instructions.

These instructions cover creating a user provided service instance. When you're _not_ using the Single Sign-On for PCF service, you only need the user provided service named `ignition-config` bound to the ignition app.

## Create Ignition-Config User Provided Service Instance
The only service required in order to run ignition is the user provided service instance that configures the application for your environment. Create a ignition-config.json file with the following attributes, modify as necessary for your environment.

Generate a Google [OAuth2 Client ID and Secret](https://console.developers.google.com/apis/credentials), and use them below:

```json
{
  "session_secret": "",
  "system_domain": "run.example.net",
  "uaa_origin": "okta",
  "api_client_id": "ignition",
  "api_client_secret": "client secret",
  "authorized_domain": "@example.net",
  "auth_variant": "openid",
  "auth_scopes": "openid,profile,email",
  "auth_url": "https://accounts.google.com",
  "client_id": "your-client-id-here",
  "client_secret": "your-client-secret-here"  
}
```

The `session_secret` can be anything unique. The `api_client_id` and `api_client_secret` are the client id and secret that ignition uses to talk to the Cloud Controller API and to UAA.

Lastly the `authorized_domain` controls what email domains are allowed to use ignition to create organizations. This is useful especially when you're using a OAuth provider that supports non-corporate domains - you don't want people creating orgs in PCF with their personal email accounts.

```bash
cf cups ignition-config -p /path/to/ignition-config.json
```

## Push Ignition to PCF
1. Build ignition `./build_local.sh`
1. Push Ignition to Cloud Foundry and bind both services

  ```
  cf push ignition -p build -c './ignition-linux' -b binary_buildpack --no-start
  cf bind-service ignition ignition-config
  cf start ignition
  ```
