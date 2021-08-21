import React, {useState} from 'react'
import axios from 'axios'
import './home.css'

function isNum(s) {
    if (s !== null && s !== '') {
        return !isNaN(s);
    }
    return false;
}

function Home() {
    // frontend input value
    const [pid, setPID] = useState('') // user input
    // backend returned value
    const [ready, setReady] = useState(0) // ready flag
    const [ret, setRet] = useState()
    
    // dynamically record the user input
    // when user changes the input text, hide the result from last input
    const _setPID = (val) => {
        if (ready) {
            setReady(0)
        }
        setPID(val)
    }

    /* 
    set ready flag for three cases:
        1: all good
        2: professor has no comment, so there is no sentiment analysis score
        3: PID does not exist
    */
    const _setRet = (data) => {
        if (data.first_name === '-1') {
            setReady(3)
        } 
        else if (data.sentiment_score_discrete === '-1') {
            setReady(2)
        }
        else {
            setReady(1)
        }
        setRet(data)
    }

    // handler of ENTER event
    const searchEnter = (e) => {
        // input validation
        if (pid !== '' && isNum(pid)) {
            // retrieve from backend API /get_prof_by_id
            axios.get('http://54.251.197.0:8080/get_prof_by_id', {
                params: {
                    input: pid,
                },
            })
            // process returned result
            .then((r) => {
                if (r.status === 200) {
                    console.log(r.data)
                    _setRet(r.data)
                } else {
                    setReady(-1)
                }
            })
            // catch error
            .catch((err) => {
                console.error(err)    
            })    
        } 
        // empty input
        else if (pid === '') {
            alert('You need to enter something')
            e.preventDefault()
        }
        // non-numbers (specifically non-positive integer) input
        else {
            alert('PID should only contain number 0 ~ 9')
            e.preventDefault()
        }
    }

    return (
        <div className = "Home">
            <h2>Hello, Welcome to AHSAR Website! Development in Progress...</h2>
            <span>
                Visit <a href = "https://github.com/Xinyu-bot/AHSAR_Web">this GitHub repository</a> for the full project: Codebase, Database, Reference List, etc. 
            </span>

            {/* Notice that this is commented out because -
            - there is no need to seperate the main function from the main page as for now
            <Link to = '/get_prof_by_id'>
                <p>Search Professor by ID</p>
            </Link>
            */}
            <section id = "main">
                <p id = "searchByID_Intro">
                    By entering the pid of a professor of your choice, 
                    get the Sentiment Analysis result of students' commentary 
                    from RateMyProfessors.com on the professor right away!
                </p>

                <p id = "searchByID_Intro">
                    For example, if the professor's URL on RateMyProfessors.com is <br/>
                    <span>https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994</span><br/>
                    Then, enter 2105994 in the text box below.
                </p>

                <p id = "searchByID_Intro">
                    Notice that Sentiment Score (discrete) is computed based on individual comments, 
                    while Sentiment Score (continuous) is computed based on all comments. 
                    In other words, the higher the discrete score is, the more individual comments are positive. 
                    The higher the continuous score is, the larger proportion of all comments are positive.
                    The sentiment score can be undeterministic, because of random tier-breaking.
                </p>
            </section>
            <input id="search-content" placeholder = "Please enter a PID" onInput={(e) => _setPID(e.target.value)} />
            <button onClick={(e) => {searchEnter(e)}}>Search</button>
            
            <div id = "result">
                {(() => {
                    if (ret !== undefined ) {
                        // evil switch case here, Oops
                        switch(ready) {
                                                    // user changes input or there is no input at all
                        case 0: return
                        // response code not 200, ERROR in backend
                        case -1: return (
                            <span>
                                <p>Something is wrong with the server, please try again later.</p>
                            </span>
                        )
                        // all good
                        case 1: return (
                            <span>
                                <p>Professor: {ret.first_name} {ret.last_name}</p>
                                <p>Difficulty Score: {ret.difficulty_score}</p>
                                <p>Quality Score: {ret.quality_score}</p>
                                <p>Sentiment Aanlysis Score (Discrete): {ret.sentiment_score_discrete}</p>
                                <p>Sentiment Analysis Score (Continuous): {ret.sentiment_score_continuous}</p>
                            </span>
                        )
                        // prof has no comment, so no sentiment analysis score
                        case 2: return (
                            <span>
                                <p>Professor: {ret.first_name} {ret.last_name}</p>
                                <p>Difficulty Score: {ret.difficulty_score}</p>
                                <p>Quality Score: {ret.quality_score}</p>
                                <p>This professor has no commentary from any student yet!</p>
                            </span>
                        )
                        // PID does not exist
                        case 3: return (
                            <span>
                                <p>Sorry, this PID is invalid. There is no such PID on RateMyProfessors.com!</p>
                            </span>
                        )
                        // make compiler happy
                        default: return
                        }
                    }
                })()}
            </div>
        </div>
    );
}

export default Home;