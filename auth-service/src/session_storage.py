from typing import Optional
import uuid
import asyncio

from .user import User

class SessionStorage:

    def __init__(self) -> None:
        self._storage = {}
        self._lock = asyncio.Lock()

    async def create_session(self, user: User) -> str:
        id_ = str(uuid.uuid4())
        async with self._lock:
            self._storage[id_] = user
        return id_

    async def delete_session(self, session_id: str):
        async with self._lock:
            if session_id in self._storage:
                del self._storage[session_id]

    async def get_user_by_session_id(self, session_id: str) -> Optional[User]:
        async with self._lock:
            return self._storage.get(session_id)