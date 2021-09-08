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
        name = name.replace(' ', '?')
        # professor info about department and school
        info = dom.find("div", {"class": "NameTitle__Title-dowf0z-1 iLYGwn"})
        # extract department and school name only
        department = info.select("span")[0].text.replace("Professor in the", "").replace("department at", "").strip().replace(' ', '`')
        school = info.select("a")[0].text.strip().replace(' ', '`')

        # Check if this professor has any rating 
        check_finder = dom.find('div', {'class': 'RatingValue__NumRatings-qw8sqy-0 jMkisx'})
        if check_finder != None and check_finder.text[0:10] == 'No ratings':
            prof = {'name': name, 'overall_score': -1, 'would_take_again': -1, 'difficulty': -1, 'comments': -1, "school": school, "department": department}
            return prof
        
        # Overall quantitative scores of this professor
        overall_score = float(dom.find('div', {'class': 'RatingValue__Numerator-qw8sqy-2 liyUjw'}).text)
        temp_finder = dom.find_all('div', {'class': 'FeedbackItem__FeedbackNumber-uof32n-1 kkESWs'})
        if len(temp_finder) == 1:
            would_take_again = 0
            difficulty = float(temp_finder[0].text)
        else:
            would_take_again = temp_finder[0].text
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
        prof = {'name': name, 'overall_score': overall_score, 'would_take_again': would_take_again, 'difficulty': difficulty, 'comments': comments, "school": school, "department": department}
        return prof
    
    except: 
        raise UrlException


# This function takes url as input, and returns a list 
# where list[0] is a list of verbal comments, list[1] is the overall score, list[2] is level of difficulty
def get_comments(user_in: str) -> list:
    url = "https://www.ratemyprofessors.com/ShowRatings.jsp?tid=" + user_in
    prof = get_prof(url)
    comments = -1
    if prof['comments'] != -1:
        comments = [x[2] for x in prof['comments']]
    return [comments, prof['overall_score'], prof['difficulty'], prof['name'], prof['would_take_again'], prof["school"], prof["department"]]

# This function takes a professor name as input, 
# and returns the url of the RMP rating page of this professor
def get_url(prof_name: str) -> list:
    name_parts = prof_name.lower().split()
    query = name_parts[0]
    for i in range(1, len(name_parts)):
        query = query + '+' + name_parts[i]

    # Get the results of searching by this professor name
    query_url = f'https://www.ratemyprofessors.com/search/teachers?query={query}'
    query_res = requests.get(query_url)
    query_dom = BeautifulSoup(query_res.text, features="html.parser")
    
    # If no professor found, return list of "-1"
    if query_dom.find('div', {'class': 'NoResultsFoundArea__StyledNoResultsFound-mju9e6-0 iManHc'}) != None:
        return None
    # Else
    else:
        prof_selector = query_dom.find_all('a', {'class': 'TeacherCard__StyledTeacherCard-syjs0d-0 dLJIlx'})
        prof_list = []
        for prof in prof_selector:
            name = prof.find('div', {'class': 'CardName__StyledCardName-sc-1gyrgim-0 cJdVEK'}).text
            department = prof.find('div', {'class': 'CardSchool__Department-sc-19lmz2k-0 haUIRO'}).text
            school = prof.find('div', {'class': 'CardSchool__School-sc-19lmz2k-1 iDlVGM'}).text
            pid = prof['href'].split('=')[-1]

            profStr = name.replace(' ', '?') + '&' + department.replace(' ', '?') + '&' + school.replace(' ', '?') + '&' + pid
            prof_list.append(profStr)

        return prof_list
