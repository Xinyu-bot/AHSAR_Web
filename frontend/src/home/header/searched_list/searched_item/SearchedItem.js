import React from 'react'
import PropTypes from 'prop-types'
import './SearchedItem.css'

SearchedItem.prototype = {
	name: PropTypes.string.isRequired,
	pid: PropTypes.string.isRequired,
	school: PropTypes.string.isRequired,
}

export default function SearchedItem(props) {
	return (
		<li className='searchedItem'>
			{props.name} &nbsp; {props.pid} &nbsp; {props.school}
		</li>
	)
}
