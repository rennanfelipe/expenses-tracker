# Requirements doc for Expenses Tracker:

This simple system aims to create a set of interfaces (such as CLI, TUI, and web UI) for an ETL process, which will load credit card statements, as CSVs files, into a normalized database to allow report generation in CLI, PDF, or HTML formats.

## Problem statement:

Manage personal expenses is a time consuming task, replicate credit card statements into a spreadsheet is time consuming, seting up a form, with solutions like Google Forns accompanied with Google Sheets, can help reducing errors and the general input time. 

With all statements put on a spread sheet the second challenge is the repots generation. Spread sheets are tricky when it comes to dates queryes, a column set to ISO date can have some cell interpreted as month-day-year, breaking formula logic.

The other solution is to trust ones data into a third party expense tracker (apps, web sites or desktop) to import and generate reports. 

But Expense trackers can be too complex to set up, too pricy to afford or they have issues during the automated import process, often categorizing expenses wrongly and even simply not loading it at all. 
I haven't got a simple software just to load data, categorize into a small set of categories and generate simple reports yet.

## Project vision:

To create a very simple-to-use tool to extract, load, and consult its data. It has to work with Brazilian banks' Inter and Bradesco CSV formats and output data according to which interface is querying (TUI, CLI web).
