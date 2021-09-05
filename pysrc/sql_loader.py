import json
import pandas as pd
import time
from sqlalchemy import create_engine

def read_SDP() -> dict:
    with open('../pysrc/SDP.json', 'r') as instream:
        SDP = json.load(instream)
        return SDP

def main() -> None:
    engine = create_engine('mysql+pymysql://root:123@localhost:3306/ahsar?charset=utf8')
    print(engine)

    a = time.time()
    SDP = read_SDP()
    print("SDP model loaded into Dictionary in {0} seconds. ".format(round(time.time() - a, 3)))

    b = time.time()
    l = []
    for school in SDP.keys():
        for department, pList in SDP[school].items():
            for p in pList:
                a = [p[0], p[1], school, department, p[2], p[3], p[4], -1, -1, 'Never']
                l.append(a)
    df = pd.DataFrame(columns = ['pid', 'prof_name', 'school', 'department', 'quality_score', 'difficulty_score', 'would_take_again', 'sentiment_socre_continuous', 'sentiment_score_discrete', 'last_update'], data = l)
    print("SDP model loaded into DataFrame in {0} seconds. ".format(round(time.time() - b, 3)))
    
    a = time.time()
    df.to_sql('prof', engine, 'ahsar')
    print("SDP model loaded into MySQL in {0} seconds. ".format(round(time.time() - a, 3)))

    return

if __name__ == "__main__":
    main()