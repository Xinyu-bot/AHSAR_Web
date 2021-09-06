# AHSAR Web
Development In Progress
正在开发中
- [AHSAR Web](#ahsar-web)
  - [Introduction 介绍](#introduction-介绍)
  - [Notice on Server Status 服务器状态](#notice-on-server-status-服务器状态)
  - [Application Structure 服务器架构](#application-structure-服务器架构)
  - [Application Setup (Locally) 本地如何运行](#application-setup-locally-本地如何运行)
  - [Public API (To be Expanded) 公开调用的API](#public-api-to-be-expanded-公开调用的api)
  - [A Bit More About NLP Server 关于NLP服务器](#a-bit-more-about-nlp-server-关于nlp服务器)
  - [License 授权证书](#license-授权证书)
  - [Project History 项目历史](#project-history-项目历史)
  - [TODO 待完成工作计划](#todo-待完成工作计划)

## Introduction 介绍
Web Application for AHSAR

Temporary Website at Tencent Cloud: http://1.14.137.215:5000/

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

AHSAR网页应用

运行在腾讯云服务器上的临时网址：http://1.14.137.215/:5000

AHSAR旨在为在美大学生提供一个量化评估教授的新角度。通过输入RateMyProfessors.com网站上教授URL末位的`tid`或者教授名字（最好是全名），用户就可以得到其他学生对于该教授的评语的情绪分析结果。得分从低到高是1到5分。

比如，纽约大学计算机系CS0002 ICP和CS480 NLP课程的授课教授Adam Meyers，他在RMP网站上的URL为`https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994`，那么可以通过以下方式进行查询：
*   选择Search by pid，然后输入2105994
*   选择Search by name，然后输入教授全名Adam Meyers，最后选择列表中`Adam Meyers 2105994 New York University ...`这个选项

情绪分析的得分，分为“连续”和“离散”两种。一般来说：
*   这两个分数会比较接近。但是少数情况两个分数会差距较大，这个时候就要多留一个心眼儿了
*   多次查询的结果可能会有一些不一致，这是因为分析数据的时候，如果出现正、负的概率相同，程序会随机决定该句子的正负性 (randomized tie-breaking)。

“离散”的情绪得分是基于单个评论的正负权重而计算的，“连续”的得分是基于所有评论汇总的正负权重而计算的。话句话说，离散分数越高，说明更多的单条评论是偏向正面的；而连续分数越高，则说明全部评论整体而言，偏向正面的比例更高。

例如，一位教授的连续得分为2.0分而离散的分为4.0分，可能意味着更多的单个评论被归类为正面的，但也许在这些正面的评论中正负的占比非常接近，而在被归类为负面的评论中，负向的占比非常高。为什么呢？这就需要真正花时间去看一下实际的评论了。也许这个教授会给很多A，但是课程从各方面而言都很糟糕...或者正好相反...

[Back to top 回到顶部](#ahsar-web)

## Notice on Server Status 服务器状态
Server might be lagging, on and off, or unstable, because:
*   ... 
*   Server is under manual deployment. This usually takes only a few minutes. 
*   Server is temporarily down because of flood requests, concurrency test, etc. Fix time is not guaranteed, but auto-recover is on the development schedule. 
*   Server is permanently down because of internal issue on Tencent Cloud side, maybe due to CPU/Memory shortage. Fix time is not guaranteed, but auto-restart is on the development schedule. 

服务器可能会卡顿、时好时坏、或者不稳定，因为以下几种原因：
*   正在手动部署最新版本的代码（自动部署功能在计划中，但不确定完成时间），这一般只需要几分钟时间。
*   服务器因为大量的访问或者并发测试而短暂挂起。修复时间不确定，但是自动恢复的功能也在计划中。
*   服务器因为腾讯云服务方面的问题而宕机，比如内存或者CPU不够用了。修复时间不确定，但是自动重启的功能也在计划中。

[Back to top 回到顶部](#ahsar-web)

## Application Structure 服务器架构
*   Frontend (Language: JavaScript, Framework: React.js) 
*   Backend
    *   HTTP Response Server, or in the project, simply called Backend Server (Language: Go, Framework: Gin) 
    *   Query & Result Cache (MiddleWare: Redis, written in C) 
    *   Partial Mirror of RMP Database (Relational Database: MySQL, written in C)
    *   NLP Server: Modified Version of AHSAR NLP project (Language: Python) 
    *   Internal Communication between Backend Server and NLP Server: Naive Socket TCP Connection
*   Language Environment Setup: 
    *   C, Go, Python, JavaScript (please choose from the most-recent/latest versions). 
    *   Check source file in this repository to find the actual Packages/Modules involved, but basically:
    *   JavaScript: in `/frontend` directory, `npm install`
    *   Python: `pip install nltk`, `pip install bs4`, and `pip install rake-nltk` should be enough
    *   Go: in `/backend` directory, `go mod download` should be enough
    *   C: should not need any extra package
*   ...
*   前端 （JavaScript语言的React.js框架）
*   后端
    *   HTTP响应服务器（下称后端服务器，Go语言的Gin框架）
    *   查询请求和响应结果的缓存（用C编写的中间件Redis）
    *   RMP网站的部分镜像数据（用C编写的关系型数据库MySQL）
    *   NLP处理服务器（下称NLP服务器，Python语言）
    *   后端服务器和NLP服务器的内部通信方式：简易的Socket TCP连接
*   语言环境
    *   C, Go, Python, JavaScript (请选用最新的几个版本之一)
    *   可以阅读代码仓库中的源码并且找到真正需要下载的包、模组，不过一般来说：
    *   JavaScript: 需要在`/frontend`目录下执行`npm install`
    *   Python: 需要执行`pip install nltk`, `pip install bs4`, 和`pip install rake-nltk`
    *   Go: 需要在`/backend`目录下执行`go mod download`自动下载需要的包
    *   C: 不需要额外的包

[Back to top 回到顶部](#ahsar-web)

## Application Setup (Locally) 本地如何运行
*   Start Redis `redis-server`
*   Start MySQL `service MySQL start`
*   Start NLP Server (in repository `/pysrc` directory, `python3 NLP_server.py`)
*   Start Backend Server (in repository `/backend` directory, `./app`, or, to recompile again, `bash run.bash`)
*   Start Frontend (in repository `/frontend` directory, for development mode `npm start`, or, for production mode first `npm run build` then follow the instruction on terminal)
*   ...
*   启动Redis `redis-server`
*   启动MySQL `service MySQL start`
*   启动NLP服务器（在`/pysrc`目录下执行`python3 NLP_Server.py`）
*   启动后端服务器（在`/backend`目录下执行`./app`，或者如果想重新编译的话，执行`bash run.bash`）
*   启动前端（在`/frontend`目录下执行`npm start`来开启开发模式的前端进程，或者如果想开启生产模式的进程，执行`npm run build`再根据命令行中显示的指引进行下一步操作）

[Back to top 回到顶部](#ahsar-web)

## Public API (To be Expanded) 公开调用的API
*   ...
*   `GET http://1.14.137.215:8080/get_pid_by_name?input=2105994`
*   `GET http://1.14.137.215:8080/get_pid_by_name?input=adam%20meyers`
*   `GET http://1.14.137.215:8080/get_schools_by_initial?initial=N`
*   `GET http://1.14.137.215:8080/get_departments_by_school?school=New%20York%20University`
*   `GET http://1.14.137.215:8080/get_prof_by_department?school=New%20York%20University&department=Computer%20Science`

*   ...
*   不翻译了，懂的都懂

[Back to top 回到顶部](#ahsar-web)

## A Bit More About NLP Server 关于NLP服务器
For the full project (including datebase of __80k labeled RMP comments__ and other imported data, codebase of __RMP scraper__ and __N-gram algorithm__, and __reference__ list for the imported data) of the NLP Server behind the screen, called __AHSAR__ *Ad-Hoc Sentiment Analysis on RateMyProfessors*, please check this [GitHub Repository](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP). Bear with the badly optimized code ^^. 

Notice that Keywords Extraction (experiment feature) is implemented by `rake-nltk` module written by GitHub user *csurfer*, which also can be found at [this GitHub repository](https://github.com/csurfer/rake-nltk)


关于NLP服务器背后的AHSAR项目（包括8万条标记过的RMP评论和其他引用数据的数据库，RMP爬虫和N-gram算法的代码库，以及引用数据的引用参考列表），我们称之为 __AHSAR__ _Ad-Hoc Sentiment Analysis on RateMyProfessors_，欢迎访问[此GitHub代码仓库](https://github.com/Xinyu-bot/NLP_SentimentAnalysis_RMP)。代码优化质量较差，敬请谅解。

关键词提取功能是由GitHub用户*csurfer*的`rake-nltk`模组完成的。欢迎访问[他的GitHub代码仓库](https://github.com/csurfer/rake-nltk)。
[Back to top 回到顶部](#ahsar-web)

## License 授权证书
Project under MIT License. Basically, feel free to adopt anything (codebase, database, reference list, paper, etc. ) from here for any usage, with no warranty, promise, or liability from the repository owners and collaborators. But a little bit of credit/reference is very appreciated. 

项目授权MIT证书。只要不犯法，随便玩。但是本代码仓库的拥有者、管理员、贡献者不对项目的内容作出任何拥有法律效力的保证和担保。如果能在您的项目中提到我们，不胜感激。

[Back to top 回到顶部](#ahsar-web)

## Project History 项目历史
*   ... 只做简单的翻译
*   2021/09/06:
    *   Backend:
        *   Use MySQL for SDP -- decoupling: 
            *   remove SDP-related code from NLP Server; also larger pool size of 10 since the removal of SDP frees some CPU and RAM resources
            *   remove Redis involvement from SDP, since MySQL connection pool and indexed query make the response time (< 5ms) short enough (as for now)
            *   add SDP-relatede code in Backend Server
        *   Keywords Extraction Feature:
            *   using `rake-nltk` module
            *   provides 10 keywords with top 10 "significant value" to the entire list of commentary of a specific professor
        *   使用MySQL提供SDP功能的支持 —— 系统解耦：
            *   从NLP服务器中移除了SDP相关代码，并且由于移除之后空余了更多系统资源，进程池扩大到10
            *   从Redis缓存中移除SDP相关代码，因为MySQL连接池和索引查询已经让响应时间(< 5ms)足够快（目前来看）
            *   SDP相关代码已经被移到后端HTTP响应服务器
        *   关键词提取功能：
            *   使用了`rake-nltk`模组
            *   从某个教授的所有评论中提取出10个最具价值的关键词组
*   2021/09/05:
    *   Deployment has been moved to Tencent Cloud at http://1.14.137.215:5000/
    *   SDP model updated with difficulty score, quality score, and would take again percentage information
    *   Professor list by searching school and department is now returned sorted by quality score then difficulty score
    *   迁移至腾讯云新地址：http://1.14.137.215:5000/
    *   全网站重爬，教授列表将以quality和difficulty两个分数进行预排序。
*   2021/09/04: 
    *   Backend:
        *   Implementation of Searching by School and Department:
            *   SDP scraper using process pool nested with thread pool
            *   Fully structured School-Department-Professor model obtained brutally from RMP website and dumped into byte file using `Python3 pickle`
            *   Three new APIs: `/get_schools_by_initial`, `/get_prof_by_department`, and `/get_departments_by_school` 
            *   related updates on involved files
            *   fix a bug where the keys in Redis could potentially be the same for differnt types of query
        *   新功能的API：根据校名+系名搜索任职教授列表
            *   进程池嵌套线程池的爬虫，用于获取SDP模型。模型通过`Python3 pickle`保存为二进制文件
            *   三个新的API，如英文部分所提及
    *   Frontend:
        *   rework and optimization
        *   home page picture
        *   in the future, frontend updates will be explained in a more detailed way
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
        *   简单来说，随机Redis缓存的TTL来减少缓存雪崩的危害，空值储存来减少缓存穿透的危害，互斥锁来解决缓存击穿（并发的、针对同一个数据的请求）的问题
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
        *   Elastic IP address: ~~http://54.251.197.0:5000/~~ __*No Longer In-Use*__
        *   Using CORS to connect Backend Server and Frontend
        *   Reduced the NLP multiprocessing pool size from 20 to ~~5~~ 3 in the deployed version to save memory usage on the free AWS EC2 ubuntu server...
        *   Notice that in a local environment, the pool size should be significantly larger to maximize the ability of hanlding concurrent requests. 
        *   加入了跨域资源共享的中间件，所以在亚马逊云服务器上可以使用生产模式的React前端进程了

*   2021/08/20:
    *   Manual Deployment on AWS EC2 Ubuntu server: 
        *   Programs running by `screen` command so accessible 24/7. 
        *   Now available at public IP address:Port ~~http://18.142.108.23:8080/~~ __*No Longer In-Use*__
        *   Frontend is unimplemented yet, currently support query with URL only
        *   Example: ~~http://18.142.108.23:8080/get_prof_by_id?input=123456~~ where `123456` is the PID. __*No Longer In-Use*__
        *   服务器24小时全天候运行

*   2021/08/19:
    *   Redis cache set with expiration limit. 
    *   AWS EC2 Ubuntu server setup. 

*   2021/08/18:
    *   Backend Server now communications with NLP Server through Naive Socket TCP Connection. 
    *   New NLP Server with Naive Socket TCP Connection and multi-processing pool. 
    *   Removal of Kafka dependency. 
    *   NLP服务器中使用进程池提升并发效率

*   2021/08/17: 
    *   First push to GitHub Repository. 
    *   Simple Redis connection included. 
    *   Project runnable with basic funcionality. 

[Back to top 回到顶部](#ahsar-web)

## TODO 待完成工作计划
Notice that this TODO list is not ordered by any factor (estimated finish time, importance, difficulty, etc.) and is not guaranteed to be implemented either:
*   ...
*   Optimization, Modularization, Robustness...
*   More features, functions, tools...
*   Continuous and automatic update SDP model slowly
*   Server auto-recovering from fatal error of NLP or Redis processes. 
*   Allow user to submit a paragraph of commentary and obtain sentiment analysis result. 
*   TCP/Redis Connection pool. 
*   Usage of Goroutine... Multi-everything! But where to use it? 
*   Human-readable domain address: maybe `www.ahsar.*` is a good name. 
*   Auto-restart Tencent Cloud Server when server is down because of internal issue: resource shortage, flood attack, etc. 


待完成事项（排序不分先后且仅供参考，完成时间无保证）
*   更多功能
*   代码优化、模块化、提升健壮性
*   慢慢地、持续地、自动化更新SDP模型
*   服务器自动恢复，应对TCP或者Redis连接断开的情况
*   用户可以上传一段评论，服务器处理之后返回该评论的情绪分析结果
*   TCP/Redis连接使用连接池（现在是用后即弃）
*   更多Goroutine的应用，并发搞起来
*   更好记的域名，比如`www.ahsar.*`
*   持续部署
*   服务器内部错误后的自动重启


[Back to top 回到顶部](#ahsar-web)