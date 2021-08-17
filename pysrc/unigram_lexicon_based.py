import os
from nltk.tokenize import word_tokenize, sent_tokenize

''' 
a = Lexicon()
word = 'hello'
sentiment = 'negative'
setattr(a, word, Token(word, 'verb', sentiment))
print(a.hello.word) --> 'hello'
Lexicon: 
self.pervert -> self.pervert.word
self.phobic
self.phony
'''
class Lexicon: 
    __slot__ = ('__dict__')
    def __init__(self): 
        pass

    def _get(self, token):
        return object.__getattribute__(self, token)

class Token: 
    __slots__ = ('word', 'sentiment')
    def __init__(self, word, sentiment):
        self.word = word
        self.sentiment = sentiment

# reading in the file
def generate_lexicon(filename: str) -> Lexicon: 
    # create the lexicon class
    unigram = Lexicon()

    # read in the file
    with open(filename, 'r') as instream: 
        for line in instream: 
            line = line.strip(os.linesep).split(',')
            word, sentiment = line[0], line[1]
            setattr(unigram, word, Token(word, sentiment))

    return unigram

def comment_parsing(comment: str, lexicon: Lexicon) -> dict: 
    comment_sentiment = []
    sentences = sent_tokenize(comment)
    # loop through each sentence
    for sentence in sentences: 
        sentence_sentiment = []
        clauses = sentence.split(',')
        for clause in clauses: 
            flag = False
            clause = word_tokenize(clause)
            # loop through the clause
            for token in clause: 
                token = token.lower()
                # try to get token info from lexicon
                if token == 'not' or token == "n't": 
                    flag = True
                if token == 'but': 
                    flag = False

                try: 
                    word_obj = lexicon._get(token)
                    _word, _sentiment = word_obj.word, word_obj.sentiment
                    
                    if flag: 
                        if _sentiment == 'positive':
                            _sentiment = 'negative'
                        elif _sentiment == 'negative': 
                            _sentiment = 'positive'
                    
                    sentence_sentiment.append(_sentiment)
                # if not in lexicon, ignore
                except (AttributeError):
                    pass 

        comment_sentiment += sentence_sentiment
    
    return comment_sentiment

def sentiment_analysis(comment_sentiment: list) -> tuple:
    positive_count = 0
    negative_count = 0
    for element in comment_sentiment: 
        element = int(element)
        # if element > 0: 
        if element == 1: 
            positive_count += 1
        # elif element < 0: 
        elif element == 0: 
            negative_count += 1
        else: 
            continue

    count_sum = positive_count + negative_count
    try: 
        # (positive, negative)
        weight = (round(positive_count / count_sum, 3), round(negative_count / count_sum, 3))
    except ZeroDivisionError: 
        weight = (0, 0)

    return tuple(weight)