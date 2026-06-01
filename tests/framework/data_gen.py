from random import randint


def get_username(len: int) -> str:
    name = []
    for i in range(len):
        name.append(chr(randint(65, 90)))
    return "".join(name)
