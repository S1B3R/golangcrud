# GolangCRUD
# Objective
This repo objective is to implement a CRUD of users REST API.<br/>
MongoDb was chosen for storing data due to its simplicity and agility to store and retrieve data, connection string and database name are environment variables.<br/>
X-API-Key is being used to secure the API and is being passed as an environment variable.<br/>

===============================
# Improvements to be made:
Retrieve database connection string as well as API Key from a service like AWS Secrets Manager or Google Cloud Secret Manager.<br/>
Create an interface to manipulate DB operations, removing the need to create an entire new CRUD for every new entity.
