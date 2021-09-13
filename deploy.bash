# run frontend
cd frontend/
serve -s build &

# run backend
cd ../
cd backend/
./app &

# run NLP Server
cd ../
cd pysrc/ # must move to ./pysrc/ directory for NLP-related files to be loaded successfully
python3 NLP_server.py &