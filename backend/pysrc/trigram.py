import unigram_lexicon_based
import pickle
import sys
import random
import string
from nltk import PorterStemmer
from nltk.tokenize import ToktokTokenizer

# Out-of-Model threshold
# if a n-gram sequence occurred lower than this threshold, 
# we consider it as not valid and back-off to a lower n-gram model
OOM_THRESHOLD = 2
tokTok = ToktokTokenizer().tokenize

# unigram backoff <-- last defense before neutralizing the individual token
# use lexicon based unigram model, because it provides better result
def unigram_backoff(unigram_model: unigram_lexicon_based.Lexicon, bigram: str) -> tuple: 
    _sum = 0
    # unpack bigram into unigrams
    unigrams = bigram.split(' ')
    for unigram in unigrams: 
        try: 
            word_obj = unigram_model._get(unigram)
            _word, _sentiment = word_obj.word, int(word_obj.sentiment)
            if _sentiment == 0: 
                polarity = -1
            elif _sentiment == 1: 
                polarity = 1
            else: 
                polarity = 0
        except AttributeError:
            polarity = 0
        
        _sum += polarity

    # generate returned value
    if _sum == 0:
        ret = (1, 1)
    elif _sum > 0: 
        ret = (1, 0)
    else: 
        ret = (0, 1)
    
    return ret

def bigram_backoff(bigram_model: dict, unigram_model: unigram_lexicon_based.Lexicon, trigram: str) -> tuple: 
    # unpack trigram into two bigrams
    trigram = trigram.split(' ')
    bigrams = [' '.join((trigram[0], trigram[1])), ' '.join((trigram[1], trigram[2]))]

    pos = 0
    neg = 0
    for bigram in bigrams: 
        # retrieve pos and neg occurrence from bigram model
        arr = bigram_model.get(bigram, (-1, -1))
        # if not found or occurrence too few, backoff to unigram model
        if sum(arr) < OOM_THRESHOLD: 
            arr = unigram_backoff(unigram_model, bigram)
        _pos, _neg = arr[0], arr[1]
        pos += _pos / (_pos + _neg)
        neg += _neg / (_pos + _neg)

    # setup returned value
    if pos > neg: 
        ret = (1, 0)
    elif neg > pos:
        ret = (0, 1)
    else: 
        ret = (1, 1)

    return ret

def analyze_trigram(sentence: str, trigram_model: dict, bigram_model: dict, unigram_model: unigram_lexicon_based.Lexicon, porterStemmer: PorterStemmer) -> tuple: 
    # clean the input
    tokens = tokTok(sentence)
    tokens = [porterStemmer.stem(token.lower()) for token in tokens]

    # lazy generate trigrams from sentence
    trigrams = [' '.join((tokens[i], tokens[i + 1], tokens[i + 2])) for i in range(len(tokens)) if i < len(tokens) - 2]

    pos, neg = 0, 0
    # loop through each trigrams of the current sentence
    for trigram in trigrams: 
        # retrieve pos and neg occurrence from trigram model
        arr = trigram_model.get(trigram, (-1, -1))
        # if not found or occurrence too few, backoff to unigram model
        if sum(arr) < OOM_THRESHOLD: 
            arr = bigram_backoff(bigram_model, unigram_model, trigram)

        _pos, _neg = arr[0], arr[1]
        pos += _pos / (_pos + _neg)
        neg += _neg / (_pos + _neg)

    count_sum = pos + neg
    try: 
        # in form of (positive, negative)
        weight = (round(pos / count_sum, 3), round(neg / count_sum, 3))
    except ZeroDivisionError: 
        weight = (0, 0)

    if weight[0] == weight[1]: 
        if random.randint(0, 1): 
            return (1, 0)
        else: 
            return (0, 1)
    return weight

# helper function to import the previously exported model bytefiles
def import_models() -> tuple: 
    with open('./pysrc/trigram.model', 'rb') as handle:
        trigram_model = pickle.load(handle)
    with open('./pysrc/bigram.model', 'rb') as handle: 
        bigram_model = pickle.load(handle)

    return (trigram_model, bigram_model)