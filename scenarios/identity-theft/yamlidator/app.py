import yaml, base64, subprocess
from flask import Flask, request
from flask_restful import Resource, Api
from cerberus import Validator

app = Flask(__name__)
api = Api(app)

# cors_header = { 'Access-Control-Allow-Origin': '*',
#                'Access-Control-Allow-Methods': 'GET,POST',
#                'Access-Control-Allow-Headers': 'Content-Type'}

class Yamlidator(Resource):
    def get(self):
        return "Yamlidator Running"
    def options(self):
        return None,200, # cors_header

    def post(self):
        try:
            data = request.form["data"]
        except Exception:
            return "Bad Request", 400, # cors_header

        try:
            b64_data = base64.b64decode(data)
        except Exception:
            return "Failed to decode request", 415, # cors_header

        prep_data = b64_data.decode("utf-8")

        try:
            yaml_data = yaml.unsafe_load(prep_data)
        except Exception:
            return yaml.YAMLError, 500, # cors_header

        schema = eval(open('pod-schema.py', 'r').read())
        v = Validator(schema)

        try:
            check = v.validate(yaml_data, schema)
        except Exception:
            return "Invalid Pod Schema", 400, # cors_header

        return check, 201, # cors_header

class Bekind(Resource):
    def get(self):
        return "Bekind Running"
    def options(self):
        return None,200, # cors_header

    def post(self):
        try:
            data = request.get_json()
        except Exception:
            return "Bad Request", 400, # cors_header

        kind = data["kind"]

        try:
            # run the command
            output = subprocess.check_output(["kubectl", "example", kind])
            decoded = str(output.decode("utf-8"))
        except Exception:
            return "Failed to process example", 500, # cors_header

        return decoded, 202, # cors_header


api.add_resource(Yamlidator, '/api/v1/schema')
api.add_resource(Bekind, '/api/v1/example')

if __name__ == '__main__':
    app.run(debug=True)