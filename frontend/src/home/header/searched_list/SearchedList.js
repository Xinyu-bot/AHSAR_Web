import React from 'react'
import SearchedItem from './searched_item/SearchedItem'
import './SearchedList.css'

export default function SearchedList(props) {
	const splitName = (item, index) => {
		const nameStr = item.split(' ')[index]
		const newStr = nameStr.replaceAll('?', ' ')//replace只replace第一个。要用replaceAll。
		return newStr
	}

	return (
		<div className='searched2'>
			{(() => {
				if (props.searchedList !== undefined) {
					return (
						<ul className='searchedList' style={props.style}>
							{props.searchedList.map((item) => (
								<SearchedItem key={item.id} name={item.name} pid={item.pid} />
							))}
						</ul>
					)
				}

				//如果要求来自Result组件的searchedListByName
				if (props.searchedListByName !== undefined) {
					return (
						<ul className='resultOfName' style={props.style}>
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
