import React from 'react'
import { useHistory } from 'react-router-dom'
import SearchedItem from '../header/searched_list/searched_item/SearchedItem'

export default function SearchResult(props) {
	let history = useHistory()
	const handleClick2 = (e) => {
		//首先利用事件委托，把事件委托给ul，要得到li里面点值，就得用e.target,而不是e.currentTarget。所以点击了谁就返回谁的innerHTML，这对英文版是没有任何问题的
		//但是如果用户翻译了页面，页面的结构就会发生改变。这时候，最小的可点击元素不是li了。因为翻译后，li节点中会被新增font节点。。。
		console.log('here', e.target.innerHTML, '-----------')
		let pid = ''
		//当用户点击了翻译页面后，页面的结构发生了改变，之前有英文的地方，都会被加上<font>标签。虽然因为事件冒泡机制，ul的click事件还是会执行，但是他获取到的innerHTML变化了。里面多了很多font标签

		if (!e.target.innerHTML.includes('&nbsp;')) {
			//当用户点击了翻译页面后, 点击了<font>pid</font>,会返回‘220’之类的
			if (!isNaN(Number(e.target.innerHTML))) {
				pid = e.target.innerHTML
			} 
			else {
				//当用户点击了翻译页面后, 点击了li里面名字区域，innerHTML会返回名字，比如‘马克亚当’。这个时候我们要得到该font节点的下一个font节点。
				if (e.target.parentNode.nextSibling) {
					let a = e.target.parentNode.nextSibling.nextSibling.firstChild //得到<font style="vertical-align: inherit;">220</font>
					pid = a.innerHTML
				}
				else{//当用户点击了翻译页面后, 点击了li里面学校区域，innerHTML会返回学校，比如‘纽约大学’。这个时候我们要得到该font节点的上一个font节点。
					let a = e.target.parentNode.previousSibling.previousSibling.firstChild //得到<font style="vertical-align: inherit;">220</font>
					pid = a.innerHTML
				}
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

	console.log('in searchresult', props)
	const splitName = (item, index) => {
		const nameStr = item.split(' ')[index]
		const newStr = nameStr.replaceAll('?', ' ') //replace只replace第一个。要用replaceAll。
		return newStr
	}
	return (
		<div className='searched2'>
			{(() => {
				//如果要求来自Result组件的searchedListByName
				if (props.searchedListByName !== undefined) {
					return (
						<ul style={{ maxWidth: '640px', margin: '0 auto' }} onClick={handleClick2} className='resultOfName'>
							{props.searchedListByName.map((item) => (
								<SearchedItem
									className='resultList'
									key={item.split(' ')[3]}
									name={splitName(item, 0)}
									filed={splitName(item, 1)}
									school={splitName(item, 2)}
									pid={item.split(' ')[3]}
								/>
							))}
						</ul>
					)
				}
			})()}
		</div>
	)
}
