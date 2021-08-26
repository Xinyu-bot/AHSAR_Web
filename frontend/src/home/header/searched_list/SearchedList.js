import React, { useState, useEffect } from 'react'
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

	return (
		<div className='searched2'>
			{(() => {
				if (a) {
					return (
						<ul onClick={handleClick1} className='searchedList' style={props.style}>
							{a.map((item) => (
								<SearchedItem key={item.id} name={item.name} pid={item.pid} />
							))}
						</ul>
					)
				}
			})()}
		</div>
	)
}
