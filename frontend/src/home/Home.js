import React, { useState, useEffect } from 'react'
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
		//pid不存在
		if (data.professor_name === '-1') {
			setRet('')
			setReady(3)

			//pid存在
		} else {
			//pid存在，但是没有评论
			setRet(data)

			const obj = {
				id: data.queryHash,
				name: data.professor_name,
				pid: data.pid,
			}
			//把返回的有效pid相关内容，加入搜索历史记录。用到了id，提高react渲染效率
			// Don't use JS array methods such as pop, push, shift, unshift
			// as these will not tell React to trigger a re-render.
			// Instead, make a copy of the array then add your new item onto the end
			// To update an item in the array use .map.
			// Assumes each array item is an object with an id.
			setSearchedList([obj, ...searchedList])
			console.log('searchedList', searchedList)

			if (data.sentiment_score_discrete === '-1') {
				if (data.difficulty_score === '-1' && data.quality_score === '-1') {
					//difficulty_score，quality_score都没有
					//console.log(4)
					setReady(4)
				} else {
					//difficulty_score，quality_score都有
					//console.log(2)
					setReady(2)
				}
				//pid存在，有评论，没有difficulty_score，没有quality_score
			} else if (data.difficulty_score === '-1' && data.quality_score === '-1') {
				setReady(5)
			} else {
				//全都有
				setReady(1)
			}
		}
	}

	const [ret, setRet] = useState()
	const [ready, setReady] = useState(0) // ready flag
	// initialize pid state
	const [pid, setPid] = useState('') // user input
	// initialize searchedList state
	const [searchedList, setSearchedList] = useState([])

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
					//上面的then block里面有问题
					console.log(err)
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

			<Article pid={pid} ret={ret} ready={ready} />
		</div>
	)
}

export default Home
