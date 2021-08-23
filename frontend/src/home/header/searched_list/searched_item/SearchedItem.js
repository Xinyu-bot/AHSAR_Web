import React from 'react'
import './SearchedItem.css'

export default function SearchedItem(props) {
	return (
		<li className='searchedItem'>
			{props.name} {props.pid}
		</li>
	)
}
