import React from 'react'
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

		// 把Header组件里，用户输入的pid传给Home组件
		props.getSearchedPid(target.value)
		/* 清空输入框里面字
		target.value = ''*/
		target.style.color = 'grey'
		/* 失去焦点 */
		target.blur()
		
	}

	/*文本全选*/
	const focus = (event) => {
		event.target.select()
	}
	return (
		<div className='header1'>
			<input onKeyUp={handleKeyUp} onFocus={focus} id='search-box' autoComplete='off' placeholder='Please enter a PID' />
			<SearchedList searchedList={props.searchedList} />
		</div>
	)
}
