import React, { useState, useEffect } from 'react'
import PropTypes from 'prop-types'
import PubSub from 'pubsub-js'
import SearchedItem from './searched_item/SearchedItem'
import './SearchedList.css'

//对接收的props进行类型以及必要性的限制

export default function SearchedList(props) {
	const [a, setA] = useState([])
	useEffect(() => {
		// componentDidMount
		const token = PubSub.subscribe('searchedList', (msg, data) => {
			console.log('pubsub', data)
			setA(data)
		})
		//componentWillUnmount
		return () => {
			PubSub.unsubscribe(token) //must unsubscribe
		}
	}, []) 

	// initialize searchedList state
	//const [searchedList, setSearchedList] = useState([]) //如果localStorage里有，就用localStorage里的
	/*useEffect(() => {
		// Check browser support
		if (typeof Storage !== 'undefined' && localStorage.getItem('searchedList') !== null) {
			// Retrieve and set
			setSearchedList(JSON.parse(localStorage.getItem('searchedList')))
		}
	}, []) // only run it once!
	
	//监听searchedList是否通过setSearchedList改变成功。如果是，把改变成功的searchedList存到localStorage里
	//must useEffect! setSearchedList会被异步执行，在主线程里面使用searchedList都不是更新过后的值。但是setSearchedList里面不能写callback，所以用useEffect
	/*useEffect(() => {
		//console.log('before', searchedList)
		localStorage.setItem('searchedList', JSON.stringify(searchedList))
	}, [searchedList]) //注意这里searchedList是一个object，地址没有变就不会触发useEffect
*/
	const handleClick1 = (e) => {
		const pid = e.target.innerHTML.split('&nbsp;')[1]
		console.log(pid) //pid是string类型，我们要数字类型
		props.getClickedPid(parseInt(pid))
	}

	const handleClick2 = (e) => {
		const pid = e.target.innerHTML.split('&nbsp;')[1]
		console.log(pid) //pid是string类型，我们要数字类型
		props.getClickedPid2(parseInt(pid))
	}

	const splitName = (item, index) => {
		const nameStr = item.split(' ')[index]
		const newStr = nameStr.replaceAll('?', ' ') //replace只replace第一个。要用replaceAll。
		return newStr
	}

	return (
		<div className='searched2'>
			{(() => {
				return (
					<ul onClick={handleClick1} className='searchedList' style={props.style}>
						{a.map((item) => (
							<SearchedItem key={item.id} name={item.name} pid={item.pid} />
						))}
					</ul>
				)

				//如果要求来自Result组件的searchedListByName
				if (props.searchedListByName !== undefined) {
					return (
						<ul onClick={handleClick2} className='resultOfName' style={props.style}>
							{props.searchedListByName.map((item) => (
								<SearchedItem
									className='resultList'
									key={item.split(' ')[3]}
									name={splitName(item, 0)}
									filed={splitName(item, 1)}
									school={splitName(item, 2)}
									pid={item.split(' ')[3]}
								/>
							))}
						</ul>
					)
				}
			})()}
		</div>
	)
}
