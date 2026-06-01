import pytest

from tests.framework.data_gen import get_username
from tests.framework.http_client import HttpClient


def test_create_user(db_cleanup, http_client):
    username = get_username(14)
    payload = {"username": username, "password": "test123@!221"}
    status, body = http_client.post_json("/api/users", payload)
    assert status == 201
    assert body is not None
    assert "id" in body
    db_cleanup.append(body)


def test_create_doc(login, http_client):
    f = open("mdFile.md")
    status, body = http_client.post_file("/api/upload", f, login["token"])
    assert status == 201
