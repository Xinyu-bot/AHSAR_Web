import React, { useState, useEffect } from 'react'
import { Switch, Route, useHistory } from 'react-router-dom'
import Search from './search/Search'
import Header from './header/Header'
import Result from './result/Result'
import Article from './result/Article'
import Content from './content/Content'
import './Home.css'
function Home() {
	return (
		<div className='Home'>
			<Switch>
				<Route path='/result' component={Result}></Route>
				<Route path='/search_by_name' component={Search}></Route>
				<Route path='/' component={Content}></Route>
			</Switch>
			<div id='footer'>
				粤ICP备2021131165号-1 | 粤ICP备2021131165号 | <a href='https://beian.miit.gov.cn/#/Integrated/recordQuery'>工信部查询</a>{' '}
			</div>
		</div>
	)
}

export default Home
