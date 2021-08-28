import logging
import os
import aiohttp_jinja2
import jinja2
from aiohttp import web

from src import forms


if __name__ == '__main__':
    logging.basicConfig(level=logging.DEBUG)
    templates_dir = os.path.join(os.path.abspath(os.path.dirname(__file__)), 'src/templates')

    app = web.Application(debug=True)

    app.add_routes(forms.routes)

    aiohttp_jinja2.setup(app, loader=jinja2.FileSystemLoader(templates_dir))

    web.run_app(app, port=8080)