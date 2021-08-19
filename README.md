# AHSAR_Web
Development In Progress

## Introduction
Web Application for AHSAR
<br></br>
Check the sentiment analysis result of student's commentary on the professor of your choice by entering the `tid` in a professor URL from RateMyProfessors.com website. 

For example, assume a "randomly-selected" professor URL is `https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994` (Salute to Professor Adam Meyers from NYU, CS0002 ICP and CS480 NLP), enter `2105994` but not the name of professor or the full URL. 

Notice that Sentiment Score (discrete) is computed based on individual comments, while Sentiment Score (continuous) is computed based on all comments.
In other words, the higher the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive.

## Application Structure (Not Finished)
* Frontend (Undecided, potentially JavaScript Vue.js/React.js framework) 
* Backend
    * Backend Server (Language: Go, Framework: Gin) 
    * Query & Result Cache (MiddleWare: Redis, written in C) 
    * Internal Communication: Naive Socket TCP Connection
    * Server: Modified Version of AHSAR NLP project (Language: Python) 
* Language Environment Setup: C, Go, Python. Check source file in this repository to find the actual Packages/Modules involved. Full list of dependencies will be listed when air in production mode. 

## Application Setup (Not Finished)
* Prepare Redis beforehand...!
* Start Redis (in Redis directory, `src/redis-server`)
* Start NLP Server (in base directory, `python3 pysrc/NLP_server.py`)
* Start Backend Server (in base directory, `bash run.bash` Notice that this is not finalized)
* Start Frontend / Access directly on browser

## Application Workflow (Not Finished)
* Frontend sends query to Backend
* Backend receives query and check if it is in Redis:
    * If in, retrieves the cached result and returns it to Frontend
    * If not in, sends the query to NLP Server through TCP Socket
* NLP Server receives the query from TCP Socket and start analyzing
* NLP Server returns the result to Backend Server through TCP Socket
* Backend updates Redis with the result and also returns result to Frontend

## NLP Server
For the full project (including datebase of __80k labeled RMP comments__ and other imported data, codebase of __RMP scraper__ and __N-gram algorithm__, and __reference__ list for the imported data) of the NLP Server behind the screen, called __AHSAR__ *Ad-Hoc Sentiment Analysis on RateMyProfessors*, please check this [Repository](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP). Bear with the badly optimized code ^^. 

## License
Project under MIT License. Basically, feel free to adopt any code from here for any usage, with no warranty, promise, or liability from the repository owners and collaborators. But a little bit of credit/reference is very appreciated. 

## Project History:
*   ...
*   2021/08/19:
    *   Redis cache set with expiration limit
*   2021/08/18:
    *   Backend Server now communications with NLP Server through Naive Socket TCP Connection. 
    *   New NLP Server with Naive Socket TCP Connection and multi-processing pool. 
    *   Removal of Kafka dependency. 
*   2021/08/17: 
    *   First push to GitHub Repository. 
    *   Simple Redis connection included. 
    *   Project runnable with basic funcionality. 

## TODO
Notice that this TODO list is not ordered by any factor (estimated finish time, importance, difficulty, etc.) and is not guaranteed to be implemented either:
*   ...
*   Choice of force-fetching latest data from RMP website instead of using cached data. 
*   TCP/Redis Connect pool. 
*   Better handling of concurrent queries on a same professor. 
*   Rundimental Frontend implementation, potentially Vue.js or React.js (JavaScript). 
*   Usage of Goroutine... Multi-everything! 
*   Air in production mode, potentially on AWS or TencentCloud. 