import React, { useState, useEffect } from 'react'
import { Switch, Route, useHistory } from 'react-router-dom'
import Search from './search/Search'
import Header from './header/Header'
import Result from './result/Result'
import Article from './result/Article'
function Home() {
	return (
		<div className='home'>
			<Header />
			<Switch>
				<Route path='/result' component={Result}></Route>
				<Route path='/search_by_name' component={Search}></Route>
			</Switch>
			<Article />
		</div>
	)
}

export default Home
