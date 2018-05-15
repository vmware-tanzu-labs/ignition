# Ignition SSO with LDAP

You must use the same identity provider for the PCF SSO Tile bound to ignition as the Pivotal Application Service. These instructions assume you're going to be using LDAP for the single sign-on method, and walk you through configuring the Single Sign-On service and service instance to bind to Ignition. The Single Sign-On service is a multi tenancy openid connect server. It can integrate with other openid connect, ldap and saml servers as identity providers.

These instructions cover creating a SSO service plan backed by LDAP, the service instance, and a user provided service instance. Ignition requires two bound services:

* ignition-config
* ignition-identity

## Install the Pivotal SSO Tile

If you have not already done so, [install the Pivotal SSO Tile](http://docs.pivotal.io/p-identity/) into Operation Manager.

## Create Dev SSO Service Plan for Ignition

Once the SSO tile is installed, you'll need to create a service plan for the ignition app. Navigate to https://p-identity.YOURSYSTEMDOMAIN and login using PAS tile admin level credentials. If your PAS account is not an admin you can use the UAA admin account creds found in the PAS tile Credentials tab.

Once signed in create a SSO service plan. This plan will be used by ignition and potentially all other applications on the platform for authentication. Click the `New Plan` button in the top left corner.

* Plan Name: `Dev SSO Plan`
* Description: `SSO Plan for application developers`
* Auth Domain: `dev-sso`
* Instance Name: `SSO Login for App Devs`

You can optionally select other users you want to be able to administer the service plan we're creating now. For organizations either check the `Enable for all Orgs` checkbox or enter the org where Ignition will be deployed too. Click `Save Plan`.

On the new plan drop down select `Manage Identity Providers`. From here click `New Identity Provider` button. Enter the following details:

* **Identity Provider Name**: Enter the name of your identity provider or domain
* **Identity Provider Description**: Enter `application developers`
* **Identity Provider Type**: `ldap`
* **Hostname**: Enter the IP or hostname without protocol of your LDAP server, the same one used to configure LDAP in PAS.
* **Port**: `389`, or `636` if using secure ldap.
* **Security Protocol**: `None` if using insecure ldap, otherwise select `SSL` or `TLS`
* **Referral**: `follow`, unless you turned this off in PAS
* **User DN**: Enter your ldap service account name, for example `ldapsvc`. This is the service account used to query LDAP.
* **Bind Password**: Enter the ldap service account password.
* **Users Search Base**: Same as PAS, for example `dc=ad,dc=pcf,dc=example,dc=com`
* **Users Search Filter**: Same as PAS, for ignition it has been tested with `sAMAccountName={0}`
* **Just in Time Provisioning**: `checked`
* **Groups Search Base**: Same as PAS, for example `dc=ad,dc=pcf,dc=example,dc=com`
* **Groups Search Filter**: Same as PAS, for example `member={0}`
* **Email Domains**: Enter your email domain

Under Advanced Settings we need to map a few user attributes from LDAP that ignition requires. Add the following User Schema Attribute to Attribute Name mappings:

* `given_name` -> `cn`
* `family_name` -> `sn`
* `email` -> `mail`

Click `Save Identity Provider`.

## Create Ignition-Identity SSO Service Instance

Create a Single Sign-On service instance named `ignition-identity` in your space.

```bash
$ cf create-service p-identity sso ignition-identity
```

Once created you'll need to bind the `ignition-identity` service to your ignition app. Once you've bound the sso service to ignition, go into Apps Manager and find the bound service instance and click the `Manage` link in the top right corner. This will take you to that service instance's specific management page.

Click the `New App` button/link and the fill in the following details:

* **App Name**: `ignition`
* **App Launch Url**: Your ignition app URL, for example: https://ignition.apps.pcf.example.com
* **Identity Providers**: Select your ldap sso identity provider, and unselect the internal user store.
* **Auth Redirect URIs**: Enter the same URL you used for the App Launch Url above.
* **Authorization**:
  * **Profile**: `profile`
  * **System Permissions**: `openid`, `user_attributes`
* **Auto Approve Scopes**: `All Selected`

Click `Save Config`.

## Create Ignition-Config User Provided Service Instance

The last part required in order to run ignition is the user provided service instance that configures the application for your environment. Create a ignition-config.json file with the following attributes, modify as necessary for your environment.

```json
{
  "session_secret": "encrypt the secure cookie used to store a user's session",
  "system_domain": "system domain of PAS",
  "uaa_origin": "ldap",
  "api_client_id": "ignition",
  "api_client_secret": "uaa client secret",
  "authorized_domain": "@example.com"
}
```

The `session_secret` can be anything unique. The `uaa_origin` is the attribute to pay attention to, we need to ensure we set that to ldap when using the SSO service backed by ldap.

The `api_client_id` and `api_client_secret` are the client id and secret that ignition uses to talk to the Cloud Controller API and to UAA.

Lastly the `authorized_domain` controls what email domains are allowed to use ignition to create organizations. This is useful especially when you're using a OAuth provider that supports non-corporate domains - you don't want people creating orgs in PCF with their personal email accounts.

```bash
cf cups ignition-config -p /path/to/ignition-config.json
```

Restage the ignition app and you should now be able to login to ignition and create organizations.

## Push Ignition to PCF

1.  Build ignition `./build_local.sh`
1.  Push Ignition to Cloud Foundry and bind both services

```
cf push ignition -p build -c './ignition-linux' -b binary_buildpack --no-start
cf bind-service ignition ignition-config
cf bind-service ignition ignition-identity
cf start ignition
```

## Links

* https://docs.pivotal.io/pivotalcf/2-0/opsguide/auth-sso.html#configure-ldap
* https://docs.pivotal.io/pivotalcf/2-0/opsguide/external-user-management.html#user-role
* https://discuss.pivotal.io/hc/en-us/articles/204140418-Configuring-LDAP-Integration-with-Pivotal-Cloud-Foundry-
