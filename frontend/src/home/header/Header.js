import React from 'react'
import axios from 'axios'
import './Header.css'
import SearchedList from './searched_list/SearchedList'

export default function Header(props) {
	const handleKeyUp = (event) => {
		const { keyCode, target } = event
		if (keyCode !== 13) return
		if (target.value.trim() === '') {
			alert('输入不能为空')
			return
		}

		console.log(target.value)
		// 把Header组件里，用户输入的pid传给Home组件
		props.getSearchedPid(target.value)
		// 清空输入框里面字
		target.value = ''
	}

	return (
		<div className='header1'>
			<div className='search-wrap'>
				<input onKeyUp={handleKeyUp} id='search-box' autoComplete='off' placeholder='Please enter a PID' />
			</div>

			<SearchedList searchedList={props.searchedList} />
		</div>
	)
}
