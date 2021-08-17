# AHSAR_Web
Development In Progress

## Introduction
Web Application for AHSAR
Check the sentiment analysis result of student's commentary on the professor of your choice by entering the `tid` in a professor URL from RateMyProfessors.com website. 

For example, assume a "randomly-selected" professor URL is `https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994` (Salute to Professor Adam Meyers from NYU, CS0002 ICP and CS480 NLP), enter `2105994` but not the name of professor or the full URL. 

Notice that Sentiment Score (discrete) is computed based on individual comments, while Sentiment Score (continuous) is computed based on all comments.
In other words, the higher the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive.

## Application Structure (Not Finished)
* Frontend (Undecided, potentially JavaScript Vue.js framework) 
* Backend Server (Language: Go, Framework: Gin) 
* Query Cache (MiddleWare: Redis, written in C) 
* Message Queue (MiddleWare: Kafka, written in Java) 
* NLP Process: Modified Version of AHSAR NLP project (Language: Python) 
* Language Environment Setup: C, Go, Python, Java. Check source file in this repository to find the actual Packages/Modules involved. Full list of dependencies will be listed when air in production mode. 

## Application Setup (Not Finished)
* Prepare Query Cache and Message Queue beforehand...!
* Start Query Cache and Message Queue
* Start NLP Process
* Start Backend Server
* Start Frontend / Access directly on browser

## Application Workflow (Not Finished)
* Frontend sends query to Backend
* Backend receives query and check if it is in Query Cache:
    * If in, retrieves the cached result and returns it to Frontend
    * If not in, sends the query to NLP Process through Message Queue
* NLP Process receives the query from Message Queue and start analyzing
* Update the Query Cache with the analyzed result and notifies Backend Server to retrieve result from Query Cache
* Backend retrieves data from Query Cache and returns it to Frontend

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
    *   TODO in mind (not ordered by any factor):
        *   Rundimental Frontend implementation, potentially Vue.js(JavaScript). 
        *   Simple Kafka connection. 
        *   NLP Process continuous modification, potentially on multi-processing / process pool / MPI. 
        *   Usage of Goroutine...Multi-everything! 
        *   Air in production mode, potentially on AWS or TencentCloud. 