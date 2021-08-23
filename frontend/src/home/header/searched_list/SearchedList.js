import React from 'react'
import SearchedItem from './searched_item/SearchedItem'
import './SearchedList.css'

export default function SearchedList(props) {
	console.log('22', props.searchedList)
	return (
		<ul className='searchedList'>
			{props.searchedList.map((item) => (
				<SearchedItem key={item.id} name={item.name} pid={item.pid} />
			))}
		</ul>
	)
}
