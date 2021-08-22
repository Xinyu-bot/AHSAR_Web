import requests
from bs4 import BeautifulSoup

class UrlException(Exception):
    pass 

# This function takes the url of an RMP rating page, 
# and returns the ratings of the prof as a dict: 
    # {name: Adam Meyers, overall_score: 3.3, would_take_again: 0.67, difficulty: 3, comments: [a_list_of_comments]}
    # "comments" stores all the comments for this professor. Each comment is stored as a list in the form of [quality, difficulty, verbal_comment].
def get_prof(url): 
    try: 
        res = requests.get(url)
        dom = BeautifulSoup(res.text, features="html.parser")

        # Professor name
        name = dom.select("div span")[2].text.strip() + ' ' + dom.select("div span")[3].text.strip()
        
        # Check if this professor has any rating 
        check_finder = dom.find('div', {'class': 'RatingValue__NumRatings-qw8sqy-0 jMkisx'})
        if check_finder != None and check_finder.text[0:10] == 'No ratings':
            prof = {'name': name, 'overall_score': None, 'would_take_again': None, 'difficulty': None, 'comments': None}
            return prof
        
        # Overall quantitative scores of this professor
        overall_score = float(dom.find('div', {'class': 'RatingValue__Numerator-qw8sqy-2 liyUjw'}).text)
        temp_finder = dom.find_all('div', {'class': 'FeedbackItem__FeedbackNumber-uof32n-1 kkESWs'})
        if len(temp_finder) == 1:
            would_take_again = 0
            difficulty = float(temp_finder[0].text)
        else:
            would_take_again = float(temp_finder[0].text.split('%')[0])
            difficulty = float(temp_finder[1].text)

        # Get all the comments of this professor
        comments_selector = dom.find('ul', {'class': 'RatingsList__RatingsUL-hn9one-0 cbdtns'}).select('li')
        quality_class = ['CardNumRating__CardNumRatingNumber-sc-17t4b9u-2 kMhQxZ', 'CardNumRating__CardNumRatingNumber-sc-17t4b9u-2 bUneqk', 'CardNumRating__CardNumRatingNumber-sc-17t4b9u-2 fJKuZx']
        difficulty_class = ['CardNumRating__CardNumRatingNumber-sc-17t4b9u-2 cDKJcc']

        comments = []  # A list of all the verbal comments of this professor
        for comment in comments_selector:
            if comment.find('div', {'class': difficulty_class[0]}) == None:
                continue
            for i in range(len(quality_class)):
                selector = comment.find('div', {'class': quality_class[i]})
                if selector != None:
                    q = float(selector.text)
                    break
                else:
                    i += 1
            d = float(comment.find('div', {'class': 'CardNumRating__CardNumRatingNumber-sc-17t4b9u-2 cDKJcc'}).text)
            vc = comment.find('div', {'class': 'Comments__StyledComments-dzzyvm-0 gRjWel'}).text

            # Get all comments
            comments.append([q, d, vc])

        # Each professor is stored as a dictionary in the following form:
        # {name: Adam Meyers, overall_score: 3.3, would_take_again: 0.67, difficulty: 3, comments: [a_list_of_comments]}
        # "comments" stores all the comments for this professor. Each comment is stored as a list in the form of [quality, difficulty, verbal_comment].
        prof = {'name': name, 'overall_score': overall_score, 'would_take_again': would_take_again, 'difficulty': difficulty, 'comments': comments}
        return prof
    
    except: 
        raise UrlException


# This function takes url as input, and returns a list 
# where list[0] is a list of verbal comments, list[1] is the overall score, list[2] is level of difficulty
def get_comments(user_in):
    url = "https://www.ratemyprofessors.com/ShowRatings.jsp?tid=" + user_in
    prof = get_prof(url)
    comments = None
    if prof['comments'] != None:
        comments = [x[2] for x in prof['comments']]
    return [comments, prof['overall_score'], prof['difficulty'], prof['name'], prof['would_take_again']]
