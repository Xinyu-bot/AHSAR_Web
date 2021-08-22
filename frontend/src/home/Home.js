import React, { useState } from 'react'
import axios from 'axios'
import Header from './header/Header'
import Article from './article/Article'
import Result from './article/result/Result'

function Home() {
	function isNum(s) {
		if (s !== null && s !== '') {
			return !isNaN(s)
		}
		return false
	}
	const _setRet = (data) => {
		if (data.first_name === '-1') {
			setReady(3)
		} else if (data.sentiment_score_discrete === '-1') {
			setReady(2)
		} else {
			setReady(1)
		}
		setRet(data)
	}

	const [ret, setRet] = useState()
	const [ready, setReady] = useState(0) // ready flag
	// initialize pid state
	const [pid, setPid] = useState('') // user input
	// initialize searchedList state
	const [searchedList, setSearchedList] = useState('')

	//Home传给Header一个函数，为了Header把pid传给Home
	const getSearchedPid = (pid) => {
		setPid(pid)

		// input validation
		if (isNum(pid)) {
			// retrieve from backend API /get_prof_by_id
			axios
				.get('http://54.251.197.0:8080/get_prof_by_id', {
					params: {
						input: pid,
					},
				})
				// process returned result
				.then((r) => {
					if (r.status === 200) {
						console.log(r.data)
						_setRet(r.data)
					} else {
						setReady(-1)
					}
				})
				// catch error
				.catch((err) => {
					console.error(err)
				})
		}

		// non-numbers (specifically non-positive integer) input
		else {
			alert('PID should only contain number 0 ~ 9')
			return
		}
	}

	return (
		<div className='home'>
			<Header getSearchedPid={getSearchedPid} searchedList={searchedList} />
			<Result />
			<Article pid={pid} ret={ret} ready={ready} />
		</div>
	)
}

export default Home
