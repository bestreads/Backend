import http.client as cl
import json
import traceback
import requests

HOST = "localhost:3000"
API_PATH = "/api/v1"
USERNAME = "test3@test.arpa"
PASSWORD = "test1234567890"

def healtcheck():
    resp = requests.get(url=f"http://{HOST}{API_PATH}/health")

    if resp.status_code != 200:
        raise Exception(f"failed the healthcheck: {resp.status_code}; body: {resp.text}")
    return

def login() -> List[dict]:
    resp = requests.post(url=f"http://{HOST}{API_PATH}/auth/login", json={"email": USERNAME, "password": PASSWORD})
    tokens = {}
    
    if resp.status_code != 200:
        raise Exception(f"failed to login, code {res.status}, msg: {body}")


    for h in resp.headers:
        if h == "Set-Cookie":
            tokens['Cookie']  = resp.headers[h]
            
    return tokens
    
def follow(tokens: dict[str]):
    headers = {"Host": HOST}
    headers['Cookie'] = tokens['Cookie']

    resp = requests.post(url=f"http://{HOST}{API_PATH}/follow?id=10", headers=headers)

    if resp.status_code != 200:
        raise Exception(f"failed to follow user: {resp.status_code}; body: {resp.text}")
    return



def green(s: str) -> str:
    return f"\033[1;49;32m{s}\033[0m"

def red(s: str) -> str:
    return f"\033[1;49;91m{s}\033[0m"

def yellow(s: str) -> str:
    return f"\033[1;49;33m{s}\033[0m"

if __name__ == "__main__":
    try:
        print(yellow("---- running healthcheck... ----"))
        healtcheck()
        print(green("---- healthcheck successful ----"))
        
        print(yellow("---- running login... ----"))
        cookies = login()
        print(green("---- login successful ----"))
        print(green(cookies))
        
        print(yellow("---- running follow... ----"))
        follow(cookies)
        print(green("---- follow successful ----"))



    except Exception as e:
        print(red("---- error ----"))
        print(red(traceback.format_exc()))
        print(red(f"testsuite ran into an error: {e}"))
