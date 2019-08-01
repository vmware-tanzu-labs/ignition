# Ignition SSO with the Pivotal Single-Sign On Service
The Pivotal SSO tile provides a way to allow subsets of your users to use
`ignition`; for example, allowing one development team access to the Ignition app
then letting others on as they onboard.

## Prepare Cloud Foundry
Install the [`Pivotal Single Sign-On for PCF`](https://network.pivotal.io/products/pivotal_single_sign-on_service)
tile.

## Create a Dev Service Plan
Once the SSO tile is installed, you'll need to create a service plan for the
ignition app. Navigate to `https://p-identity.YOUR-SYSTEM-DOMAIN` and log in using
PAS tile admin level credentials. If your PAS account is not an admin you can
use the UAA `Admin Credentials` found in the PAS tile Credentials tab in Ops
Manager.

Once logged in, create a SSO service plan. This plan will be used by ignition
and potentially all other applications on the platform for authentication. Click
the `New Plan` button in the top left corner.

* Plan Name: `Dev SSO Plan`
* Description: `SSO Plan for application developers`
* Auth Domain: `dev-sso`
* Instance Name: `SSO Login for App Devs`

You can optionally select other users you want to be able to administer the
service plan you're creating now. For organizations either check the `Enable for
all Orgs` checkbox or enter the org where Ignition will be deployed too. Click
`Save Plan`.

This tutorial discusses creating an identity plan backed by an LDAP server. To
use SSO, follow the [official docs](https://docs.pivotal.io/p-identity/1-8/configure-id-providers.html#config-saml-prov)
and keep in mind some of the values below.

On the new plan drop down select `Manage Identity Providers`, and then choose
`New Identity Provider`. Enter the following details:
* **Identity Provider Name**: Enter the name of your identity provider or domain
* **Identity Provider Description**: Enter `application developers`
* **Identity Provider Type**: `ldap`
* **Hostname**: Enter the IP or hostname of your LDAP server. It must be the same
  LDAP server configured in the PAS tile's "Authentication and Enterprise SSO"
  section
* **Port**: `389`, or `636` if using secure ldap.
* **Security Protocol**: `None` if using insecure ldap, otherwise select `SSL`
  or `TLS`
* **Referral**: `follow`, unless you turned this off in PAS
* **User DN**: Enter your ldap service account name, for example `ldapsvc`.
  This is the service account used to query LDAP.
* **Bind Password**: Enter the ldap service account password.
* **Users Search Base**: Same as PAS, for example `dc=ad,dc=pcf,dc=example,dc=com`
* **Users Search Filter**: Same as PAS, for ignition it has been tested with
  `sAMAccountName={0}`
* **Just in Time Provisioning**: `checked`
* **Groups Search Base**: Same as PAS, for example `dc=ad,dc=pcf,dc=example,dc=com`
* **Groups Search Filter**: Same as PAS, for example `member={0}`
* **Email Domains**: Enter your email domain

Under Advanced Settings you need to map a few user attributes from LDAP that
ignition requires. Add the following User Schema Attribute to Attribute Name
mappings:

* `given_name` -> `cn`
* `family_name` -> `sn`
* `email` -> `mail`

Click `Save Identity Provider`.

## Create Ignition-Identity SSO Service Instance
Create a Single Sign-On service instance named `ignition-identity` in your space.

```bash
$ cf target -o ignition -s production
api endpoint:   https://api.run.example.net
api version:    2.131.0
user:           admin
org:            ignition
space:          production

$ cf create-service p-identity sso ignition-identity
```

## Create the App for Ignition
1. Follow step 1 (and **only** step 1) from the [main installation
   guide](./README.md#deploy-ignition)
1. In the directory from the previous step, run
   ```shell
   $ cf push ignition --no-start
   $ cf bind-service ignition ignition-identity
   ```
1. Go to the Apps Manager at `https://apps.YOUR-SYSTEM-DOMAIN` and navigate to the
   `ignition` app. Click on the `Services` tab at the top. Find `ignition-identity`
   and click the `Manage` link.
1. Click the `New App` button/link and the fill in the following details:
   * **App Name**: `ignition`
   * **App Launch Url**: Your ignition app URL of `https://ignition.YOUR-APPS-DOMAIN` (for example,
   `https://ignition.apps.example.net`)
   * **Identity Providers**: Select your ldap sso identity provider, and unselect
   the internal user store.
   * **Auth Redirect URIs**: Enter the same URL you used for the App Launch Url
   above.
   * **Authorization**:
     * **Profile**: `profile`
     * **System Permissions**: `openid`, `user_attributes`
   * **Auto Approve Scopes**: `All Selected`

   Click `Save Config`.

## Modify Ignition Config for SSO
1. Complete the steps located
   [here](./README.md#create-the-ignition-config-user-provided-service)
1. Modify `ignition-config.json` from this:
   ```json
   {
     "session_secret": "YOUR_SESSION_SECRET",
     "system_domain": "YOUR-SYSTEM-DOMAIN",
     "api_client_id": "ignition",
     "api_client_secret": "UAA_IGNITION_CLIENT_SECRET",
     "authorized_domain": "@example.net",
   }
   ```
   to look like this
   ```json
   {
     "session_secret": "YOUR_SESSION_SECRET",
     "system_domain": "YOUR-SYSTEM-DOMAIN",
     "api_client_id": "ignition",
     "api_client_secret": "UAA_IGNITION_CLIENT_SECRET",
     "authorized_domain": "@example.net",
     "uaa_origin": "ldap"
   }
   ```

Return to the [main installation
instructions](./README.md#finish-the-json-and-create-the-service-in-pas)
and finish the instructions on that page.

## Links
* https://docs.pivotal.io/pivotalcf/2-5/opsguide/auth-sso.html#configure-ldap
* https://docs.pivotal.io/pivotalcf/2-5/opsguide/external-user-management.html#user-role
* https://discuss.pivotal.io/hc/en-us/articles/204140418-Configuring-LDAP-Integration-with-Pivotal-Cloud-Foundry-
