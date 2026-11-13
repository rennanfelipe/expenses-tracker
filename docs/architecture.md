# Architecture

The architecture is a traditional MVC:

Model: Classes to read the CSV files and persist data on the DB.

View: TUI, CLI and Web classes to present data and forms to load more data to the database.

Controller: Classes responsible for organizing CSV data from the model into coherent data for the DB. Each bank statement will have its own Controller that knows how to read its specific CSV format, like Inter and Bradesco.atement will have it's own Controller that knows how to read it's specific CSV format, like Inter and Bradesco.