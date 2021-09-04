import React, { useState, useEffect } from 'react'
import axios from 'axios'
import qs from 'querystring'
import SearchResult from './SearchResult'
import Header from '../header/Header'

export default function Search(props) {
	// 初始化state
	const [ret, setRet] = useState()
	const [ready, setReady] = useState(0) // ready flag
	// initialize pid state
	const [pid, setPid] = useState(props.location.search) // user input

	useEffect(() => {
		console.log('props', props)
		//接收search参数
		const { search } = props.location
		const { name } = qs.parse(search.slice(1))
		getSearchedName(name)
	}, [props.location.search]) //must 监听history对象里面的内容（也就是url是否变了）

	const _setRetByName = (data) => {
		//name不存在
		if (data.hasResult === 'false') {
			setReady(6)
			setRet('')

			//name存在
		} else {
			setReady(7)
			setRet(data.ret)
		}
	}

	function hasNumber(myString) {
		return /\d/.test(myString)
	}
	//Home传给Header一个函数，为了Header把name传给Home
	const getSearchedName = (name) => {
		setPid(name)
		setRet('loading')
		if (!hasNumber(name)) {
			// retrieve from backend API /get_prof_by_id
			axios
				.get('http://1.14.137.215:8080/get_pid_by_name', {
					params: {
						input: name,
					},
				})
				// process returned result3w4
				.then((r) => {
					if (r.status === 200) {
						//正确的数据从服务器返回，Header重新渲染。
						console.log(r.data)
						_setRetByName(r.data)
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
		} else {
			alert('Name should not contain number')
			setRet('non-numbers')
			return
		}
	}

	return (
		<div className='result'>
			<Header />
			{(() => {
				if (ret !== undefined) {
					// evil switch case here, Oops

					if (ret === 'loading') {
						return (
							<div>
								<h2 id='result'>Result for {pid}</h2>
								<span>Loading data...</span>
								<span className='loader'></span>
							</div>
						)
					} else if (ret === 'non-numbers') {
						return <p>Please enter a valid pid or name</p>
					}
					//如果用户输入内容符合条件，并且收到了服务器的返回内容。
					switch (ready) {
						// user changes input or there is no input at all
						case 0:
							return
						// response code not 200, ERROR in backend
						case -1:
							return (
								<span>
									<p>Something is wrong with the server, please try again later.</p>
								</span>
							)

						//for search by name
						case 6:
							return (
								<span>
									<h2 id='result'>Result for {pid}</h2>
									<p>Sorry, this name is invalid. There is no such name on RateMyProfessors.com!</p>
								</span>
							)
						case 7:
							return (
								<span>
									<h2 id='result'>Result for {pid}</h2>

									<SearchResult searchedListByName={ret} />
								</span>
							)

						// make compiler happy
						default:
							return
					}
				}
			})()}
		</div>
	)
}
