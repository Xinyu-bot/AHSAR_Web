# kill old processes
ps aux | grep './app' | awk '{print $2}' | xargs kill -9 
ps aux | grep 'python3 NLP_server.py' | awk '{print $2}' | xargs kill -9 
ps aux | grep 'serve -l 443 --ssl-cert ./ahsar.pem --ssl-key ./ahsar.key -s build' | awk '{print $2}' | xargs kill -9