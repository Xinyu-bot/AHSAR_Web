import React from 'react'
import { useHistory } from 'react-router-dom'
import SearchedItem from '../header/searched_list/searched_item/SearchedItem'

export default function SearchResult(props) {
	let history = useHistory()
	const handleClick2 = (e) => {
		const pid = e.target.innerHTML.split('&nbsp;')[1]
		console.log(pid) //pid是string类型，我们要数字类型
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
						<ul onClick={handleClick2} className='resultOfName'>
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
