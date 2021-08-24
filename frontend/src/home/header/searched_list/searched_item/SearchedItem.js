import React from 'react'
import './SearchedItem.css'

export default function SearchedItem(props) {
	return (
		<li className='searchedItem'>
			<span id='name'>{props.name} </span>   &nbsp;  <span id='pid'>{props.pid}</span>  &nbsp;   <span id='school'>{props.school}</span>
		</li>
	)
}
