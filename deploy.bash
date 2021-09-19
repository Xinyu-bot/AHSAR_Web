# run frontend
cd frontend/
nohup serve -l 443 -s build &

# run NLP Server
cd ../
cd pysrc/ # must move to ./pysrc/ directory for NLP-related files to be loaded successfully
nohup python3 NLP_server.py 2>&1 &
sleep 5s

# run backend
cd ../
cd backend/
nohup ./app 2>&1 &