import React from 'react'
import './App.css'
import RMP from './sentiment_analysis/rmp'
import Home from './home/Home'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'

function App() {
	return (
		<div className='App'>
			<Router>
				<Switch>
					<Route path='/get_prof_by_id'>
						<RMP />
					</Route>

					<Route path='/'>
						<Home />
					</Route>
				</Switch>
			</Router>
		</div>
	)
}

export default App
