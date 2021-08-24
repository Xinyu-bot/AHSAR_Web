import React, { useState, useEffect } from 'react'
import './Header.css'
import SearchedList from './searched_list/SearchedList'

export default function Header(props) {
	useEffect(() => {
		document.getElementById('pid').checked = true
	}, []) // only run it once!

	const [show, setShow] = useState(false)
	const [searchBy, setSearchBy] = useState('pid')

	const handleClick = (event) => {
		setShow(!show)
	}
	const handleClick2 = (event) => {
		document.getElementById('pname').checked = false
		event.target.checked = true
		document.getElementsByClassName('search-box')[0].setAttribute('placeholder', 'Please enter a PID eg: 2105994') //getElementsByClassName返回的是一个伪数组。。

		setSearchBy('pid')
	}
	const handleClick3 = (event) => {
		document.getElementById('pid').checked = false
		event.target.checked = true
		document.getElementsByClassName('search-box')[0].setAttribute('placeholder', 'Please enter a name eg: Adam Meyers')

		setSearchBy('name')
	}

	const handleKeyUp = (event) => {
		const { keyCode, target } = event
		if (keyCode !== 13) return
		//检查用户输入的是否只是空格或者根本没输入
		if (target.value.trim() === '') {
			alert('输入不能为空 Input cannot be empty')
			return
		}
		if (searchBy === 'pid') {
			// 把Header组件里，用户输入的pid传给Home组件。把用户输入string前后的空格去掉。
			console.log(searchBy)

			props.getSearchedPid(target.value.trim())
		} else {
			console.log(searchBy)
			props.getSearchedName(target.value.trim())
		}

		/* 清空输入框里面字
		target.value = ''*/
		target.style.color = 'grey'
		/* 失去焦点 */
		target.blur()
	}

	//Header传给SearchedList一个函数，为了SearchedList把pid传给Header,然后Header再传给Home
	const getClickedPid = (pid) => {
		console.log('here', pid)
		props.getSearchedPid(pid)
	}

	return (
		<div className='header1'>
			<div onClick={handleClick} className='history'>
				<span style={{ color: show ? '#ccc' : 'white' }}>Search History</span> {/* 文字不能被选中。被选中时，成为灰色。 */}
			</div>
			<div className='search'>
				<input onKeyUp={handleKeyUp} className='search-box' autoComplete='off' placeholder='Please enter a PID eg: 2105994' />
			</div>
			<form className='search_by'>
				Search by: &nbsp;
				<input type='radio' id='pid' name='searchby' value='pid' onClick={handleClick2} /> {/*name的值要一样才能单选 */}
				<label for='pid'>pid</label>
				&nbsp;
				<input type='radio' id='pname' name='searchby' value='pname' onClick={handleClick3} />
				<label for='pname'>name</label>
			</form>
			{/* 点击搜索历史节点，显示SearchedList */}
			<SearchedList getClickedPid={getClickedPid} searchedList={props.searchedList} style={{ display: show ? 'block' : 'none' }} />{' '}
			{/*style在component上不起作用，把style传给SearchedList组件，在里面的ul节点加上这个style */}
		</div>
	)
}
