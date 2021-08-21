import React from 'react'
import './Article.css'
function Article() {
	return (
		<div className='article'>
			<h2>Hello, Welcome to AHSAR Website! Development in Progress...</h2>
			<span>
				Visit <a href='https://github.com/Xinyu-bot/AHSAR_Web'>this GitHub repository</a> for the full project: Codebase, Database, Reference List, etc.
			</span>

			{/* Notice that this is commented out because -
            - there is no need to seperate the main function from the main page as for now
            <Link to = '/get_prof_by_id'>
                <p>Search Professor by ID</p>
            </Link>
            */}
			<section id='main'>
				<p id='searchByID_Intro'>
					By entering the pid of a professor of your choice, get the Sentiment Analysis result of students' commentary from RateMyProfessors.com on the professor right away!
				</p>

				<p id='searchByID_Intro'>
					For example, if the professor's URL on RateMyProfessors.com is <br />
					<span>https://www.ratemyprofessors.com/ShowRatings.jsp?tid=2105994</span>
					<br />
					Then, enter 2105994 in the text box below.
				</p>

				<p id='searchByID_Intro'>
					Notice that Sentiment Score (discrete) is computed based on individual comments, while Sentiment Score (continuous) is computed based on all comments. In other words, the higher
					the discrete score is, the more individual comments are positive. The higher the continuous score is, the larger proportion of all comments are positive. The sentiment score can be
					undeterministic, because of random tier-breaking.
				</p>
			</section>
		</div>
	)
}
export default Article
