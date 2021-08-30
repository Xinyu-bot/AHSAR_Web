import React, { useState, useEffect } from 'react'
import { Switch, Route, useHistory } from 'react-router-dom'
import Search from './search/Search'
import Header from './header/Header'
import Result from './result/Result'
import Article from './result/Article'
import Content from './content/Content'

function Home() {
	return (
		<div className='home'>
			<Switch>
				<Route path='/result' component={Result}></Route>
				<Route path='/search_by_name' component={Search}></Route>
				<Route path='/' component={Content}></Route>
			</Switch>
		</div>
	)
}

export default Home
