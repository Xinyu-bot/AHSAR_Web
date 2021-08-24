import React from 'react'
import './Result.css'
import SearchedList from '../../../home/header/searched_list/SearchedList.js'
export default function Result(props) {
	return (
		<div className='result'>
			{(() => {
				if (props.ret !== undefined) {
					// evil switch case here, Oops

					if (props.ret === 'loading') {
						return (
							<div>
								<h2 id='result'>Result for {props.pid}</h2>
								<span>Loading data...</span>
								<span className='loader'></span>
							</div>
						)
					} else if (props.ret === 'non-numbers') {
						return <div></div>
					}
					//如果用户输入内容符合条件，并且收到了服务器的返回内容。
					switch (props.ready) {
						// user changes input or there is no input at all
						case 0:
							return
						// response code not 200, ERROR in backend
						case -1:
							return (
								<span>
									<p>Something is wrong with the server, please try again later.</p>
								</span>
							)
						// all good
						case 1:
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Professor: {props.ret.professor_name}</p>
									<p>PID: {props.ret.pid}</p>
									<p>Difficulty Score: {props.ret.difficulty_score}</p>
									<p>Quality Score: {props.ret.quality_score}</p>
									<p>Would Take Again: {props.ret.would_take_again}</p>
									<p>Sentiment Analysis Score (Discrete): {props.ret.sentiment_score_discrete}</p>
									<p>Sentiment Analysis Score (Continuous): {props.ret.sentiment_score_continuous}</p>
								</span>
							)
						// prof has no comment, so no sentiment analysis score, but the professor has quality score and difficulty score
						//pid存在，没有评论, diffuculty_score，quality_score都有
						case 2:
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Professor: {props.ret.professor_name}</p>
									<p>Difficulty Score: {props.ret.difficulty_score}</p>
									<p>Quality Score: {props.ret.quality_score}</p>
									<p>Would Take Again: {props.ret.would_take_again}</p>
									<p>(This professor has no commentary from any student yet!)</p>
									<p>Sentiment Analysis Score (Discrete): N/A</p>
									<p>Sentiment Analysis Score (Continuous): N/A</p>
								</span>
							)
						// PID does not exist
						case 3:
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Sorry, this PID is invalid. There is no such PID on RateMyProfessors.com!</p>
								</span>
							)

						//pid exist, but the professor has no difficulty_score and no quality score
						case 4:
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Professor: {props.ret.professor_name}</p>
									<p>(This professor has no difficulty score and no quality score yet!)</p>
									<p>Difficulty Score: N/A</p>
									<p>Quality Score: N/A</p>
									<p>Would Take Again: N/A</p>
									<p>(This professor has no comment yet!)</p>
									<p>Sentiment Analysis Score (Discrete): N/A</p>
									<p>Sentiment Analysis Score (Continuous): N/A</p>
								</span>
							)
						case 5: //pid存在，有评论，没有difficulty_score，没有quality_score
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Professor: {props.ret.professor_name}</p>
									<p>(This professor has no difficulty score and no quality score yet!)</p>
									<p>Difficulty Score: N/A</p>
									<p>Quality Score: N/A</p>
									<p>Would Take Again: {props.ret.would_take_again}</p>
									<p>Sentiment Analysis Score (Discrete): {props.ret.sentiment_score_discrete}</p>
									<p>Sentiment Analysis Score (Continuous): {props.ret.sentiment_score_continuous}</p>
								</span>
							)

						//for search by name
						case 6:
							return (
								<span>
									<h2 id='result'>Result for </h2>
									<p>Sorry, this name is invalid. There is no such name on RateMyProfessors.com!</p>
								</span>
							)
						case 7:
							return (
								<span>
									<h2 id='result'>Result for </h2>

									<SearchedList searchedListByName={props.ret} />
								</span>
							)
						// make compiler happy
						default:
							return
					}
				}
			})()}
		</div>
	)
}
