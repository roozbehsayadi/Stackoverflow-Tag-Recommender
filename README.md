# Recommender System with Stackoverflow's Data

This project is written to recommend tags based on a provided one.

## Running the Project

- Create a database with name "stackoverflow\_recommender" on postgresql.
- Run migration files with `migrate.sh` file. See the file's content for additional information.
- Dump data in database with `loaddata.sh` script.
- Run calculate_similarities.go. The application's output will be written to similarities.txt.
