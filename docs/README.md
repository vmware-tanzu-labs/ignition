# How to Install Ignition

## Prepare Cloud Foundry

1. Install the
[`Pivotal Application Service`](https://network.pivotal.io/products/elastic-runtime)
(PAS) tile

## Create necessary artifacts
1. Create an `ignition` org: `cf create-org ignition`. This will grant your user
   `OrgManager` permissions.
1. Create a `production` space in the ignition org: `cf create-space production
   -o ignition`. This will grant your user `SpaceManager` and `SpaceDeveloper`
   permissions. This is where the `ignition` app will live.
1. Create an `ignition` quota: `cf create-quota ignition -m 10G -i -1 -r 1000
   -s 100 --allow-paid-service-plans -a -1 --reserved-route-ports 1`. When
   `ignition` creates orgs for users, it will set this quota on the newly created
   org.

## Create the `ignition` UAA client
Any UAA commands will be run with the `uaa` CLI, which can be found
[here](https://github.com/cloudfoundry-incubator/uaa-cli). Equivalent `uaac`
commands exist but are out of scope for this document.

The `ignition` client is an OAuth2 client with admin privileges so that it can
create orgs for users. To create the client, run the following command:

```shell
$ uaa create-client ignition -s <client-secret> \
    --authorized_grant_types client_credentials \
    --scope cloud_controller.admin,scim.write,scim.read \
    --authorities cloud_controller.admin,scim.write,scim.read
```

## Create the Ignition Config User Provided Service
This user provided service instance configures ignition for your environment.
Create a file called `ignition-config.json`, and include the following required
attributes:
```json
{
  "session_secret": "REQUIRED",
  "system_domain": "run.example.net",
  "api_client_id": "ignition",
  "api_client_secret": "REQUIRED",
  "authorized_domain": "@example.net",
}
```
Please see the [glossary](./config-options.md) for definitions for available fields.

## Choose Your Authentication Method
Before creating the service in PAS, you must choose which authentication method
you wish to use, and further configure the JSON file. Choose the appropriate link
for your authentication method:
1. [Pivotal SSO Tile](./sso.md)
1. [Internal PAS UAA](./internal_uaa.md)
1. [External OpenID Connect Provider](./oidc.md)

## Finish the JSON and Create the Service in PAS
Once you have set your authentication method, add any [optional fields](./config-options.md)
you need for your deployment. Create the service in PAS by running the commands:
```shell
$ cf target -o ignition -s production
api endpoint:   https://api.run.example.net
api version:    2.131.0
user:           admin
org:            ignition
space:          production

$ cf create-user-provided-service ignition-config -p ignition-config.json
Creating user provided service ignition-config in org ignition / space production as admin...
OK

```

## Deploy `ignition`
1. Download the [latest release](https://github.com/pivotalservices/ignition/releases/latest)
   from Github and expand it **into its own directory**.
1. In that directory, create a file called `manifest.yml` that looks like this:
   ```yaml
   ---
   applications:
   - name: ignition
     memory: 128M
     instances: 2
     buildpacks: [binary_buildpack]
     command: ./ignition
     services:
     - ignition-config
     - ignition-identity # Only include this line if you chose SSO as your auth method
   ```
1. Deploy the app to the correct org and space:
   ```shell
   $ cf target -o ignition -s production
   api endpoint:   https://api.run.example.net
   api version:    2.131.0
   user:           admin
   org:            ignition
   space:          production

   $ cf push
   Pushing from manifest to org ignition / space production as admin...

   Getting app info...
   Creating app with these attributes...
   + name:         ignition
     path:         /tmp/ignition
     buildpacks:
   +   binary_buildpack
   + command:      ./ignition
   + instances:    2
   + memory:       128M
     services:
   +   ignition-config
     routes:
   +   ignition.apps.example.net

   Creating app ignition...
   Mapping routes...
   Binding services...
   Comparing local files to remote cache...
   Packaging files to upload...
   Uploading files...
    347.70 KiB / 347.70 KiB [=======================================] 100.00% 1s

   Waiting for API to complete processing files...

   name:              ignition
   requested state:   started
   routes:            ignition.apps.example.net
   last uploaded:     Wed 15 May 09:28:24 MDT 2019
   stack:             cflinuxfs3
   buildpacks:        binary

   type:           web
   instances:      2/2
   memory usage:   128M
        state     since                  cpu    memory          disk          details
   #0   running   2019-05-15T15:28:43Z   0.2%   12.6M of 128M   14.1M of 1G
   #1   running   2019-05-15T15:28:33Z   0.2%   12.9M of 128M   14.1M of 1G
   ```