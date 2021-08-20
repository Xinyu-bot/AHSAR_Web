import React, {useState} from 'react'
import axios from 'axios'
import { Link } from 'react-router-dom'

function RMP() {
    const [pid_text, setPID] = useState('')
    const [send, setSend] = useState(false)
    
    const searchEnter = (e) => {
        setSend(!send)
        setPID('my post text')
        console.log('send:', send)
        console.log('pid_text:', pid_text)
        document.getElementById('myTextarea').value = pid_text
        console.log(pid_text)

        if (pid_text !== '') {
            axios.get('/get_prof_by_id', {
                //send along the pid_text user typed
                params: {
                    input: pid_text,
                },
            })
            .then((response) => {
                if(response.status === 200){
                    
                }
            })
            .catch(function (err) {
                console.error(err)    
            })
            
        } else {
            alert('You need to write some text')
            e.preventDefault()
        }
    }

    return (
        <div className = "RMP">
            <Link to = '/'>
                <h1>back</h1>
            </Link>
            <button onClick={(e) => {searchEnter(e)}} id='send_button'>Search</button>
            <textarea placeholder = "Please enter a PID" onInput={(e) => setPID(e.target.value)} />
            
        </div>
    );
}

export default RMP;
