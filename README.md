# AHSAR Web
Development In Progress
- [AHSAR Web](#ahsar-web)
  - [Introduction](#introduction)
  - [Notice on Server Status](#notice-on-server-status)
  - [Application Structure](#application-structure)
  - [Application Setup (Locally)](#application-setup-locally)
  - [Public API (To be Expanded)](#public-api-to-be-expanded)
  - [Application Workflow](#application-workflow)
  - [A Bit More About NLP Server](#a-bit-more-about-nlp-server)
  - [License](#license)
  - [Project History](#project-history)
  - [TODO](#todo)

## Introduction
Web Application for AHSAR

Temporary Website under Production Mode on AWS: http://54.251.197.0:5000/

AHSAR intends to provide students with a different perspective on quatitatively evaluating their professors. Check the sentiment analysis result (in scale of 1 to 5) of other students' commentary on the professor of your choice by entering the `tid`(in AHSAR, it is called PID) in a professor URL from RateMyProfessors.com website, or the name (preferably full name) of professor. 

For example, assume a "randomly-selected" professor URL is `https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994` (Salute to Professor Adam Meyers from NYU, CS0002 ICP and CS480 NLP): 
*   select __Search by pid__, and enter `2105994`. 
*   select __Search by name__, and enter `Adam Meyers`; then choose the entry of `Adam Meyers 2105994 New York University ...`

Sentiment Score (continuous) and Sentiment Score (discrete) are usually:
*   close to each other in numbers, but sometimes they differ a lot... _be cautious when it happens_. 
*   can be undeterministic in different queries, because of randomized tie-breaking. 

Notice that Sentiment Score (discrete) is computed based on the positive and negative weight of individual comments, while Sentiment Score (continuous) is computed based on the positive and negative weight of all comments.
In other words, the higher the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive. 

For example, a professor having Sentiment Score (continuous) of 2.0 and Sentiment Score (discrete) of 4.0 might imply that more individual comments are classified as positive, but maybe the positive comments have really close weight on positivity and negativity, while the negative comments significantly skew to negativity. Why? Check the actual comments to find the reason. _Maybe the professor gives easy A, but the course is bad in many ways... or the other way around..._
 
[Back to top](#ahsar_web)

## Notice on Server Status
Server might be lagging, on and off, or unstable, because:
*   ... 
*   Server is under manual deployment. This usually takes only a few minutes. 
*   Server is down because of flood requests, concurrency test, etc. Fix time is not guaranteed, but auto-restart is on the development schedule. 
*   Server is down internally on AWS side, because of CPU/Memory shortage (this is a free version of EC2 Server, so the hardware is weak). Fix time is not guaranteed.

[Back to top](#ahsar_web)

## Application Structure
*   Frontend (Language: JavaScript, Framework: React.js) 
*   Backend
    * Backend Server (Language: Go, Framework: Gin) 
    * Query & Result Cache (MiddleWare: Redis, written in C) 
    * Internal Communication: Naive Socket TCP Connection
    * Server: Modified Version of AHSAR NLP project (Language: Python) 
*   Language Environment Setup: 
    *   C, Go, Python. 
    *   Check source file in this repository to find the actual Packages/Modules involved. 
    *   Python: `pip install nltk` and `pip install bs4` should be enough
    *   Go: in `/backend` directory, `go mod download` should be enough
    *   C: should not need any extra package
    *   Full list of preparation procedure will be listed when air with Continuous Deployment on AWS. 

[Back to top](#ahsar_web)

## Application Setup (Locally)
*   Start Redis (if compiled from source, in redis directory `/src/redis-server`; otherwise, `sudo service redis-server`)
*   Start NLP Server (in repository `/pysrc` directory, `python3 NLP_server.py`)
*   Start Backend Server (in repository `/backend` directory, `./app`, or, to recompile again, `bash run.bash`)
*   Start Frontend (in repository `/frontend` directory, for development mode `npm start`, or, for production mode first `npm run build` then follow the instruction on terminal)

[Back to top](#ahsar_web)

## Public API (To be Expanded)
*   ...
*   `GET http://54.251.197.0:5000/get_prof_by_id?input={numeric}&noCache={bool}`:
    *   Get sentiment analysis result by PID
    *   parameters: `input` should be numeric PID and is __*mandatory*__, `noCache` means "do not use cached result but fetch latest data" and is __*optional*__. 
    *   sample usage: `/get_pid_by_name?input=2105994&noCache=true` --> returns sentiment analysis result of newly fetched professor data of PID 2105994 on RMP website
*   `GET http://54.251.197.0:5000/get_pid_by_name?input={alphabetic}&noCache={bool}`:
    *   Get list of Professor Info and PID by Name
    *   parameters: `input` should be alphabetic professor name and is __*mandatory*__, `noCache` means "do not use cached result but fetch latest data" and is __*optional*__. 
    *   sample usage: `/get_pid_by_name?input=adam%20meyers&noCache=true` --> returns newly fetched professor entries on RMP website with name like "adam meyers", and user can choose from the entries

[Back to top](#ahsar_web)

## Application Workflow 
*   Frontend sends query to Backend
*   Backend receives query and check if it is in Redis:
    *   If in, retrieves the cached result and returns it to Frontend
    *   If not in, sends the query to NLP Server through TCP Socket
*   NLP Server receives the query from TCP Socket and start analyzing
*   NLP Server returns the result to Backend Server through TCP Socket
*   Backend updates Redis with the result and also returns result to Frontend

[Back to top](#ahsar_web)

## A Bit More About NLP Server
For the full project (including datebase of __80k labeled RMP comments__ and other imported data, codebase of __RMP scraper__ and __N-gram algorithm__, and __reference__ list for the imported data) of the NLP Server behind the screen, called __AHSAR__ *Ad-Hoc Sentiment Analysis on RateMyProfessors*, please check this [GitHub Repository](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP). Bear with the badly optimized code ^^. 

[Back to top](#ahsar_web)

## License
Project under MIT License. Basically, feel free to adopt anything (codebase, database, reference list, paper, etc. ) from here for any usage, with no warranty, promise, or liability from the repository owners and collaborators. But a little bit of credit/reference is very appreciated. 

[Back to top](#ahsar_web)

## Project History
*   ...

*   2021/08/24:
    *   Backend:
        *   Redis usage rework: 
            *   Simulated Mutex when concurrent queries are on the same PID or name
            *   Notice that Go-gin Framework automatically create goroutine for each individual incoming query from the frontend (and actually every API calling)
            *   Only one goroutine will acquire the mutex, get new data from NLP Server, and update  Redis with the data
            *   The rest of goroutines will detect that the mutex has been released, and then fetch cached data from Redis
            *   Better handling of concurrent queries as the same query will only need one process in NLP Server, and the other processes can work on queries with different PID or name. 
        *   Change to go-redis/v8
        *   Usage of randomized TTL in range of 1 to 6 hours instead of fixed 6 hours, so that it is much less likely to happen when majority of cached data suddenly go expired at the same time and Server gets flooded with queries on those cached data. 
    *   Frontend:
        *   Options of __Search by PID__ and __Search by Name__
        *   Display adaption for portable devices

*   2021/08/23:
    *   Backend: 
        *   API rework:  
            *   Search by PID API `/get_prof_by_id`:
                *   parameters: `input` should be numeric PID and is __*mandatory*__, `noCache` means "do not use cached result but fetch latest data" and is __*optional*__. 
                *   sample usage: `/get_pid_by_name?input=2105994&noCache=true` --> returns sentiment analysis result of newly fetched professor data of PID 2105994 on RMP website
            *   Search by Name API `/get_pid_by_name`:
                *   parameters: `input` should be alphabetic professor name and is __*mandatory*__, `noCache` means "do not use cached result but fetch latest data" and is __*optional*__. 
                *   sample usage: `/get_pid_by_name?input=adam%20meyers&noCache=true` --> returns newly fetched professor entries on RMP website with name like "adam meyers", and user can choose from the entries
                *   Frontend part for this API has not done yet
        *   Split API handler functions into different files for better modulability. 
    *   Frontend:
        *   Use localStorage to store the latest 10 query history

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

[Back to top](#ahsar_web)

## TODO
Notice that this TODO list is not ordered by any factor (estimated finish time, importance, difficulty, etc.) and is not guaranteed to be implemented either:
*   ...
*   Optimization, Modularization, Robustness...
*   Search by school / department, Hey, why doesn't RMP website provide this feature? 
*   Concurrent Scraper written in Go and split up the python scraper and NLP analyzer. 
*   Server auto-recovering from fatal error of NLP or Redis processes. 
*   Login / Registration (don't see a reason for that at the moment). 
*   Allow user to submit a paragraph of commentary and obtain sentiment analysis result. 
*   More features / functions / tools...
*   TCP/Redis Connection pool. 
*   Usage of Goroutine... Multi-everything! But where to use it? 
*   Human-readable domain address: maybe `www.ahsar.*` is a good name. 
*   Air with Continuous Deployment on AWS. 
*   Auto-restart AWS Server when server is down because of internal issue: resource shortage, flood attack, etc. 

[Back to top](#ahsar_web)