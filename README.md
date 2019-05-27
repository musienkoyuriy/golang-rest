# Too Many Strains

The goal of this challenge is to take a subset of data, migrate it over to a MySQL database instance and build an API that can interact with the data.

With this document, you should have also received a JSON file titled: `strains.json`.

1. Given a blank MySQL database, build tables to store the information in `strains.json` in a relational manner.

2. Using `Go`, write an import script that will be able to process `strains.json` and upload the information into the database designed in step 1.

    - This script should create the table or tables needed to store the data
    - This script should then upload the data
    - This script should run by itself using `go run`. Suggestion is to place in a separate folder in your repository

3. Build an API that can interact with the MySQL database using `Go`. Implement the following functionality. It is up to you how many different handlers are needed.
    - Create new strain
    - Edit an existing strain by ID
    - Remove a strain
    - Search for strains by name 
    - Search for strains by race (Available races: Sativa, Indica, and Hybrid)
    - Search for strains by effect
    - Search for strains by flavor


Things to note:
- You may use any external library you wish, but please use `dep` (https://golang.github.io/dep/) and include the `vendor` directory in your repository.

- The resulting API should be able to handle multiple requests at once.

- The resulting API should run on `localhost` and accept requests on port `:8888`.

- For the MySQL DB, you can use Docker to install a local MySQL container and instructions have been provided for you in the folder `Docker`.

- Your code should be properly documented using `godoc` formatting (https://blog.golang.org/godoc-documenting-go-code)



## Up & Running

1. Execute schema.sql outside the docker container to create database with tables

```
mysql --host=127.0.0.1 --user=root -p < schema.sql
```

2. Run script which will fill in database with data from `strains.json` file

```
go run cmd/importToDB/importToDB.go
```

3. Run the app

```
go run main.go model.go utils.go
```

It will do these actions:

- Creates a database connection
- initializes an http routes


## Notes

There are several notes and things to enhance for this REST API:

- For production we need to envisage an SQL Injection
- For production we should use foreign keys to database constistence
- Also there is another way to store the data which presented in `strains.json`

We can use a `SET` operator to store the data collections in one column

```
CREATE TABLE strains (
id NOT NULL AUTO_INCREMENT,
name  VARCHAR(50) NOT NULL,
race VARCHAR(50) NOT NULL,
flavors SET('Earthy', 'Chemical', 'Pine'),
effects_positive SET('Relaxed', 'Hungry', 'Happy', 'Sleepy'),
effects_negative SET('Dizzy'),
effects_medical SET('Depression', 'Insomnia', 'Pain', 'Stress', 'Lack of Appetite');
```

But I implemented this task using the relational model with multiple tables due to task requirements
