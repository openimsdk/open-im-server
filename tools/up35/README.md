#  README for OpenIM Server Data Conversion Tool

## Overview

This tool is part of the OpenIM Server suite, specifically designed for data conversion between MySQL and MongoDB databases. It handles the migration of various data types, including user information, friendships, group memberships, and more from a MySQL database to MongoDB, ensuring data consistency and integrity during the transition.

## Features

+ **Configurable Database Connections:** Supports connections to both MySQL and MongoDB, configurable through a YAML file.
+ **Data Conversion Tasks:** Converts a range of data models, including user profiles, friend requests, group memberships, and logs.
+ **Version Control:** Maintains data versioning, ensuring only necessary migrations are performed.
+ **Error Handling:** Robust error handling for database connectivity and query execution.

## Requirements

+ Go programming environment
+ MySQL and MongoDB servers
+ OpenIM Server dependencies installed

## Installation

1. Ensure Go is installed and set up on your system.
2. Clone the OpenIM Server repository.
3. Navigate to the directory containing this tool.
4. Install required dependencies.

## Configuration

+ Configuration is managed through a YAML file specified at runtime.
+ Set up the MySQL and MongoDB connection parameters in the config file.

## Usage

To run the tool, use the following command from the terminal:

```go
make build BINS="up35"
```

Where `path/to/config.yaml` is the path to your configuration file.

## Functionality

The main functions of the script include:

+ `InitConfig(path string)`: Reads and parses the YAML configuration file.
+ `GetMysql()`: Establishes a connection to the MySQL database.
+ `GetMongo()`: Establishes a connection to the MongoDB database.
+ `Main(path string)`: Orchestrates the data migration process.
+ `SetMongoDataVersion(db *mongo.Database, curver string)`: Updates the data version in MongoDB after migration.
+ `NewTask(...)`: Generic function to handle the migration of different data types.
+ `insertMany(coll *mongo.Collection, objs []any)`: Inserts multiple records into a MongoDB collection.
+ `getColl(obj any)`: Retrieves the MongoDB collection associated with a given object.
+ `convert struct`: Contains methods for converting MySQL models to MongoDB models.

## Notes

+ Ensure that the MySQL and MongoDB instances are accessible and that the credentials provided in the config file are correct.
+ It is advisable to backup databases before running the migration to prevent data loss.

## Contributing

Contributions to improve the tool or address issues are welcome. Please follow the project's contribution guidelines.

## License

Refer to the project's license document for usage and distribution rights.