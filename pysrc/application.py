import trigram
import scraper
import unigram_lexicon_based
from nltk import PorterStemmer
import sys

''' Modulability is everything '''
def load_models() -> tuple: 
    # import models and setup for sentiment analysis process
    unigram_model = unigram_lexicon_based.generate_lexicon("../pysrc/unigram_lexicon_extended.csv")
    trigram_model, bigram_model = trigram.import_models()
    porterStemmer = PorterStemmer()

    return trigram_model, bigram_model, unigram_model, porterStemmer

#def fetch_prof_info() -> tuple: 
#    return scraper.get_comments(sys.argv[1])

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

    # display the results
    print(name, quality_score, difficulty_score, sentiment_score_disc, sentiment_score_cont)

    return name, quality_score, difficulty_score, sentiment_score_disc, sentiment_score_cont

''' main function '''
# main function of the application
def main() -> None: 
    # load models
    trigram_model, bigram_model, unigram_model, porterStemmer = load_models()

    # instruction info
    # print("Notice that Sentiment Score (discrete) is computed based on individual comments, \nwhile Sentiment Score (continuous) is computed based on all comments. ")
    # print("In other words, the higher the discrete score is, the more individual comments are positive. \nThe higher the continuous score is, the larger proportion of all comments are positive. \n")
    try:
        comments, quality_score, difficulty_score, name = scraper.get_comments(sys.argv[1]) #fetch_prof_info()
    except scraper.UrlException:
        return 
    if comments is None: 
        print("Professor {0} does not have any comment. ".format(name))
    else:
        name, quality_score, difficulty_score, sentiment_score_disc, sentiment_score_cont =\
            analyze_sentiment(comments, trigram_model, bigram_model, unigram_model, porterStemmer, quality_score, difficulty_score, name)
    
    '''
    # read user's input
    flag = 'y'
    while flag == 'y': 
        # actual analysis 
        try: 
            # get professor info
            comments, quality_score, difficulty_score, name = fetch_prof_info()
            # check info validity
            if comments is None: 
                print("Professor {0} does not have any comment. ".format(name))
                pass
            # sentiment analysis
            analyze_sentiment(comments, trigram_model, bigram_model, unigram_model, porterStemmer, quality_score, difficulty_score, name)
        except (AssertionError, TypeError) as e:
            pass

        # ask user for next round
        flag = input("Continue to next professor? [y/n] ")
        while flag != 'y' and flag != 'n': 
            flag = input("Invalid entry. Enter again... [y/n] ")
    '''    

    return name, quality_score, difficulty_score, sentiment_score_disc, sentiment_score_cont

# always a good practice
if __name__ == '__main__': 
    main()

