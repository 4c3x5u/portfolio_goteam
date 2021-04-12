from rest_framework.response import Response
from rest_framework.exceptions import ErrorDetail
from ..models import User, Team
import bcrypt

forbidden_response = Response({
    'auth': ErrorDetail(string="Authentication failure.",
                        code='not_authenticated')
}, 403)


def validate_username(username):
    if not username:
        return forbidden_response
    try:
        User.objects.get(username=username)
    except User.DoesNotExist:
        return forbidden_response


def validate_token(token, username, password):
    if not token:
        return forbidden_response
    try:
        tokens_match = bcrypt.checkpw(
            bytes(username, 'utf-8') + password,
            bytes(token, 'utf-8'))
        if not tokens_match:
            return forbidden_response
    except ValueError:
        return forbidden_response


def validate_team_id(team_id):
    if not team_id:
        return Response({
            'team_id': ErrorDetail(string='Team ID cannot be empty.',
                                   code='blank')
        }, 400)
    try:
        Team.objects.get(id=team_id)
    except Team.DoesNotExist:
        return Response({
            'team_id': ErrorDetail(string='Team not found.',
                                   code='not_found')
        }, 404)

