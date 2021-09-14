import React from 'react'
import './App.css'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'
import Home from './home/Home'

function App() {
	return (
		<div className='app'>
			<Router>
				<Switch>
					<Route path='/'>
						<Home />
					</Route>
				</Switch>
			</Router>
		</div>
	)
}

export default App
