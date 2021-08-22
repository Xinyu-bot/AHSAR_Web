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
		if (data.professor_name === '-1') {
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
		//console.log('here')
		//下面axios发送请求，返回ret有延迟。在ret state被set之前，把ret state的内容设置为loading，触发一次render渲染页面。ret被返回，在setRet，再触发render重新渲染
		//用户输入了非空的内容，Header重新渲染，显示正在加载。
		setRet('loading')

		// input validation
		if (isNum(pid)) {
			// retrieve from backend API /get_prof_by_id
			axios
				.get('http://54.251.197.0:8080/get_prof_by_id', {
					params: {
						input: pid,
					},
				})
				// process returned result3w4
				.then((r) => {
					if (r.status === 200) {
						//正确的数据从服务器返回，Header重新渲染。
						_setRet(r.data)
					} else {
						setReady(-1)
					}
				})
				// catch error
				.catch((err) => {
					//服务器错误
					setReady(-1)
				})
		}

		// non-numbers (specifically non-positive integer) input
		else {
			alert('PID should only contain number 0 ~ 9')
			//用户输入的内容包含数字，让Hedaer重新渲染，返回一个空div
			setRet('non-numbers')
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
