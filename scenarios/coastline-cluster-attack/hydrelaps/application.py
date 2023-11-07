#!/usr/bin/python3

from flask import Flask
from flask_rest_jsonapi_next import Api, ResourceList, ResourceDetail
from flask_rest_jsonapi_next.exceptions import ObjectNotFound
from flask_sqlalchemy import SQLAlchemy
from marshmallow_jsonapi.flask import Schema
from marshmallow_jsonapi import fields

# Create a new Flask application
app = Flask(__name__)
app.config['DEBUG'] = True

# Set up SQLAlchemy
app.config['SQLALCHEMY_DATABASE_URI'] = 'sqlite:////host/users.db'
db = SQLAlchemy(app)

# Define a class for the Artist table
class Users(db.Model):
    id = db.Column(db.Integer, primary_key=True)
    ip = db.Column(db.String)
    hostname = db.Column(db.String)
    mac_address = db.Column(db.String)
    user_agent = db.Column(db.String)

# Create the table
with app.app_context():
    db.create_all()

# Logical data abstraction
class UserSchema(Schema):
    class Meta:
        type_ = 'user'
        self_view = 'user_detail'
        self_view_kwargs = {'id': '<id>'}
        self_view_many = 'user_list'

    id = fields.Integer(primary_key=True)
    ip = fields.Str()
    hostname = fields.Str()
    mac_address = fields.Str()
    user_agent = fields.Str()

# Create resource managers
class UserList(ResourceList):
    schema = UserSchema
    data_layer = {'session': db.session,
                   'model': Users}

class UserDetail(ResourceDetail):
    schema = UserSchema
    data_layer = {'session': db.session,
                   'model': Users}

# Create Endpoints
api = Api(app)
api.route(UserList, 'user_list', '/users')
api.route(UserDetail, 'user_detail', '/users/<int:id>')

if __name__ == '__main__':
    app.run(debug=True)