import React from 'react'
import PropTypes from 'prop-types'
import './Article.css'
import Result from './result/Result'

Article.prototype = {
	getSearchedPid: PropTypes.func,
	//pid can be string and number 点击搜索历史li和searchbyname的li时传来的是number，其他时候是string
	//ret can be string and object
	ready: PropTypes.number.isRequired,
}
function Article(props) {
	const getClickedPid = (pid) => {
		console.log('here article', pid)
		props.getSearchedPid(pid)
	}
	return (
		<div className='article'>
			{/*jsx里面只能放js表达式，包括执行一个函数。如果pid不为空，显示下面的节点 */}
			{(() => {
				if (props.pid !== '') {
					return (
						<div>
							<Result getClickedPid={getClickedPid} ret={props.ret} pid={props.pid} ready={props.ready} />
						</div>
					)
				}
			})()}

			{/* Notice that this is commented out because -
            - there is no need to seperate the main function from the main page as for now
            <Link to = '/get_prof_by_id'>
                <p>Search Professor by ID</p>
            </Link>
            */}
			<section id='main'>
				<span>
					Visit <a href='https://github.com/Xinyu-bot/AHSAR_Web'>this GitHub repository</a> for the full project: Codebase, Database, Reference List, etc.
				</span>
				<p id='searchByID_Intro'>
				AHSAR intends to provide students with a different perspective on quatitatively evaluating their professors. 
				Check the sentiment analysis result (in scale of 1 to 5) of other students' commentary on the professor of your choice right away, 
				by entering the tid (in AHSAR, it is called PID) inside a professor URL from RateMyProfessors.com website, 
				or the name (preferably full name) of a professor. 
				</p>

				<p id='searchByID_Intro'>
					For example, if the professor's name is Adam Meyers, and his URL on RateMyProfessors.com is <br />
					<span>https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994</span>
					<br />
					Then, select Search by pid, and enter "2105994"; <br />
					Or, select Search by name, enter Adam Meyers, and select the entry of "Adam Meyers 2105994 New York University ..."
				</p>

				<p id='searchByID_Intro'>
					Sentiment Score (continuous) and Sentiment Score (discrete) are usually:
					<br /> 1. close to each other in numbers, but sometimes they differ a lot... be cautious when it happens. 
					<br /> 2. can be undeterministic in different queries, because of randomized tie-breaking. 
					
					<br />Notice that Sentiment Score (discrete) is computed based on individual comments, while Sentiment Score (continuous) is computed based on all comments. In other words, the higher
					the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive. The sentiment score can be
					undeterministic, because of random tie-breaking. 
					<br />
					For example, a professor having Sentiment Score (continuous) of 2.0 and Sentiment Score (discrete) of 4.0 might imply that more individual comments are classified as positive, 
					but maybe the positive comments have really close weight on positivity and negativity, while the negative comments significantly skew to negativity. 
					Why? Check the actual comments to find the reason. Maybe the professor gives easy A, but the course is bad in many ways... or the other way around...
				</p>
			</section>
		</div>
	)
}
export default Article
