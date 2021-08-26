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
		// Check browser support
		if (typeof Storage !== 'undefined' && localStorage.getItem('searchedList') !== null) {
			// Retrieve and set
			setA(JSON.parse(localStorage.getItem('searchedList')))
		}
		const token = PubSub.subscribe('searchedList', (msg, data) => {
			console.log('pubsub', data)
			setA(data)
		})
		//componentWillUnmount
		return () => {
			PubSub.unsubscribe(token) //must unsubscribe
		}
	}, [])

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
				if (a ) {
					return (
						<ul onClick={handleClick1} className='searchedList' style={props.style}>
							{a.map((item) => (
								<SearchedItem key={item.id} name={item.name} pid={item.pid} />
							))}
						</ul>
					)
				}
				/*
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
				}*/
			})()}
		</div>
	)
}
