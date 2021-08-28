import uuid
import logging
import aiohttp_jinja2
from aiohttp import web, hdrs
from .session_storage import SessionStorage
from .user_storage import UserStorage

routes = web.RouteTableDef()
logging.basicConfig(level=logging.DEBUG)

session_storage = SessionStorage()
user_storage = UserStorage()

@routes.route(hdrs.METH_GET, '/auth/logout')
async def login(req: web.Request):
    session_id = req.cookies.get('session_id')
    if session_id is None:
        raise web.HTTPUnauthorized()
    await session_storage.delete_session(session_id)
    return web.Response()


@routes.route(hdrs.METH_GET, '/auth/login')
async def login(req: web.Request):
    return aiohttp_jinja2.render_template('login.html', req, {})


@routes.route(hdrs.METH_POST, '/auth/login')
async def login(req: web.Request):
    data = await req.post()
    username = data.get('username')
    password = data.get('password')

    user = await user_storage.get(username)
    if user is None or not user.check_password(password):
        return web.Response(status=403, text='Username or password is invalid')

    session_id = await session_storage.create_session(user)
    response = web.Response(status=200, text='Logged in')
    response.set_cookie('session_id', str(session_id), httponly=True)
    return response


@routes.route(hdrs.METH_GET, '/auth/register')
async def register(req: web.Request):
    return aiohttp_jinja2.render_template('login.html', req, {})


@routes.route(hdrs.METH_POST, '/auth/register')
async def register(req: web.Request):
    data = await req.post()
    username = data.get('username')
    password = data.get('password')
    logging.info(f"Username = {username}, password = {password}")
    is_created = await user_storage.add(username, password)
    if not is_created:
        return web.Response(status=409, text = "Such login already exists")
    
    return web.Response(status=201, text="Registered")


@routes.route(hdrs.METH_ANY, '/auth')
async def auth(req: web.Request):
    session_id = req.cookies.get('session_id')
    user = await session_storage.get_user_by_session_id(session_id)
    if user is None:
        raise web.HTTPUnauthorized()
    
    logging.info(f"SessionId = {session_id}, user = {user}")

    headers = {"x-user-id": str(user.id)}
    return web.Response(headers=headers)


@routes.route(hdrs.METH_GET, '/auth/backdoor')
async def backdoor(req: web.Request):
    for session_id, user in session_storage._storage.items():
        logging.info(f"Session_id = {session_id}, user = {user}")

    for login, user in user_storage._storage.items():
        logging.info(f"Login = {login}, user = {user}")
    
    return web.Response()