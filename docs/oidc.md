# Ignition SSO with an External OpenID Connect Provider
While you can use ignition without the [Single Sign-On for PCF](./sso.md) service,
we recommend you use either SSO or the [Internal UAA](./internal_uaa.md) to
authenticate and authorize users for the application. It is possible that an
external OpenID Connect compliant provider can be used directly, but it is not
recommended.

When you want to use Google (or an equivalent OpenID Connect provider) to
authenticate users (i.e. you do _not_ want to use the SSO service or the internal UAA),
follow these instructions.

1. Complete the steps located
   [here](./README.md#create-the-ignition-config-user-provided-service)
1. Generate a Google [OAuth2 Client ID and Secret](https://console.developers.google.com/apis/credentials),
   and add them to the existing json like below:
   ```json
   {
     "session_secret": "YOUR_SESSION_SECRET",
     "system_domain": "YOUR-SYSTEM-DOMAIN",
     "api_client_id": "ignition",
     "api_client_secret": "UAA_IGNITION_CLIENT_SECRET",
     "authorized_domain": "@example.net",
   }
   ```
   should become
   ```json
   {
     "session_secret": "YOUR_SESSION_SECRET",
     "system_domain": "YOUR-SYSTEM-DOMAIN",
     "api_client_id": "ignition",
     "api_client_secret": "UAA_IGNITION_CLIENT_SECRET",
     "authorized_domain": "@example.net",
     "uaa_origin": "google",
     "auth_variant": "openid",
     "auth_scopes": "openid,profile,email",
     "auth_url": "https://accounts.google.com",
     "client_id": "your-google-client-id-here",
     "client_secret": "your-google-client-secret-here"
   }
   ```
Return to the [main installation
instructions](./README.md#finish-the-json-and-create-the-service-in-pas)
and finish the instructions on that page.
