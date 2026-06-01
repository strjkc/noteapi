from tests.framework.data_gen import get_username
from tests.framework.http_client import HttpClient
import os


# TODO: random pass
def create_user(http_client: HttpClient):
    payload = {"username": get_username(15), "password": "TestPass1234"}
    status, body = http_client.post_json("/api/users", payload)
    if status != 201 or body is None:
        print("User not created")
        print(f"data: {payload}")
        exit()
    return payload, body


def login(http_client: HttpClient, user):
    status, body = http_client.post_json("/api/login", user)
    # should be 200
    if status != 201 or body is None:
        print("User not logged in")
        print(f"data: {user}")
        exit()
    return body

def remove_file_from_storage(name):
    try:
        os.remove(name)
    except FileNotFoundError:
        print("File already removed")
    except Exception as e:
        print(f"Remove file from disk: {e}")
