import React from 'react'
import './App.css'
import Result from './home/result/Result'
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom'
import Content from './home/content/Content'
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
