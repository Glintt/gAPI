## Groups

gAPI allows to associate APIs to an application group.

These groups don't have any logic associated, they exists simply to better organize your APIs.

#### Create a group

1. Open the dashboard
2. Navbar button "Apps"
3. Create Application
4. Fill in the form and save

#### Associate API with Application Group

1. Open the dashboard
2. Open API information
3. API configuration section
4. From "Associated Application Group", select the group you want to associate the API to
5. Associate button

#### Database table

When using relational database, the table where groups are stored is: **_gapi_services_groups_**. This table contains a name, an identifier and a reachable status.

Each service contains a column/attribute which is a foreign key to the group's identifier.
