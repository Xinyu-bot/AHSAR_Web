import React, { useState, useEffect } from 'react'
import { useHistory } from 'react-router-dom'
import PropTypes from 'prop-types'
import { MenuOutlined, RotateLeftOutlined } from '@ant-design/icons'
import { DownOutlined } from '@ant-design/icons'
import { CaretDownOutlined } from '@ant-design/icons'
import './Header.css'
import SearchedList from './searched_list/SearchedList'

//对接收的props进行类型以及必要性的限制
Header.propTypes = {
	getSearchedPid: PropTypes.func,
	getSearchedName: PropTypes.func,
	searchedList: PropTypes.array,
}
export default function Header(props) {
	let history = useHistory()
	useEffect(() => {
		document.getElementById('pid').checked = true
	}, []) // only run it once!

	const [show, setShow] = useState(false)
	const [showName, setShowName] = useState(false)
	const [searchBy, setSearchBy] = useState('pid')

	const [category, setCategory] = useState('pid')

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

	const handleClick4 = (event) => {
		setShowName(!showName)
	}
	//我服了，翻译成中文后，不管之前这个字有没有被font包裹，都会新加font标签。。
	//英文版我们用event.target=event.currentTarget,因为最小点击节点就是span。然而翻译成中文版就不一样了。最小点击节点是<font>姓名</font>。
	//所以我们这里用currentTarget，这样的话中文版点击<span>内任意一点currentTarget都是<span>。然后我们在span里面找font再找姓名/进程号就行了
	const handleClick5 = (event) => {
		let current = event.currentTarget.innerHTML
		if (current.includes('<font')) {
			current = current.split('>')[2].split('<')[0]
			if (current === '姓名') {
				//手动改（英文版的不用改，current状态更新了它会自动改。但是这个翻译版，他不会自动改。。什么毛病。所以我这里手动改。
				event.currentTarget.previousSibling.firstChild.firstChild.innerHTML = `${current}`
				current = 'name'
			} else {
				//手动改 为翻译版。。
				event.currentTarget.previousSibling.firstChild.firstChild.innerHTML = `${current}`
				current = 'pid'
			}
		}
		console.log(current)
		let result = current === 'name' ? 'name' : 'pid'
		setCategory(result)
		let a = result === 'name' ? 'Please enter a name eg: Adam Meyers' : 'Please enter a PID eg: 2105994' //这里要用result，不用category（还没更新呢）
		document.getElementsByClassName('search-box')[0].setAttribute('placeholder', a)
		setSearchBy(result)
		setShowName(!showName)
	}

	function isNum(s) {
		if (s !== null && s !== '') {
			return !isNaN(s)
		}
		return false
	}

	const handleKeyUp = (event) => {
		const { keyCode, target } = event
		if (keyCode !== 13) return
		//检查用户输入的是否只是空格或者根本没输入
		if (target.value.trim() === '') {
			alert('输入不能为空 Input cannot be empty')
			return
		}
		if (isNum(searchBy)) {
			alert('PID should only contain number 0 ~ 9')
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

	//Header传给SearchedList一个函数，为了SearchedList把pid传给Header,然后Header再传给Home
	const getClickedPid = (pid) => {
		console.log('here', pid)
		props.getSearchedPid(pid)
	}

	return (
		<div className='header1-wrap'>
			<div className='header1'>
				<form className='search_by tablet'>
					Search by: &nbsp;
					<input type='radio' id='pid' name='searchby' value='pid' onClick={handleClick2} /> {/*name的值要一样才能单选 */}
					<label for='pid'>pid</label>
					&nbsp;
					<input type='radio' id='pname' name='searchby' value='pname' onClick={handleClick3} />
					<label for='pname'>name</label>
				</form>
				<div className='search_by mobile'>
					<span onClick={handleClick4}>
						{category}
						<DownOutlined />
					</span>

					<span className='select_name' onClick={handleClick5} style={{ display: showName ? 'block' : 'none' }}>
						{category === 'name' ? 'pid' : 'name'}
					</span>
				</div>
				<div className='search'>
					<input onKeyUp={handleKeyUp} className='search-box' autoComplete='off' placeholder='Please enter a PID eg: 2105994' />
				</div>
				{/*style在component上不起作用，把style传给SearchedList组件，在里面的ul节点加上这个style */}
				<div onClick={handleClick} className='history'>
					<span style={{ color: show ? 'rgba(0, 0, 0, 0.7)' : 'black' }} className='history-span tablet'>
						Search History
						<span className='downoutlined'>
							<CaretDownOutlined />
						</span>
					</span>
					{/* 文字不能被选中。被选中时，成为灰色。 */}
					<span className='mobile'>
						<MenuOutlined style={{ color: 'black', fontSize: '17px' }} />
					</span>
				</div>
				{/* 点击搜索历史节点，显示SearchedList */}
				<SearchedList className='history-list' getClickedPid={getClickedPid} searchedList={props.searchedList} style={{ display: show ? 'block' : 'none' }} />{' '}
			</div>
		</div>
	)
}
