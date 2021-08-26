import React, { useState, useEffect } from 'react'
import { Route, useHistory } from 'react-router-dom'

import axios from 'axios'
import Header from './header/Header'
import Result from './result/Result'
import Article from './result/Article'
function Home() {
	return (
		<div className='home'>
			<Header />
			<Route path='/result' component={Result}></Route> 
			<Article />
		</div>
	)
}

export default Home
