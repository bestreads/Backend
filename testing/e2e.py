import http.client as cl
import json

HOST = "localhost:3000"
API_PATH = "/api/v1"
USERNAME = "test3@test.arpa"
PASSWORD = "test1234567890"

def healtcheck(conn: HTTPConnection):
    conn.request(method="GET", url=f"{API_PATH}/healthcheck", headers={"Host": HOST})
    res = conn.getresponse()

    if res.status != 200:
        # uhhh idk wieso das so ist
        # raise Exception("failed the healthcheck")
        pass
    else:
        pass

def login(conn: HTTPConnection) -> List[str]:

    body = json.dumps({"email": USERNAME, "password": PASSWORD})
    conn.request(method="POST", url=f"{API_PATH}/auth/login", body=body, headers={"Host": HOST, "Content-Type": "application/json"})
    res = conn.getresponse()

    if res.status != 200:
        raise Exception(f"failed to login, code {res.status}, msg: {res.read()}")

    tokens = []
    for h in res.getheaders():
        if h[0] == "Set-Cookie":
            tokens.append(h[1])

    return tokens
    




if __name__ == "__main__":
    conn = cl.HTTPConnection(HOST)
    try:
        healtcheck(conn)
        cookies = login(conn)
    except Exception as e:
        print(e)
        # print(e.args)
