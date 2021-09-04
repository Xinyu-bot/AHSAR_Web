import requests
from bs4 import BeautifulSoup
import multiprocessing
from concurrent.futures import ThreadPoolExecutor, as_completed
import time
import json
from tqdm import tqdm


class UrlException(Exception):
    pass 

def get_prof(url, pid) -> tuple: 
    try: 
        # get html body
        res = requests.get(url)
        # get html dom tree parsed
        dom = BeautifulSoup(res.text, features="html.parser")

        # professor name
        name = dom.select("div span")[2].text.strip() + ' ' + dom.select("div span")[3].text.strip()
        # professor info about department and school
        info = dom.find("div", {"class": "NameTitle__Title-dowf0z-1 iLYGwn"})
        # extract department and school name only
        department = info.select("span")[0].text.replace("Professor in the", "").replace("department at", "").strip()
        school = info.select("a")[0].text.strip()

        # printing to stdout under massive concurrent/parallel processing will negatively affect the running speed
        #print(pid, name, department, school) 
        return (pid, name, department, school)
    
    except: 
        return

# hmmmm
def work(start: int, end: int) -> list:
    ret = []

    executor = ThreadPoolExecutor(max_workers = 10)
    tasks = [executor.submit(get_prof, "https://www.ratemyprofessors.com/ShowRatings.jsp?tid=" + str(i), i) for i in range(start, end)]
    for future in as_completed(tasks):
        res = future.result()
        if res:
            ret.append(res)

    return ret

# main function
def main() -> None:
    # initilization
    POOL_SIZE = 20
    START = 1000000
    END = 2000000
    SEP = (END - START) // POOL_SIZE
    pool = multiprocessing.Pool(processes = POOL_SIZE)
    out = []
    d = {}

    pbar = tqdm(total = POOL_SIZE)
    def update(*a):
        pbar.update()

    st = time.time()
    # async apply work
    for i in range(POOL_SIZE):
        ret = pool.apply_async(work, (START + (i * SEP), START + ((i + 1) * SEP)), callback = update)
        out.append(ret)

    pool.close()
    pool.join()

    # async fetch result
    for i in range(len(out)):
        out[i] = out[i].get()
    et = time.time()

    #print("Work done in {0} seconds. ".format(round(et - st, 3)))

    for entries in out:
        # dispatch
        for entry in entries:
            pid, name, department, school = entry
            ds = d.get(school, {})
            dd = ds.get(department, [])
            dd.append((pid, name))
            d[school] = ds
            d[school][department] = dd

    with open("out.json", "a") as outstream:
        outstream.write(json.dumps(d))

    return

# always a good practice
if __name__ == "__main__":
    main()