import React, { useState, useEffect } from 'react'
import { useHistory } from 'react-router-dom'
import PubSub from 'pubsub-js'
import SearchedItem from './searched_item/SearchedItem'
import './SearchedList.css'

//对接收的props进行类型以及必要性的限制

export default function SearchedList(props) {
	let history = useHistory()
	const [a, setA] = useState([])
	useEffect(() => {
		// componentDidMount
		// Check browser support
		if (typeof Storage !== 'undefined' && localStorage.getItem('searchedList') !== null) {
			// Retrieve and set
			setA(JSON.parse(localStorage.getItem('searchedList')))
		}
		const token = PubSub.subscribe('searchedList', (msg, data) => {
			console.log('pubsub', data)
			setA(data)
		})
		//componentWillUnmount
		return () => {
			PubSub.unsubscribe(token) //must unsubscribe
		}
	}, [])

	const handleClick1 = (e) => {
		//首先利用事件委托，把事件委托给ul，要得到li里面点值，就得用e.target,而不是e.currentTarget。所以点击了谁就返回谁的innerHTML，这对英文版是没有任何问题的
		//但是如果用户翻译了页面，页面的结构就会发生改变。这时候，最小的可点击元素不是li了。因为翻译后，li节点中会被新增font节点。。。
		console.log('here', e.target.innerHTML, '-----------')
		let pid = ''
		//当用户点击了翻译页面后，页面的结构发生了改变，之前有英文的地方，都会被加上<font>标签。虽然因为事件冒泡机制，ul的click事件还是会执行，但是他获取到的innerHTML变化了。里面多了很多font标签

		if (!e.target.innerHTML.includes('&nbsp;')) {
			//当用户点击了翻译页面后, 点击了<font>pid</font>,会返回‘220’之类的
			if (!isNaN(Number(e.target.innerHTML))) {
				pid = e.target.innerHTML
			} //当用户点击了翻译页面后, 点击了li里面名字区域，innerHTML会返回名字，比如‘马克亚当’。这个时候我们要得到该font节点的下一个font节点。
			else {
				let a = e.target.parentNode.nextSibling.nextSibling.firstChild //得到<font style="vertical-align: inherit;">220</font>
				pid = a.innerHTML
			}
		}
		//当用户点击了翻译页面后, 点击了li里面除font以外的区域，innerHTML会返回即含有名字的font+含有pid的font
		//比如‘<font style="vertical-align: inherit;"><font style="vertical-align: inherit;">马克亚当</font></font> &nbsp; <font style="vertical-align: inherit;"><font style="vertical-align: inherit;">511543</font></font> &nbsp;’
		else if (e.target.innerHTML.split('&nbsp;')[1].includes('<font')) {
			pid = e.target.innerHTML.split('&nbsp;')[1].split('">')[2].split('<')[0]
		}

		//用户没有点翻译，innerHTML就是“adam 220”之类的
		else if (e.target.innerHTML.split('&nbsp;')[1] && !e.target.innerHTML.split('&nbsp;')[1].includes('<font')) {
			pid = e.target.innerHTML.split('&nbsp;')[1]
		}

		console.log('pid', pid) //pid是string类型，我们要数字类型
		history.push({
			pathname: `/result`,
			search: `?pid=${parseInt(pid)}`,
		})
	}

	return (
		<div className='searched2'>
			{(() => {
				if (a) {
					return (
						<ul onClick={handleClick1} className='searchedList' style={props.style}>
							{a.map((item) => (
								<SearchedItem key={item.id} name={item.name} pid={item.pid} />
							))}
						</ul>
					)
				}
			})()}
		</div>
	)
}
