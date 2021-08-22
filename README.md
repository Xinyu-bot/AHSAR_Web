# AHSAR_Web
Development In Progress

## Introduction
Web Application for AHSAR

Temporary Website under Production Mode: http://54.251.197.0:5000/

Check the sentiment analysis result of student's commentary on the professor of your choice by entering the `tid` in a professor URL from RateMyProfessors.com website. 

For example, assume a "randomly-selected" professor URL is `https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994` (Salute to Professor Adam Meyers from NYU, CS0002 ICP and CS480 NLP), enter `2105994` but not the name of professor or the full URL. 

Notice that Sentiment Score (discrete) is computed based on individual comments, while Sentiment Score (continuous) is computed based on all comments.
In other words, the higher the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive.

## Notice on Server Status
Server might be lagging, on and off, or unstable, because:
*   Server is under manual deployment. This usually takes only a few minutes. 
*   Server is down because of flood requests, concurrency test, etc. Fix time is not guaranteed, but auto-restart is on the development schedule. 
*   Server is down internally on AWS side, because of CPU/Memory shortage (this is a free version of EC2 Server, so the hardware is weak). Fix time is not guaranteed.
*   ...

## Application Structure
* Frontend (Language: JavaScript, Framework: React.js) 
* Backend
    * Backend Server (Language: Go, Framework: Gin) 
    * Query & Result Cache (MiddleWare: Redis, written in C) 
    * Internal Communication: Naive Socket TCP Connection
    * Server: Modified Version of AHSAR NLP project (Language: Python) 
* Language Environment Setup: 
  * C, Go, Python. 
  * Check source file in this repository to find the actual Packages/Modules involved. 
    * Python: `pip install nltk` and `pip install bs4` should be enough
    * Go: in `/backend` directory, `go mod download` should be enough
    * C: should not need any extra package
  * Full list of preparation procedure will be listed when air with Continuous Deployment on AWS. 

## Application Setup (Locally)
* Start Redis (if compiled from source, in redis directory `/src/redis-server`; otherwise, `sudo service redis-server`)
* Start NLP Server (in repository `/pysrc` directory, `python3 NLP_server.py`)
* Start Backend Server (in repository `/backend` directory, `./app`, or, to recompile again, `bash run.bash`)
* Start Frontend (in repository `/frontend` directory, for development mode `npm start`, or, for production mode first `npm run build` then follow the instruction on terminal)

## Public API (To be Expanded)
* `GET http://54.251.197.0:5000/get_prof_by_id?input={pid}` where pid is the user input


## Application Workflow 
* Frontend sends query to Backend
* Backend receives query and check if it is in Redis:
    * If in, retrieves the cached result and returns it to Frontend
    * If not in, sends the query to NLP Server through TCP Socket
* NLP Server receives the query from TCP Socket and start analyzing
* NLP Server returns the result to Backend Server through TCP Socket
* Backend updates Redis with the result and also returns result to Frontend

## NLP Server
For the full project (including datebase of __80k labeled RMP comments__ and other imported data, codebase of __RMP scraper__ and __N-gram algorithm__, and __reference__ list for the imported data) of the NLP Server behind the screen, called __AHSAR__ *Ad-Hoc Sentiment Analysis on RateMyProfessors*, please check this [GitHub Repository](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP). Bear with the badly optimized code ^^. 

## License
Project under MIT License. Basically, feel free to adopt anything (codebase, database, reference list, paper, etc. ) from here for any usage, with no warranty, promise, or liability from the repository owners and collaborators. But a little bit of credit/reference is very appreciated. 

## Project History
*   ...
*   2021/08/22:
    *   Bug fix: 
        *   `professor name` and `would take again` is now fetched from RMP website and displayed on AHSAR Web correctly. 
    *   Frontend rework: 
        *   optimization on rendering process
        *   modularization - split home page file into smaller components
        *   css refinement
*   2021/08/21:
    *   Manual Deployment continued: 
        *   Elastic IP address: http://54.251.197.0:5000/ <-- notice that this `:5000` port means that the React App is in Production Mode
        *   Using CORS to connect Backend Server and Frontend
        *   Reduced the NLP multiprocessing pool size from 20 to ~~5~~ 3 in the deployed version to save memory usage on the free AWS EC2 ubuntu server...
        *   Notice that in a local environment, the pool size should be significantly larger to maximize the ability of hanlding concurrent requests. 
*   2021/08/20:
    *   Manual Deployment on AWS EC2 Ubuntu server: 
        *   Programs running by `screen` command so accessible 24/7. 
        *   Now available at public IP address:Port ~~http://18.142.108.23:8080/~~ __*No Longer In-Use*__
        *   Frontend is unimplemented yet, currently support query with URL only
        *   Example: ~~http://18.142.108.23:8080/get_prof_by_id?input=123456~~ where `123456` is the PID. __*No Longer In-Use*__
*   2021/08/19:
    *   Redis cache set with expiration limit. 
    *   AWS EC2 Ubuntu server setup. 
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
*   Optimization, Modularization, Robustness
*   Using an `.env` file to store the backend API, or maybe no needed because AWS security groop will take care of it
*   Choice of force-fetching latest data from RMP website instead of using cached data. 
*   Allow user to submit a paragraph of commentary and obtain sentiment analysis result. 
*   More features / functions / tools
*   TCP/Redis Connection pool. 
*   Better handling of concurrent queries on a same professor. 
*   Rundimental Frontend implementation, potentially Vue.js or React.js (JavaScript). 
*   Usage of Goroutine... Multi-everything! 
*   Human-readable domain address: maybe `www.ahsar.*`
*   Air with Continuous Deployment on AWS. 
*   Auto-restart AWS Server when server is down because of internal issue: resource shortage, flood attack, TCP/Redis connection failure, etc. 