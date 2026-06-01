import sqlite3


class DbClient:
    def __init__(self, db_name):
        self.db_name = db_name
        self.connection = sqlite3.connect(db_name)

    def remove_user(self, id):
        try:
            curr = self.connection.cursor()
            query = f"delete from users where id = {id}"
            curr.execute(query)
            self.connection.commit()
        except Exception as e:
            print(f"DbClient: exception: {e}")

    def remove_file(self, id):
        try:
            curr = self.connection.cursor()
            query = f"update files set deleted = 1 where name = '{id}'"
            curr.execute(query)
            self.connection.commit()
        except Exception as e:
            print(f"DbClient: exception: {e}")

    def __get(self, query):
        try:
            curr = self.connection.cursor()
            res = curr.execute(query)
            res = res.fetchone()
            return res
        except Exception as e:
            print(f"DbClient: exception: {e}")

    def get_file(self, id):
        query = f"select * from files where name = '{id}' order by created_at desc limit 1"
        try:
            return self.__get(query)
        except Exception as e:
            print("Exception")
            exit()

    def get_user(self, id):
        query = f"select * from users where username = '{id}'"
        try:
            return self.__get(query)
        except Exception as e:
            print("Exception")
            exit()

    def close(self):
        self.connection.close()
