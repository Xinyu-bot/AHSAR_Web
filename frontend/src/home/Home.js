import React, { useState, useEffect } from 'react'
import axios from 'axios'
import Header from './header/Header'
import Article from './article/Article'

function Home() {
	useEffect(() => {
		// Check browser support
		if (typeof Storage !== 'undefined' && localStorage.getItem('searchedList') !== null) {
			// Retrieve and set
			setSearchedList(JSON.parse(localStorage.getItem('searchedList')))
		}
	}, []) // only run it once!

	// 初始化state
	const [ret, setRet] = useState()
	const [ready, setReady] = useState(0) // ready flag
	// initialize pid state
	const [pid, setPid] = useState('') // user input
	// initialize searchedList state
	const [searchedList, setSearchedList] = useState([]) //如果localStorage里有，就用localStorage里的

	//监听searchedList是否通过setSearchedList改变成功。如果是，把改变成功的searchedList存到localStorage里
	//must useEffect! setSearchedList会被异步执行，在主线程里面使用searchedList都不是更新过后的值。但是setSearchedList里面不能写callback，所以用useEffect
	useEffect(() => {
		//console.log('before', searchedList)
		localStorage.setItem('searchedList', JSON.stringify(searchedList))
	}, [searchedList]) //注意这里searchedList是一个object，地址没有变就不会触发useEffect

	//如果用户搜索的pid有效，把从后端发回的与该pid有关的数据作为一个object的形式，用setSearchedList添加到searchedList
	function addToSearchedList(data) {
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
		//数据去重
		const flag = searchedList.some((item, index) => {
			//重复了
			if (item.pid === data.pid) {
				//删除当前元素，因为和pid重复了
				searchedList.splice(index, 1)
				//把当前pid添加到index[0]
				/*searchedList.unshift(obj)
				//触发重新渲染，让新元素显示在最上面。searchedList还是同一个object地址，所以不会触发useEffect[searchedList]...
				*/
				let newArr = [obj, ...searchedList] //必须return一个新的地址
				setSearchedList(newArr)

				return true //第一次true，some就结束迭代
			} else {
				return false
			}
		})
		//没有重复
		if (!flag) {
			//添加当前pid对应的obj到searchedList
			const newArr = [obj, ...searchedList]
			//检查是否多于10个元素，如果多于10个删掉。
			if (newArr.length > 10) {
				newArr.splice(10)
			}
			setSearchedList(newArr) //触发重新渲染
		}
	}

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
			//把符合条件的data以object的形式添加到searchedList
			addToSearchedList(data)

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

	const _setRetByName = (data) => {
		//name不存在
		if (data.hasResult === 'false') {
			setRet('')
			setReady(6)

			//name存在
		} else {
			setRet(data.ret)
			setReady(7)
		}
	}

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
				.get('http://54.251.197.0:8080/get_pid_by_name', {
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
		<div className='home'>
			<Header getSearchedPid={getSearchedPid} getSearchedName={getSearchedName} searchedList={searchedList} />

			<Article pid={pid} ret={ret} ready={ready} />
		</div>
	)
}

export default Home
