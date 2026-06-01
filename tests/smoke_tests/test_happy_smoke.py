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


def test_login(http_client: HttpClient, create_user):
    creds, _ = create_user
    status, body = http_client.post_json("/api/login", creds)
    # should be 200
    assert status == 201
    assert body is not None
    assert "token" in body


def test_spell_check(http_client: HttpClient):
    f = open("mdFileMistakes.md")
    status, body = http_client.post_file("/api/spellcheck/eng", f)
    assert status == 200
    assert body is not None
    # we should assert the number of errors


def test_upload_file(login, http_client: HttpClient, remove_files):
    f = open("mdFile.md")
    status, body = http_client.post_file("/api/upload", f, login["token"])
    assert status == 201
    remove_files.append("mdFile.md")

def test_convert_to_html(remove_files, login, http_client: HttpClient, create_file):
    _, name = create_file
    #this should be the full name of the file, whatever it is
    status, body = http_client.post_json(f"/app/convert-html/mdFile", None, login["token"])
    assert status == 201
    remove_files.append("mdFile.html")

def test_get_html(login, http_client: HttpClient, convert_file):
    _, name = convert_file
    status, body = http_client.get(f"/app/get-html/mdFile", login["token"])
    assert status == 200
    # we should assert the body somehow


def test_remove_file(login, http_client: HttpClient, create_file, db_client):
    _, name = create_file
    status, body = http_client.delete(f"/api/files/{name}", login["token"])
    res = db_client.get_file(name)
    # assert that the file doesn't exist on disk
    assert status == 204
    assert res[-1] == 1


def test_remove_user(login, db_client: DbClient, http_client: HttpClient):
    status, body = http_client.delete(f"/api/users", login["token"])
    res = db_client.get_user(login["username"])
    assert status == 204
    assert res is None
