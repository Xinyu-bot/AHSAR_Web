import React, { useState, useEffect } from 'react'
import { useHistory, NavLink } from 'react-router-dom'
import './Content.css'
import Article from '../result/Article'
export default function Content() {
	const [searchBy, setSearchBy] = useState('pid')
	const [category, setCategory] = useState('pid')
	let history = useHistory()

	const handleKeyUp = (event) => {
		const { keyCode, target } = event
		if (keyCode !== 13) return
		//检查用户输入的是否只是空格或者根本没输入
		if (target.value.trim() === '') {
			return
		}

		if (searchBy === 'pid') {
			// 把Header组件里，用户输入的pid传给Home组件。把用户输入string前后的空格去掉。
			console.log(searchBy, 'huiche')

			//props.getSearchedPid(target.value.trim())
			history.push({
				pathname: `/result`,
				search: `?pid=${target.value.trim()}`,
			})
		} else {
			//search by name
			console.log(searchBy)
			//props.getSearchedName(target.value.trim())

			history.push({
				pathname: `/search_by_name`,
				search: `?name=${target.value.trim()}`,
			})
		}
		/* 清空输入框里面字
		target.value = ''*/
		target.style.color = 'grey'
		/* 失去焦点 */
		target.blur()
	}

	const handleClick5 = () => {
		console.log(searchBy)
		let result = searchBy === 'name' ? 'pid' : 'name'
		setCategory(result)
		let a = result === 'name' ? 'Please enter a name eg: Adam Meyers' : 'Please enter a PID eg: 2105994' //这里要用result，不用category（还没更新呢）
		document.getElementsByClassName('search-box')[0].setAttribute('placeholder', a)
		setSearchBy(result)
	}
	return (
		<div className='content'>
			<div className='header-wrap'>
				<div className='header2'>
					<a id='home' href='#'>
						Home
					</a>
					<a id='about' href='#main'>
						ABOUT
					</a>
				</div>
			</div>

			<div className='content-wrap'>
				<div className='search-wrap'>
					<div className='search'>
						<div>Search professor by {category}</div>

						<input onKeyUp={handleKeyUp} className='search-box' autoComplete='off' placeholder='Please enter a PID eg: 2105994' />

						<div className='select_name' onClick={handleClick5}>
							I want to search professor by {category === 'name' ? 'pid' : 'name'}
						</div>
					</div>
				</div>
			</div>

			<Article />
		</div>
	)
}
