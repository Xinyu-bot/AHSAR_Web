# kill old processes
ps aux | grep './app' | awk '{print $2}' | xargs kill -9 
ps aux | grep 'python3 NLP_server.py' | awk '{print $2}' | xargs kill -9 
ps aux | grep 'serve -s build' | awk '{print $2}' | xargs kill -9