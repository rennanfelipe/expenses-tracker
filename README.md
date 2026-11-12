# Expenses Tracker - ET!

       .""".
      /      \
     |        |
     | (0) (0)|
      \  \_/ /
       `.___.'
      /  |  \
     /   |   \
        / \
       |   |

ET! is a simple personal expense tracking application.


It's a simple exprense tracker that read CSV files from credit card statements, normalize it and impor it into a corherent SQL database to generate reports.

The initial goal is to implement a CLI tool; afterward, it will create a TUI or web version.

## Tech stack

**Language:** Go Lang.
**Database:** MariaDB.
**User Interface:** CLI (initial), TUI/WEB (later).

## Architecture

The architecture is a traditional MVC:

Model: Classes to read the CSV files and persist data on the DB.

View: TUI classes to present data and forms to load more data to the database.

Controller: Classes responsible for organizing CSV data from the model into coherent data for the DB. Each bank statement will have its own Controller that knows how to read its specific CSV format, like Inter and Bradesco.atement will have it's own Controller that knows how to read it's specific CSV format, like Inter and Bradesco.

## Build

TBD

## Installation

TBD

## Docs

Docs for this project are present in the those folders:

```shell
.
README.md                       # Project entripoint document
AGENTS.md                       # Instructions for AI coding assistants.
├── ai                          # Describes how an AI coding assistant should work within the project.
│   ├── architecture-rules.md   # How architecture should the approached.
│   ├── coding-standards.md     # How code should be written..
│   └── testing-guidelines.md   # How correctness is verified.
└── docs                        # Describes the system: what it does, how it works, and why decisions were made
    ├── architecture.md         # How the system is structured.
    ├── data-model.md           # Entities and relationships.
    ├── decisions               # Why particular technical decisions were made.
    └── requirements.md         # What the system must accomplish.
``` 