import React, { useState } from 'react'
import axios from 'axios'
import './Home.css'
import Header from './header/Header'
import Article from './article/Article'
import Result from './result/Result'

function Home() {
	//fake data
	const data = [
		{ id: '001', firstName: 'Lin', lastName: 'He' },
		{ id: '002', firstName: 'Anna', lastName: 'Zhang' },
		{ id: '003', firstName: 'Jessica', lastName: 'Zhou' },
	]

	// initialize pid state
	const [pid, setPid] = useState('') // user input
	// initialize searchedList state
	const [searchedList, setSearchedList] = useState('')

	//Home传给Header一个函数，为了Header把pid传给Home
	const getSearchedPid = (pid) => {
		setPid(pid)
	}

	return (
		<div className='home'>
			<Header getSearchedPid={getSearchedPid} searchedList={searchedList} />
			<Result />
			<Article pid={pid} />
		</div>
	)
}

export default Home
