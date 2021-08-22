import React from 'react'
import './Result.css'

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
						// prof has no comment, so no sentiment analysis score
						case 2:
							return (
								<span>
									<h2 id='result'>Result for {props.pid}</h2>
									<p>Professor: {props.ret.professor_name}</p>
									<p>Difficulty Score: {props.ret.difficulty_score}</p>
									<p>Quality Score: {props.ret.quality_score}</p>
									<p>This professor has no commentary from any student yet!</p>
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
						// make compiler happy
						default:
							return
					}
				}
			})()}
		</div>
	)
}
