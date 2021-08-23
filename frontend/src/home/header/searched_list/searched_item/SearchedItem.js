import React from 'react'
import './SearchedItem.css'

export default function SearchedItem(props) {
    console.log('33',props)
	return (
		<li className='searchedItem'>
			{props.name} {props.pid}
		</li>
	)
}
