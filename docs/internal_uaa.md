# Ignition SSO with the Internal PAS UAA
Using the internal UAA allows a seamless integration between ignition and the PAS
system as a whole, with no external dependencies required.

1. Complete the steps located
   [here](./README.md#create-the-ignition-config-user-provided-service)
1. Create a second client used for logging in to ignition. Run a command similar
   to this:
   ```shell
   $ uaa create-client ignition-login -s <client-secret> \
       --authorized_grant_types password,authorization_code \
       --scope openid,profile,email \
       --authorities openid,profile,email,user_attributes \
       --redirect_uri "https://ignition.YOURAPPSDOMAIN" # or whatever URL you will use for ignition
   ```
1. Next, you need to modify the client to autoapprove all scopes:
   * Get the client JSON so it can be modified:
     ```shell
     $ uaa curl /oauth/clients/ignition-login > client.json
     ```
   * Modify `client.json` and change `"autoapprove": []` to `"autoapprove":
   ["openid","profile","email"]`
   * Update the client in the UAA:
     ```shell
     $ uaa curl -X PUT -H "content-type: application/json" \
         -H "accept: application/json" -d $(cat client.json) \
         /oauth/clients/ignition-login
     ```
1. Modify `ignition-config.json` from this:
   ```json
   {
     "session_secret": "REQUIRED",
     "system_domain": "run.example.net",
     "api_client_id": "ignition",
     "api_client_secret": "REQUIRED",
     "authorized_domain": "@example.net",
   }
   ```
   to this:
   ```json
   {
     "session_secret": "REQUIRED",
     "system_domain": "run.example.net",
     "api_client_id": "ignition",
     "api_client_secret": "REQUIRED",
     "authorized_domain": "@example.net",
     "uaa_origin": "<variant>",
     "auth_variant": "openid",
     "auth_scopes": "openid,profile,email",
     "auth_url": "https://login.YOURSYSTEMDOMAIN",
     "client_id": "ignition-login",
     "client_secret": "REQUIRED"
   }
   ```
   For `uaa_origin`, the value should be one of the following:
   * `uaa` if you are using purely internal users
   * `ldap` if PAS is configured to use LDAP
   * If PAS is using SAML, `uaa_origin` should be the value of the `Provider Name`
     field in the `Authentication and Enterprise SSO` section of the PAS Tile
     configuration in Ops Manager.

Return to the [main installation
instructions](./README.md#finish-the-json-and-create-the-service-in-pas)
and finish the instructions on that page.
