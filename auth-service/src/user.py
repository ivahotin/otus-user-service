from dataclasses import dataclass
from uuid import UUID

@dataclass
class User:
    id: str
    login: str
    password: str

    def check_password(self, password: str) -> bool:
        return self.password == password