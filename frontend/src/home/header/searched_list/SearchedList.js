import React from 'react'
import SearchedItem from './searched_item/SearchedItem'
import './SearchedList.css'

export default function SearchedList(props) {
	return (
		<div className='searched2'>
			<ul className='searchedList' style={props.style}>
				{props.searchedList.map((item) => (
					<SearchedItem key={item.id} name={item.name} pid={item.pid} />
				))}
			</ul>
		</div>
	)
}
