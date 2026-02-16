import http.client as cl
import json
import traceback
import requests
import datetime

class Co:
    """ ANSI color codes """
    BLACK = "\033[0;30m"
    RED = "\033[0;31m"
    GREEN = "\033[0;32m"
    BROWN = "\033[0;33m"
    BLUE = "\033[0;34m"
    PURPLE = "\033[0;35m"
    CYAN = "\033[0;36m"
    LIGHT_GRAY = "\033[0;37m"
    DARK_GRAY = "\033[1;30m"
    LIGHT_RED = "\033[1;31m"
    LIGHT_GREEN = "\033[1;32m"
    YELLOW = "\033[1;33m"
    LIGHT_BLUE = "\033[1;34m"
    LIGHT_PURPLE = "\033[1;35m"
    LIGHT_CYAN = "\033[1;36m"
    LIGHT_WHITE = "\033[1;37m"
    BOLD = "\033[1m"
    FAINT = "\033[2m"
    ITALIC = "\033[3m"
    UNDERLINE = "\033[4m"
    BLINK = "\033[5m"
    NEGATIVE = "\033[7m"
    CROSSED = "\033[9m"
    END = "\033[0m"

#############
#   setup   # 
#############

HOST = "localhost:3000"
API_PATH = "/api/v1"

# login-test
USERNAME = "test3@test.arpa"
PASSWORD = "test1234567890"

# follow-tests
FOLLOW_ID = 10

#############
#   tests   # 
#############
_MY_ID = None


def healtcheck():
    resp = requests.get(url=f"http://{HOST}{API_PATH}/health")

    if resp.status_code != 200:
        raise Exception(f"failed the healthcheck: {resp.status_code}; body: {resp.text}")
    else:
        log(f"{resp.text}", 3)

def login() -> List[dict]:
    resp = requests.post(url=f"http://{HOST}{API_PATH}/auth/login", json={"email": USERNAME, "password": PASSWORD})
    tokens = {}
    
    if resp.status_code != 200:
        raise Exception(f"failed to login, code {resp.status_code}, msg: {body}")


    for h in resp.headers:
        if h == "Set-Cookie":
            tokens['Cookie']  = resp.headers[h]
    if len(tokens) == 0:
        raise Exception("no tokens found :(")


    log(f"tokens: {len(tokens)}", 3)

    return tokens

def get_user(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.get(url=f"http://{HOST}{API_PATH}/user", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to get user information: {resp.status_code}; body: {resp.text}")
    else:
        global _MY_ID
        _MY_ID = int(json.loads(resp.text)['userId'])
        log(f"logged in as: {_MY_ID}", 3)



def follow(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.post(url=f"http://{HOST}{API_PATH}/follow?id={FOLLOW_ID}", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to follow user: {resp.status_code}; body: {resp.text}")
    else:
        log(f"{resp.text}", 3)

def unfollow(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.delete(url=f"http://{HOST}{API_PATH}/follow?id={FOLLOW_ID}", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to follow user: {resp.status_code}; body: {resp.text}")
    else:
        log(f"{resp.text}", 3)

def get_followers(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.get(url=f"http://{HOST}{API_PATH}/user/{FOLLOW_ID}/followers", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to get follower list: {resp.status_code}; body: {resp.text}")
    else:
        log(f"{FOLLOW_ID} is being followed by: {resp.text}", 3)

def get_following(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.get(url=f"http://{HOST}{API_PATH}/user/{_MY_ID}/following", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to get follower list: {resp.status_code}; body: {resp.text}")
    else:
        log(f"{_MY_ID} is following: {resp.text}", 3)

    


def log(msg: str, type: int):
    now = datetime.datetime.now()
    if type == 0:
        # OK
        print(f"{Co.END}{str(now)} | [{Co.LIGHT_GREEN}  OK  {Co.END}] {Co.END}{msg}")
    elif type == 1:
        # ERROR
        print(f"{Co.END}{str(now)} | [{Co.RED} FAIL {Co.END}] {Co.END}{msg}")

    elif type == 2:
        # running
        print(f"{Co.END}{str(now)} | [{Co.LIGHT_WHITE} EXEC {Co.END}] {Co.END}{msg}")
    elif type == 3:
        print(f"{Co.END}{str(now)} | [{Co.LIGHT_WHITE} INFO {Co.END}] {Co.END}{msg}")
    else:
        print("what")
    
if __name__ == "__main__":
    try:
        log("running healthchck...", 2)
        healtcheck()
        log("healthcheck successful", 0)

        log("running login...", 2)
        cookies = login()
        log("login successful", 0)

        log("running getuser...", 2)
        get_user(cookies)
        log("getuser successful", 0)


        log("running follow...", 2)
        follow(cookies)
        log("follow successful", 0)


        log("running get_followers...", 2)
        get_followers(cookies)
        log("get_followers successful", 0)

        log("running get_following...", 2)
        get_following(cookies)
        log("get_following successful", 0)

        log("running unfollow...", 2)
        unfollow(cookies)
        log("unfollow successful", 0)


    except Exception as e:
        log("An error occured", 1)
        log(f"testsuite ran into an error: {e}", 1)
