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
    * Query Cache (MiddleWare: Redis, written in C) 
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
* Backend updates Redis
* Backend retrieves data from Redis and returns it to Frontend

## NLP Algorithm Process
For the full implementation and reference of the NLP Algorithm Process behind, called __AHSAR__ *Ad-Hoc Sentiment Analysis on RateMyProfessors*, please check this [Repository](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP). Bear with the badly optimized code ^^. 

## License
Project under MIT License. Basically, feel free to adopt any code from here for any usage, with no warranty, promise, or liability from the repository owners and collaborators. But a little bit of credit/reference is very appreciated. 

## History and TODOs:
*   2021/08/17: 
    *   Done:
        *   First push to GitHub Repository. 
        *   Simple Redis connection included. 
        *   Project runnable with basic funcionality. 
*   2021/08/18:
    *   Done:
        *   Backend Server now communications with NLP Server through Naive Socket TCP Connection. 
        *   New NLP Server with Naive Socket TCP Connection and multi-processing pool. 
        *   Removal of Kafka dependency. 
<br></br>
*   TODO in mind (not ordered by any factor):
    *   ...
    *   Redis cache set with expiration limit, 
    *   TCP Socket pool, maybe 
    *   Better handling of concurrency: 10 concurrent request of same professor should only launch the scraper and calculation process once and the rest should read from Redis, not ten times. 
    *   Rundimental Frontend implementation, potentially Vue.js or React.js (JavaScript). 
    *   Usage of Goroutine...Multi-everything! 
    *   Air in production mode, potentially on AWS or TencentCloud. 