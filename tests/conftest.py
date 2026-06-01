import os

import pytest

import tests.framework.helpers as helpers
from tests.framework.db_client import DbClient
from tests.framework.http_client import HttpClient


@pytest.fixture(scope="module")
def db_client():
    client = DbClient("notes.db")
    yield client
    client.close()


@pytest.fixture(scope="module")
def http_client():
    client = HttpClient(os.getenv("BASE_URL"))
    yield client


@pytest.fixture
def db_cleanup(db_client):
    users = []
    yield users
    for user in users:
        db_client.remove_user(user["id"])

@pytest.fixture
def create_user(http_client, db_client):
    credentials, created_user = helpers.create_user(http_client)
    yield credentials, created_user
    db_client.remove_user(created_user["id"])


@pytest.fixture
def login(http_client, create_user):
    creds, _ = create_user
    payload = {"username": creds["username"], "password": creds["password"]}
    logged_in_user = helpers.login(http_client, payload)
    yield logged_in_user


@pytest.fixture
def create_file(http_client, login, db_client):
    name = "mdFile.md"
    f = open(name)
    status, body = http_client.post_file("/api/upload", f, login["token"])
    yield status, name
    helpers.remove_file_from_storage(".storage/mdFile.md")
    db_client.remove_file(name)

@pytest.fixture
def convert_file(db_client, http_client, create_file, login):
    _, name = create_file
    #this should be the full name of the file, whatever it is
    status, body = http_client.post_json(f"/app/convert-html/mdFile", None, login["token"])
    yield status, body
    helpers.remove_file_from_storage(".storage/mdFile.html")
    db_client.remove_file("mdFile.html")


@pytest.fixture
def remove_files(db_client):
    files = []
    yield files
    for file in files:
        helpers.remove_file_from_storage(f".storage/{file}")
        db_client.remove_file(file)
