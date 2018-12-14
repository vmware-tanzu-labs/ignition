# Ignition Installation

## Prepare Cloud Foundry

1. Install the [`Pivotal Application Service`](https://network.pivotal.io/products/elastic-runtime) (PAS) tile
1. Install the [`Pivotal Single Sign-On for PCF`](https://network.pivotal.io/products/pivotal_single_sign-on_service) tile

## Create A p-identity Plan That Authenticates and Authorizes Your Users



## Prepare Orgs + Spaces

1. Create an `ignition` org: `cf create-org ignition`
  * Grant access to the ignition org to yourself: `cf set-org-role you@example.net ignition OrgManager`
1. Create a `production` space in the ignition org: `cf create-space production -o ignition`
  * Grant Space Manager access to the production space to yourself: `cf set-space-role you@example.net ignition production SpaceManager`
  * Grant Space Developer access to the production space to yourself: `cf set-space-role you@example.net ignition production SpaceDeveloper`
  * Grant Space Auditor access to the production space to yourself: `cf set-space-role you@example.net ignition production SpaceAuditor`
1. Create an `ignition` quota: `cf create-quota ignition -m 10G -i -1 -r 1000 -s 100 --allow-paid-service-plans -a -1 --reserved-route-ports 1`

## Create The `ignition-config` User Provided Service Instance

The user provided service instance configures ignition for your environment. Create an `ignition-config.json` file with the following attributes, modifying as necessary for your environment:

```json
{
  "session_secret": "",
  "system_domain": "run.example.net",
  "uaa_origin": "ldap",
  "api_client_id": "ignition",
  "api_client_secret": "",
  "authorized_domain": "@pivotal.io",
  "company_name": "Pivotal",
  "iso_segment_name": "iso-01"
}
```

Here's a reference of the available values:

* `session_secret`: The session secret is used to secure the cookie used to store a user's session. You should randomly generate the contents of this value and limit access to it.
* `system_domain`: The system domain is
* `uaa_origin`: This is used when creating UAA users while giving users access to your PAS deployment. The values are typically:
  * `uaa` for users that are authenticated by the UAA deployment (i.e. you are not using an external identity provider)
  * `ldap` for users that are authenticated by the LDAP provider (e.g. Active Directory)
  * `{OIDC provider alias}` for users authenticated via an OIDC provider
  * `{SAML provider alias}` for users authenticated via a SAML identity provider (e.g. `okta`)
* `api_client_id`: This is typically `ignition`.
* `api_client_secret`: This is the client secret created for the `ignition` client.
* `authorized_domain`: This is the email domain that valid users belong to (e.g. `pivotal.io`).
* `auth_variant`: This is `p-identity` by default. Only change this if you have a specific reason to.
* `auth_scopes`: This is `openid,profile,user_attributes` by default. Only change this if you have a specific reason to.
* `auth_servicename`: This is `ignition-identity` by default. Change this if you have a different `p-identity` service instance name.
* `auth_url`: This is supplied by the `ignition-identity` service instance.
* `client_id`: This is supplied by the `ignition-identity` service instance.
* `client_secret`: This is supplied by the `ignition-identity` service instance.
* `skip_tls_validation`:
* `org_prefix`
* `org_count_update_interval`:
* `space_name`:
* `quota_name`:
* `iso_segment_name`:

* `uaa create-client ignition -s <client-secret> --authorized_grant_types client_credentials --scope cloud_controller.admin,scim.write,scim.read --authorities cloud_controller.admin,scim.write,scim.read`

```yml
---
applications:
- name: ignition
  memory: 128M
  instances: 2
  buildpack: binary_buildpack
  command: ./ignition
  services:
    - ignition-config
    - ignition-identity
```
