from typing import Optional
import asyncio
import uuid

from .user import User

class UserStorage:

    def __init__(self) -> None:
        self._storage = {}
        self._lock = asyncio.Lock()
    
    async def add(self, login: str, password: str) -> bool:
        id_ = str(uuid.uuid4())
        user = User(id_, login, password)

        if login in self._storage:
            return False
        
        async with self._lock:
            self._storage[login] = user

        return True

    async def get(self, login: str) -> Optional[User]:
        async with self._lock:
            return self._storage.get(login)