import trigram
import scraper
import unigram_lexicon_based
import socket
import multiprocessing
from nltk import PorterStemmer
from time import time

# load the pre-trained models
def load_models() -> tuple: 
    # import models and setup for sentiment analysis process
    unigram_model = unigram_lexicon_based.generate_lexicon("../pysrc/unigram_lexicon_extended.csv")
    trigram_model, bigram_model = trigram.import_models()
    porterStemmer = PorterStemmer()
    
    return trigram_model, bigram_model, unigram_model, porterStemmer

# analyze sentiment given comments retrieved from RMP
def analyze_sentiment(comments: list, trigram_model: dict, bigram_model: dict, unigram_model: dict, 
porterStemmer: PorterStemmer, quality_score: float, difficulty_score: float, name: str) -> None: 
    pos, neg, count = 0, 0, 0
    weight = [0, 0]
    # loop through the comments
    for comment in comments: 
        # update the counts
        _pos, _neg = trigram.analyze_trigram(comment, trigram_model, bigram_model, unigram_model, porterStemmer)
        if _pos > _neg: 
            pos += 1
        elif _neg > _pos: 
            neg += 1
        else: 
            pass
        weight[0] += _pos
        weight[1] += _neg
        count += 1
    # compute two sentiment scores
    sentiment_score_disc = round(3.0 + (4.0 * ((pos / (pos + neg)) - 0.5)), 1)
    sentiment_score_cont = round(3.0 + (4.0 * ((weight[0] / count) - 0.5)), 1)

    return name, str(quality_score), str(difficulty_score), str(sentiment_score_disc), str(sentiment_score_cont)

# retrieve comments from RMP, analyze sentimenet, and return result to Backend Server
def analyze(comm_socket: socket.socket, pid: str) -> None:
    func_start = time()
    ret = None
    try:
        comments, quality_score, difficulty_score, name, would_take_again = scraper.get_comments(pid) 
    except scraper.UrlException:
        ret = ("-1", "-1", "-1", "-1", "-1", "-1", "-1") 
    else:
        if comments == -1: 
            ret = (name, quality_score, difficulty_score, "-1", "-1", would_take_again, pid)
        else:
            ret = list(analyze_sentiment(
                comments, trigram_model, bigram_model, unigram_model, 
                porterStemmer, quality_score, difficulty_score, name
                ))
            ret.append(would_take_again)
            ret.append(pid)
    finally:
        ret = [str(x) for x in ret]
        print(ret)
        comm_socket.send(" ".join(ret).encode())
        print("Task pid#{0} done in {1} seconds".format(pid, round(time() - func_start, 3)))
        comm_socket.close()
    return

def fetchPID(comm_socket: socket.socket, name: str) -> None:
    func_start = time()
    ret = None
    try:
        ret = scraper.get_url(name)
    except:
        pass
    finally:
        if ret == None or ret == []:
            comm_socket.send("-1 -1".encode())
        else:
            comm_socket.send(' '.join(ret).encode())
        print("Task name#{0} done in {1} seconds".format(name, round(time() - func_start, 3)))
        comm_socket.close()
    return

''' gloabl variables '''
gv_start = time()
# load models
trigram_model, bigram_model, unigram_model, porterStemmer = load_models()
# start TCP server socket
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.bind(("localhost", 5005))
sock.listen()
# process pool
pool = multiprocessing.Pool(processes = 10)

# main function of the application
def main() -> None: 
    print("NLP Server environment loaded in {0} seconds".format(round(time() - gv_start, 3)))
    print("NLP Server running")

    # read query from Backend Server, process the query, and return result
    while True: 
        (comm_socket, client_addr) = sock.accept()
        data = comm_socket.recv(128).decode()
        mode, query = data[0], data[1:]
        print("Received mode:", mode + ", query:", query)
        if mode == "0":
            pool.apply_async(analyze, (comm_socket, query, ))
        elif mode == "1":
            pool.apply_async(fetchPID, (comm_socket, query, ))

    # quit elegantly
    pool.close()
    pool.join()
    return

# always a good practice
if __name__ == '__main__': 
    main()
